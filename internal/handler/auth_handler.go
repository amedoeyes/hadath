package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/amedoeyes/hadath/internal/dto"
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
	var req dto.AuthSignUpRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	err = h.service.SignUp(r.Context(), &req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthSignInRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	res, err := h.service.SignIn(r.Context(), &req)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	err := h.service.SignOut(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
