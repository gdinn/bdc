// internal/router/router.go
package router

import (
	"net/http"

	"bdc/internal/handlers"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	// Handlers
	oauthHandler *handlers.OAuthHandler
}

func NewRouter(
	oauthHandler *handlers.OAuthHandler,
) *Router {
	return &Router{
		oauthHandler: oauthHandler,
	}
}

func (r *Router) Setup() http.Handler {
	router := chi.NewRouter()

	// Rotas públicas
	router.Route("/api/v1", func(v1 chi.Router) {

		// OAuth/OIDC (Cognito)
		v1.Route("/oauth", func(oauth chi.Router) {
			oauth.Get("/login", r.oauthHandler.Login)
			oauth.Get("/callback", r.oauthHandler.Callback)
			oauth.Get("/logout", r.oauthHandler.Logout)
		})

		// Rotas protegidas
		v1.Group(func(protected chi.Router) {
			// Claims do token (para debug/demonstração)
			protected.Get("/auth/claims", r.oauthHandler.ShowClaims)
		})
	})

	// Arquivos estáticos (se necessário)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Backend Go - v1.0.0"))
	})

	// Rota 404
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Route not found"}`))
	})

	return router
}
