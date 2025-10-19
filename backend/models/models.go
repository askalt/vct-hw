package models

import (
	"time"
)

// Event describes a posted event.
type Event struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Date        string   `json:"date"`
}

// Registration associates an user with an event.
type Registration struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
