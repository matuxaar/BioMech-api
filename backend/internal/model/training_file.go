package model

import "time"

type TrainingFile struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	DeviceID     string    `json:"device_id,omitempty"`
	OriginalName string    `json:"original_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	Label        string    `json:"label,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type FileUploadResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}
