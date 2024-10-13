package service

import (
	"context"

	"github.com/amedoeyes/hadath/internal/api"
	"github.com/amedoeyes/hadath/internal/api/request"
	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/amedoeyes/hadath/internal/database/repository"
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

func (s *EventService) Create(ctx context.Context, req *request.EventRequest) error {
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

func (s *EventService) List(ctx context.Context) ([]response.EventResponse, error) {
	events, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]response.EventResponse, 0, len(events))
	for _, event := range events {
		response = append(response, event.ToResponse())
	}

	return response, nil
}

func (s *EventService) Get(ctx context.Context) response.EventResponse {
	event := ctx.Value("event").(*model.Event)
	return event.ToResponse()
}

func (s *EventService) Update(ctx context.Context, req *request.EventRequest) error {
	err := validator.Get().Struct(req)
	if err != nil {
		return err
	}

	user := ctx.Value("user").(*model.User)
	event := ctx.Value("event").(*model.Event)

	if event.User.ID != user.ID {
		return api.ErrUnauthorized
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
		return api.ErrUnauthorized
	}

	return s.repo.Delete(ctx, event.ID)
}
