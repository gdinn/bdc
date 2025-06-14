package models

import "github.com/golang-jwt/jwt"

type ClaimsPage struct {
	AccessToken string
	Claims      jwt.MapClaims
}
