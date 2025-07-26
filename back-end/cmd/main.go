package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"bdc/api"
	"bdc/internal/database"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
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

	// TESTE
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	cred, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		print("deu ruim", err.Error())
	}
	print(cred.SecretAccessKey)

	updateInfo := cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserAttributes: []types.AttributeType{
			{Name: stringPtr("custom:role"), Value: stringPtr("EXTERNAL_USER")},
		},
		UserPoolId: stringPtr(""),
		Username:   stringPtr(""),
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)

	output, err := client.AdminUpdateUserAttributes(context.Background(), &updateInfo)
	if err != nil {
		// Tratar o erro
		log.Printf("Erro ao atualizar atributos: %v", err)
		print(output)
		return
	}

	// Fim teste

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func stringPtr(s string) *string {
	return &s
}
