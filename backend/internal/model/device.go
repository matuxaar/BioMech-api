package model

import "time"

type DeviceType string

const (
	DeviceTypeProsthetic DeviceType = "prosthetic"
	DeviceTypeSensor     DeviceType = "sensor"
)

type Device struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Type       DeviceType `json:"type"`
	Name       string     `json:"name"`
	HWVersion  string     `json:"hw_version"`
	CreatedAt  time.Time  `json:"created_at"`
}

type CreateDeviceRequest struct {
	Type      DeviceType `json:"type" binding:"required,oneof=prosthetic sensor"`
	Name      string     `json:"name" binding:"required"`
	HWVersion string     `json:"hw_version" binding:"required"`
}
