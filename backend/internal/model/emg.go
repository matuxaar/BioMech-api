package model

import "time"

type EMGSession struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	DeviceID  string     `json:"device_id"`
	Label     string     `json:"label"`
	StartedAt time.Time  `json:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type EMGSample struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Timestamp time.Time `json:"timestamp"`
	Channel1  float64   `json:"channel_1"`
	Channel2  float64   `json:"channel_2"`
	Channel3  float64   `json:"channel_3"`
	Channel4  float64   `json:"channel_4"`
	Channel5  float64   `json:"channel_5"`
	Channel6  float64   `json:"channel_6"`
	Channel7  float64   `json:"channel_7"`
	Channel8  float64   `json:"channel_8"`
	Metadata  string    `json:"metadata,omitempty"`
}

type CreateEMGSessionRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Label    string `json:"label"`
}

type AddSampleRequest struct {
	Timestamp time.Time `json:"timestamp" binding:"required"`
	Channel1  float64   `json:"channel_1"`
	Channel2  float64   `json:"channel_2"`
	Channel3  float64   `json:"channel_3"`
	Channel4  float64   `json:"channel_4"`
	Channel5  float64   `json:"channel_5"`
	Channel6  float64   `json:"channel_6"`
	Channel7  float64   `json:"channel_7"`
	Channel8  float64   `json:"channel_8"`
	Metadata  string    `json:"metadata"`
}
