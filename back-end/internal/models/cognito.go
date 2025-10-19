package models

import "github.com/golang-jwt/jwt/v5"

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
