package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"gorm.io/gorm"
)

type warehouseRepository struct {
	db *gorm.DB
}

func (repository *warehouseRepository) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	model := toWarehouseModel(*warehouse)
	if err := repository.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	*warehouse = toDomainWarehouse(model)
	return nil
}

func (repository *warehouseRepository) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	result := repository.db.WithContext(ctx).Model(&warehouseModel{}).
		Where("id = ?", warehouse.ID).
		Updates(map[string]any{
			"name":       warehouse.Name,
			"code":       warehouse.Code,
			"location":   warehouse.Location,
			"is_active":  warehouse.IsActive,
			"updated_at": gorm.Expr("now()"),
		})

	return notFoundOnZeroRows(result)
}

func (repository *warehouseRepository) Delete(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&warehouseModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *warehouseRepository) GetByID(ctx context.Context, id domain.ID) (*domain.Warehouse, error) {
	var model warehouseModel
	if err := repository.db.WithContext(ctx).First(&model, uint64(id)).Error; err != nil {
		return nil, mapError(err)
	}

	warehouse := toDomainWarehouse(model)
	return &warehouse, nil
}

func (repository *warehouseRepository) GetByCode(ctx context.Context, code string) (*domain.Warehouse, error) {
	var model warehouseModel
	if err := repository.db.WithContext(ctx).Where("code = ?", code).First(&model).Error; err != nil {
		return nil, mapError(err)
	}

	warehouse := toDomainWarehouse(model)
	return &warehouse, nil
}

func (repository *warehouseRepository) List(ctx context.Context, filter application.WarehouseFilter) ([]domain.Warehouse, error) {
	query := repository.db.WithContext(ctx).Model(&warehouseModel{}).Order("id ASC")
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Code != "" {
		query = query.Where("code = ?", filter.Code)
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ? OR location ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var models []warehouseModel
	if err := applyListOptions(query, filter.ListOptions).Find(&models).Error; err != nil {
		return nil, err
	}

	warehouses := make([]domain.Warehouse, 0, len(models))
	for _, model := range models {
		warehouses = append(warehouses, toDomainWarehouse(model))
	}

	return warehouses, nil
}

func (repository *warehouseRepository) SetActive(ctx context.Context, id domain.ID, isActive bool) error {
	result := repository.db.WithContext(ctx).Model(&warehouseModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_active":  isActive,
			"updated_at": gorm.Expr("now()"),
		})

	return notFoundOnZeroRows(result)
}

func (repository *warehouseRepository) HasStock(ctx context.Context, id domain.ID) (bool, error) {
	var count int64
	err := repository.db.WithContext(ctx).Model(&stockItemModel{}).
		Where("warehouse_id = ? AND quantity > 0", id).
		Count(&count).Error

	return count > 0, err
}

func (repository *warehouseRepository) HasActiveOperations(ctx context.Context, id domain.ID) (bool, error) {
	var count int64
	err := repository.db.WithContext(ctx).Raw(`
SELECT COUNT(*) FROM (
	SELECT id FROM supplies WHERE warehouse_id = ? AND status NOT IN ('Completed', 'Cancelled')
	UNION ALL
	SELECT id FROM outbound_requests WHERE warehouse_id = ? AND status NOT IN ('Completed', 'Cancelled')
	UNION ALL
	SELECT id FROM transfers WHERE (source_warehouse_id = ? OR target_warehouse_id = ?) AND status NOT IN ('Completed', 'Cancelled')
) active_operations`, id, id, id, id).Scan(&count).Error

	return count > 0, err
}
