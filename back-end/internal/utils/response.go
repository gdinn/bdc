package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func SendErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

func IsEmailAlreadyExistsError(err error) bool {
	return strings.Contains(err.Error(), "email already exists")
}
