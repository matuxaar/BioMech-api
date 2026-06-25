package service

import (
	"context"

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
