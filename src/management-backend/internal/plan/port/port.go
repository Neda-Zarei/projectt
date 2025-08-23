package port

import (
	"context"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/domain"
)

type Service interface {
	CreatePlan(ctx context.Context, plan *domain.Plan) error
	GetPlanByID(ctx context.Context, id uint) (*domain.Plan, error)
	GetPlanByName(ctx context.Context, name string) (*domain.Plan, error)
	UpdatePlan(ctx context.Context, plan *domain.Plan) error
	DeletePlan(ctx context.Context, id uint) error
	ListPlans(ctx context.Context, limit, offset int) ([]*domain.Plan, error)
	TogglePlanActive(ctx context.Context, id uint) error
}
