package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"gorm.io/gorm"
)

type transferRepository struct {
	db *gorm.DB
}

func (repository *transferRepository) Create(ctx context.Context, transfer *domain.Transfer) error {
	model := toTransferModel(*transfer)
	if err := repository.db.WithContext(ctx).Omit("Items").Create(&model).Error; err != nil {
		return err
	}

	*transfer = toDomainTransfer(model)
	return nil
}

func (repository *transferRepository) Update(ctx context.Context, transfer *domain.Transfer) error {
	result := repository.db.WithContext(ctx).Model(&transferModel{}).
		Where("id = ?", transfer.ID).
		Updates(map[string]any{
			"source_warehouse_id":  transfer.SourceWarehouseID,
			"target_warehouse_id":  transfer.TargetWarehouseID,
			"status":               string(transfer.Status),
			"created_by_user_id":   transfer.CreatedByUserID,
			"confirmed_by_user_id": domainIDPtrToUint64(transfer.ConfirmedByUserID),
			"cancelled_by_user_id": domainIDPtrToUint64(transfer.CancelledByUserID),
			"confirmed_at":         transfer.ConfirmedAt,
			"cancelled_at":         transfer.CancelledAt,
			"updated_at":           gorm.Expr("now()"),
		})

	return notFoundOnZeroRows(result)
}

func (repository *transferRepository) Delete(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&transferModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *transferRepository) GetByID(ctx context.Context, id domain.ID) (*domain.Transfer, error) {
	var model transferModel
	if err := repository.db.WithContext(ctx).Preload("Items").First(&model, uint64(id)).Error; err != nil {
		return nil, mapError(err)
	}

	transfer := toDomainTransfer(model)
	return &transfer, nil
}

func (repository *transferRepository) List(ctx context.Context, filter application.TransferFilter) ([]domain.Transfer, error) {
	query := repository.db.WithContext(ctx).Model(&transferModel{}).Preload("Items").Order("id ASC")
	if filter.WarehouseID != nil {
		query = query.Where("source_warehouse_id = ? OR target_warehouse_id = ?", *filter.WarehouseID, *filter.WarehouseID)
	}
	if filter.SourceWarehouseID != nil {
		query = query.Where("source_warehouse_id = ?", *filter.SourceWarehouseID)
	}
	if filter.TargetWarehouseID != nil {
		query = query.Where("target_warehouse_id = ?", *filter.TargetWarehouseID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}

	var models []transferModel
	if err := applyListOptions(query, filter.ListOptions).Find(&models).Error; err != nil {
		return nil, err
	}

	transfers := make([]domain.Transfer, 0, len(models))
	for _, model := range models {
		transfers = append(transfers, toDomainTransfer(model))
	}

	return transfers, nil
}

func (repository *transferRepository) AddItem(ctx context.Context, item *domain.TransferItem) error {
	model := toTransferItemModel(*item)
	if err := repository.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	*item = toDomainTransferItem(model)
	return nil
}

func (repository *transferRepository) UpdateItem(ctx context.Context, item *domain.TransferItem) error {
	result := repository.db.WithContext(ctx).Model(&transferItemModel{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"transfer_id": item.TransferID,
			"product_id":  item.ProductID,
			"quantity":    item.Quantity,
		})

	return notFoundOnZeroRows(result)
}

func (repository *transferRepository) DeleteItem(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&transferItemModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *transferRepository) ListItems(ctx context.Context, transferID domain.ID) ([]domain.TransferItem, error) {
	var models []transferItemModel
	if err := repository.db.WithContext(ctx).Where("transfer_id = ?", transferID).Order("id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]domain.TransferItem, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainTransferItem(model))
	}

	return items, nil
}
