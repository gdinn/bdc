package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"bdc/api"
	"bdc/internal/database"
	"bdc/internal/repositories"
	"bdc/internal/services"

	"github.com/aws/aws-sdk-go-v2/config"
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

	// TESTE
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	cognitoRepository := repositories.NewCognitoRepository(cfg)

	_ = services.NewCognitoService(cognitoRepository)

	// FIM TESTE

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
