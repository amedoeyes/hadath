package service

import (
	"context"
	"testing"

	"github.com/amedoeyes/hadath/internal/api/request"
	"github.com/amedoeyes/hadath/internal/database/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuthService(t *testing.T) {
	apiKeyRepo := repository.NewAPIKeyRepository()
	userRepo := repository.NewUserRepository()
	service := NewAuthService(apiKeyRepo, userRepo)

	validSignUpRequest := request.AuthSignUpRequest{
		Name:     "Cat",
		Email:    "cat@meow.purr",
		Password: "meowmeowmeow",
	}

	validSignInRequest := request.AuthSignInRequest{
		Email:    "cat@meow.purr",
		Password: "meowmeowmeow",
	}

	t.Run("SignUp", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()

		tests := []struct {
			name    string
			request request.AuthSignUpRequest
			wantErr bool
		}{
			{"ValidRequest", validSignUpRequest, false},
			{"EmptyName", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Name = ""; return r }(), true},
			{"EmptyEmail", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Email = ""; return r }(), true},
			{"EmptyPassword", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Password = ""; return r }(), true},
			{"InvalidEmail", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Email = "catmeow.purr"; return r }(), true},
			{"ShortPassword", func() request.AuthSignUpRequest { r := validSignUpRequest; r.Password = "meow"; return r }(), true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := service.SignUp(ctx, &tt.request)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("SignIn", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()

		err := service.SignUp(ctx, &validSignUpRequest)
		assert.NoError(t, err)

		tests := []struct {
			name    string
			request request.AuthSignInRequest
			wantErr bool
		}{
			{"ValidRequest", validSignInRequest, false},
			{"EmptyEmail", func() request.AuthSignInRequest { r := validSignInRequest; r.Email = ""; return r }(), true},
			{"EmptyPassword", func() request.AuthSignInRequest { r := validSignInRequest; r.Password = ""; return r }(), true},
			{"WrongEmail", func() request.AuthSignInRequest { r := validSignInRequest; r.Email = "wrong@email.com"; return r }(), true},
			{"WrongPassword", func() request.AuthSignInRequest { r := validSignInRequest; r.Password = "wrongpassword"; return r }(), true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				user, err := service.SignIn(ctx, &tt.request)
				if tt.wantErr {
					assert.Error(t, err)
					assert.Nil(t, user)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, validSignUpRequest.Name, user.Name)
					assert.Equal(t, validSignUpRequest.Email, user.Email)
				}
			})
		}
	})

	t.Run("SignOut", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()

		signUpReq := validSignUpRequest
		err := service.SignUp(ctx, &signUpReq)
		assert.NoError(t, err)

		user, err := service.SignIn(ctx, &request.AuthSignInRequest{Email: signUpReq.Email, Password: signUpReq.Password})
		assert.NoError(t, err)

		apiKey, err := apiKeyRepo.GetByKey(ctx, HashAPIKey(user.APIKey))
		assert.NoError(t, err)

		ctx = context.WithValue(ctx, "apiKey", apiKey)

		err = service.SignOut(ctx)
		assert.NoError(t, err)

		_, err = apiKeyRepo.GetByKey(ctx, HashAPIKey(user.APIKey))
		assert.Error(t, err, "API key should not exist after sign out")
	})
}
