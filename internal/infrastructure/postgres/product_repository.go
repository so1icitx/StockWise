package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func (repository *productRepository) Create(ctx context.Context, product *domain.Product) error {
	model := toProductModel(*product)
	if err := repository.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	*product = toDomainProduct(model)
	return nil
}

func (repository *productRepository) Update(ctx context.Context, product *domain.Product) error {
	result := repository.db.WithContext(ctx).Model(&productModel{}).
		Where("id = ?", product.ID).
		Updates(map[string]any{
			"name":                product.Name,
			"sku":                 product.SKU,
			"category_id":         product.CategoryID,
			"unit_of_measure":     product.UnitOfMeasure,
			"min_stock_threshold": product.MinStockThreshold,
			"is_active":           product.IsActive,
			"updated_at":          gorm.Expr("now()"),
		})

	return notFoundOnZeroRows(result)
}

func (repository *productRepository) Delete(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&productModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *productRepository) GetByID(ctx context.Context, id domain.ID) (*domain.Product, error) {
	var model productModel
	if err := repository.db.WithContext(ctx).First(&model, uint64(id)).Error; err != nil {
		return nil, mapError(err)
	}

	product := toDomainProduct(model)
	return &product, nil
}

func (repository *productRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var model productModel
	if err := repository.db.WithContext(ctx).Where("sku = ?", sku).First(&model).Error; err != nil {
		return nil, mapError(err)
	}

	product := toDomainProduct(model)
	return &product, nil
}

func (repository *productRepository) List(ctx context.Context, filter application.ProductFilter) ([]domain.Product, error) {
	query := repository.db.WithContext(ctx).Model(&productModel{}).Order("id ASC")
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.SKU != "" {
		query = query.Where("sku = ?", filter.SKU)
	}
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Search != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	var models []productModel
	if err := applyListOptions(query, filter.ListOptions).Find(&models).Error; err != nil {
		return nil, err
	}

	products := make([]domain.Product, 0, len(models))
	for _, model := range models {
		products = append(products, toDomainProduct(model))
	}

	return products, nil
}

func (repository *productRepository) SetActive(ctx context.Context, id domain.ID, isActive bool) error {
	result := repository.db.WithContext(ctx).Model(&productModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"is_active":  isActive,
			"updated_at": gorm.Expr("now()"),
		})

	return notFoundOnZeroRows(result)
}

func (repository *productRepository) HasConfirmedMovements(ctx context.Context, id domain.ID) (bool, error) {
	var count int64
	err := repository.db.WithContext(ctx).Raw(`
SELECT COUNT(*) FROM (
	SELECT si.id FROM supply_items si
	JOIN supplies s ON s.id = si.supply_id
	WHERE si.product_id = ? AND s.status = 'Completed'
	UNION ALL
	SELECT ori.id FROM outbound_request_items ori
	JOIN outbound_requests ore ON ore.id = ori.outbound_request_id
	WHERE ori.product_id = ? AND ore.status = 'Completed'
	UNION ALL
	SELECT ti.id FROM transfer_items ti
	JOIN transfers t ON t.id = ti.transfer_id
	WHERE ti.product_id = ? AND t.status = 'Completed'
) confirmed_movements`, id, id, id).Scan(&count).Error

	return count > 0, err
}
