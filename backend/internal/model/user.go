package model

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Nickname     string    `json:"nickname,omitempty"`
	DisplayName  string    `json:"display_name,omitempty"`
	PhotoURL     string    `json:"photo_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	Nickname    *string `json:"nickname"`
	DisplayName *string `json:"display_name"`
	PhotoURL    *string `json:"photo_url"`
}

type ProfileResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Nickname    string `json:"nickname,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	PhotoURL    string `json:"photo_url,omitempty"`
	DeviceCount int    `json:"device_count"`
}
