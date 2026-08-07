package bootstrap

import (
	"context"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/device"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/iot"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/pkg/event"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/product"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/workorder"
)

// Container holds all the services for the application.
type Container struct {
	ProductSvc   *product.Service
	DeviceSvc    *device.Service
	WorkorderSvc *workorder.Service
}

// NewContainer creates and wires up all the application services.
func NewContainer(ctx context.Context, bus event.Bus) *Container {
	// Product Domain
	productRepo := product.NewMemoryRepository()
	productSvc := product.NewService(productRepo)

	// Device Domain
	deviceRepo := device.NewMemoryRepository()
	deviceSvc := device.NewService(deviceRepo, bus)

	// IoT Domain (internal service)
	iotSvc := iot.NewService(bus)
	iotHandler := iot.NewMQTTHandler(iotSvc)
	iotConsumer := iot.NewConsumer(iotHandler)

	// WorkOrder Domain
	workorderRepo := workorder.NewMemoryRepository()
	workorderSvc := workorder.NewService(workorderRepo, productRepo, bus)


	// Start background services
	go iotConsumer.Run(ctx)

	return &Container{
		ProductSvc:   productSvc,
		DeviceSvc:    deviceSvc,
		WorkorderSvc: workorderSvc,
	}
}