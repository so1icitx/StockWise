package application

import (
	"context"

	"github.com/so1icitx/StockWise/internal/domain"
)

// ProductService handles product-related business operations.
type ProductService struct {
	provider RepositoryProvider
}

// NewProductService creates a product service.
func NewProductService(provider RepositoryProvider) *ProductService {
	return &ProductService{provider: provider}
}

// Create creates an active product with a unique SKU.
func (service *ProductService) Create(ctx context.Context, input CreateProductInput) (*domain.Product, error) {
	name, err := requireText(input.Name, "name")
	if err != nil {
		return nil, err
	}
	sku, err := requireText(input.SKU, "sku")
	if err != nil {
		return nil, err
	}
	unitOfMeasure, err := requireText(input.UnitOfMeasure, "unit of measure")
	if err != nil {
		return nil, err
	}
	if input.MinStockThreshold < 0 {
		return nil, validationError("minimum stock threshold cannot be negative")
	}

	repos := repositories(service.provider)
	if _, err := requireActiveCategory(ctx, repos.Categories, input.CategoryID); err != nil {
		return nil, err
	}

	existing, err := repos.Products.GetBySKU(ctx, sku)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if existing != nil {
		return nil, conflictError(ErrDuplicateSKU, "product sku must be unique")
	}

	product := domain.Product{
		Name:              name,
		SKU:               sku,
		CategoryID:        input.CategoryID,
		UnitOfMeasure:     unitOfMeasure,
		MinStockThreshold: input.MinStockThreshold,
		IsActive:          true,
	}
	if err := repos.Products.Create(ctx, &product); err != nil {
		return nil, err
	}

	return &product, nil
}

// Update updates product details while preserving SKU uniqueness.
func (service *ProductService) Update(ctx context.Context, id domain.ID, input UpdateProductInput) (*domain.Product, error) {
	if err := requireID(id, "product id"); err != nil {
		return nil, err
	}
	name, err := requireText(input.Name, "name")
	if err != nil {
		return nil, err
	}
	sku, err := requireText(input.SKU, "sku")
	if err != nil {
		return nil, err
	}
	unitOfMeasure, err := requireText(input.UnitOfMeasure, "unit of measure")
	if err != nil {
		return nil, err
	}
	if input.MinStockThreshold < 0 {
		return nil, validationError("minimum stock threshold cannot be negative")
	}

	repos := repositories(service.provider)
	product, err := repos.Products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.IsActive {
		if _, err := requireActiveCategory(ctx, repos.Categories, input.CategoryID); err != nil {
			return nil, err
		}
	} else if err := requireID(input.CategoryID, "category id"); err != nil {
		return nil, err
	}

	existing, err := repos.Products.GetBySKU(ctx, sku)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if existing != nil && existing.ID != id {
		return nil, conflictError(ErrDuplicateSKU, "product sku must be unique")
	}

	product.Name = name
	product.SKU = sku
	product.CategoryID = input.CategoryID
	product.UnitOfMeasure = unitOfMeasure
	product.MinStockThreshold = input.MinStockThreshold
	product.IsActive = input.IsActive

	if err := repos.Products.Update(ctx, product); err != nil {
		return nil, err
	}

	return repos.Products.GetByID(ctx, id)
}

// Delete hard-deletes unused products and soft-deactivates products with stock or confirmed movements.
func (service *ProductService) Delete(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "product id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	if _, err := repos.Products.GetByID(ctx, id); err != nil {
		return err
	}

	hasConfirmedMovements, err := repos.Products.HasConfirmedMovements(ctx, id)
	if err != nil {
		return err
	}
	if hasConfirmedMovements {
		return repos.Products.SetActive(ctx, id, false)
	}

	stockRows, err := repos.Stock.GetByProductAcrossWarehouses(ctx, id)
	if err != nil {
		return err
	}
	if len(stockRows) > 0 {
		return repos.Products.SetActive(ctx, id, false)
	}

	return repos.Products.Delete(ctx, id)
}

// Activate marks a product as active when its category is active.
func (service *ProductService) Activate(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "product id"); err != nil {
		return err
	}

	repos := repositories(service.provider)
	product, err := repos.Products.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if _, err := requireActiveCategory(ctx, repos.Categories, product.CategoryID); err != nil {
		return err
	}

	return repos.Products.SetActive(ctx, id, true)
}

// Deactivate marks a product as inactive.
func (service *ProductService) Deactivate(ctx context.Context, id domain.ID) error {
	if err := requireID(id, "product id"); err != nil {
		return err
	}

	return repositories(service.provider).Products.SetActive(ctx, id, false)
}

// GetByID returns a product by identifier.
func (service *ProductService) GetByID(ctx context.Context, id domain.ID) (*domain.Product, error) {
	if err := requireID(id, "product id"); err != nil {
		return nil, err
	}

	return repositories(service.provider).Products.GetByID(ctx, id)
}

// List returns products matching the filter.
func (service *ProductService) List(ctx context.Context, filter ProductFilter) ([]domain.Product, error) {
	return repositories(service.provider).Products.List(ctx, filter)
}
