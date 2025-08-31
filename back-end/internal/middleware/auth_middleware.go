package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CognitoJWKS representa a estrutura das chaves públicas do Cognito
type CognitoJWKS struct {
	Keys []CognitoJWK `json:"keys"`
}

type CognitoJWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
}

// UserClaims represents JWT token claims (available properties)
type UserClaims struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Sub      string `json:"sub"`
	TokenUse string `json:"token_use"`
	jwt.RegisteredClaims
}

// AuthMiddleware manages cognito auth
type AuthMiddleware struct {
	cognitoRegion  string
	cognitoPoolID  string
	jwksURL        string
	publicKeys     map[string]*rsa.PublicKey
	lastKeyRefresh time.Time
	keyRefreshTTL  time.Duration
}

type authMiddleWareKey string

const (
	authMiddleWareKeyUserClaims authMiddleWareKey = "user_claims"
)

func NewAuthMiddleware() *AuthMiddleware {
	region := os.Getenv("AWS_REGION")
	poolID := os.Getenv("AWS_COGNITO_USER_POOL_ID")

	if region == "" || poolID == "" {
		panic("AWS_REGION and AWS_COGNITO_USER_POOL_ID must be set")
	}

	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", region, poolID)

	return &AuthMiddleware{
		cognitoRegion: region,
		cognitoPoolID: poolID,
		jwksURL:       jwksURL,
		publicKeys:    make(map[string]*rsa.PublicKey),
		keyRefreshTTL: time.Hour * 24, // Atualizar chaves a cada 24 horas
	}
}

func (am *AuthMiddleware) ValidateToken(tokenString string) (*UserClaims, error) {

	// Extracts token without prefix
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse token without validation do extract kid (key ID)
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify if it is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Extracts kid from header
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		// Gets respective public key
		publicKey, err := am.getPublicKey(kid)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key: %w", err)
		}

		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	if claims.TokenUse != "access" {
		return nil, fmt.Errorf("invalid token use: expected 'access', got '%s'", claims.TokenUse)
	}

	// Validate issuer
	expectedIssuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", am.cognitoRegion, am.cognitoPoolID)
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: expected '%s', got '%s'", expectedIssuer, claims.Issuer)
	}

	return claims, nil
}

// getPublicKey obtains public cognito key by kid
func (am *AuthMiddleware) getPublicKey(kid string) (*rsa.PublicKey, error) {
	// Verificar se precisamos atualizar as chaves
	if time.Since(am.lastKeyRefresh) > am.keyRefreshTTL || am.publicKeys[kid] == nil {
		if err := am.refreshPublicKeys(); err != nil {
			return nil, fmt.Errorf("failed to refresh public keys: %w", err)
		}
	}

	publicKey, exists := am.publicKeys[kid]
	if !exists {
		return nil, fmt.Errorf("public key not found for kid: %s", kid)
	}

	return publicKey, nil
}

func (am *AuthMiddleware) refreshPublicKeys() error {
	resp, err := http.Get(am.jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch JWKS: status code %d", resp.StatusCode)
	}

	var jwks CognitoJWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	// Limpar chaves antigas
	am.publicKeys = make(map[string]*rsa.PublicKey)

	// Converter JWKs para chaves RSA
	for _, jwk := range jwks.Keys {
		if jwk.Kty != "RSA" || jwk.Use != "sig" {
			continue
		}

		publicKey, err := am.jwkToRSAPublicKey(jwk)
		if err != nil {
			continue // Pular chaves inválidas
		}

		am.publicKeys[jwk.Kid] = publicKey
	}

	am.lastKeyRefresh = time.Now()
	return nil
}

func (am *AuthMiddleware) jwkToRSAPublicKey(jwk CognitoJWK) (*rsa.PublicKey, error) {
	// Decode N (modulus)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	// Decode E (exponent)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	// Convert bytes to big.Int
	n := new(big.Int).SetBytes(nBytes)
	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

func (am *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
			return
		}

		claims, err := am.ValidateToken(authHeader)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "Invalid token: %s"}`, err.Error()), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authMiddleWareKeyUserClaims, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserClaimsFromContext extrai as claims do usuário do contexto
func GetUserClaimsFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(authMiddleWareKeyUserClaims).(*UserClaims)
	return claims, ok
}
