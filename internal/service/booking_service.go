package service

import (
	"context"

	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/amedoeyes/hadath/internal/database/repository"
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{
		repo: repo,
	}
}

func (s *BookingService) Create(ctx context.Context) error {
	user := ctx.Value("user").(*model.User)
	event := ctx.Value("event").(*model.Event)

	return s.repo.Create(ctx, user.ID, event.ID)
}

func (s *BookingService) ListByCurrentUser(ctx context.Context) ([]response.EventResponse, error) {
	user := ctx.Value("user").(*model.User)
	events, err := s.repo.ListByUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	response := make([]response.EventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, event.ToResponse())
	}

	return response, nil
}

func (s *BookingService) ListByEvent(ctx context.Context) ([]response.UserResponse, error) {
	event := ctx.Value("event").(*model.Event)
	users, err := s.repo.ListByEvent(ctx, event.ID)
	if err != nil {
		return nil, err
	}

	response := make([]response.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, user.ToResponse())
	}

	return response, nil
}

func (s *BookingService) Delete(ctx context.Context) error {
	user := ctx.Value("user").(*model.User)
	event := ctx.Value("event").(*model.Event)

	return s.repo.Delete(ctx, user.ID, event.ID)
}
