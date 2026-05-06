package application

import (
	"context"

	"github.com/so1icitx/StockWise/internal/domain"
)

// StockStatus describes a stock item with its product-derived state.
type StockStatus struct {
	StockItem domain.StockItem
	Product   domain.Product
	State     domain.StockState
}

// StockService handles stock read operations and stock state calculations.
type StockService struct {
	provider RepositoryProvider
}

// NewStockService creates a stock service.
func NewStockService(provider RepositoryProvider) *StockService {
	return &StockService{provider: provider}
}

// GetByWarehouse returns all stock rows for a warehouse.
func (service *StockService) GetByWarehouse(ctx context.Context, warehouseID domain.ID) ([]domain.StockItem, error) {
	if err := requireID(warehouseID, "warehouse id"); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	if _, err := repos.Warehouses.GetByID(ctx, warehouseID); err != nil {
		return nil, err
	}

	return repos.Stock.GetByWarehouse(ctx, warehouseID)
}

// GetByProductAcrossWarehouses returns all stock rows for a product.
func (service *StockService) GetByProductAcrossWarehouses(ctx context.Context, productID domain.ID) ([]domain.StockItem, error) {
	if err := requireID(productID, "product id"); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	if _, err := repos.Products.GetByID(ctx, productID); err != nil {
		return nil, err
	}

	return repos.Stock.GetByProductAcrossWarehouses(ctx, productID)
}

// GetByWarehouseAndProduct returns one stock row for a warehouse and product pair.
func (service *StockService) GetByWarehouseAndProduct(ctx context.Context, warehouseID domain.ID, productID domain.ID) (*domain.StockItem, error) {
	if err := requireID(warehouseID, "warehouse id"); err != nil {
		return nil, err
	}
	if err := requireID(productID, "product id"); err != nil {
		return nil, err
	}

	repos := repositories(service.provider)
	return repos.Stock.GetByWarehouseAndProduct(ctx, warehouseID, productID)
}

// GetTotalQuantityForProduct returns total stock for a product across all warehouses.
func (service *StockService) GetTotalQuantityForProduct(ctx context.Context, productID domain.ID) (int64, error) {
	if err := requireID(productID, "product id"); err != nil {
		return 0, err
	}

	repos := repositories(service.provider)
	if _, err := repos.Products.GetByID(ctx, productID); err != nil {
		return 0, err
	}

	return repos.Stock.GetTotalQuantityForProduct(ctx, productID)
}

// GetLowStock returns stock rows that are low or out of stock.
func (service *StockService) GetLowStock(ctx context.Context) ([]StockStatus, error) {
	repos := repositories(service.provider)
	stockItems, err := repos.Stock.GetLowStock(ctx)
	if err != nil {
		return nil, err
	}

	statuses := make([]StockStatus, 0, len(stockItems))
	for _, stockItem := range stockItems {
		product, err := repos.Products.GetByID(ctx, stockItem.ProductID)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, StockStatus{
			StockItem: stockItem,
			Product:   *product,
			State:     stockItem.StateForProduct(*product),
		})
	}

	return statuses, nil
}

// GetStockStatus returns the low-stock state for one warehouse and product pair.
func (service *StockService) GetStockStatus(ctx context.Context, warehouseID domain.ID, productID domain.ID) (*StockStatus, error) {
	repos := repositories(service.provider)
	stockItem, err := service.GetByWarehouseAndProduct(ctx, warehouseID, productID)
	if err != nil {
		return nil, err
	}
	product, err := repos.Products.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return &StockStatus{
		StockItem: *stockItem,
		Product:   *product,
		State:     stockItem.StateForProduct(*product),
	}, nil
}
