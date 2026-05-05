package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"gorm.io/gorm"
)

type stockRepository struct {
	db *gorm.DB
}

func (repository *stockRepository) GetByWarehouse(ctx context.Context, warehouseID domain.ID) ([]domain.StockItem, error) {
	var models []stockItemModel
	if err := repository.db.WithContext(ctx).
		Where("warehouse_id = ?", warehouseID).
		Order("product_id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return toDomainStockItems(models), nil
}

func (repository *stockRepository) GetByProductAcrossWarehouses(ctx context.Context, productID domain.ID) ([]domain.StockItem, error) {
	var models []stockItemModel
	if err := repository.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("warehouse_id ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	return toDomainStockItems(models), nil
}

func (repository *stockRepository) GetByWarehouseAndProduct(ctx context.Context, warehouseID domain.ID, productID domain.ID) (*domain.StockItem, error) {
	var model stockItemModel
	err := repository.db.WithContext(ctx).
		Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).
		First(&model).Error
	if err != nil {
		return nil, mapError(err)
	}

	stockItem := toDomainStockItem(model)
	return &stockItem, nil
}

func (repository *stockRepository) Upsert(ctx context.Context, stockItem *domain.StockItem) error {
	var model stockItemModel
	err := repository.db.WithContext(ctx).Raw(`
INSERT INTO stock_items (warehouse_id, product_id, quantity)
VALUES (?, ?, ?)
ON CONFLICT (warehouse_id, product_id) DO UPDATE
SET quantity = EXCLUDED.quantity,
	updated_at = now()
RETURNING id, warehouse_id, product_id, quantity, created_at, updated_at`,
		stockItem.WarehouseID,
		stockItem.ProductID,
		stockItem.Quantity,
	).Scan(&model).Error
	if err != nil {
		return err
	}

	*stockItem = toDomainStockItem(model)
	return nil
}

func (repository *stockRepository) Increment(ctx context.Context, warehouseID domain.ID, productID domain.ID, quantity int64) (*domain.StockItem, error) {
	var model stockItemModel
	err := repository.db.WithContext(ctx).Raw(`
INSERT INTO stock_items (warehouse_id, product_id, quantity)
VALUES (?, ?, ?)
ON CONFLICT (warehouse_id, product_id) DO UPDATE
SET quantity = stock_items.quantity + EXCLUDED.quantity,
	updated_at = now()
RETURNING id, warehouse_id, product_id, quantity, created_at, updated_at`,
		warehouseID,
		productID,
		quantity,
	).Scan(&model).Error
	if err != nil {
		return nil, err
	}

	stockItem := toDomainStockItem(model)
	return &stockItem, nil
}

func (repository *stockRepository) Decrement(ctx context.Context, warehouseID domain.ID, productID domain.ID, quantity int64) (*domain.StockItem, error) {
	var model stockItemModel
	err := repository.db.WithContext(ctx).Raw(`
UPDATE stock_items
SET quantity = quantity - ?,
	updated_at = now()
WHERE warehouse_id = ?
	AND product_id = ?
	AND quantity >= ?
RETURNING id, warehouse_id, product_id, quantity, created_at, updated_at`,
		quantity,
		warehouseID,
		productID,
		quantity,
	).Scan(&model).Error
	if err != nil {
		return nil, err
	}
	if model.ID == 0 {
		return nil, application.ErrNotFound
	}

	stockItem := toDomainStockItem(model)
	return &stockItem, nil
}

func (repository *stockRepository) GetTotalQuantityForProduct(ctx context.Context, productID domain.ID) (int64, error) {
	var total int64
	err := repository.db.WithContext(ctx).Model(&stockItemModel{}).
		Where("product_id = ?", productID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&total).Error

	return total, err
}

func (repository *stockRepository) GetLowStock(ctx context.Context) ([]domain.StockItem, error) {
	var models []stockItemModel
	err := repository.db.WithContext(ctx).Table("stock_items").
		Select("stock_items.*").
		Joins("JOIN products ON products.id = stock_items.product_id").
		Where("products.is_active = true").
		Where("stock_items.quantity <= products.min_stock_threshold").
		Order("stock_items.warehouse_id ASC, stock_items.product_id ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	return toDomainStockItems(models), nil
}

func toDomainStockItems(models []stockItemModel) []domain.StockItem {
	stockItems := make([]domain.StockItem, 0, len(models))
	for _, model := range models {
		stockItems = append(stockItems, toDomainStockItem(model))
	}

	return stockItems
}
