package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/user/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.AdminUser) error
	GetByID(ctx context.Context, id uint) (*domain.AdminUser, error)
	GetByEmail(ctx context.Context, email string) (*domain.AdminUser, error)
	Update(ctx context.Context, user *domain.AdminUser) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int, filters map[string]string) ([]*domain.AdminUser, error)
	ToggleActive(ctx context.Context, id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.AdminUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*domain.AdminUser, error) {
	var user domain.AdminUser
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.AdminUser, error) {
	var user domain.AdminUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *domain.AdminUser) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.AdminUser{}, id).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int, filters map[string]string) ([]*domain.AdminUser, error) {
	query := r.db.WithContext(ctx).Model(&domain.AdminUser{})

	if name, ok := filters["name"]; ok {
		query = query.Where("first_name LIKE ? OR last_name LIKE ?",
			fmt.Sprintf("%%%s%%", name),
			fmt.Sprintf("%%%s%%", name))
	}

	if email, ok := filters["email"]; ok {
		query = query.Where("email LIKE ?", fmt.Sprintf("%%%s%%", email))
	}

	if phone, ok := filters["phone"]; ok {
		query = query.Where("phone LIKE ?", fmt.Sprintf("%%%s%%", phone))
	}

	var users []*domain.AdminUser
	err := query.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

func (r *userRepository) ToggleActive(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&domain.AdminUser{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).Error
}
