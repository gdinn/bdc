package services

import (
	"bdc/internal/domain"
	"bdc/internal/models"
	"bdc/internal/repositories"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

type CognitoService struct {
	cognitoRepository *repositories.CognitoRepository
}

const (
	ErrMsgDeleteUserCognito = "failed to delete user in cognito"
	ErrMsgCreateUserCognito = "failed to create user in cognito"
	ErrMsgUpdateUserCognito = "failed to update user in cognito"
)

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
		return fmt.Errorf("%s: %w", ErrMsgUpdateUserCognito, err)
	}

	// Atualizar no Cognito
	err = cs.cognitoRepository.UpdateUserInCognito(user)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrMsgUpdateUserCognito, err)
	}

	log.Printf("User %s successfully updated in Cognito", user.Email)
	return nil
}

// DeleteUserFromCognito remove usuário do Cognito
func (cs *CognitoService) DeleteUserFromCognito(email string) error {
	if email == "" {
		return fmt.Errorf("%s: %w", ErrMsgDeleteUserCognito, domain.ErrEmailEmpty)
	}

	// Verificar se existe
	_, err := cs.cognitoRepository.GetUserFromCognito(email)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrMsgDeleteUserCognito, err)
	}

	// Deletar
	err = cs.cognitoRepository.DeleteUser(email)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrMsgDeleteUserCognito, err)
	}

	log.Printf("User %s successfully deleted from Cognito", email)
	return nil
}

// DisableUserInCognito desabilita usuário no Cognito
func (cs *CognitoService) DisableUserInCognito(email string) error {
	err := cs.cognitoRepository.DisableUser(email)
	if err != nil {
		return err
	}

	log.Printf("User %s successfully disabled in Cognito", email)
	return nil
}

// EnableUserInCognito habilita usuário no Cognito
func (cs *CognitoService) EnableUserInCognito(email string) error {
	err := cs.cognitoRepository.EnableUser(email)
	if err != nil {
		return err
	}

	log.Printf("User %s successfully enabled in Cognito", email)
	return nil
}

func (cs *CognitoService) GetUserInCognito(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	user, err := cs.cognitoRepository.GetUserFromCognito(username)

	if err != nil {
		return nil, err
	}

	return user, nil
}
