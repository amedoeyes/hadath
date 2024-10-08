package service

import (
	"context"
	"time"

	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/google/uuid"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Create(ctx context.Context, user_id uuid.UUID, name, description, address string, startTime, endTime time.Time) error {
	return s.repo.Create(ctx, user_id, name, description, address, startTime, endTime)
}

func (s *EventService) GetAll(ctx context.Context) ([]model.Event, error) {
	return s.repo.GetAll(ctx)
}

func (s *EventService) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EventService) UpdateByID(ctx context.Context, id uuid.UUID, name, description, address string, startTime, endTime time.Time) error {
	return s.repo.UpdateByID(ctx, id, name, description, address, startTime, endTime)
}

func (s *EventService) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteByID(ctx, id)
}
