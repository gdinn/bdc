package api

import (
	"encoding/json"
	"net/http"

	"bdc/internal/database"
	"bdc/internal/models"

	"github.com/gorilla/mux"
)

func CreateProfileHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição: "+err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" {
		http.Error(w, "Nome e Email são campos obrigatórios", http.StatusBadRequest)
		return
	}

	err = database.SaveUser(&user)
	if err != nil {
		http.Error(w, "Erro ao salvar o usuário no banco de dados: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Perfil criado com sucesso",
	})
}

func ConfigRoutes(router *mux.Router) {
	router.HandleFunc("/profile", CreateProfileHandler).Methods("POST")
}
