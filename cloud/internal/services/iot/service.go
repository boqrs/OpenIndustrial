package iot

import (
	"github.com/boqrs/OpenIndustrial/cloud/internal/services/kernel/security"
)

type service struct {
	repo     Repository
	security security.Service

	// Cloud -> Device.
	commandPublisher MQTTCommandPublisher

	// Device -> Cloud.
	statusSubscriber MQTTStatusSubscriber
}

func NewService(
	repo Repository,
	securitySvc security.Service,
	commandPublisher MQTTCommandPublisher,
	statusSubscriber MQTTStatusSubscriber,
) Service {
	return &service{
		repo:             repo,
		security:         securitySvc,
		commandPublisher: commandPublisher,
		statusSubscriber: statusSubscriber,
	}
}

var _ Service = (*service)(nil)
