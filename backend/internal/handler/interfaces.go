package handler

import (
	"context"
	"io"

	"github.com/matuxaar/BioMech-api/internal/model"
)

type AuthService interface {
	SyncUser(ctx context.Context, firebaseUID, email string) (*model.User, error)
	GetProfile(ctx context.Context, userID string) (*model.ProfileResponse, error)
	UpdateProfile(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error)
}

type DeviceService interface {
	Create(ctx context.Context, userID string, req *model.CreateDeviceRequest) (*model.Device, error)
	ListByUser(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.Device], error)
	GetByID(ctx context.Context, userID, id string) (*model.Device, error)
	Update(ctx context.Context, userID, id string, req *model.UpdateDeviceRequest) (*model.Device, error)
	Delete(ctx context.Context, userID, id string) error
	GetActions(ctx context.Context, userID, id string) (*model.DeviceActionsResponse, error)
}

type EMGService interface {
	StartSession(ctx context.Context, userID string, req *model.CreateEMGSessionRequest) (*model.EMGSession, error)
	EndSession(ctx context.Context, userID, id string) error
	ListSessions(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.EMGSession], error)
	GetSession(ctx context.Context, userID, id string) (*model.EMGSession, error)
	AddSample(ctx context.Context, userID, sessionID string, req *model.AddSampleRequest) (*model.EMGSample, error)
	AddSamplesBatch(ctx context.Context, userID, sessionID string, samples []model.AddSampleRequest) error
	GetSamples(ctx context.Context, userID, sessionID string, page, limit int) (*model.PaginatedResponse[model.EMGSample], error)
}

type TrainingService interface {
	CreateJob(ctx context.Context, userID string, req *model.CreateTrainingJobRequest) (*model.TrainingJob, error)
	StartTraining(ctx context.Context, jobID string) error
	ListJobs(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.TrainingJob], error)
	GetJob(ctx context.Context, userID, id string) (*model.TrainingJob, error)
	Predict(ctx context.Context, samples []model.EMGSample) (*model.PredictResponse, error)
	ProcessUpload(ctx context.Context, userID, deviceID, label string, file io.Reader) (*model.EMGSession, error)
	UpdateJobStatus(ctx context.Context, id, status, modelPath string, accuracy float64, errMsg string) error
}

type StatsService interface {
	GetDashboardStats(ctx context.Context, userID string) (*model.DashboardStats, error)
}

type TrainingFileService interface {
	Upload(ctx context.Context, userID, deviceID, label, filename string, file io.Reader, size int64) (*model.TrainingFile, error)
	List(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.TrainingFile], error)
	Get(ctx context.Context, id, userID string) (*model.TrainingFile, error)
	Delete(ctx context.Context, id, userID string) error
}
