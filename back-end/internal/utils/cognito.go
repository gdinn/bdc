package utils

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func GetUserAttribute(user *cognitoidentityprovider.AdminGetUserOutput, attributeName string) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user is nil")
	}

	if user.UserAttributes == nil {
		return "", fmt.Errorf("user attributes is nil")
	}

	for _, attr := range user.UserAttributes {
		if attr.Name != nil && *attr.Name == attributeName {
			if attr.Value != nil {
				return *attr.Value, nil
			}
			return "", fmt.Errorf("attribute '%s' exists but has nil value", attributeName)
		}
	}

	return "", fmt.Errorf("attribute '%s' not found", attributeName)
}
