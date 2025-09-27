package interfaces

import (
	"bdc/internal/models"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type CognitoRepositoryInterface interface {
	UpdateUserInCognito(user *models.User) error
	GetUserFromCognito(username string) (*cognitoidentityprovider.AdminGetUserOutput, error)
	DeleteUser(email string) error
	DisableUser(email string) error
	EnableUser(email string) error
}
