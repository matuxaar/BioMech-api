package model

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name,omitempty"`
	PhotoURL     string    `json:"photo_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	DisplayName *string `json:"display_name"`
	PhotoURL    *string `json:"photo_url"`
}
