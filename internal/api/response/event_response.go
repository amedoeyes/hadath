package response

import (
	"time"

	"github.com/google/uuid"
)

type EventResponse struct {
	ID          uuid.UUID    `json:"id"`
	User        UserResponse `json:"user"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Address     string       `json:"address"`
	StartTime   time.Time    `json:"start_time"`
	EndTime     time.Time    `json:"end_time"`
}
