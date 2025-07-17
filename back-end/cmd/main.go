package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gorm.io/gorm/logger"

	"bdc/api"
	"bdc/internal/models"
)

func main() {
	// 1. Executar migrations com usuário privilegiado
	if err := runMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// 2. Conectar com usuário da aplicação para operações normais
	db, err := connectApplicationDatabase()
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

// runMigrations executa as migrations usando usuário privilegiado
func runMigrations() error {
	log.Println("Running database migrations...")

	// Conectar com usuário de migration
	migrationDB, err := connectMigrationDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect migration database: %w", err)
	}

	// Executar migrations
	if err := migrationDB.AutoMigrate(
		&models.User{},
		&models.Apartment{},
		&models.Vehicle{},
		&models.Pet{},
		&models.Bicycle{},
		&models.UserApartment{},
	); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// connectMigrationDatabase conecta usando usuário com privilégios de DDL
func connectMigrationDatabase() (*gorm.DB, error) {
	host := os.Getenv("PSQL_DB_HOST")
	port := os.Getenv("PSQL_DB_PORT")
	user := os.Getenv("PSQL_DB_MIGRATIONS_USER")
	password := os.Getenv("PSQL_DB_MIGRATIONS_PASSWORD")
	dbname := os.Getenv("PSQL_DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=bdc",
		host, port, user, password, dbname)

	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info, // Mostra todas as queries SQL
			},
		),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectApplicationDatabase conecta usando usuário da aplicação (operações normais)
func connectApplicationDatabase() (*gorm.DB, error) {
	host := os.Getenv("PSQL_DB_HOST")
	port := os.Getenv("PSQL_DB_PORT")
	user := os.Getenv("PSQL_DB_USER")
	password := os.Getenv("PSQL_DB_PASSWORD")
	dbname := os.Getenv("PSQL_DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=bdc",
		host, port, user, password, dbname)

	// Configuração do logger do GORM (opcional)
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info, // Mostra todas as queries SQL
			},
		),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
