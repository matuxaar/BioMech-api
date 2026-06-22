package model

import "time"

type TrainingStatus string

const (
	TrainingStatusPending    TrainingStatus = "pending"
	TrainingStatusRunning    TrainingStatus = "running"
	TrainingStatusCompleted  TrainingStatus = "completed"
	TrainingStatusFailed     TrainingStatus = "failed"
)

type TrainingJob struct {
	ID           string         `json:"id"`
	UserID       string         `json:"user_id"`
	SessionIDs   []string       `json:"session_ids"`
	Status       TrainingStatus `json:"status"`
	ModelPath    string         `json:"model_path,omitempty"`
	Accuracy     float64        `json:"accuracy,omitempty"`
	ErrorMessage string         `json:"error_message,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CreateTrainingJobRequest struct {
	SessionIDs []string `json:"session_ids" binding:"required,min=1"`
}
