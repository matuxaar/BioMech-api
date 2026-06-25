package service

import (
	"context"
	"errors"

	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
)

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
	return s.deviceRepo.FindByID(ctx, id)
}

func (s *DeviceService) Update(ctx context.Context, userID, deviceID string, req *model.UpdateDeviceRequest) (*model.Device, error) {
	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
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
		return nil, err
	}
	if device.UserID != userID {
		return nil, errors.New("access denied")
	}

	actions := []model.DeviceAction{
		{Name: "Rest", Emoji: "\u270B", Accuracy: 0.95},
		{Name: "Fist", Emoji: "\u270A", Accuracy: 0.92},
		{Name: "Open", Emoji: "\uD83D\uDD90\uFE0F", Accuracy: 0.88},
		{Name: "Pinch", Emoji: "\uD83E\uDD1F", Accuracy: 0.75},
		{Name: "Point", Emoji: "\u261D\uFE0F", Accuracy: 0.70},
	}

	return &model.DeviceActionsResponse{
		DeviceID: deviceID,
		Actions:  actions,
	}, nil
}
