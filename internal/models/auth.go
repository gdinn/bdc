package models

import "github.com/golang-jwt/jwt/v4"

type ClaimsPage struct {
	AccessToken string
	Claims      jwt.MapClaims
}
