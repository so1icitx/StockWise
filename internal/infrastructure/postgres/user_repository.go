package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (repository *userRepository) Create(ctx context.Context, user *domain.User) error {
	model := toUserModel(*user)
	if err := repository.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	*user = toDomainUser(model)
	return nil
}

func (repository *userRepository) Update(ctx context.Context, user *domain.User) error {
	result := repository.db.WithContext(ctx).Model(&userModel{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"name":       user.Name,
			"email":      user.Email,
			"role":       string(user.Role),
			"is_active":  user.IsActive,
			"updated_at": gorm.Expr("now()"),
		})

	return notFoundOnZeroRows(result)
}

func (repository *userRepository) Delete(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&userModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *userRepository) GetByID(ctx context.Context, id domain.ID) (*domain.User, error) {
	var model userModel
	if err := repository.db.WithContext(ctx).First(&model, uint64(id)).Error; err != nil {
		return nil, mapError(err)
	}

	user := toDomainUser(model)
	return &user, nil
}

func (repository *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model userModel
	if err := repository.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		return nil, mapError(err)
	}

	user := toDomainUser(model)
	return &user, nil
}

func (repository *userRepository) List(ctx context.Context, filter application.UserFilter) ([]domain.User, error) {
	query := repository.db.WithContext(ctx).Model(&userModel{}).Order("id ASC")
	if filter.Role != nil {
		query = query.Where("role = ?", string(*filter.Role))
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var models []userModel
	if err := applyListOptions(query, filter.ListOptions).Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(models))
	for _, model := range models {
		users = append(users, toDomainUser(model))
	}

	return users, nil
}
