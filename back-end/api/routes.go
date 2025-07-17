package api

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"bdc/internal/handlers"
	"bdc/internal/repositories"
	"bdc/internal/services"
)

func SetupRoutes(db *gorm.DB) *mux.Router {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Setup router
	router := mux.NewRouter()

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// User routes
	userRoutes := api.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("", userHandler.CreateUser).Methods("POST")

	return router
}
