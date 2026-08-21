package factory

import (
	"context"

	"github.com/google/uuid"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context,entity *model.Factory) error
	GetByUUID(ctx context.Context,id uuid.UUID) (*model.Factory, error)
	GetByResourceID(ctx context.Context,resourceID uuid.UUID) (*model.Factory, error)
	GetByCode(ctx context.Context,code string) (*model.Factory, error)
	Update(	ctx context.Context,entity *model.Factory) error
	Delete(ctx context.Context,id uuid.UUID) error
}

type Service interface {
	CreateFactory(ctx context.Context,req *CreateFactoryRequest) (*FactoryResponse, error)
	GetFactory(ctx context.Context,factoryID uuid.UUID) (*FactoryResponse, error)
	UpdateFactory(ctx context.Context,factoryID uuid.UUID,req *UpdateFactoryRequest) (*FactoryResponse, error)
	DeleteFactory(ctx context.Context,factoryID uuid.UUID) error
	CreateTopologyNode(ctx context.Context,req *CreateTopologyNodeRequest) (*TopologyNodeResponse, error)
	UpdateTopologyNode(ctx context.Context,resourceID uuid.UUID,req *UpdateTopologyNodeRequest) (*TopologyNodeResponse, error)
	MoveTopologyNode(ctx context.Context,req *MoveTopologyNodeRequest) error
	DeleteTopologyNode(ctx context.Context,resourceID uuid.UUID) error
	GetTopology(ctx context.Context,factoryID uuid.UUID) (*FactoryTopologyResponse, error)
}