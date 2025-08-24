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

func (r *userRepository) GetByOauthID(ctx context.Context, oauthID string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("oauth_id = ?", oauthID).First(&user).Error
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
		Update("active", gorm.Expr("NOT active")).Error
}

func (r *userRepository) List(ctx context.Context, filter *domain.UserFilter, limit, offset int) (*domain.PaginatedUsers, error) {
	var users []*domain.User
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.User{})

	// Apply filters
	if filter != nil {
		if filter.Name != "" {
			query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
		}
		if filter.Email != "" {
			query = query.Where("email ILIKE ?", "%"+filter.Email+"%")
		}
		if filter.Phone != "" {
			query = query.Where("phone ILIKE ?", "%"+filter.Phone+"%")
		}
		if filter.UserID != 0 {
			query = query.Where("id = ?", filter.UserID)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Get paginated records
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	return &domain.PaginatedUsers{
		Users: users,
		Total: total,
		Size:  int64(limit),
		Page:  int64(offset/limit + 1),
	}, nil
}
