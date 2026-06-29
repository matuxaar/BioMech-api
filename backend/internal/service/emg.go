package service

import (
	"context"
	"time"

	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
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
	session, err := s.emgRepo.CreateSession(ctx, userID, req)
	if err != nil {
		return nil, err
	}
	s.deviceRepo.UpdateLastRecordingAt(ctx, req.DeviceID, time.Now())
	return session, nil
}

func (s *EMGService) EndSession(ctx context.Context, userID, sessionID string) error {
	session, err := s.emgRepo.FindSessionByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if session.UserID != userID {
		return ErrAccessDenied
	}
	return s.emgRepo.EndSession(ctx, sessionID)
}

func (s *EMGService) ListSessions(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.EMGSession], error) {
	total, err := s.emgRepo.CountSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	sessions, err := s.emgRepo.FindSessionsByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}
	result := model.NewPaginatedResponse(sessions, total, page, limit)
	return &result, nil
}

func (s *EMGService) GetSession(ctx context.Context, userID, id string) (*model.EMGSession, error) {
	session, err := s.emgRepo.FindSessionByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, ErrAccessDenied
	}
	return session, nil
}

func (s *EMGService) checkSessionOwnership(ctx context.Context, userID, sessionID string) (*model.EMGSession, error) {
	session, err := s.emgRepo.FindSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if session.UserID != userID {
		return nil, ErrAccessDenied
	}
	return session, nil
}

func (s *EMGService) AddSample(ctx context.Context, userID, sessionID string, req *model.AddSampleRequest) (*model.EMGSample, error) {
	session, err := s.checkSessionOwnership(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	sample, err := s.emgRepo.AddSample(ctx, sessionID, req)
	if err != nil {
		return nil, err
	}
	s.deviceRepo.UpdateLastRecordingAt(ctx, session.DeviceID, time.Now())
	return sample, nil
}

func (s *EMGService) AddSamplesBatch(ctx context.Context, userID, sessionID string, samples []model.AddSampleRequest) error {
	session, err := s.checkSessionOwnership(ctx, userID, sessionID)
	if err != nil {
		return err
	}
	if err := s.emgRepo.AddSamplesBatch(ctx, sessionID, samples); err != nil {
		return err
	}
	s.deviceRepo.UpdateLastRecordingAt(ctx, session.DeviceID, time.Now())
	return nil
}

func (s *EMGService) GetSamples(ctx context.Context, userID, sessionID string, page, limit int) (*model.PaginatedResponse[model.EMGSample], error) {
	if _, err := s.checkSessionOwnership(ctx, userID, sessionID); err != nil {
		return nil, err
	}
	total, err := s.emgRepo.CountSamplesBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	samples, err := s.emgRepo.FindSamplesBySessionID(ctx, sessionID, page, limit)
	if err != nil {
		return nil, err
	}
	result := model.NewPaginatedResponse(samples, total, page, limit)
	return &result, nil
}
