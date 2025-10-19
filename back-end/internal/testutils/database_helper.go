package testutils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"bdc/internal/models"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDatabase configura um container PostgreSQL para testes E2E
func SetupTestDatabase(t *testing.T) *gorm.DB {
	ctx := context.Background()

	// 1. Criar container PostgreSQL
	postgresContainer, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "postgres:15",
				ExposedPorts: []string{"5432/tcp"},
				Env: map[string]string{
					"POSTGRES_DB":               "bdc_test",
					"POSTGRES_USER":             "test",
					"POSTGRES_PASSWORD":         "test",
					"POSTGRES_HOST_AUTH_METHOD": "trust", // Para testes, sem senha
				},
				WaitingFor: wait.ForAll(
					// Esperar a segunda ocorrência da mensagem (PostgreSQL emite 2x)
					wait.ForLog("database system is ready to accept connections").
						WithOccurrence(2).
						WithStartupTimeout(20*time.Second),
					// Esperar a porta estar disponível
					wait.ForListeningPort("5432/tcp").
						WithStartupTimeout(20*time.Second),
				),
			},
			Started: true,
		})
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// 2. Adicionar cleanup para parar o container ao final do teste
	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	})

	// 3. Obter host e porta mapeada do container
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	// 4. Construir DSN para conexão
	dsn := fmt.Sprintf(
		"host=%s port=%s user=test password=test dbname=bdc_test sslmode=disable",
		host,
		mappedPort.Port(),
	)

	// 5. Conectar ao banco de dados
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Silenciar logs em testes
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// 6. Executar migrations
	if err := db.AutoMigrate(
		&models.User{},
		&models.Apartment{},
		&models.Vehicle{},
		&models.Pet{},
		&models.Bicycle{},
		&models.UserApartment{},
	); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// 7. Criar schema bdc se necessário
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS bdc").Error; err != nil {
		t.Logf("Warning: Failed to create schema (may already exist): %v", err)
	}

	return db
}

// cleanupTestDatabase limpa todos os dados do banco de teste
func CleanupTestDatabase(t *testing.T, db *gorm.DB) {
	// Truncar todas as tabelas na ordem correta (respeitando foreign keys)
	tables := []string{
		"user_apartments",
		"bicycles",
		"pets",
		"vehicles",
		"apartments",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Logf("Warning: Failed to truncate table %s: %v", table, err)
		}
	}
}

// seedTestData popula o banco com dados de teste
func SeedTestData(t *testing.T, db *gorm.DB) {
	// Exemplo de seed básico
	testUsers := []models.User{
		{
			Name:     "Test User 1",
			Email:    "test1@example.com",
			Phone:    "11999999991",
			Type:     models.UserTypeResident,
			AgeGroup: models.UserAgeGroupAdult,
			Role:     models.UserRoleCommon,
		},
		{
			Name:     "Test Manager",
			Email:    "manager@example.com",
			Phone:    "11999999992",
			Type:     models.UserTypeResident,
			AgeGroup: models.UserAgeGroupAdult,
			Role:     models.UserRoleManager,
		},
	}

	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to seed test user: %v", err)
		}
	}
}
