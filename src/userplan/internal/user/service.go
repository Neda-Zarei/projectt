package user

import (
	"context"

	userD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/domain"
	userP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/port"
)

type service struct {
	repo userP.Repo
}

func New(r userP.Repo) userP.Service {
	return &service{repo: r}
}

func (s *service) ListUsers(ctx context.Context, uf *userD.UserFilter) (*userD.PaginatedUsers, error) {
	const defaultLimit = 20
	return s.repo.List(ctx, uf, defaultLimit, 0)
}

func (s *service) CreateUser(ctx context.Context, u *userD.User) error {
	return s.repo.Create(ctx, u)
}

func (s *service) UpdateUser(ctx context.Context, u *userD.User) error {
	return s.repo.Update(ctx, u)
}

func (s *service) SetUserActive(ctx context.Context, ua *userD.UserActivation) error {
	return s.repo.ToggleActive(ctx, ua.UserID)
}

func (s *service) GetByID(ctx context.Context, id uint) (*userD.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*userD.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) GetByOauthID(ctx context.Context, oauthID string) (*userD.User, error) {
	return s.repo.GetByOauthID(ctx, oauthID)
}
