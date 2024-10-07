package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.SignUp(r.Context(), request.Name, request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var response struct {
		ID     uint32 `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		APIKey string `json:"api_key"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, apiKey, err := h.service.SignIn(r.Context(), request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.ID = user.ID
	response.Name = user.Name
	response.Email = user.Email
	response.APIKey = apiKey

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	err := h.service.SignOut(r.Context(), r.Context().Value("apiKey").(*model.APIKey).ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
