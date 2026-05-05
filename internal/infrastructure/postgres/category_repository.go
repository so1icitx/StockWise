package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func (repository *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	model := toCategoryModel(*category)
	if err := repository.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	*category = toDomainCategory(model)
	return nil
}

func (repository *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	result := repository.db.WithContext(ctx).Model(&categoryModel{}).
		Where("id = ?", category.ID).
		Updates(map[string]any{
			"name":        category.Name,
			"description": category.Description,
			"is_active":   category.IsActive,
			"updated_at":  gorm.Expr("now()"),
		})

	return notFoundOnZeroRows(result)
}

func (repository *categoryRepository) Delete(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&categoryModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *categoryRepository) GetByID(ctx context.Context, id domain.ID) (*domain.Category, error) {
	var model categoryModel
	if err := repository.db.WithContext(ctx).First(&model, uint64(id)).Error; err != nil {
		return nil, mapError(err)
	}

	category := toDomainCategory(model)
	return &category, nil
}

func (repository *categoryRepository) GetByName(ctx context.Context, name string) (*domain.Category, error) {
	var model categoryModel
	if err := repository.db.WithContext(ctx).Where("name = ?", name).First(&model).Error; err != nil {
		return nil, mapError(err)
	}

	category := toDomainCategory(model)
	return &category, nil
}

func (repository *categoryRepository) List(ctx context.Context, filter application.CategoryFilter) ([]domain.Category, error) {
	query := repository.db.WithContext(ctx).Model(&categoryModel{}).Order("id ASC")
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var models []categoryModel
	if err := applyListOptions(query, filter.ListOptions).Find(&models).Error; err != nil {
		return nil, err
	}

	categories := make([]domain.Category, 0, len(models))
	for _, model := range models {
		categories = append(categories, toDomainCategory(model))
	}

	return categories, nil
}

func (repository *categoryRepository) HasActiveProducts(ctx context.Context, id domain.ID) (bool, error) {
	var count int64
	err := repository.db.WithContext(ctx).Model(&productModel{}).
		Where("category_id = ? AND is_active = true", id).
		Count(&count).Error

	return count > 0, err
}
