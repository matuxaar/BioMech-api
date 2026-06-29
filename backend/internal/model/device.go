package model

import "time"

type DeviceType string

const (
	DeviceTypeProsthetic DeviceType = "prosthetic"
	DeviceTypeSensor     DeviceType = "sensor"
)

type Device struct {
	ID                 string     `json:"id"`
	UserID             string     `json:"user_id"`
	Type               DeviceType `json:"type"`
	Name               string     `json:"name"`
	HWVersion          string     `json:"hw_version"`
	BLEServiceUUID     string     `json:"ble_service_uuid,omitempty"`
	BLECommandCharUUID string     `json:"ble_command_char_uuid,omitempty"`
	BLEStatusCharUUID  string     `json:"ble_status_char_uuid,omitempty"`
	BLEEMGCharUUID     string     `json:"ble_emg_char_uuid,omitempty"`
	LastRecordingAt    *time.Time `json:"last_recording_at,omitempty"`
	LastTrainingAt     *time.Time `json:"last_training_at,omitempty"`
	UpdatedAt          time.Time  `json:"updated_at"`
	CreatedAt          time.Time  `json:"created_at"`
}

type CreateDeviceRequest struct {
	Type               DeviceType `json:"type" binding:"required,oneof=prosthetic sensor"`
	Name               string     `json:"name" binding:"required"`
	HWVersion          string     `json:"hw_version" binding:"required"`
	BLEServiceUUID     string     `json:"ble_service_uuid"`
	BLECommandCharUUID string     `json:"ble_command_char_uuid"`
	BLEStatusCharUUID  string     `json:"ble_status_char_uuid"`
	BLEEMGCharUUID     string     `json:"ble_emg_char_uuid"`
}

type UpdateDeviceRequest struct {
	Name               *string     `json:"name"`
	HWVersion          *string     `json:"hw_version"`
	Type               *DeviceType `json:"type" binding:"omitempty,oneof=prosthetic sensor"`
	BLEServiceUUID     *string     `json:"ble_service_uuid"`
	BLECommandCharUUID *string     `json:"ble_command_char_uuid"`
	BLEStatusCharUUID  *string     `json:"ble_status_char_uuid"`
	BLEEMGCharUUID     *string     `json:"ble_emg_char_uuid"`
}

type DeviceAction struct {
	Name       string  `json:"name"`
	Emoji      string  `json:"emoji"`
	ActionCode int     `json:"action_code"`
	Accuracy   float64 `json:"accuracy"`
}

type DeviceActionsResponse struct {
	DeviceID string         `json:"device_id"`
	Actions  []DeviceAction `json:"actions"`
}
