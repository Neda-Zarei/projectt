package repository

import (
	"context"

	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"
)

type priceRepository struct {
	db *gorm.DB
}

func NewPriceRepository(db *gorm.DB) planP.PriceRepository {
	return &priceRepository{db: db}
}

func (r *priceRepository) Create(ctx context.Context, price *domain.Price) error {
	return r.db.WithContext(ctx).Create(price).Error
}

func (r *priceRepository) GetByPlanID(ctx context.Context, planID uint) ([]*domain.Price, error) {
	var prices []*domain.Price
	err := r.db.WithContext(ctx).Where("plan_id = ?", planID).Find(&prices).Error
	return prices, err
}

func (r *priceRepository) GetByPlanIDAndMonth(ctx context.Context, planID uint, month int) (*domain.Price, error) {
	var price domain.Price
	err := r.db.WithContext(ctx).Where("plan_id = ? AND month = ?", planID, month).First(&price).Error
	return &price, err
}

func (r *priceRepository) Update(ctx context.Context, price *domain.Price) error {
	return r.db.WithContext(ctx).Save(price).Error
}

func (r *priceRepository) Delete(ctx context.Context, planID uint, month int) error {
	return r.db.WithContext(ctx).Where("plan_id = ? AND month = ?", planID, month).Delete(&domain.Price{}).Error
}
