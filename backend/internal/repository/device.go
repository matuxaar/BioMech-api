package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matuxaar/BioMech-api/internal/model"
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

func (r *DeviceRepository) Update(ctx context.Context, id string, req *model.UpdateDeviceRequest) (*model.Device, error) {
	if req.Name != nil {
		_, err := r.db.Exec(ctx, `UPDATE devices SET name = $1 WHERE id = $2`, *req.Name, id)
		if err != nil {
			return nil, err
		}
	}
	if req.HWVersion != nil {
		_, err := r.db.Exec(ctx, `UPDATE devices SET hw_version = $1 WHERE id = $2`, *req.HWVersion, id)
		if err != nil {
			return nil, err
		}
	}
	if req.Type != nil {
		_, err := r.db.Exec(ctx, `UPDATE devices SET type = $1 WHERE id = $2`, *req.Type, id)
		if err != nil {
			return nil, err
		}
	}
	return r.FindByID(ctx, id)
}

func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM devices WHERE id = $1`, id)
	return err
}
