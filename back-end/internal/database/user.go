package database

import (
	"bdc/internal/models"
	"database/sql"
	"fmt"
)

func SaveUser(user *models.User) error {
	db, err := sql.Open("postgres", "user=seu-usuario password=sua-senha dbname=seu-banco sslmode=disable")
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	query := `
		INSERT INTO users (name, email, phone, birth_date, type, age_group, is_manager, is_advisor, is_legal_rep)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	err = db.QueryRow(query, user.Name, user.Email, user.Phone, user.BirthDate, user.Type, user.AgeGroup, user.IsManager, user.IsAdvisor, user.IsLegalRep).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("erro ao inserir usuário no banco de dados: %v", err)
	}

	return nil
}
