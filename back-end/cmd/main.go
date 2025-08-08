package main

import (
	"log"
	"net/http"
	"os"

	"bdc/api"
	"bdc/internal/database"
)

func main() {
	// 1. Executar migrations com usuário privilegiado
	if err := database.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 2. Conectar com usuário da aplicação para operações normais
	db, err := database.ConnectApplicationDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to application database: %v", err)
	}

	// Setup routes
	router := api.SetupRoutes(db)

	// Start server
	port := os.Getenv("PSQL_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
