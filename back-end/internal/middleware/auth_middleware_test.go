package middleware

import (
	"bdc/internal/models"
	"bdc/internal/testutils"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Mock HTTP server para simular o endpoint JWKS do Cognito
func createMockJWKSServer(jwks *models.CognitoJWKS, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		if jwks != nil {
			json.NewEncoder(w).Encode(jwks)
		}
	}))
}

// Gera uma chave RSA para testes
func generateTestRSAKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// Converte chave RSA para formato JWK
func rsaKeyToJWK(pubKey *rsa.PublicKey, kid string) models.CognitoJWK {
	n := base64.RawURLEncoding.EncodeToString(pubKey.N.Bytes())

	// Converter expoente para bytes
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

// Cria um token JWT válido para testes
func createTestJWT(privateKey *rsa.PrivateKey, kid string, claims *models.UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	return token.SignedString(privateKey)
}

func TestNewAuthMiddleware(t *testing.T) {
	// Arrange
	ma := testutils.GetMockAuthData()
	mockJwksURL := testutils.GetJwksUrl(ma)
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)
	t.Setenv("AWS_COGNITO_JWKS_URL", mockJwksURL)

	// Act
	middleware := NewAuthMiddleware()

	// Assert
	if middleware == nil {
		t.Error("Expected middleware to be created, but got nil")
	}

	if middleware.cognitoRegion != ma.AWS_REGION {
		t.Errorf("Expected region to be %s, but got %s", ma.AWS_REGION, middleware.cognitoRegion)
	}

	if middleware.cognitoPoolID != ma.AWS_COGNITO_USER_POOL_ID {
		t.Errorf("Expected poolID to be %s, but got %s", ma.AWS_COGNITO_USER_POOL_ID, middleware.cognitoPoolID)
	}

	if middleware.jwksURL != mockJwksURL {
		t.Errorf("Expected jwksURL to be %s, but got %s", mockJwksURL, middleware.jwksURL)
	}
}

func TestNewAuthMiddleware_MissingEnvironmentVariables(t *testing.T) {
	testutils.UnsetEnvironmentVars()
	defer func() {
		// Recover from panic
		if r := recover(); r == nil {
			t.Error("Expected panic when environment variables are missing")
		}
	}()

	// Act & Assert
	NewAuthMiddleware() // Should panic
}

func TestValidateToken_Success(t *testing.T) {
	// Arrange
	privateKey, err := generateTestRSAKey()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	kid := "test-kid"
	jwk := rsaKeyToJWK(&privateKey.PublicKey, kid)
	jwks := &models.CognitoJWKS{Keys: []models.CognitoJWK{jwk}}

	server := createMockJWKSServer(jwks, http.StatusOK)
	defer server.Close()

	// Set up environment variables
	ma := testutils.GetMockAuthData()
	mockIssuer := testutils.GetIssuerUrl(ma)
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()
	middleware.jwksURL = server.URL

	// Create valid claims
	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "COMMON",
		TokenUse: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: mockIssuer,
		},
	}

	tokenString, err := createTestJWT(privateKey, kid, claims)
	if err != nil {
		t.Fatalf("Failed to create test JWT: %v", err)
	}

	// Act
	result, err := middleware.ValidateToken("Bearer " + tokenString)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if result == nil {
		t.Error("Expected claims to be returned, but got nil")
	}

	if result.Email != "test@example.com" {
		t.Errorf("Expected email to be test@example.com, but got %s", result.Email)
	}
}

func TestValidateToken_InvalidSigningMethod(t *testing.T) {
	// Arrange
	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()

	// Create token with HMAC instead of RSA (invalid signing method)
	claims := &models.UserClaims{
		Email:    "test@example.com",
		TokenUse: "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("secret"))

	// Act
	_, err := middleware.ValidateToken(tokenString)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid signing method, but got nil")
	}

	if !strings.Contains(err.Error(), "unexpected signing method") {
		t.Errorf("Expected error message about signing method, but got: %v", err)
	}
}

func TestValidateToken_MissingKid(t *testing.T) {
	// Arrange
	privateKey, err := generateTestRSAKey()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()

	claims := &models.UserClaims{
		Email:    "test@example.com",
		TokenUse: "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// Not setting kid in header
	tokenString, _ := token.SignedString(privateKey)

	// Act
	_, err = middleware.ValidateToken(tokenString)

	// Assert
	if err == nil {
		t.Error("Expected error for missing kid, but got nil")
	}

	if !strings.Contains(err.Error(), "kid not found") {
		t.Errorf("Expected error message about missing kid, but got: %v", err)
	}
}

func TestValidateToken_InvalidTokenUse(t *testing.T) {
	// Arrange
	privateKey, err := generateTestRSAKey()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	kid := "test-kid"
	jwk := rsaKeyToJWK(&privateKey.PublicKey, kid)
	jwks := &models.CognitoJWKS{Keys: []models.CognitoJWK{jwk}}

	server := createMockJWKSServer(jwks, http.StatusOK)
	defer server.Close()

	ma := testutils.GetMockAuthData()
	mockIssuer := testutils.GetJwksUrl(ma)
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()
	middleware.jwksURL = server.URL

	// Create claims with invalid token_use
	claims := &models.UserClaims{
		Email:    "test@example.com",
		TokenUse: "id", // Should be "access"
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: mockIssuer,
		},
	}

	tokenString, err := createTestJWT(privateKey, kid, claims)
	if err != nil {
		t.Fatalf("Failed to create test JWT: %v", err)
	}

	// Act
	_, err = middleware.ValidateToken(tokenString)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid token use, but got nil")
	}

	if !strings.Contains(err.Error(), "invalid token use") {
		t.Errorf("Expected error message about invalid token use, but got: %v", err)
	}
}

func TestValidateToken_InvalidIssuer(t *testing.T) {
	// Arrange
	privateKey, err := generateTestRSAKey()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	kid := "test-kid"
	jwk := rsaKeyToJWK(&privateKey.PublicKey, kid)
	jwks := &models.CognitoJWKS{Keys: []models.CognitoJWK{jwk}}

	server := createMockJWKSServer(jwks, http.StatusOK)
	defer server.Close()

	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()
	middleware.jwksURL = server.URL

	// Create claims with invalid issuer
	claims := &models.UserClaims{
		Email:    "test@example.com",
		TokenUse: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "https://invalid-issuer.com", // Invalid issuer
		},
	}

	tokenString, err := createTestJWT(privateKey, kid, claims)
	if err != nil {
		t.Fatalf("Failed to create test JWT: %v", err)
	}

	// Act
	_, err = middleware.ValidateToken(tokenString)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid issuer, but got nil")
	}

	if !strings.Contains(err.Error(), "invalid issuer") {
		t.Errorf("Expected error message about invalid issuer, but got: %v", err)
	}
}

func TestRefreshPublicKeys_Success(t *testing.T) {
	// Arrange
	privateKey, err := generateTestRSAKey()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	kid := "test-kid"
	jwk := rsaKeyToJWK(&privateKey.PublicKey, kid)
	jwks := &models.CognitoJWKS{Keys: []models.CognitoJWK{jwk}}

	server := createMockJWKSServer(jwks, http.StatusOK)
	defer server.Close()

	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()
	middleware.jwksURL = server.URL

	// Act
	err = middleware.refreshPublicKeys()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(middleware.publicKeys) != 1 {
		t.Errorf("Expected 1 public key, but got %d", len(middleware.publicKeys))
	}

	if _, exists := middleware.publicKeys[kid]; !exists {
		t.Error("Expected public key to be stored with correct kid")
	}
}

func TestRefreshPublicKeys_HTTPError(t *testing.T) {
	// Arrange
	server := createMockJWKSServer(nil, http.StatusInternalServerError)
	defer server.Close()

	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()
	middleware.jwksURL = server.URL

	// Act
	err := middleware.refreshPublicKeys()

	// Assert
	if err == nil {
		t.Error("Expected error for HTTP error, but got nil")
	}

	if !strings.Contains(err.Error(), "status code 500") {
		t.Errorf("Expected error message about status code, but got: %v", err)
	}
}

func TestRequireAuth_Success(t *testing.T) {
	// Arrange
	privateKey, err := generateTestRSAKey()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	kid := "test-kid"
	jwk := rsaKeyToJWK(&privateKey.PublicKey, kid)
	jwks := &models.CognitoJWKS{Keys: []models.CognitoJWK{jwk}}

	server := createMockJWKSServer(jwks, http.StatusOK)
	defer server.Close()

	ma := testutils.GetMockAuthData()
	mockIssuer := testutils.GetIssuerUrl(ma)
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()
	middleware.jwksURL = server.URL

	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "COMMON",
		TokenUse: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: mockIssuer,
		},
	}

	tokenString, err := createTestJWT(privateKey, kid, claims)
	if err != nil {
		t.Fatalf("Failed to create test JWT: %v", err)
	}

	// Create test handler
	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true

		// Verify claims are in context
		userClaims, ok := GetUserClaimsFromContext(r.Context())
		if !ok {
			t.Error("Expected claims in context, but not found")
		}
		if userClaims.Email != "test@example.com" {
			t.Errorf("Expected email test@example.com, but got %s", userClaims.Email)
		}

		w.WriteHeader(http.StatusOK)
	})

	// Create request with valid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()

	// Act
	middleware.RequireAuth(testHandler).ServeHTTP(w, req)

	// Assert
	if !handlerCalled {
		t.Error("Expected handler to be called, but it wasn't")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, but got %d", w.Code)
	}
}

func TestRequireAuth_MissingAuthHeader(t *testing.T) {
	// Arrange
	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when auth header is missing")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	// Not setting Authorization header
	w := httptest.NewRecorder()

	// Act
	middleware.RequireAuth(testHandler).ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, but got %d", w.Code)
	}

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "Authorization header required" {
		t.Errorf("Expected specific error message, but got: %s", response["error"])
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	// Arrange
	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called when token is invalid")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	// Act
	middleware.RequireAuth(testHandler).ServeHTTP(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, but got %d", w.Code)
	}

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if !strings.Contains(response["error"], "Invalid token") {
		t.Errorf("Expected error about invalid token, but got: %s", response["error"])
	}
}

func TestGetUserClaimsFromContext_Success(t *testing.T) {
	// Arrange
	claims := &models.UserClaims{
		Email:    "test@example.com",
		Username: "testuser",
		Role:     "COMMON",
	}

	ctx := context.WithValue(context.Background(), authMiddleWareKeyUserClaims, claims)

	// Act
	result, ok := GetUserClaimsFromContext(ctx)

	// Assert
	if !ok {
		t.Error("Expected claims to be found in context, but they weren't")
	}

	if result == nil {
		t.Error("Expected claims to be returned, but got nil")
	}

	if result.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, but got %s", result.Email)
	}
}

func TestGetUserClaimsFromContext_NotFound(t *testing.T) {
	// Arrange
	ctx := context.Background() // Empty context

	// Act
	result, ok := GetUserClaimsFromContext(ctx)

	// Assert
	if ok {
		t.Error("Expected claims not to be found in context, but they were")
	}

	if result != nil {
		t.Error("Expected nil claims, but got a value")
	}
}

func TestJwkToRSAPublicKey_Success(t *testing.T) {
	// Arrange
	privateKey, err := generateTestRSAKey()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()
	jwk := rsaKeyToJWK(&privateKey.PublicKey, "test-kid")

	// Act
	publicKey, err := middleware.jwkToRSAPublicKey(jwk)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if publicKey == nil {
		t.Error("Expected public key to be returned, but got nil")
	}

	// Verify the public key matches
	if publicKey.N.Cmp(privateKey.PublicKey.N) != 0 {
		t.Error("Expected modulus to match original key")
	}

	if publicKey.E != privateKey.PublicKey.E {
		t.Error("Expected exponent to match original key")
	}
}

func TestJwkToRSAPublicKey_InvalidBase64(t *testing.T) {
	// Arrange
	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()

	jwk := models.CognitoJWK{
		Kid: "test-kid",
		Kty: "RSA",
		Use: "sig",
		N:   "invalid-base64!", // Invalid base64
		E:   "AQAB",
		Alg: "RS256",
	}

	// Act
	_, err := middleware.jwkToRSAPublicKey(jwk)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid base64, but got nil")
	}

	if !strings.Contains(err.Error(), "failed to decode modulus") {
		t.Errorf("Expected error about modulus decoding, but got: %v", err)
	}
}

func TestKeyRefreshTTL(t *testing.T) {
	// Arrange
	privateKey, err := generateTestRSAKey()
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	kid := "test-kid"
	jwk := rsaKeyToJWK(&privateKey.PublicKey, kid)
	jwks := &models.CognitoJWKS{Keys: []models.CognitoJWK{jwk}}

	server := createMockJWKSServer(jwks, http.StatusOK)
	defer server.Close()

	ma := testutils.GetMockAuthData()
	t.Setenv("AWS_REGION", ma.AWS_REGION)
	t.Setenv("AWS_COGNITO_USER_POOL_ID", ma.AWS_COGNITO_USER_POOL_ID)

	middleware := NewAuthMiddleware()
	middleware.jwksURL = server.URL

	// Set a very short TTL for testing
	middleware.keyRefreshTTL = time.Millisecond * 100

	// Act - First call should fetch keys
	_, err = middleware.getPublicKey(kid)
	if err != nil {
		t.Fatalf("First getPublicKey failed: %v", err)
	}

	// Wait for TTL to expire
	time.Sleep(time.Millisecond * 150)

	// Act - Second call should refresh keys due to TTL
	_, err = middleware.getPublicKey(kid)

	// Assert
	if err != nil {
		t.Errorf("Expected no error after TTL refresh, but got: %v", err)
	}
}
