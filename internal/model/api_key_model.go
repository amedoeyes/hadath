package model

import (
	"github.com/google/uuid"
)

type APIKey struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Key    string
}
