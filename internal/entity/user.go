package entity

import "time"

// User -.
type User struct {
	CreatedAt   time.Time `json:"created_at"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email"`
	ID          int       `json:"-"`
}
