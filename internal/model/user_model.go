package model

import (
	"github.com/amedoeyes/hadath/internal/dto"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Name     string
	Email    string
	Password string
}

func (u *User) ToResponse() dto.UserResponse {
	return dto.UserResponse{
		ID:   u.ID,
		Name: u.Name,
	}
}
