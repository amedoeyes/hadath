package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amedoeyes/hadath/internal/api/request"
	"github.com/amedoeyes/hadath/internal/database/repository"
	"github.com/amedoeyes/hadath/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler(t *testing.T) {
	apiKeyRepo := repository.NewAPIKeyRepository()
	userRepo := repository.NewUserRepository()
	authService := service.NewAuthService(apiKeyRepo, userRepo)
	handler := NewAuthHandler(authService)

	validSignUpRequest := request.AuthSignUpRequest{
		Name:     "Cat",
		Email:    "cat@meow.purr",
		Password: "meowmeowmeow",
	}

	validSignInRequest := request.AuthSignInRequest{
		Email:    "cat@meow.purr",
		Password: "meowmeowmeow",
	}

	r := chi.NewRouter()
	r.Post("/signup", handler.SignUp)
	r.Post("/signin", handler.SignIn)
	r.Post("/signout", handler.SignOut)

	t.Run("SignUp", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()

		tests := []struct {
			name       string
			request    request.AuthSignUpRequest
			wantStatus int
			wantErr    bool
		}{
			{"ValidRequest", validSignUpRequest, http.StatusCreated, false},
			{"EmptyName", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Name = ""; return r }(), http.StatusBadRequest, true},
			{"EmptyEmail", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Email = ""; return r }(), http.StatusBadRequest, true},
			{"EmptyPassword", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Password = ""; return r }(), http.StatusBadRequest, true},
			{"InvalidEmail", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Email = "catmeow.purr"; return r }(), http.StatusBadRequest, true},
			{"ShortPassword", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Password = "meow"; return r }(), http.StatusBadRequest, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				reqBody, _ := json.Marshal(tt.request)
				req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
				req = req.WithContext(ctx)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				r.ServeHTTP(rr, req)

				assert.Equal(t, tt.wantStatus, rr.Code)
				if tt.wantErr {
					assert.NotEqual(t, http.StatusCreated, rr.Code)
				} else {
					assert.Equal(t, http.StatusCreated, rr.Code)
				}
			})
		}
	})

	t.Run("SignIn", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()

		err := authService.SignUp(ctx, &validSignUpRequest)
		assert.NoError(t, err)

		tests := []struct {
			name       string
			request    request.AuthSignInRequest
			wantStatus int
			wantErr    bool
		}{
			{"ValidRequest", validSignInRequest, http.StatusOK, false},
			{"EmptyEmail", func() request.AuthSignInRequest { r := validSignInRequest; r.Email = ""; return r }(), http.StatusBadRequest, true},
			{"EmptyPassword", func() request.AuthSignInRequest { r := validSignInRequest; r.Password = ""; return r }(), http.StatusBadRequest, true},
			{"WrongEmail", func() request.AuthSignInRequest { r := validSignInRequest; r.Email = "wrong@email.com"; return r }(), http.StatusUnauthorized, true},
			{"WrongPassword", func() request.AuthSignInRequest { r := validSignInRequest; r.Password = "wrongpassword"; return r }(), http.StatusUnauthorized, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				reqBody, _ := json.Marshal(tt.request)
				req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))
				req = req.WithContext(ctx)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				r.ServeHTTP(rr, req)

				assert.Equal(t, tt.wantStatus, rr.Code)
				if tt.wantErr {
					assert.NotEqual(t, http.StatusOK, rr.Code)
				} else {
					assert.Equal(t, http.StatusOK, rr.Code)

					var responseBody map[string]interface{}
					err := json.NewDecoder(rr.Body).Decode(&responseBody)
					assert.NoError(t, err)

					assert.Equal(t, validSignUpRequest.Email, responseBody["email"])
				}
			})
		}
	})

	t.Run("SignOut", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()

		err := authService.SignUp(ctx, &validSignUpRequest)
		assert.NoError(t, err)

		user, err := authService.SignIn(ctx, &validSignInRequest)
		assert.NoError(t, err)

		apiKey, err := apiKeyRepo.GetByKey(ctx, service.HashAPIKey(user.APIKey))
		assert.NoError(t, err)

		ctx = context.WithValue(ctx, "apiKey", apiKey)

		req := httptest.NewRequest(http.MethodPost, "/signout", nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

		_, err = apiKeyRepo.GetByKey(ctx, service.HashAPIKey(user.APIKey))
		assert.Error(t, err, "API key should not exist after sign out")
	})
}
