package repository

import (
	"context"

	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/domain"
	userP "hamgit.ir/arcaptcha/arcaptcha-dumbledore/userplan/internal/user/port"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) userP.Repo {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *userRepository) GetByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("external_id = ?", externalID).First(&user).Error
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) ToggleActive(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).Error
}
