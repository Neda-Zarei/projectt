package port

import (
	"context"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/domain"
)

type Service interface {
	ListUsers(context.Context, *domain.UserFilter) (*domain.PaginatedUsers, error)
	CreateUser(context.Context, *domain.User) error
	UpdateUser(context.Context, *domain.User) error
	SetUserActive(context.Context, *domain.UserActivation) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByOauthID(ctx context.Context, oauthID string) (*domain.User, error)
}

type Repo interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	GetByOauthID(ctx context.Context, oauthID string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	ToggleActive(ctx context.Context, id uint) error
	List(ctx context.Context, filter *domain.UserFilter, limit, offset int) (*domain.PaginatedUsers, error)
}
