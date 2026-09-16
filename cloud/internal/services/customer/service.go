package customer

import (
	"context"
	"errors"
	"strings"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
)

var (
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrCustomerCodeExists = errors.New("customer code already exists")
	ErrCustomerInvalid    = errors.New("invalid customer")
	ErrCustomerInactive   = errors.New("customer is inactive")
)

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) Create(
	ctx context.Context,
	req *CreateRequest,
) (*Response, error) {
	if req == nil {
		return nil, ErrCustomerInvalid
	}

	code := strings.TrimSpace(req.Code)
	name := strings.TrimSpace(req.Name)
	customerType := strings.TrimSpace(req.Type)

	if code == "" || name == "" || customerType == "" {
		return nil, ErrCustomerInvalid
	}

	existing, err := s.repository.GetByCode(ctx, code)
	if err == nil && existing != nil {
		return nil, ErrCustomerCodeExists
	}

	customer := &model.Customer{
		Code: code,
		Name: name,
		Type: model.CustomerType(customerType),

		ContactName:  strings.TrimSpace(req.ContactName),
		ContactEmail: strings.TrimSpace(req.ContactEmail),
		ContactPhone: strings.TrimSpace(req.ContactPhone),
		Address:      strings.TrimSpace(req.Address),

		Status: model.CustomerStatusActive,
	}

	if err := s.repository.Create(ctx, customer); err != nil {
		return nil, err
	}

	return toResponse(customer), nil
}

func (s *service) GetByID(
	ctx context.Context,
	id uint,
) (*Response, error) {
	customer, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	return toResponse(customer), nil
}

func (s *service) GetByCode(
	ctx context.Context,
	code string,
) (*Response, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, ErrCustomerInvalid
	}

	customer, err := s.repository.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	return toResponse(customer), nil
}

func (s *service) List(
	ctx context.Context,
	req *ListRequest,
) (*ListResponse, error) {
	offset := 0
	limit := 20

	if req != nil {
		if req.Offset > 0 {
			offset = req.Offset
		}

		if req.Limit > 0 {
			limit = req.Limit
		}
	}

	customers, total, err := s.repository.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	items := make([]*Response, 0, len(customers))

	for _, customer := range customers {
		items = append(items, toResponse(customer))
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
	if req == nil {
		return nil, ErrCustomerInvalid
	}

	customer, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, ErrCustomerNotFound
	}

	if customer.Status != model.CustomerStatusActive {
		return nil, ErrCustomerInactive
	}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, ErrCustomerInvalid
		}
		customer.Name = name
	}

	if req.Type != nil {
		customer.Type = model.CustomerType(strings.TrimSpace(*req.Type))
	}

	if req.ContactName != nil {
		customer.ContactName = strings.TrimSpace(*req.ContactName)
	}

	if req.ContactEmail != nil {
		customer.ContactEmail = strings.TrimSpace(*req.ContactEmail)
	}

	if req.ContactPhone != nil {
		customer.ContactPhone = strings.TrimSpace(*req.ContactPhone)
	}

	if req.Address != nil {
		customer.Address = strings.TrimSpace(*req.Address)
	}

	if err := s.repository.Update(ctx, customer); err != nil {
		return nil, err
	}

	return toResponse(customer), nil
}

func (s *service) Deactivate(
	ctx context.Context,
	id uint,
) error {
	customer, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if customer == nil {
		return ErrCustomerNotFound
	}

	if customer.Status == model.CustomerStatusInactive {
		return nil
	}

	customer.Status = model.CustomerStatusInactive

	return s.repository.Update(ctx, customer)
}

func toResponse(customer *model.Customer) *Response {
	if customer == nil {
		return nil
	}

	return &Response{
		ID: customer.ID,

		Code: customer.Code,
		Name: customer.Name,
		Type: customer.Type,

		ContactName:  customer.ContactName,
		ContactEmail: customer.ContactEmail,
		ContactPhone: customer.ContactPhone,
		Address:      customer.Address,

		Status: customer.Status,

		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
}
