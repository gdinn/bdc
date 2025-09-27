package interfaces

import "bdc/internal/models"

type UserRepositoryInterface interface {
	Create(user *models.User) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	IsEmailExists(email string) (bool, error)
	GetByID(id uint) (*models.User, error)
}
