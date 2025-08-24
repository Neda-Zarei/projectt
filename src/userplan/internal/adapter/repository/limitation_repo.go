package repository

import (
	"context"

	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"
)

type limitationRepository struct {
	db *gorm.DB
}

func NewLimitationRepository(db *gorm.DB) planP.LimitationRepository {
	return &limitationRepository{db: db}
}

func (r *limitationRepository) Create(ctx context.Context, limitation *domain.Limitation) error {
	return r.db.WithContext(ctx).Create(limitation).Error
}

func (r *limitationRepository) GetByID(ctx context.Context, id uint) (*domain.Limitation, error) {
	var limitation domain.Limitation
	err := r.db.WithContext(ctx).First(&limitation, id).Error
	return &limitation, err
}

func (r *limitationRepository) GetByTitle(ctx context.Context, title string) (*domain.Limitation, error) {
	var limitation domain.Limitation
	err := r.db.WithContext(ctx).Where("title = ?", title).First(&limitation).Error
	return &limitation, err
}

func (r *limitationRepository) List(ctx context.Context) ([]*domain.Limitation, error) {
	var limitations []*domain.Limitation
	err := r.db.WithContext(ctx).Find(&limitations).Error
	return limitations, err
}

func (r *limitationRepository) Update(ctx context.Context, limitation *domain.Limitation) error {
	return r.db.WithContext(ctx).Save(limitation).Error
}

func (r *limitationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Limitation{}, id).Error
}

func (r *limitationRepository) AssignToPlan(ctx context.Context, planLimitation *domain.PlanLimitation) error {
	return r.db.WithContext(ctx).Create(planLimitation).Error
}

func (r *limitationRepository) GetPlanLimitations(ctx context.Context, planID uint) ([]*domain.PlanLimitation, error) {
	var planLimitations []*domain.PlanLimitation
	err := r.db.WithContext(ctx).
		Preload("Limitation").
		Where("plan_id = ?", planID).
		Find(&planLimitations).Error
	return planLimitations, err
}

func (r *limitationRepository) UpdatePlanLimitation(ctx context.Context, planLimitation *domain.PlanLimitation) error {
	return r.db.WithContext(ctx).
		Where("plan_id = ? AND limitation_id = ?", planLimitation.PlanID, planLimitation.LimitationID).
		Updates(planLimitation).Error
}

func (r *limitationRepository) RemoveFromPlan(ctx context.Context, planID, limitationID uint) error {
	return r.db.WithContext(ctx).
		Where("plan_id = ? AND limitation_id = ?", planID, limitationID).
		Delete(&domain.PlanLimitation{}).Error
}
