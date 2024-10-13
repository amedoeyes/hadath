package model

import (
	"time"

	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	User        User
	Name        string
	Description string
	Address     string
	StartTime   time.Time
	EndTime     time.Time
}

func (e *Event) ToResponse() response.EventResponse {
	return response.EventResponse{
		ID:          e.ID,
		User:        e.User.ToResponse(),
		Name:        e.Name,
		Description: e.Description,
		Address:     e.Address,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
	}
}
