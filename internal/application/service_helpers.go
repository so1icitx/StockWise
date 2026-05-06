package application

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/so1icitx/StockWise/internal/domain"
)

func repositories(provider RepositoryProvider) Repositories {
	return provider.Repositories()
}

func nowUTC() time.Time {
	return time.Now().UTC()
}

func clean(value string) string {
	return strings.TrimSpace(value)
}

func requireText(value string, field string) (string, error) {
	cleaned := clean(value)
	if cleaned == "" {
		return "", validationError(field + " is required")
	}

	return cleaned, nil
}

func requireID(id domain.ID, field string) error {
	if id.IsZero() {
		return validationError(field + " is required")
	}

	return nil
}

func ensurePositiveQuantity(quantity int64) error {
	if !domain.IsPositiveQuantity(quantity) {
		return validationError("quantity must be positive")
	}

	return nil
}

func ensureNonNegativePrice(priceCents int64) error {
	if !domain.IsNonNegativePriceCents(priceCents) {
		return validationError("price cannot be negative")
	}

	return nil
}

func ensureNotFoundOrNil(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNotFound) {
		return nil
	}

	return err
}

func isNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func requireActiveUser(ctx context.Context, repository UserRepository, id domain.ID) (*domain.User, error) {
	if err := requireID(id, "user id"); err != nil {
		return nil, err
	}

	user, err := repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, businessRuleError(ErrBusinessRule, "inactive user cannot perform operations")
	}

	return user, nil
}

func requireActiveWarehouse(ctx context.Context, repository WarehouseRepository, id domain.ID) (*domain.Warehouse, error) {
	if err := requireID(id, "warehouse id"); err != nil {
		return nil, err
	}

	warehouse, err := repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !warehouse.CanBeUsedForOperations() {
		return nil, businessRuleError(ErrInactiveWarehouse, "inactive warehouse cannot be used for operations")
	}

	return warehouse, nil
}

func requireActiveCategory(ctx context.Context, repository CategoryRepository, id domain.ID) (*domain.Category, error) {
	if err := requireID(id, "category id"); err != nil {
		return nil, err
	}

	category, err := repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !category.CanBeUsedForProducts() {
		return nil, businessRuleError(ErrInactiveCategory, "inactive category cannot be assigned to active products")
	}

	return category, nil
}

func requireActiveProduct(ctx context.Context, repository ProductRepository, id domain.ID) (*domain.Product, error) {
	if err := requireID(id, "product id"); err != nil {
		return nil, err
	}

	product, err := repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !product.CanBeUsedForOperations() {
		return nil, businessRuleError(ErrInactiveProduct, "inactive product cannot be used for operations")
	}

	return product, nil
}
