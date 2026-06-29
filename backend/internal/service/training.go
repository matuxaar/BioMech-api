package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
)

type TrainingService struct {
	trainingRepo *repository.TrainingRepository
	emgRepo      *repository.EMGRepository
	deviceRepo   *repository.DeviceRepository
	mlClient     *MLClient
}

func NewTrainingService(trainingRepo *repository.TrainingRepository, emgRepo *repository.EMGRepository, deviceRepo *repository.DeviceRepository, mlClient *MLClient) *TrainingService {
	return &TrainingService{trainingRepo: trainingRepo, emgRepo: emgRepo, deviceRepo: deviceRepo, mlClient: mlClient}
}

func (s *TrainingService) CreateJob(ctx context.Context, userID string, req *model.CreateTrainingJobRequest) (*model.TrainingJob, error) {
	return s.trainingRepo.Create(ctx, userID, req)
}

func (s *TrainingService) ListJobs(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.TrainingJob], error) {
	total, err := s.trainingRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	jobs, err := s.trainingRepo.FindByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}
	result := model.NewPaginatedResponse(jobs, total, page, limit)
	return &result, nil
}

func (s *TrainingService) GetJob(ctx context.Context, userID, id string) (*model.TrainingJob, error) {
	job, err := s.trainingRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if job.UserID != userID {
		return nil, ErrAccessDenied
	}
	return job, nil
}

func (s *TrainingService) UpdateJobStatus(ctx context.Context, id, status, modelPath string, accuracy float64, errMsg string) error {
	return s.trainingRepo.UpdateStatus(ctx, id, model.TrainingStatus(status), modelPath, accuracy, errMsg)
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
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in StartTraining goroutine: %v\n%s", r, debug.Stack())
				s.trainingRepo.UpdateStatus(context.Background(), jobID, model.TrainingStatusFailed, "", 0, fmt.Sprintf("panic: %v", r))
			}
		}()
		result, err := s.mlClient.Train(job)
		if err != nil {
			s.trainingRepo.UpdateStatus(context.Background(), jobID, model.TrainingStatusFailed, "", 0, err.Error())
			return
		}
		status := model.TrainingStatusCompleted
		if result.Status == "failed" {
			status = model.TrainingStatusFailed
		}
		s.trainingRepo.UpdateStatus(context.Background(), jobID, status, result.ModelPath, result.Accuracy, "")
		if status == model.TrainingStatusCompleted {
			deviceIDs, err := s.emgRepo.FindDeviceIDsBySessionIDs(context.Background(), job.SessionIDs)
			if err == nil {
				now := time.Now()
				for _, did := range deviceIDs {
					s.deviceRepo.UpdateLastTrainingAt(context.Background(), did, now)
				}
			}
		}
	}()

	return nil
}

func (s *TrainingService) Predict(ctx context.Context, samples []model.EMGSample) (*model.PredictResponse, error) {
	predictions, err := s.mlClient.Predict(samples)
	if err != nil {
		return nil, err
	}
	return &model.PredictResponse{Predictions: predictions}, nil
}

func (s *TrainingService) ProcessUpload(ctx context.Context, userID, deviceID, label string, file io.Reader) (*model.EMGSession, error) {
	session, err := s.emgRepo.CreateSession(ctx, userID, &model.CreateEMGSessionRequest{
		DeviceID: deviceID,
		Label:    label,
	})
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, errors.New("invalid CSV: missing header row")
	}

	if len(headers) < 9 {
		return nil, errors.New("invalid CSV: expected at least 9 columns (timestamp + 8 channels)")
	}

	var samples []model.AddSampleRequest
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("invalid CSV: " + err.Error())
		}

		ts, err := time.Parse(time.RFC3339, strings.TrimSpace(record[0]))
		if err != nil {
			ts = time.Now()
		}

		sample := model.AddSampleRequest{
			Timestamp: ts,
		}
		for i := 0; i < 8 && i+1 < len(record); i++ {
			val, _ := strconv.ParseFloat(strings.TrimSpace(record[i+1]), 64)
			switch i {
			case 0:
				sample.Channel1 = val
			case 1:
				sample.Channel2 = val
			case 2:
				sample.Channel3 = val
			case 3:
				sample.Channel4 = val
			case 4:
				sample.Channel5 = val
			case 5:
				sample.Channel6 = val
			case 6:
				sample.Channel7 = val
			case 7:
				sample.Channel8 = val
			}
		}
		samples = append(samples, sample)
	}

	if len(samples) == 0 {
		return nil, errors.New("CSV contains no data rows")
	}

	if err := s.emgRepo.AddSamplesBatch(ctx, session.ID, samples); err != nil {
		return nil, err
	}

	s.deviceRepo.UpdateLastRecordingAt(ctx, deviceID, time.Now())

	return session, nil
}
