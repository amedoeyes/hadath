package service

import (
	"context"

	"github.com/amedoeyes/hadath/internal/api"
	"github.com/amedoeyes/hadath/internal/api/request"
	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/amedoeyes/hadath/internal/database/repository"
	"github.com/amedoeyes/hadath/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	apiKeyRepo *repository.APIKeyRepository
	userRepo   *repository.UserRepository
}

func NewAuthService(apiKeyRepo *repository.APIKeyRepository, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		apiKeyRepo: apiKeyRepo,
		userRepo:   userRepo,
	}
}

func (s *AuthService) SignUp(ctx context.Context, req *request.AuthSignUpRequest) error {
	err := validator.Get().Struct(req)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.userRepo.Create(ctx, req.Name, req.Email, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SignIn(ctx context.Context, req *request.AuthSignInRequest) (*response.AuthResponse, error) {
	err := validator.Get().Struct(req)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, api.ErrUnauthorized
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	key, err := GenerateAPIKey()
	if err != nil {
		return nil, err
	}

	err = s.apiKeyRepo.Create(ctx, user.ID, HashAPIKey(key))
	if err != nil {
		return nil, err
	}

	response := &response.AuthResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		APIKey: key,
	}

	return response, err
}

func (s *AuthService) SignOut(ctx context.Context) error {
	apiKey := ctx.Value("apiKey").(*model.APIKey)
	return s.apiKeyRepo.Delete(ctx, apiKey.ID)
}
