package testutils

import (
	"bdc/internal/models"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// setupMockJWKSServer cria servidor mock para JWKS do Cognito
func SetupMockJWKSServer(t *testing.T) (*rsa.PrivateKey, *httptest.Server) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	jwk := RsaKeyToJWK(&privateKey.PublicKey, "test-kid")
	jwks := &models.CognitoJWKS{
		Keys: []models.CognitoJWK{jwk},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	}))

	// Atualizar JWKS URL no middleware (via variável de ambiente ou config)
	// Nota: Em testes reais, você precisaria de uma forma de injetar a URL do mock
	t.Setenv("AWS_COGNITO_JWKS_URL", server.URL)

	return privateKey, server
}

// rsaKeyToJWK converte chave RSA para formato JWK
func RsaKeyToJWK(pubKey *rsa.PublicKey, kid string) models.CognitoJWK {
	n := base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes())
	eBytes := big.NewInt(int64(pubKey.E)).Bytes()
	e := base64.RawURLEncoding.EncodeToString(eBytes)

	return models.CognitoJWK{
		Kid: kid,
		Kty: "RSA",
		Use: "sig",
		N:   n,
		E:   e,
		Alg: "RS256",
	}
}

// createValidToken cria um token JWT válido para testes
func CreateValidToken(t *testing.T, privateKey *rsa.PrivateKey, email, username string, role models.UserRole) (string, *models.UserClaims) {
	claims := &models.UserClaims{ // cycling
		Email:    email,
		Username: username,
		Role:     string(role),
		Sub:      "user-sub-" + username,
		TokenUse: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    GetIssuerUrl(GetMockAuthData()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-kid"

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	return tokenString, claims
}

func UnsetEnvironmentVars() {
	os.Unsetenv("AWS_REGION")
	os.Unsetenv("AWS_COGNITO_USER_POOL_ID")
	os.Unsetenv("AWS_COGNITO_JWKS_URL")
}
