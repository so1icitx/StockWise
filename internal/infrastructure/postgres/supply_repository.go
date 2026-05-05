package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"gorm.io/gorm"
)

type supplyRepository struct {
	db *gorm.DB
}

func (repository *supplyRepository) Create(ctx context.Context, supply *domain.Supply) error {
	model := toSupplyModel(*supply)
	if err := repository.db.WithContext(ctx).Omit("Items").Create(&model).Error; err != nil {
		return err
	}

	*supply = toDomainSupply(model)
	return nil
}

func (repository *supplyRepository) Update(ctx context.Context, supply *domain.Supply) error {
	result := repository.db.WithContext(ctx).Model(&supplyModel{}).
		Where("id = ?", supply.ID).
		Updates(map[string]any{
			"warehouse_id":         supply.WarehouseID,
			"status":               string(supply.Status),
			"created_by_user_id":   supply.CreatedByUserID,
			"confirmed_by_user_id": domainIDPtrToUint64(supply.ConfirmedByUserID),
			"cancelled_by_user_id": domainIDPtrToUint64(supply.CancelledByUserID),
			"confirmed_at":         supply.ConfirmedAt,
			"cancelled_at":         supply.CancelledAt,
			"updated_at":           gorm.Expr("now()"),
		})

	return notFoundOnZeroRows(result)
}

func (repository *supplyRepository) Delete(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&supplyModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *supplyRepository) GetByID(ctx context.Context, id domain.ID) (*domain.Supply, error) {
	var model supplyModel
	if err := repository.db.WithContext(ctx).Preload("Items").First(&model, uint64(id)).Error; err != nil {
		return nil, mapError(err)
	}

	supply := toDomainSupply(model)
	return &supply, nil
}

func (repository *supplyRepository) List(ctx context.Context, filter application.SupplyFilter) ([]domain.Supply, error) {
	query := repository.db.WithContext(ctx).Model(&supplyModel{}).Preload("Items").Order("id ASC")
	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}

	var models []supplyModel
	if err := applyListOptions(query, filter.ListOptions).Find(&models).Error; err != nil {
		return nil, err
	}

	supplies := make([]domain.Supply, 0, len(models))
	for _, model := range models {
		supplies = append(supplies, toDomainSupply(model))
	}

	return supplies, nil
}

func (repository *supplyRepository) AddItem(ctx context.Context, item *domain.SupplyItem) error {
	model := toSupplyItemModel(*item)
	if err := repository.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	*item = toDomainSupplyItem(model)
	return nil
}

func (repository *supplyRepository) UpdateItem(ctx context.Context, item *domain.SupplyItem) error {
	result := repository.db.WithContext(ctx).Model(&supplyItemModel{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"supply_id":        item.SupplyID,
			"product_id":       item.ProductID,
			"quantity":         item.Quantity,
			"unit_price_cents": item.UnitPriceCents,
		})

	return notFoundOnZeroRows(result)
}

func (repository *supplyRepository) DeleteItem(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&supplyItemModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *supplyRepository) ListItems(ctx context.Context, supplyID domain.ID) ([]domain.SupplyItem, error) {
	var models []supplyItemModel
	if err := repository.db.WithContext(ctx).Where("supply_id = ?", supplyID).Order("id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]domain.SupplyItem, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainSupplyItem(model))
	}

	return items, nil
}
