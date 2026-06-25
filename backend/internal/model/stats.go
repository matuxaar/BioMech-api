package model

type DashboardStats struct {
	DeviceCount        int     `json:"device_count"`
	TotalTrainings     int     `json:"total_trainings"`
	CompletedTrainings int     `json:"completed_trainings"`
	AverageAccuracy    float64 `json:"average_accuracy"`
	TopMovements       []string `json:"top_movements"`
}
