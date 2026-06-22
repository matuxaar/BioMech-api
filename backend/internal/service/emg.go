package service

import (
	"context"

	"github.com/motvii/desertacia/internal/model"
	"github.com/motvii/desertacia/internal/repository"
)

type EMGService struct {
	emgRepo    *repository.EMGRepository
	deviceRepo *repository.DeviceRepository
}

func NewEMGService(emgRepo *repository.EMGRepository, deviceRepo *repository.DeviceRepository) *EMGService {
	return &EMGService{emgRepo: emgRepo, deviceRepo: deviceRepo}
}

func (s *EMGService) StartSession(ctx context.Context, userID string, req *model.CreateEMGSessionRequest) (*model.EMGSession, error) {
	device, err := s.deviceRepo.FindByID(ctx, req.DeviceID)
	if err != nil {
		return nil, err
	}
	if device.UserID != userID {
		return nil, ErrAccessDenied
	}
	return s.emgRepo.CreateSession(ctx, userID, req)
}

func (s *EMGService) EndSession(ctx context.Context, sessionID string) error {
	return s.emgRepo.EndSession(ctx, sessionID)
}

func (s *EMGService) ListSessions(ctx context.Context, userID string) ([]model.EMGSession, error) {
	return s.emgRepo.FindSessionsByUserID(ctx, userID)
}

func (s *EMGService) GetSession(ctx context.Context, id string) (*model.EMGSession, error) {
	return s.emgRepo.FindSessionByID(ctx, id)
}

func (s *EMGService) AddSample(ctx context.Context, sessionID string, req *model.AddSampleRequest) (*model.EMGSample, error) {
	return s.emgRepo.AddSample(ctx, sessionID, req)
}

func (s *EMGService) AddSamplesBatch(ctx context.Context, sessionID string, samples []model.AddSampleRequest) error {
	return s.emgRepo.AddSamplesBatch(ctx, sessionID, samples)
}

func (s *EMGService) GetSamples(ctx context.Context, sessionID string) ([]model.EMGSample, error) {
	return s.emgRepo.FindSamplesBySessionID(ctx, sessionID)
}
