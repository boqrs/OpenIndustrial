package event

// IoT Domain
const (
	EventIoTDeviceOnline  = "iot.device.online"
	EventIoTDeviceOffline = "iot.device.offline"
	EventIoTPointChanged  = "iot.point.changed"
)

// MES Domain
const (
	EventMESWorkOrderCreated   = "mes.workorder.created"
	EventMESWorkOrderStarted   = "mes.workorder.started"
	EventMESWorkOrderCompleted = "mes.workorder.completed"
	EventMESStationStarted     = "mes.station.started"
	EventMESStationCompleted   = "mes.station.completed"
)

// Product Domain
const (
	EventProductCreated         = "product.created"
	EventProductLifecycleChanged = "product.lifecycle.changed"
	EventProductShipped         = "product.shipped"
	EventProductDelivered       = "product.delivered"
)

// Quality Domain
const (
	EventQualityInspectionStarted = "quality.inspection.started"
	EventQualityInspectionPassed  = "quality.inspection.pass"
	EventQualityInspectionFailed  = "quality.inspection.fail"
)