package repository

import (
	"context"

	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"
)

type planRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) planP.PlanRepository {
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

func (r *planRepository) ToggleActive(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&domain.Plan{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).Error
}

func (r *planRepository) ListActive(ctx context.Context) ([]*domain.Plan, error) {
	var plans []*domain.Plan
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&plans).Error
	return plans, err
}

func (r *planRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Plan{}, id).Error
}

func (r *planRepository) GetByTitle(ctx context.Context, title string) (*domain.Plan, error) {
	var plan domain.Plan
	err := r.db.WithContext(ctx).Where("title = ?", title).First(&plan).Error
	return &plan, err
}

func (r *planRepository) List(ctx context.Context, includeInactive bool) ([]*domain.Plan, error) {
	var plans []*domain.Plan
	query := r.db.WithContext(ctx)

	if !includeInactive {
		query = query.Where("custom = ? OR payg = ?", true, true)
	}

	err := query.Find(&plans).Error
	return plans, err
}
