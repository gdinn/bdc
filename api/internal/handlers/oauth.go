package handlers

import (
	"bdc/internal/models"
	"bdc/internal/services"
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type OAuthHandler struct {
	oauthService *services.OAuthService
}

func NewOAuthHandler(oauthService *services.OAuthService) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
	}
}

func (h *OAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Implementação do login OAuth (snippet AWS)
	state := "state" // Replace with a secure random string in production
	url := h.oauthService.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")

	// Exchange the authorization code for a token
	rawToken, err := h.oauthService.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tokenString := rawToken.AccessToken

	// Parse the token (do signature verification for your use case in production)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return
	}

	// Check if the token is valid and extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusBadRequest)
		return
	}

	// Prepare data for rendering the template
	pageData := models.ClaimsPage{
		AccessToken: tokenString,
		Claims:      claims,
	}

	// Define the HTML template
	tmpl := `
    <html>
        <body>
            <h1>User Information</h1>
            <h1>JWT Claims</h1>
            <p><strong>Access Token:</strong> {{.AccessToken}}</p>
            <ul>
                {{range $key, $value := .Claims}}
                    <li><strong>{{$key}}:</strong> {{$value}}</li>
                {{end}}
            </ul>
            <a href="/logout">Logout</a>
        </body>
    </html>`

	// Parse and execute the template
	t := template.Must(template.New("claims").Parse(tmpl))
	t.Execute(w, pageData)
}

func (h *OAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Implementação do logout OAuth
	// Here, you would clear the session or cookie if stored.
	http.Redirect(w, r, "/", http.StatusFound)
}
