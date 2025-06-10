package router

type Route struct {
	Method      string
	Pattern     string
	Handler     string
	Middlewares []string
	Description string
}

var Routes = []Route{
	// OAuth routes
	{Method: "GET", Pattern: "/oauth/login", Handler: "OAuthLogin", Description: "OAuth login redirect"},
	{Method: "GET", Pattern: "/oauth/logout", Handler: "OAuthLogout", Description: "OAuth logout redirect"},
	{Method: "GET", Pattern: "/oauth/callback", Handler: "OAuthCallback", Description: "OAuth callback"},
	{Method: "GET", Pattern: "/oauth/claims", Handler: "OAuthShowClaims", Description: "OAuth claims"},
}
