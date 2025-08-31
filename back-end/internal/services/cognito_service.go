package services

import (
	"bdc/internal/models"
	"bdc/internal/repositories"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type CognitoService struct {
	cognitoRepository *repositories.CognitoRepository
}

func NewCognitoService(cognitoRepo *repositories.CognitoRepository) *CognitoService {
	return &CognitoService{
		cognitoRepository: cognitoRepo,
	}
}

// UpdateUserInCognito atualiza atributos do usuário no Cognito com sincronização
func (cs *CognitoService) UpdateUserInCognito(user *models.User) error {
	// Verificar se o usuário existe no Cognito primeiro
	_, err := cs.cognitoRepository.GetUserFromCognito(user.Email)
	if err != nil {
		return fmt.Errorf("failed to update user in Cognito. User does not exist: %w", err)
	}

	// Atualizar no Cognito
	err = cs.cognitoRepository.UpdateUserInCognito(user)
	if err != nil {
		return fmt.Errorf("failed to update user in Cognito: %w", err)
	}

	log.Printf("User %s successfully updated in Cognito", user.Email)
	return nil
}

// DeleteUserFromCognito remove usuário do Cognito
func (cs *CognitoService) DeleteUserFromCognito(email string) error {
	// Verificar se existe
	_, err := cs.cognitoRepository.GetUserFromCognito(email)
	if err != nil {
		return fmt.Errorf("user not found in Cognito: %w", err)
	}

	// Deletar
	err = cs.cognitoRepository.DeleteUser(email)
	if err != nil {
		return fmt.Errorf("failed to delete user from Cognito: %w", err)
	}

	log.Printf("User %s successfully deleted from Cognito", email)
	return nil
}

// DisableUserInCognito desabilita usuário no Cognito
func (cs *CognitoService) DisableUserInCognito(email string) error {
	err := cs.cognitoRepository.DisableUser(email)
	if err != nil {
		return fmt.Errorf("failed to disable user in Cognito: %w", err)
	}

	log.Printf("User %s successfully disabled in Cognito", email)
	return nil
}

// EnableUserInCognito habilita usuário no Cognito
func (cs *CognitoService) EnableUserInCognito(email string) error {
	err := cs.cognitoRepository.EnableUser(email)
	if err != nil {
		return fmt.Errorf("failed to enable user in Cognito: %w", err)
	}

	log.Printf("User %s successfully enabled in Cognito", email)
	return nil
}

func (cs *CognitoService) GetUserInCognito(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	user, err := cs.cognitoRepository.GetUserFromCognito(username)

	if err != nil {
		return nil, fmt.Errorf("failed to get user in Cognito: %w", err)
	}

	return user, nil
}
