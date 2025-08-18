package repository

import (
	"context"

	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/domain"
)

type PlanRepository interface {
	Create(ctx context.Context, plan *domain.Plan) error
	GetByID(ctx context.Context, id uint) (*domain.Plan, error)
	GetByName(ctx context.Context, name string) (*domain.Plan, error)
	Update(ctx context.Context, plan *domain.Plan) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*domain.Plan, error)
}

type planRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) PlanRepository {
	return &planRepository{db: db}
}

func (r *planRepository) Create(ctx context.Context, plan *domain.Plan) error {
	return r.db.WithContext(ctx).Create(plan).Error
}

func (r *planRepository) GetByID(ctx context.Context, id uint) (*domain.Plan, error) {
	var plan domain.Plan
	err := r.db.WithContext(ctx).First(&plan, id).Error
	return &plan, err
}

func (r *planRepository) GetByName(ctx context.Context, name string) (*domain.Plan, error) {
	var plan domain.Plan
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&plan).Error
	return &plan, err
}

func (r *planRepository) Update(ctx context.Context, plan *domain.Plan) error {
	return r.db.WithContext(ctx).Save(plan).Error
}

func (r *planRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Plan{}, id).Error
}

func (r *planRepository) List(ctx context.Context, limit, offset int) ([]*domain.Plan, error) {
	var plans []*domain.Plan
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&plans).Error
	return plans, err
}
