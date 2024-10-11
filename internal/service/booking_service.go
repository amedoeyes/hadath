package service

import (
	"context"

	"github.com/amedoeyes/hadath/internal/dto"
	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/amedoeyes/hadath/internal/validator"
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{
		repo: repo,
	}
}

func (s *BookingService) Create(ctx context.Context, req *dto.BookingRequest) error {
	err := validator.Get().Struct(req)
	if err != nil {
		return err
	}

	user := ctx.Value("user").(*model.User)

	return s.repo.Create(ctx, user.ID, req.EventID)
}

func (s *BookingService) ListByCurrentUser(ctx context.Context) ([]dto.EventResponse, error) {
	user := ctx.Value("user").(*model.User)
	events, err := s.repo.ListByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.EventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, event.ToResponse())
	}

	return response, nil
}

func (s *BookingService) ListByEvent(ctx context.Context) ([]dto.UserResponse, error) {
	event := ctx.Value("event").(*model.Event)
	users, err := s.repo.ListByEvent(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	return response, nil
}

func (s *BookingService) Delete(ctx context.Context, req *dto.BookingRequest) error {
	err := validator.Get().Struct(req)
	if err != nil {
		return err
	}

	user := ctx.Value("user").(*model.User)

	return s.repo.Delete(ctx, user.ID, req.EventID)
}
