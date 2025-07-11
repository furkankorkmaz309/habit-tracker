package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password" validate:"containsany=1234567890"`
	CreatedAt time.Time `json:"created_at"`
}
