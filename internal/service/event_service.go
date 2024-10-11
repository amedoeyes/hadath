package service

import (
	"context"
	"errors"

	"github.com/amedoeyes/hadath/internal/dto"
	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/amedoeyes/hadath/internal/validator"
)

type EventService struct {
	repo *repository.EventRepository
}

func NewEventService(repo *repository.EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Create(ctx context.Context, req *dto.EventRequest) error {
	err := validator.Get().Struct(req)
	if err != nil {
		return err
	}

	user := ctx.Value("user").(*model.User)

	return s.repo.Create(
		ctx,
		user.ID,
		req.Name,
		req.Description,
		req.Address,
		req.StartTime,
		req.EndTime,
	)
}

func (s *EventService) List(ctx context.Context) ([]dto.EventResponse, error) {
	events, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.EventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, event.ToResponse())
	}

	return response, nil
}

func (s *EventService) Get(ctx context.Context) dto.EventResponse {
	event := ctx.Value("event").(*model.Event)
	return event.ToResponse()
}

func (s *EventService) Update(ctx context.Context, req *dto.EventRequest) error {
	err := validator.Get().Struct(req)
	if err != nil {
		return err
	}

	user := ctx.Value("user").(*model.User)
	event := ctx.Value("event").(*model.Event)

	if event.User.ID != user.ID {
		return errors.New("Unauthorized")
	}

	return s.repo.Update(
		ctx,
		event.ID,
		req.Name,
		req.Description,
		req.Address,
		req.StartTime,
		req.EndTime,
	)
}

func (s *EventService) Delete(ctx context.Context) error {
	user := ctx.Value("user").(*model.User)
	event := ctx.Value("event").(*model.Event)

	if event.User.ID != user.ID {
		return errors.New("Unauthorized")
	}

	return s.repo.Delete(ctx, event.ID)
}
