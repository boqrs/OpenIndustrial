package bootstrap

import (
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/infrastructure/repository"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/lifecycle"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/productinstance"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/trace"
	"gorm.io/gorm"
)

type Container struct {
	ProductInstance *productinstance.Service
	Lifecycle       *lifecycle.Service
	Trace           *trace.Service
}

func NewContainer(db *gorm.DB) *Container {
	productRepo := repository.NewProductInstanceRepository(db)
	lifecycleRepo := repository.NewLifecycleRepository(db)

	return &Container{
		ProductInstance: productinstance.NewService(
			productRepo,
		),
		Lifecycle: lifecycle.NewService(
			lifecycleRepo,
		),
		Trace: trace.NewService(
			lifecycleRepo,
		),
	}
}