package service

import (
	"context"

	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) GetById(ctx context.Context, id uint32) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, id uint32) error {
	return s.repo.Delete(ctx, id)
}
