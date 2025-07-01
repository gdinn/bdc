package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bdc/internal/config"
	"bdc/internal/handlers"
	"bdc/internal/router"
	"bdc/internal/services"
)

func main() {
	// Carregar configurações
	oauthCfg := config.LoadOAuthConfig()

	// Inicializar serviços
	oauthService, err := services.NewOAuthService(oauthCfg)
	if err != nil {
		log.Fatalf("Failed to initialize OAuth service: %v", err)
	}

	// Inicializar handlers
	oauthHandler := handlers.NewOAuthHandler(oauthService)

	// Configurar rotas
	r := router.NewRouter(
		oauthHandler,
	)

	// Criar servidor HTTP
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", "8080"),
		Handler:      r.Setup(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Iniciar servidor em goroutine
	go func() {
		log.Print("Starting server", "port", 8080)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Print("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Print("Server exited")
}
