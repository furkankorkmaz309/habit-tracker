package models

import "time"

type Habit struct {
	ID          int
	Title       string
	Description string
	Frequency   string
	CreatedAt   time.Time
	UserID      int
}
