package models

import (
	"time"
)

// Event describes a posted event.
type Event struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
}

// Registration associates an user with an event.
type Registration struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
