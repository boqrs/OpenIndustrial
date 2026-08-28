package routing

import (
	"context"

	"github.com/google/uuid"
	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

// Repository defines the persistence interface for routings and operations
type Repository interface {
	// Routing methods
	CreateRouting(ctx context.Context, entity *model.Routing) error
	GetRoutingByID(ctx context.Context, tenantID uuid.UUID, id uint) (*model.Routing, error)
	GetRoutingByNameAndVersion(ctx context.Context, tenantID uuid.UUID, productID uint, name string, version int) (*model.Routing, error)
	ListRoutings(ctx context.Context, tenantID uuid.UUID, productID *uint, status *model.RoutingStatus) ([]*model.Routing, error)
	UpdateRouting(ctx context.Context, entity *model.Routing) error
	DeactivateOtherRoutings(ctx context.Context, tenantID uuid.UUID, productID uint, exceptRoutingID uint) error
	CountOperations(ctx context.Context, tenantID uuid.UUID, routingID uint) (int64, error)

	// Operation methods
	CreateOperation(ctx context.Context, entity *model.RoutingOperation) error
	GetOperation(ctx context.Context, tenantID uuid.UUID, routingID uint, operationID uint) (*model.RoutingOperation, error)
	GetOperationByCode(ctx context.Context, tenantID uuid.UUID, routingID uint, code string) (*model.RoutingOperation, error)
	GetOperationBySequence(ctx context.Context, tenantID uuid.UUID, routingID uint, sequence int) (*model.RoutingOperation, error)
	ListOperations(ctx context.Context, tenantID uuid.UUID, routingID uint) ([]*model.RoutingOperation, error)
	UpdateOperation(ctx context.Context, entity *model.RoutingOperation) error
	DeleteOperation(ctx context.Context, tenantID uuid.UUID, routingID uint, operationID uint) error
}


// --- Service Interface ---
type Service interface {
	// Routing methods
	CreateRouting(ctx context.Context, req *CreateRoutingRequest) (*RoutingResponse, error)
	GetRouting(ctx context.Context, id uint) (*RoutingResponse, error)
	ListRoutings(ctx context.Context, productID *uint, status *model.RoutingStatus) ([]*RoutingResponse, error)
	UpdateRouting(ctx context.Context, id uint, req *UpdateRoutingRequest) (*RoutingResponse, error)
	ActivateRouting(ctx context.Context, id uint) error
	DeactivateRouting(ctx context.Context, id uint) error

	// Operation methods
	AddOperation(ctx context.Context, routingID uint, req *CreateOperationRequest) (*OperationResponse, error)
	ListOperations(ctx context.Context, routingID uint) ([]*OperationResponse, error)
	UpdateOperation(ctx context.Context, routingID uint, operationID uint, req *UpdateOperationRequest) (*OperationResponse, error)
	DeleteOperation(ctx context.Context, routingID, operationID uint) error
}
