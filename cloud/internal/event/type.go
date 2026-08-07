package event

type Type string

const (
	// IoT events
	TypeDeviceOnline  Type = "iot.device.online"
	TypeDeviceOffline Type = "iot.device.offline"
	TypePointChanged  Type = "iot.point.changed"
	TypeAlarmTriggered Type = "iot.alarm.triggered"

	// MES events
	TypeWorkOrderCreated   Type = "mes.workorder.created"
	TypeWorkOrderStarted   Type = "mes.workorder.started"
	TypeWorkOrderCompleted Type = "mes.workorder.completed"

	// Product
	TypeProductCreated          Type = "product.created"
	TypeProductLifecycleChanged Type = "product.lifecycle.changed"
)