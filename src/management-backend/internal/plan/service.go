package plan

import (
	"context"
	"errors"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/domain"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/port"
)

var (
	ErrPlanNotFound      = errors.New("plan not found")
	ErrPlanAlreadyExists = errors.New("plan already exists")
)

type service struct {
	repo port.Repository
}

func NewService(repo port.Repository) port.Service {
	return &service{repo: repo}
}

func (s *service) CreatePlan(ctx context.Context, plan *domain.Plan) error {
	existing, _ := s.repo.GetByName(ctx, plan.Name)
	if existing != nil {
		return ErrPlanAlreadyExists
	}

	return s.repo.Create(ctx, plan)
}

func (s *service) GetPlanByID(ctx context.Context, id uint) (*domain.Plan, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetPlanByName(ctx context.Context, name string) (*domain.Plan, error) {
	return s.repo.GetByName(ctx, name)
}

func (s *service) UpdatePlan(ctx context.Context, plan *domain.Plan) error {
	return s.repo.Update(ctx, plan)
}

func (s *service) DeletePlan(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) ListPlans(ctx context.Context, limit, offset int) ([]*domain.Plan, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *service) TogglePlanActive(ctx context.Context, id uint) error {
	plan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrPlanNotFound
	}

	plan.IsActive = !plan.IsActive
	return s.repo.Update(ctx, plan)
}
