package models

import "time"

type Habit struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Frequency   string    `json:"frequency"`
	Day         string    `json:"day"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      int       `json:"user_id"`
}
