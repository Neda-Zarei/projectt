package plan

import (
	"context"

	planD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/domain"
	planP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/plan/port"
)

type service struct {
	repo planP.Repo
}

func New(r planP.Repo) planP.Service { return &service{repo: r} }

func (s *service) AssignPlan(ctx context.Context, req *planD.AssignPlanRequest) error {
	return s.repo.AssignPlan(ctx, &planD.UserPlan{
		UserID: req.UserID,
		PlanID: req.PlanID,
	})
}

func (s *service) GetUserPlan(ctx context.Context, userID uint) (*planD.UserPlan, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) RenewUserPlan(ctx context.Context, req *planD.RenewPlanRequest) error {
	return s.repo.RenewPlan(ctx, req.UserID, req.EndDate)
}

// Plan management methods
func (s *service) Create(ctx context.Context, plan *planD.Plan) error {
	return s.repo.Create(ctx, plan)
}

func (s *service) GetByID(ctx context.Context, id uint) (*planD.Plan, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByName(ctx context.Context, name string) (*planD.Plan, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *service) Update(ctx context.Context, plan *planD.Plan) error {
	return s.repo.Update(ctx, plan)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) ToggleActive(ctx context.Context, id uint) error {
	return s.repo.ToggleActive(ctx, id)
}

func (s *service) ListActive(ctx context.Context) ([]*planD.Plan, error) {
	return s.repo.ListActive(ctx)
}

// Expiration management methods
func (s *service) ExpirePlans(ctx context.Context) error {
	return s.repo.ExpirePlans(ctx)
}

func (s *service) GetExpiringPlans(ctx context.Context, daysThreshold int) ([]*planD.UserPlan, error) {
	return s.repo.GetExpiringPlans(ctx, daysThreshold)
}
