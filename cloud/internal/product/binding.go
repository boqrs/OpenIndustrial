package product

import "time"

// ProductDeviceBinding creates a link between a specific product instance (by SN)
// and a device that is part of it or monitors it.
type ProductDeviceBinding struct {
	ID                string     `json:"id"`
	ProductInstanceID string     `json:"productInstanceId"`
	DeviceID          string     `json:"deviceId"`
	StartTime         time.Time  `json:"startTime"`
	EndTime           *time.Time `json:"endTime,omitempty"`
}