package main

import (
	"log"
	"net/http"

	"bdc/api"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Config routes
	api.ConfigRoutes(router)

	_, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration, %v", err)
	}

	// Start server
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
