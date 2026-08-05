package gateway

import (
	"github.com/OpenGongChang/OpenIndustrial/runtime/driver"
)

/*
GatewayConfig

边缘网关整体配置
*/
type GatewayConfig struct {
	ID   string
	Name string

	Drivers    []driver.Config
	Publishers []PublisherConfig
}