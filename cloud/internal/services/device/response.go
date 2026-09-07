package device

import "github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"

type DeviceResponse struct {
	ID         uint `json:"id"`
	ResourceID uint `json:"resource_id"`
	ProductID  uint `json:"product_id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
	HardwareID   string `json:"hardware_id"`
	WorkOrderID       uint `json:"work_order_id"`
	ExecutionID       uint `json:"execution_id"`
	ExecutionResultID uint `json:"execution_result_id"`
	Status model.DeviceStatus `json:"status"`
	ParentResourceID *uint `json:"parent_resource_id"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	LastOnlineAt *string `json:"last_online_at,omitempty"`
}

type ListDevicesResponse struct {
	Items    []*DeviceResponse `json:"items"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}