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
