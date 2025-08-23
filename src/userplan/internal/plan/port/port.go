package port

import (
	"context"
	"time"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
)

type Service interface {
	AssignPlan(ctx context.Context, req *domain.AssignPlanRequest) error
	GetUserPlan(ctx context.Context, userID uint) (*domain.UserPlan, error)
	RenewUserPlan(ctx context.Context, req *domain.RenewPlanRequest) error
}

type Repo interface {
	PlanRepository
	UserPlanRepository
}

type PlanRepository interface {
	Create(ctx context.Context, plan *domain.Plan) error
	GetByID(ctx context.Context, id uint) (*domain.Plan, error)
	GetByName(ctx context.Context, name string) (*domain.Plan, error)
	Update(ctx context.Context, plan *domain.Plan) error
	ToggleActive(ctx context.Context, id uint) error
	ListActive(ctx context.Context) ([]*domain.Plan, error)
}

type UserPlanRepository interface {
	AssignPlan(ctx context.Context, userPlan *domain.UserPlan) error
	GetByUserID(ctx context.Context, userID uint) (*domain.UserPlan, error)
	RenewPlan(ctx context.Context, userID uint, newEndDate time.Time) error
	CancelPlan(ctx context.Context, userID uint) error
	GetHistory(ctx context.Context, userID uint) ([]*domain.PlanHistory, error)
	RecordHistory(ctx context.Context, history *domain.PlanHistory) error
}
