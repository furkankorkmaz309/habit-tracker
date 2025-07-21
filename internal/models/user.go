package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password" validate:"containsany=1234567890"`
	CreatedAt time.Time `json:"created_at"`
}

type Email struct {
	Email string `json:"email" validate:"email"`
}
