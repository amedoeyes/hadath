package service

import (
	"context"

	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/google/uuid"
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{
		repo: repo,
	}
}

func (s *BookingService) Create(ctx context.Context, userID uuid.UUID, eventID uuid.UUID) error {
	return s.repo.Create(ctx, userID, eventID)
}

func (s *BookingService) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]model.Booking, error) {
	return s.repo.GetAllByUserID(ctx, userID)
}

func (s *BookingService) GetAllByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Booking, error) {
	return s.repo.GetAllByEventID(ctx, eventID)
}

func (s *BookingService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteByID(ctx, id)
}
