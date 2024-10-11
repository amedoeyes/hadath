package dto

import (
	"time"

	"github.com/google/uuid"
)

type EventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Address     string    `json:"address" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
}

type EventResponse struct {
	ID          uuid.UUID    `json:"id"`
	User        UserResponse `json:"user"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Address     string       `json:"address"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     time.Time    `json:"end_time"`
}
