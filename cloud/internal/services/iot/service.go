package iot

import (
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/iot/adapter"
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
)

type service struct {
	repo     Repository
	security security.Service

	messageAdapter adapter.DeviceMessageAdapter
}

func NewService(
	repo Repository,
	securitySvc security.Service,
	messageAdapter adapter.DeviceMessageAdapter,
) Service {
	return &service{
		repo:           repo,
		security:       securitySvc,
		messageAdapter: messageAdapter,
	}
}

var _ Service = (*service)(nil)
