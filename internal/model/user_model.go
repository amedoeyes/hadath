package model

import "time"

type User struct {
	ID        uint32
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
