package services

import (
	"bdc/internal/config"
	"context"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	provider     *oidc.Provider
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func (s *OAuthService) AuthCodeURL(state string) string {
	return s.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *OAuthService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauth2Config.Exchange(ctx, code)
}

func NewOAuthService(cfg *config.OAuthConfig) (*OAuthService, error) {
	provider, err := oidc.NewProvider(context.Background(), cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       cfg.Scopes,
	}

	verifier := provider.Verifier(&oidc.Config{ // Não está na config exemplo da aws..
		ClientID: cfg.ClientID,
	})

	return &OAuthService{
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
	}, nil
}
