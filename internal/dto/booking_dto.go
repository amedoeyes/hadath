package dto

import "github.com/google/uuid"

type BookingRequest struct {
	EventID uuid.UUID `json:"event_id" validate:"required"`
}
