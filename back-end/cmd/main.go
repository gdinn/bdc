package main

import (
	"log"
	"net/http"

	"bdc/api"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Config routes
	api.ConfigRoutes(router)

	// Start server
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
