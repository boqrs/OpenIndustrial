package gateway

import (
	"time"

	"github.com/OpenGongChang/OpenIndustrial/gateway/runtime/driver"
)

/*
GatewayState 表示网关运行状态
*/
type GatewayState uint8

const (
	StateCreated GatewayState = iota
	StateInitialized
	StateRunning
	StateStopped
	StateError
)

/*
EventEnvelope

Gateway内部统一事件

driver.Event 会经过这里
*/
type EventEnvelope struct {
	Source string
	Event  driver.Event
	Timestamp time.Time
}