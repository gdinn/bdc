package repositories

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"bdc/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) IsEmailExists(email string) (bool, error) {
	if email == "" {
		return false, fmt.Errorf("email cannot be empty")
	}

	email = strings.TrimSpace(strings.ToLower(email))

	var count int64
	err := r.db.Model(&models.User{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}

	return count > 0, nil
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("deleted_at IS NULL").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
