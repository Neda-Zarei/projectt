package user

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/user/domain"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/user/port"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type service struct {
	repo port.Repository
}

func NewService(repo port.Repository) port.Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, user *domain.AdminUser, password string) error {
	existing, _ := s.repo.GetByEmail(ctx, user.Email)
	if existing != nil {
		return ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.IsActive = true
	user.Role = "admin"

	return s.repo.Create(ctx, user)
}

func (s *service) Authenticate(ctx context.Context, email, password string) (*domain.AdminUser, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	user.LastLogin = time.Now()
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, id uint) (*domain.AdminUser, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (*domain.AdminUser, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *service) UpdateUser(ctx context.Context, user *domain.AdminUser) error {
	return s.repo.Update(ctx, user)
}

func (s *service) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) ListUsers(ctx context.Context, limit, offset int, filters map[string]string) ([]*domain.AdminUser, error) {
	return s.repo.List(ctx, limit, offset, filters)
}

func (s *service) ToggleUserActive(ctx context.Context, id uint) error {
	return s.repo.ToggleActive(ctx, id)
}

func (s *service) ChangePassword(ctx context.Context, id uint, currentPassword, newPassword string) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return ErrInvalidPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.Update(ctx, user)
}
