package database

import (
	"bdc/internal/models"
	"fmt"
	"log"
)

// runMigrations executa as migrations usando usuário privilegiado
func RunMigrations() error {
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
