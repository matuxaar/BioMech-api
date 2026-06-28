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
		ID:                uuid.New().String(),
		UserID:            userID,
		Type:              req.Type,
		Name:              req.Name,
		HWVersion:         req.HWVersion,
		BLEServiceUUID:    req.BLEServiceUUID,
		BLECommandCharUUID: req.BLECommandCharUUID,
		BLEStatusCharUUID: req.BLEStatusCharUUID,
		BLEEMGCharUUID:    req.BLEEMGCharUUID,
		CreatedAt:         time.Now(),
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO devices (id, user_id, type, name, hw_version,
		 ble_service_uuid, ble_command_char_uuid, ble_status_char_uuid, ble_emg_char_uuid, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		device.ID, device.UserID, device.Type, device.Name, device.HWVersion,
		nullIfEmpty(device.BLEServiceUUID), nullIfEmpty(device.BLECommandCharUUID),
		nullIfEmpty(device.BLEStatusCharUUID), nullIfEmpty(device.BLEEMGCharUUID),
		device.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (r *DeviceRepository) FindByUserID(ctx context.Context, userID string) ([]model.Device, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, type, name, hw_version,
		 COALESCE(ble_service_uuid, ''), COALESCE(ble_command_char_uuid, ''),
		 COALESCE(ble_status_char_uuid, ''), COALESCE(ble_emg_char_uuid, ''), created_at
		 FROM devices WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device
	for rows.Next() {
		var d model.Device
		if err := rows.Scan(&d.ID, &d.UserID, &d.Type, &d.Name, &d.HWVersion,
			&d.BLEServiceUUID, &d.BLECommandCharUUID, &d.BLEStatusCharUUID, &d.BLEEMGCharUUID, &d.CreatedAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil
}

func (r *DeviceRepository) FindByID(ctx context.Context, id string) (*model.Device, error) {
	device := &model.Device{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, type, name, hw_version,
		 COALESCE(ble_service_uuid, ''), COALESCE(ble_command_char_uuid, ''),
		 COALESCE(ble_status_char_uuid, ''), COALESCE(ble_emg_char_uuid, ''), created_at
		 FROM devices WHERE id = $1`, id,
	).Scan(&device.ID, &device.UserID, &device.Type, &device.Name, &device.HWVersion,
		&device.BLEServiceUUID, &device.BLECommandCharUUID, &device.BLEStatusCharUUID, &device.BLEEMGCharUUID, &device.CreatedAt)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (r *DeviceRepository) Update(ctx context.Context, id string, req *model.UpdateDeviceRequest) (*model.Device, error) {
	if req.Name != nil {
		if _, err := r.db.Exec(ctx, `UPDATE devices SET name = $1 WHERE id = $2`, *req.Name, id); err != nil {
			return nil, err
		}
	}
	if req.HWVersion != nil {
		if _, err := r.db.Exec(ctx, `UPDATE devices SET hw_version = $1 WHERE id = $2`, *req.HWVersion, id); err != nil {
			return nil, err
		}
	}
	if req.Type != nil {
		if _, err := r.db.Exec(ctx, `UPDATE devices SET type = $1 WHERE id = $2`, *req.Type, id); err != nil {
			return nil, err
		}
	}
	if req.BLEServiceUUID != nil {
		if _, err := r.db.Exec(ctx, `UPDATE devices SET ble_service_uuid = $1 WHERE id = $2`, *req.BLEServiceUUID, id); err != nil {
			return nil, err
		}
	}
	if req.BLECommandCharUUID != nil {
		if _, err := r.db.Exec(ctx, `UPDATE devices SET ble_command_char_uuid = $1 WHERE id = $2`, *req.BLECommandCharUUID, id); err != nil {
			return nil, err
		}
	}
	if req.BLEStatusCharUUID != nil {
		if _, err := r.db.Exec(ctx, `UPDATE devices SET ble_status_char_uuid = $1 WHERE id = $2`, *req.BLEStatusCharUUID, id); err != nil {
			return nil, err
		}
	}
	if req.BLEEMGCharUUID != nil {
		if _, err := r.db.Exec(ctx, `UPDATE devices SET ble_emg_char_uuid = $1 WHERE id = $2`, *req.BLEEMGCharUUID, id); err != nil {
			return nil, err
		}
	}
	return r.FindByID(ctx, id)
}

func (r *DeviceRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM devices WHERE id = $1`, id)
	return err
}

func (r *DeviceRepository) GetActions(ctx context.Context, deviceID string) ([]model.DeviceAction, error) {
	rows, err := r.db.Query(ctx,
		`SELECT name, emoji, action_code, accuracy
		 FROM device_actions WHERE device_id = $1 ORDER BY action_code ASC`, deviceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actions []model.DeviceAction
	for rows.Next() {
		var a model.DeviceAction
		if err := rows.Scan(&a.Name, &a.Emoji, &a.ActionCode, &a.Accuracy); err != nil {
			return nil, err
		}
		actions = append(actions, a)
	}
	return actions, nil
}

func (r *DeviceRepository) GetUserActionStats(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.Query(ctx,
		`SELECT DISTINCT da.name FROM device_actions da
		 JOIN devices d ON d.id = da.device_id
		 WHERE d.user_id = $1 ORDER BY da.name`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		names = append(names, n)
	}
	return names, nil
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
