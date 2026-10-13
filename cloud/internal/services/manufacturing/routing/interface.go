package routing

import (
	"context"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

type Repository interface {
	CreateRouting(ctx context.Context, entity *model.Routing) error
	ProductBelongsToTenant(ctx context.Context, tenantID, productID uint) (bool, error)
	GetRoutingByID(ctx context.Context, tenantID, id uint) (*model.Routing, error)
	GetRoutingByNameAndVersion(ctx context.Context, tenantID, productID uint, name string, version int) (*model.Routing, error)
	ListRoutings(ctx context.Context, tenantID uint, productID *uint, status *model.RoutingStatus) ([]*model.Routing, error)
	UpdateRouting(ctx context.Context, entity *model.Routing) error
	DeactivateOtherRoutings(ctx context.Context, tenantID, productID, exceptRoutingID uint) error
	CountOperations(ctx context.Context, tenantID, routingID uint) (int64, error)

	CreateOperation(ctx context.Context, entity *model.RoutingOperation) error
	GetOperation(ctx context.Context, tenantID, routingID, operationID uint) (*model.RoutingOperation, error)
	GetOperationByCode(ctx context.Context, tenantID, routingID uint, code string) (*model.RoutingOperation, error)
	GetOperationBySequence(ctx context.Context, tenantID, routingID uint, sequence int) (*model.RoutingOperation, error)
	ListOperations(ctx context.Context, tenantID, routingID uint) ([]*model.RoutingOperation, error)
	UpdateOperation(ctx context.Context, entity *model.RoutingOperation) error
	DeleteOperation(ctx context.Context, tenantID, routingID, operationID uint) error
}

type Service interface {
	CreateRouting(ctx context.Context, req *CreateRoutingRequest) (*RoutingResponse, error)
	GetRouting(ctx context.Context, id uint) (*RoutingResponse, error)
	ListRoutings(ctx context.Context, productID *uint, status *model.RoutingStatus) ([]*RoutingResponse, error)
	UpdateRouting(ctx context.Context, id uint, req *UpdateRoutingRequest) (*RoutingResponse, error)
	ActivateRouting(ctx context.Context, id uint) error
	DeactivateRouting(ctx context.Context, id uint) error

	AddOperation(ctx context.Context, routingID uint, req *CreateOperationRequest) (*OperationResponse, error)
	ListOperations(ctx context.Context, routingID uint) ([]*OperationResponse, error)
	UpdateOperation(ctx context.Context, routingID, operationID uint, req *UpdateOperationRequest) (*OperationResponse, error)
	DeleteOperation(ctx context.Context, routingID, operationID uint) error
}
