package user

import (
	"context"

	userD "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/domain"
	userP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/port"
)

type service struct {
	repo userP.Repo
}

func New(r userP.Repo) userP.Service { return &service{repo: r} }

func (s *service) ListUsers(ctx context.Context, uf *userD.UserFilter) (*userD.PaginatedUsers, error) {
	return &userD.PaginatedUsers{}, nil
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
