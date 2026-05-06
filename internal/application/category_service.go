package application

import (
	"context"

	"github.com/so1icitx/StockWise/internal/domain"
)

// CategoryService handles category-related business operations.
type CategoryService struct {
	provider RepositoryProvider
}

// NewCategoryService creates a category service.
func NewCategoryService(provider RepositoryProvider) *CategoryService {
	return &CategoryService{provider: provider}
}

// Create creates an active category with a unique name.
func (service *CategoryService) Create(ctx context.Context, input CreateCategoryInput) (*domain.Category, error) {
	name, err := requireText(input.Name, "name")
	if err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	existing, err := repos.Categories.GetByName(ctx, name)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if existing != nil {
		return nil, conflictError(ErrConflict, "category name already exists")
	}

	category := domain.Category{
		Name:        name,
		Description: clean(input.Description),
		IsActive:    true,
	}
	if err := repos.Categories.Create(ctx, &category); err != nil {
		return nil, err
	}

	return &category, nil
}

// Update updates category details.
func (service *CategoryService) Update(ctx context.Context, id domain.ID, input UpdateCategoryInput) (*domain.Category, error) {
	if err := requireID(id, "category id"); err != nil {
		return nil, err
	}
	name, err := requireText(input.Name, "name")
	if err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	category, err := repos.Categories.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing, err := repos.Categories.GetByName(ctx, name)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if existing != nil && existing.ID != id {
		return nil, conflictError(ErrConflict, "category name already exists")
	}

	category.Name = name
	category.Description = clean(input.Description)
	category.IsActive = input.IsActive

	if err := repos.Categories.Update(ctx, category); err != nil {
		return nil, err
	}

	return repos.Categories.GetByID(ctx, id)
}

// Delete deletes a category only when it has no active products.
func (service *CategoryService) Delete(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "category id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	category, err := repos.Categories.GetByID(ctx, id)
	if err != nil {
		return err
	}

	hasActiveProducts, err := repos.Categories.HasActiveProducts(ctx, id)
	if err != nil {
		return err
	}
	if hasActiveProducts {
		return businessRuleError(ErrDeleteRestricted, "category with active products cannot be deleted")
	}

	products, err := repos.Products.List(ctx, ProductFilter{CategoryID: &id, ListOptions: ListOptions{Limit: 1}})
	if err != nil {
		return err
	}
	if len(products) > 0 {
		category.IsActive = false
		return repos.Categories.Update(ctx, category)
	}

	return repos.Categories.Delete(ctx, id)
}

// GetByID returns a category by identifier.
func (service *CategoryService) GetByID(ctx context.Context, id domain.ID) (*domain.Category, error) {
	if err := requireID(id, "category id"); err != nil {
		return nil, err
	}

	return repositories(service.provider).Categories.GetByID(ctx, id)
}

// List returns categories matching the filter.
func (service *CategoryService) List(ctx context.Context, filter CategoryFilter) ([]domain.Category, error) {
	return repositories(service.provider).Categories.List(ctx, filter)
}
