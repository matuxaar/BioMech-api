package service

import (
	"context"
	"errors"

	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
)

var ErrAccessDenied = errors.New("access denied")

type TrainingService struct {
	trainingRepo *repository.TrainingRepository
	mlClient     *MLClient
}

func NewTrainingService(trainingRepo *repository.TrainingRepository, mlClient *MLClient) *TrainingService {
	return &TrainingService{trainingRepo: trainingRepo, mlClient: mlClient}
}

func (s *TrainingService) CreateJob(ctx context.Context, userID string, req *model.CreateTrainingJobRequest) (*model.TrainingJob, error) {
	return s.trainingRepo.Create(ctx, userID, req)
}

func (s *TrainingService) ListJobs(ctx context.Context, userID string) ([]model.TrainingJob, error) {
	return s.trainingRepo.FindByUserID(ctx, userID)
}

func (s *TrainingService) GetJob(ctx context.Context, id string) (*model.TrainingJob, error) {
	return s.trainingRepo.FindByID(ctx, id)
}

func (s *TrainingService) StartTraining(ctx context.Context, jobID string) error {
	job, err := s.trainingRepo.FindByID(ctx, jobID)
	if err != nil {
		return err
	}

	if err := s.trainingRepo.UpdateStatus(ctx, jobID, model.TrainingStatusRunning, "", 0, ""); err != nil {
		return err
	}

	go func() {
		if err := s.mlClient.Train(job); err != nil {
			s.trainingRepo.UpdateStatus(context.Background(), jobID, model.TrainingStatusFailed, "", 0, err.Error())
			return
		}
		s.trainingRepo.UpdateStatus(context.Background(), jobID, model.TrainingStatusCompleted, "models/"+jobID+".h5", 0.95, "")
	}()

	return nil
}
