package productinstance

import "context"

type Repository interface {
	Create(
		ctx context.Context,
		instance *ProductInstance,
	) error

	GetBySN(
		ctx context.Context,
		sn string,
	) (
		*ProductInstance,
		error,
	)

	Update(
		ctx context.Context,
		instance *ProductInstance,
	) error
}