package port

import (
	"context"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/admin/domain"
)

type Service interface {
	CreateUser(ctx context.Context, user *domain.AdminUser, password string) error
	Authenticate(ctx context.Context, email, password string) (*domain.AdminUser, error)
	GetUserByID(ctx context.Context, id uint) (*domain.AdminUser, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.AdminUser, error)
	UpdateUser(ctx context.Context, user *domain.AdminUser) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, limit, offset int, filters map[string]string) ([]*domain.AdminUser, error)
	ToggleUserActive(ctx context.Context, id uint) error
	ChangePassword(ctx context.Context, id uint, currentPassword, newPassword string) error
}

type Repository interface {
	Create(ctx context.Context, user *domain.AdminUser) error
	GetByID(ctx context.Context, id uint) (*domain.AdminUser, error)
	GetByEmail(ctx context.Context, email string) (*domain.AdminUser, error)
	Update(ctx context.Context, user *domain.AdminUser) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int, filters map[string]string) ([]*domain.AdminUser, error)
	ToggleActive(ctx context.Context, id uint) error
}
