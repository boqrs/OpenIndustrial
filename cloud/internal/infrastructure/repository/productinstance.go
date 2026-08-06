package repository

import (
	"context"

	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/productinstance"
	"gorm.io/gorm"
)

type ProductInstanceRepository struct {
	db *gorm.DB
}

func NewProductInstanceRepository(db *gorm.DB) *ProductInstanceRepository {
	return &ProductInstanceRepository{
		db: db,
	}
}

func (r *ProductInstanceRepository) Create(
	ctx context.Context,
	entity *productinstance.ProductInstance,
) error {
	model := ProductInstanceModel{
		ID:        entity.ID,
		SN:        entity.SN,
		ProductID: entity.ProductID,
		OrgID:     entity.OrgID,
		State:     entity.State,
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *ProductInstanceRepository) GetBySN(
	ctx context.Context,
	sn string,
) (*productinstance.ProductInstance, error) {
	var model ProductInstanceModel

	err := r.db.WithContext(ctx).Where("sn = ?", sn).First(&model).Error

	if err != nil {
		return nil, err
	}

	return &productinstance.ProductInstance{
		ID:        model.ID,
		SN:        model.SN,
		ProductID: model.ProductID,
		OrgID:     model.OrgID,
		State:     model.State,
	}, nil
}

func (r *ProductInstanceRepository) Update(
	ctx context.Context,
	entity *productinstance.ProductInstance,
) error {
	model := ProductInstanceModel{
		ID:        entity.ID,
		SN:        entity.SN,
		ProductID: entity.ProductID,
		OrgID:     entity.OrgID,
		State:     entity.State,
	}
	return r.db.WithContext(ctx).Save(&model).Error
}