package application

import (
	"context"

	"github.com/so1icitx/StockWise/internal/domain"
)

// MovementService reads product and warehouse movement history.
type MovementService struct {
	provider RepositoryProvider
}

// NewMovementService creates a movement service.
func NewMovementService(provider RepositoryProvider) *MovementService {
	return &MovementService{provider: provider}
}

// ListByProduct returns completed movement history for a product.
func (service *MovementService) ListByProduct(ctx context.Context, productID domain.ID, options ListOptions) ([]MovementRecord, error) {
	if err := requireID(productID, "product id"); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	if _, err := repos.Products.GetByID(ctx, productID); err != nil {
		return nil, err
	}

	return repos.Movements.ListByProduct(ctx, productID, options)
}

// ListByWarehouse returns completed movement history for a warehouse.
func (service *MovementService) ListByWarehouse(ctx context.Context, warehouseID domain.ID, options ListOptions) ([]MovementRecord, error) {
	if err := requireID(warehouseID, "warehouse id"); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	if _, err := repos.Warehouses.GetByID(ctx, warehouseID); err != nil {
		return nil, err
	}

	return repos.Movements.ListByWarehouse(ctx, warehouseID, options)
}
