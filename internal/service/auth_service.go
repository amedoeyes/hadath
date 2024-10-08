package service

import (
	"context"

	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/amedoeyes/hadath/internal/utility"
	"github.com/google/uuid"
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

func (s *AuthService) SignUp(ctx context.Context, name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.userRepo.Create(ctx, name, email, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) SignIn(ctx context.Context, email, password string) (*model.User, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", err
	}

	apiKey, err := utility.GenerateAPIKey()
	if err != nil {
		return nil, "", err
	}

	err = s.apiKeyRepo.Create(ctx, user.ID, utility.HashAPIKey(apiKey))
	if err != nil {
		return nil, "", err
	}

	return user, apiKey, err
}

func (s *AuthService) SignOut(ctx context.Context, id uuid.UUID) error {
	return s.apiKeyRepo.DeleteByID(ctx, id)
}
