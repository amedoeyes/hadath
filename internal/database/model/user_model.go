package model

import (
	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
}

func (u *User) ToResponse() response.UserResponse {
	return response.UserResponse{
		ID:   u.ID,
		Name: u.Name,
	}
}
