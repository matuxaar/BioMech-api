package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/motvii/desertacia/internal/model"
)

type DeviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(ctx context.Context, userID string, req *model.CreateDeviceRequest) (*model.Device, error) {
	device := &model.Device{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      req.Type,
		Name:      req.Name,
		HWVersion: req.HWVersion,
		CreatedAt: time.Now(),
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO devices (id, user_id, type, name, hw_version, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		device.ID, device.UserID, device.Type, device.Name, device.HWVersion, device.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (r *DeviceRepository) FindByUserID(ctx context.Context, userID string) ([]model.Device, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, type, name, hw_version, created_at
		 FROM devices WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.UserID, &d.Type, &d.Name, &d.HWVersion, &d.CreatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil
}

func (r *DeviceRepository) FindByID(ctx context.Context, id string) (*model.Device, error) {
	device := &model.Device{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, type, name, hw_version, created_at
		 FROM devices WHERE id = $1`, id,
	).Scan(&device.ID, &device.UserID, &device.Type, &device.Name, &device.HWVersion, &device.CreatedAt)
	if err != nil {
		return nil, err
	}
	return device, nil
}
