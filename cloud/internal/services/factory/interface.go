package factory

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	Create(ctx context.Context,entity *model.Factory) error
	GetByID(ctx context.Context,id uint) (*model.Factory, error)
	GetByResourceID(ctx context.Context,resourceID uint) (*model.Factory, error)
	GetByCode(ctx context.Context,code string) (*model.Factory, error)
	Update(	ctx context.Context,entity *model.Factory) error
	Delete(ctx context.Context,id uint) error
}

type Service interface {
	CreateFactory(ctx context.Context,req *CreateFactoryRequest) (*FactoryResponse, error)
	GetFactory(ctx context.Context,factoryID uint) (*FactoryResponse, error)
	UpdateFactory(ctx context.Context,factoryID uint,req *UpdateFactoryRequest) (*FactoryResponse, error)
	DeleteFactory(ctx context.Context,factoryID uint) error
	CreateTopologyNode(ctx context.Context,req *CreateTopologyNodeRequest) (*TopologyNodeResponse, error)
	UpdateTopologyNode(ctx context.Context,resourceID uint,req *UpdateTopologyNodeRequest) (*TopologyNodeResponse, error)
	MoveTopologyNode(ctx context.Context,req *MoveTopologyNodeRequest) error
	DeleteTopologyNode(ctx context.Context,resourceID uint) error
	GetTopology(ctx context.Context,factoryID uint) (*FactoryTopologyResponse, error)
}