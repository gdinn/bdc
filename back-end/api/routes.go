package api

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"

	"bdc/internal/handlers"
	"bdc/internal/repositories"
	"bdc/internal/services"
)

func SetupRoutes(db *gorm.DB) http.Handler {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repositories.NewUserRepository(db)
	cognitoRepository := repositories.NewCognitoRepository(cfg)

	userService := services.NewUserService(userRepo)
	cognitoService := services.NewCognitoService(cognitoRepository)

	userHandler := handlers.NewUserHandler(userService, cognitoService)

	r := mux.NewRouter()

	r.Use(loggingMiddleware)

	api := r.PathPrefix("/api/v1").Subrouter()

	// ====================
	// PUBLIC ROUTES - No auth
	// ====================
	public := api.PathPrefix("/public").Subrouter()
	public.HandleFunc("/health", healthCheckHandler).Methods("GET")
	public.HandleFunc("/version", versionHandler).Methods("GET")

	// ====================
	// PROTECTED ROUTES
	// ====================
	api.HandleFunc("/users", userHandler.CreateUserWithAuth()).Methods("POST")

	// ====================
	// CORS CONFIGURATION
	// ====================
	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:4200", // Angular dev server
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	})

	return c.Handler(r)
}

// loggingMiddleware adds loging for all reqs
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// healthCheckHandler verify if the API is up
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "service": "BDC API"}`))
}

// versionHandler returns api version
func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"version": "1.0.0", "service": "BDC API"}`))
}
