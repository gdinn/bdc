package utils

import (
	"bdc/internal/domain"
	"bdc/internal/middleware"
	"bdc/internal/models"
)

func ValidateManagerRole(claims *middleware.UserClaims) error {
	if claims.Role != string(models.UserRoleManager) {
		return domain.ErrManagerRoleRequired
	}
	return nil
}

func ValidateAdvisorRole(claims *middleware.UserClaims) error {
	if claims.Role != string(models.UserRoleAdvisor) {
		return domain.ErrAdvisorRoleRequired
	}
	return nil
}

func ValidateManagerOrAdvisorRole(claims *middleware.UserClaims) error {
	if claims.Role != string(models.UserRoleManager) && (claims.Role != string(models.UserRoleAdvisor)) {
		return domain.ErrManagerOrAdvisorRoleRequired
	}
	return nil
}
