package response

import "github.com/google/uuid"

type AuthResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	APIKey string    `json:"api_key"`
}
