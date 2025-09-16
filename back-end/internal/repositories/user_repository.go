package repositories

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"bdc/internal/domain"
	"bdc/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

const (
	ErrCheckingEmailExistence = "error checking email existence"
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) IsEmailExists(email string) (bool, error) {
	if email == "" {
		return false, fmt.Errorf("%s: %w", ErrCheckingEmailExistence, domain.ErrEmailEmpty)
	}

	email = strings.TrimSpace(strings.ToLower(email))

	var count int64
	err := r.db.Model(&models.User{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("%s: %w", ErrCheckingEmailExistence, err)
	}

	return count > 0, nil
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("deleted_at IS NULL").First(&user, id).Error
	if err != nil {
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}
	return &user, nil
}
