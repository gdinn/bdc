package repositories

import (
	"bdc/internal/models"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoRepository struct {
	client     *cognitoidentityprovider.Client
	userPoolId string
	timeout    time.Duration
}

// NewCognitoRepository cria uma nova instância do repositório Cognito
func NewCognitoRepository(cfg aws.Config) *CognitoRepository {
	return &CognitoRepository{
		client:     cognitoidentityprovider.NewFromConfig(cfg),
		userPoolId: os.Getenv("AWS_COGNITO_USER_POOL_ID"),
	}
}

// UpdateUserInCognito atualiza atributos do usuário no Cognito
func (cr *CognitoRepository) UpdateUserInCognito(user *models.User) error {
	userAttributes := []types.AttributeType{
		{
			Name:  aws.String("name"),
			Value: aws.String(user.Name),
		},
		{
			Name:  aws.String("email"), // Verificar como é feita a confirmação desse email depois...
			Value: aws.String(user.Email),
		},
		{
			Name:  aws.String("custom:role"),
			Value: aws.String(string(user.Role)),
		},
	}

	input := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId:     aws.String(cr.userPoolId),
		Username:       aws.String(user.Email),
		UserAttributes: userAttributes,
	}

	_, err := cr.client.AdminUpdateUserAttributes(context.Background(), input)
	if err != nil {
		return fmt.Errorf("CognitoRepository: failed to update user in cognito: %w", err)
	}

	return nil
}

// GetUserFromCognito busca um usuário no Cognito pelo username
func (cr *CognitoRepository) GetUserFromCognito(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(cr.userPoolId),
		Username:   aws.String(username),
	}

	result, err := cr.client.AdminGetUser(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from cognito: %w", err)
	}

	return result, nil
}

// DeleteUser exclui um usuário do Cognito
func (cr *CognitoRepository) DeleteUser(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), cr.timeout)
	defer cancel()

	input := &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(cr.userPoolId),
		Username:   aws.String(email),
	}

	_, err := cr.client.AdminDeleteUser(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete user from cognito: %w", err)
	}

	return nil
}

// DisableUser desabilita um usuário no Cognito
func (cr *CognitoRepository) DisableUser(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), cr.timeout)
	defer cancel()

	input := &cognitoidentityprovider.AdminDisableUserInput{
		UserPoolId: aws.String(cr.userPoolId),
		Username:   aws.String(email),
	}

	_, err := cr.client.AdminDisableUser(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to disable user in cognito: %w", err)
	}

	return nil
}

// EnableUser habilita um usuário no Cognito
func (cr *CognitoRepository) EnableUser(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), cr.timeout)
	defer cancel()

	input := &cognitoidentityprovider.AdminEnableUserInput{
		UserPoolId: aws.String(cr.userPoolId),
		Username:   aws.String(email),
	}

	_, err := cr.client.AdminEnableUser(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to enable user in cognito: %w", err)
	}

	return nil
}
