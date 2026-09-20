package salesorder

import (
	"context"
	"errors"
	"strings"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	customerSrv "github.com/boqrs/OpenIndustrial/cloud/internal/services/customer"
	"gorm.io/gorm"
)

var (
	ErrSalesOrderInvalid         = errors.New("invalid sales order")
	ErrSalesOrderNotFound        = errors.New("sales order not found")
	ErrSalesOrderNoExists        = errors.New("sales order number already exists")
	ErrSalesOrderEmptyItems      = errors.New("sales order must contain at least one item")
	ErrSalesOrderInvalidQuantity = errors.New("sales order item quantity must be greater than zero")
	ErrSalesOrderNotDraft        = errors.New("sales order is not in draft status")
	ErrSalesOrderCancelled       = errors.New("sales order is cancelled")
	ErrCustomerNotFound          = errors.New("customer not found")
	ErrCustomerInactive          = errors.New("customer is inactive")
)

type service struct {
	uow        UnitOfWork
	repository Repository
	customers  customerSrv.Repository
}

func NewService(
	uow UnitOfWork,
	repository Repository,
	customers customerSrv.Repository,
) Service {
	return &service{
		uow:        uow,
		repository: repository,
		customers:  customers,
	}
}

func (s *service) Create(
	ctx context.Context,
	req *CreateRequest,
) (*Response, error) {
	if req == nil {
		return nil, ErrSalesOrderInvalid
	}

	orderNo := strings.TrimSpace(req.OrderNo)
	if orderNo == "" {
		return nil, ErrSalesOrderInvalid
	}

	if req.CustomerID == 0 {
		return nil, ErrSalesOrderInvalid
	}

	if len(req.Items) == 0 {
		return nil, ErrSalesOrderEmptyItems
	}

	for _, item := range req.Items {
		if item.ProductID == 0 || item.OrderedQuantity <= 0 {
			return nil, ErrSalesOrderInvalidQuantity
		}
	}

	customer, err := s.customers.GetByID(
		ctx,
		req.CustomerID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCustomerNotFound
		}

		return nil, err
	}

	if customer.Status != model.CustomerStatusActive {
		return nil, ErrCustomerInactive
	}

	existing, err := s.repository.GetByOrderNo(
		ctx,
		orderNo,
	)

	if err == nil && existing != nil {
		return nil, ErrSalesOrderNoExists
	}

	if err != nil &&
		!errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var result *Response

	err = s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			order := &model.SalesOrder{
				CustomerID:  req.CustomerID,
				OrderNo:     orderNo,
				Status:      model.SalesOrderStatusDraft,
				OrderDate:   req.OrderDate,
				Description: strings.TrimSpace(req.Description),
			}

			if order.OrderDate.IsZero() {
				return ErrSalesOrderInvalid
			}

			if err := s.repository.Create(
				txCtx,
				order,
			); err != nil {
				return err
			}

			items := make(
				[]*model.SalesOrderItem,
				0,
				len(req.Items),
			)

			for _, itemReq := range req.Items {
				items = append(
					items,
					&model.SalesOrderItem{
						SalesOrderID:    order.ID,
						ProductID:       itemReq.ProductID,
						OrderedQuantity: itemReq.OrderedQuantity,
					},
				)
			}

			if err := s.repository.CreateItemsTx(
				txCtx,
				items,
			); err != nil {
				return err
			}

			order.Items = items
			result = toResponse(order)

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetByID(
	ctx context.Context,
	id uint,
) (*Response, error) {
	if id == 0 {
		return nil, ErrSalesOrderInvalid
	}

	order, err := s.repository.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSalesOrderNotFound
		}

		return nil, err
	}

	items, err := s.repository.ListItems(
		ctx,
		order.ID,
	)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return toResponse(order), nil
}

func (s *service) GetByOrderNo(
	ctx context.Context,
	orderNo string,
) (*Response, error) {
	orderNo = strings.TrimSpace(orderNo)

	if orderNo == "" {
		return nil, ErrSalesOrderInvalid
	}

	order, err := s.repository.GetByOrderNo(
		ctx,
		orderNo,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSalesOrderNotFound
		}

		return nil, err
	}

	items, err := s.repository.ListItems(
		ctx,
		order.ID,
	)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return toResponse(order), nil
}

func (s *service) List(
	ctx context.Context,
	req *ListRequest,
) (*ListResponse, error) {
	if req == nil {
		req = &ListRequest{}
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	var status *model.SalesOrderStatus

	if req.Status != nil &&
		strings.TrimSpace(*req.Status) != "" {

		value := model.SalesOrderStatus(
			strings.TrimSpace(*req.Status),
		)

		status = &value
	}

	orders, total, err := s.repository.List(
		ctx,
		status,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	items := make(
		[]*Response,
		0,
		len(orders),
	)

	for _, order := range orders {
		orderItems, err := s.repository.ListItems(
			ctx,
			order.ID,
		)
		if err != nil {
			return nil, err
		}

		order.Items = orderItems

		items = append(
			items,
			toResponse(order),
		)
	}

	return &ListResponse{
		Items: items,
		Total: total,
	}, nil
}

func (s *service) Update(
	ctx context.Context,
	id uint,
	req *UpdateRequest,
) (*Response, error) {
	if id == 0 || req == nil {
		return nil, ErrSalesOrderInvalid
	}

	order, err := s.repository.GetByID(
		ctx,
		id,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSalesOrderNotFound
		}

		return nil, err
	}

	if order.Status != model.SalesOrderStatusDraft {
		return nil, ErrSalesOrderNotDraft
	}

	if req.OrderDate != nil {
		if req.OrderDate.IsZero() {
			return nil, ErrSalesOrderInvalid
		}

		order.OrderDate = *req.OrderDate
	}

	if req.Description != nil {
		order.Description = strings.TrimSpace(
			*req.Description,
		)
	}

	if err := s.repository.Update(
		ctx,
		order,
	); err != nil {
		return nil, err
	}

	items, err := s.repository.ListItems(
		ctx,
		order.ID,
	)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return toResponse(order), nil
}

func (s *service) Confirm(
	ctx context.Context,
	id uint,
) error {
	if id == 0 {
		return ErrSalesOrderInvalid
	}

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			order, err := s.repository.GetByIDForUpdateTx(
				txCtx,
				id,
			)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrSalesOrderNotFound
				}

				return err
			}

			if order.Status != model.SalesOrderStatusDraft {
				return ErrSalesOrderNotDraft
			}

			items, err := s.repository.ListItems(
				txCtx,
				order.ID,
			)
			if err != nil {
				return err
			}

			if len(items) == 0 {
				return ErrSalesOrderEmptyItems
			}

			order.Status = model.SalesOrderStatusConfirmed

			return s.repository.UpdateTx(
				txCtx,
				order,
			)
		},
	)
}

func (s *service) Cancel(
	ctx context.Context,
	id uint,
) error {
	if id == 0 {
		return ErrSalesOrderInvalid
	}

	return s.uow.Execute(
		ctx,
		func(txCtx context.Context) error {
			order, err := s.repository.GetByIDForUpdateTx(
				txCtx,
				id,
			)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrSalesOrderNotFound
				}

				return err
			}

			if order.Status == model.SalesOrderStatusCancelled {
				return nil
			}

			if order.Status != model.SalesOrderStatusDraft &&
				order.Status != model.SalesOrderStatusConfirmed {
				return ErrSalesOrderNotDraft
			}

			order.Status = model.SalesOrderStatusCancelled

			return s.repository.UpdateTx(
				txCtx,
				order,
			)
		},
	)
}

func toResponse(
	order *model.SalesOrder,
) *Response {
	items := make(
		[]*ItemResponse,
		0,
		len(order.Items),
	)

	for _, item := range order.Items {
		items = append(
			items,
			&ItemResponse{
				ID:              item.ID,
				SalesOrderID:    item.SalesOrderID,
				ProductID:       item.ProductID,
				OrderedQuantity: item.OrderedQuantity,
				CreatedAt:       item.CreatedAt,
				UpdatedAt:       item.UpdatedAt,
			},
		)
	}

	return &Response{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		OrderNo:     order.OrderNo,
		Status:      order.Status,
		OrderDate:   order.OrderDate,
		Description: order.Description,
		Items:       items,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}