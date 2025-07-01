package config

import (
	"os"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	IssuerURL    string
	Scopes       []string
}

func LoadOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		IssuerURL:    os.Getenv("OAUTH_ISSUER_URL"),
		Scopes:       []string{"openid", "phone", "email"},
	}
}
