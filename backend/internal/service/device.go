package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
)

var ErrDeviceNotFound = errors.New("device not found")

type DeviceService struct {
	deviceRepo *repository.DeviceRepository
}

func NewDeviceService(deviceRepo *repository.DeviceRepository) *DeviceService {
	return &DeviceService{deviceRepo: deviceRepo}
}

func (s *DeviceService) Create(ctx context.Context, userID string, req *model.CreateDeviceRequest) (*model.Device, error) {
	return s.deviceRepo.Create(ctx, userID, req)
}

func (s *DeviceService) ListByUser(ctx context.Context, userID string) ([]model.Device, error) {
	return s.deviceRepo.FindByUserID(ctx, userID)
}

func (s *DeviceService) GetByID(ctx context.Context, id string) (*model.Device, error) {
	device, err := s.deviceRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	return device, nil
}

func (s *DeviceService) Update(ctx context.Context, userID, deviceID string, req *model.UpdateDeviceRequest) (*model.Device, error) {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	if device.UserID != userID {
		return nil, errors.New("access denied")
	}
	return s.deviceRepo.Update(ctx, deviceID, req)
}

func (s *DeviceService) Delete(ctx context.Context, userID, deviceID string) error {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrDeviceNotFound
		}
		return err
	}
	if device.UserID != userID {
		return errors.New("access denied")
	}
	return s.deviceRepo.Delete(ctx, deviceID)
}

func (s *DeviceService) GetActions(ctx context.Context, userID, deviceID string) (*model.DeviceActionsResponse, error) {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDeviceNotFound
		}
		return nil, err
	}
	if device.UserID != userID {
		return nil, errors.New("access denied")
	}

	actions, err := s.deviceRepo.GetActions(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	return &model.DeviceActionsResponse{
		DeviceID: deviceID,
		Actions:  actions,
	}, nil
}
