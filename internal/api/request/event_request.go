package request

import "time"

type EventRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Address     string    `json:"address" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required,gtfield=StartTime"`
}
