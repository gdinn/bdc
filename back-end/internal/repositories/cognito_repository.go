package repositories

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"

	"bdc/internal/models"
)

type CognitoRepository struct {
	client     *cognitoidentityprovider.Client
	userPoolId string
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
			Name:  aws.String("custom:role"),
			Value: aws.String(string(user.Role)),
		},
	}

	// Adiciona telefone se fornecido
	if user.Phone != "" {
		userAttributes = append(userAttributes, types.AttributeType{
			Name:  aws.String("phone_number"),
			Value: aws.String(user.Phone),
		})
	}

	input := &cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserPoolId:     aws.String(cr.userPoolId),
		Username:       aws.String(user.Email),
		UserAttributes: userAttributes,
	}

	_, err := cr.client.AdminUpdateUserAttributes(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to update user in cognito: %w", err)
	}

	return nil
}

// GetUserFromCognito busca um usuário no Cognito pelo email
func (cr *CognitoRepository) GetUserFromCognito(email string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(cr.userPoolId),
		Username:   aws.String(email),
	}

	result, err := cr.client.AdminGetUser(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from cognito: %w", err)
	}

	return result, nil
}

// stringPtr retorna um ponteiro para string
func stringPtr(s string) *string {
	return &s
}
