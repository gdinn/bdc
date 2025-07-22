package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// connectMigrationDatabase conecta usando usuário com privilégios de DDL
func connectMigrationDatabase() (*gorm.DB, error) {
	dsn := buildDSN(
		os.Getenv("PSQL_DB_HOST"),
		os.Getenv("PSQL_DB_PORT"),
		os.Getenv("PSQL_DB_MIGRATIONS_USER"),
		os.Getenv("PSQL_DB_MIGRATIONS_PASSWORD"),
		os.Getenv("PSQL_DB_NAME"),
	)

	return connect(dsn)
}

// connectApplicationDatabase conecta usando usuário da aplicação (operações normais)
func ConnectApplicationDatabase() (*gorm.DB, error) {
	dsn := buildDSN(
		os.Getenv("PSQL_DB_HOST"),
		os.Getenv("PSQL_DB_PORT"),
		os.Getenv("PSQL_DB_USER"),
		os.Getenv("PSQL_DB_PASSWORD"),
		os.Getenv("PSQL_DB_NAME"),
	)

	return connect(dsn)
}

// connect estabelece conexão com configurações padrão
func connect(dsn string) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
			},
		),
	}

	return gorm.Open(postgres.Open(dsn), gormConfig)
}

// buildDSN constrói a string de conexão
func buildDSN(host, port, user, password, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable search_path=bdc",
		host, port, user, password, dbname)
}
