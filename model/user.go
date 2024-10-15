package model

import (
	"time"
)

type User struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type Credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
