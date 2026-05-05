package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"gorm.io/gorm"
)

type outboundRequestRepository struct {
	db *gorm.DB
}

func (repository *outboundRequestRepository) Create(ctx context.Context, request *domain.OutboundRequest) error {
	model := toOutboundRequestModel(*request)
	if err := repository.db.WithContext(ctx).Omit("Items").Create(&model).Error; err != nil {
		return err
	}

	*request = toDomainOutboundRequest(model)
	return nil
}

func (repository *outboundRequestRepository) Update(ctx context.Context, request *domain.OutboundRequest) error {
	result := repository.db.WithContext(ctx).Model(&outboundRequestModel{}).
		Where("id = ?", request.ID).
		Updates(map[string]any{
			"warehouse_id":         request.WarehouseID,
			"status":               string(request.Status),
			"created_by_user_id":   request.CreatedByUserID,
			"approved_by_user_id":  domainIDPtrToUint64(request.ApprovedByUserID),
			"executed_by_user_id":  domainIDPtrToUint64(request.ExecutedByUserID),
			"cancelled_by_user_id": domainIDPtrToUint64(request.CancelledByUserID),
			"approved_at":          request.ApprovedAt,
			"executed_at":          request.ExecutedAt,
			"cancelled_at":         request.CancelledAt,
			"updated_at":           gorm.Expr("now()"),
		})

	return notFoundOnZeroRows(result)
}

func (repository *outboundRequestRepository) Delete(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&outboundRequestModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *outboundRequestRepository) GetByID(ctx context.Context, id domain.ID) (*domain.OutboundRequest, error) {
	var model outboundRequestModel
	if err := repository.db.WithContext(ctx).Preload("Items").First(&model, uint64(id)).Error; err != nil {
		return nil, mapError(err)
	}

	request := toDomainOutboundRequest(model)
	return &request, nil
}

func (repository *outboundRequestRepository) List(ctx context.Context, filter application.OutboundRequestFilter) ([]domain.OutboundRequest, error) {
	query := repository.db.WithContext(ctx).Model(&outboundRequestModel{}).Preload("Items").Order("id ASC")
	if filter.WarehouseID != nil {
		query = query.Where("warehouse_id = ?", *filter.WarehouseID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}

	var models []outboundRequestModel
	if err := applyListOptions(query, filter.ListOptions).Find(&models).Error; err != nil {
		return nil, err
	}

	requests := make([]domain.OutboundRequest, 0, len(models))
	for _, model := range models {
		requests = append(requests, toDomainOutboundRequest(model))
	}

	return requests, nil
}

func (repository *outboundRequestRepository) AddItem(ctx context.Context, item *domain.OutboundRequestItem) error {
	model := toOutboundRequestItemModel(*item)
	if err := repository.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	*item = toDomainOutboundRequestItem(model)
	return nil
}

func (repository *outboundRequestRepository) UpdateItem(ctx context.Context, item *domain.OutboundRequestItem) error {
	result := repository.db.WithContext(ctx).Model(&outboundRequestItemModel{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"outbound_request_id": item.OutboundRequestID,
			"product_id":          item.ProductID,
			"quantity":            item.Quantity,
		})

	return notFoundOnZeroRows(result)
}

func (repository *outboundRequestRepository) DeleteItem(ctx context.Context, id domain.ID) error {
	result := repository.db.WithContext(ctx).Delete(&outboundRequestItemModel{}, uint64(id))
	return notFoundOnZeroRows(result)
}

func (repository *outboundRequestRepository) ListItems(ctx context.Context, outboundRequestID domain.ID) ([]domain.OutboundRequestItem, error) {
	var models []outboundRequestItemModel
	if err := repository.db.WithContext(ctx).Where("outbound_request_id = ?", outboundRequestID).Order("id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	items := make([]domain.OutboundRequestItem, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainOutboundRequestItem(model))
	}

	return items, nil
}
