package port

import (
	"context"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
)

type Service interface {
	AssignPlan(ctx context.Context, req *domain.AssignPlanRequest) error
	GetUserPlan(ctx context.Context, userID uint) (*domain.UserPlan, error)
	RenewUserPlan(ctx context.Context, req *domain.RenewPlanRequest) error
	CancelUserPlan(ctx context.Context, userID uint) error
	GetUserPlanHistory(ctx context.Context, userID uint) ([]*domain.UserPlan, error)

	CreatePlan(ctx context.Context, plan *domain.Plan) error
	GetPlanByID(ctx context.Context, id uint) (*domain.Plan, error)
	GetPlanByTitle(ctx context.Context, title string) (*domain.Plan, error)
	UpdatePlan(ctx context.Context, plan *domain.Plan) error
	DeletePlan(ctx context.Context, id uint) error
	ListPlans(ctx context.Context, includeInactive bool) ([]*domain.Plan, error)

	SetPlanPrice(ctx context.Context, planID uint, months int, price int) error
	GetPlanPrices(ctx context.Context, planID uint) ([]*domain.Price, error)

	CreateLimitation(ctx context.Context, limitation *domain.Limitation) error
	AssignLimitationToPlan(ctx context.Context, planID, limitationID uint, value int) error
	GetPlanLimitations(ctx context.Context, planID uint) ([]*domain.PlanLimitation, error)

	ExpirePlans(ctx context.Context) error
	GetExpiringPlans(ctx context.Context, daysThreshold int) ([]*domain.UserPlan, error)
}

type PlanRepository interface {
	Create(ctx context.Context, plan *domain.Plan) error
	GetByID(ctx context.Context, id uint) (*domain.Plan, error)
	GetByTitle(ctx context.Context, title string) (*domain.Plan, error)
	Update(ctx context.Context, plan *domain.Plan) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, includeInactive bool) ([]*domain.Plan, error)
}

type UserPlanRepository interface {
	Create(ctx context.Context, userPlan *domain.UserPlan) error
	GetByUserID(ctx context.Context, userID uint) (*domain.UserPlan, error)
	GetActiveByUserID(ctx context.Context, userID uint) (*domain.UserPlan, error)
	Update(ctx context.Context, userPlan *domain.UserPlan) error
	GetUserHistory(ctx context.Context, userID uint) ([]*domain.UserPlan, error)
	ExpirePlans(ctx context.Context) error
	GetExpiringPlans(ctx context.Context, daysThreshold int) ([]*domain.UserPlan, error)
	SoftDelete(ctx context.Context, id uint) error
}

type PriceRepository interface {
	Create(ctx context.Context, price *domain.Price) error
	GetByPlanID(ctx context.Context, planID uint) ([]*domain.Price, error)
	GetByPlanIDAndMonth(ctx context.Context, planID uint, month int) (*domain.Price, error)
	Update(ctx context.Context, price *domain.Price) error
	Delete(ctx context.Context, planID uint, month int) error
}

type LimitationRepository interface {
	Create(ctx context.Context, limitation *domain.Limitation) error
	GetByID(ctx context.Context, id uint) (*domain.Limitation, error)
	GetByTitle(ctx context.Context, title string) (*domain.Limitation, error)
	List(ctx context.Context) ([]*domain.Limitation, error)
	Update(ctx context.Context, limitation *domain.Limitation) error
	Delete(ctx context.Context, id uint) error

	//plan limitation operations
	AssignToPlan(ctx context.Context, planLimitation *domain.PlanLimitation) error
	GetPlanLimitations(ctx context.Context, planID uint) ([]*domain.PlanLimitation, error)
	UpdatePlanLimitation(ctx context.Context, planLimitation *domain.PlanLimitation) error
	RemoveFromPlan(ctx context.Context, planID, limitationID uint) error
}
