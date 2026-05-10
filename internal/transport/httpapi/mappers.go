package httpapi

import (
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

func toUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toUserResponses(users []domain.User) []UserResponse {
	responses := make([]UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, toUserResponse(user))
	}

	return responses
}

func toWarehouseResponse(warehouse domain.Warehouse) WarehouseResponse {
	return WarehouseResponse{
		ID:        warehouse.ID,
		Name:      warehouse.Name,
		Code:      warehouse.Code,
		Location:  warehouse.Location,
		IsActive:  warehouse.IsActive,
		CreatedAt: warehouse.CreatedAt,
		UpdatedAt: warehouse.UpdatedAt,
	}
}

func toWarehouseResponses(warehouses []domain.Warehouse) []WarehouseResponse {
	responses := make([]WarehouseResponse, 0, len(warehouses))
	for _, warehouse := range warehouses {
		responses = append(responses, toWarehouseResponse(warehouse))
	}

	return responses
}

func toCategoryResponse(category domain.Category) CategoryResponse {
	return CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func toCategoryResponses(categories []domain.Category) []CategoryResponse {
	responses := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, toCategoryResponse(category))
	}

	return responses
}

func toProductResponse(product domain.Product) ProductResponse {
	return ProductResponse{
		ID:                product.ID,
		Name:              product.Name,
		SKU:               product.SKU,
		CategoryID:        product.CategoryID,
		UnitOfMeasure:     product.UnitOfMeasure,
		MinStockThreshold: product.MinStockThreshold,
		IsActive:          product.IsActive,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
	}
}

func toProductResponses(products []domain.Product) []ProductResponse {
	responses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, toProductResponse(product))
	}

	return responses
}

func toStockItemResponse(stockItem domain.StockItem) StockItemResponse {
	return StockItemResponse{
		ID:          stockItem.ID,
		WarehouseID: stockItem.WarehouseID,
		ProductID:   stockItem.ProductID,
		Quantity:    stockItem.Quantity,
		CreatedAt:   stockItem.CreatedAt,
		UpdatedAt:   stockItem.UpdatedAt,
	}
}

func toStockItemResponses(stockItems []domain.StockItem) []StockItemResponse {
	responses := make([]StockItemResponse, 0, len(stockItems))
	for _, stockItem := range stockItems {
		responses = append(responses, toStockItemResponse(stockItem))
	}

	return responses
}

func toStockStatusResponse(status application.StockStatus) StockStatusResponse {
	return StockStatusResponse{
		StockItem: toStockItemResponse(status.StockItem),
		Product:   toProductResponse(status.Product),
		State:     status.State,
	}
}

func toStockStatusResponses(statuses []application.StockStatus) []StockStatusResponse {
	responses := make([]StockStatusResponse, 0, len(statuses))
	for _, status := range statuses {
		responses = append(responses, toStockStatusResponse(status))
	}

	return responses
}

func toSupplyResponse(supply domain.Supply) SupplyResponse {
	items := make([]SupplyItemResponse, 0, len(supply.Items))
	for _, item := range supply.Items {
		items = append(items, toSupplyItemResponse(item))
	}

	return SupplyResponse{
		ID:                supply.ID,
		WarehouseID:       supply.WarehouseID,
		Status:            supply.Status,
		CreatedByUserID:   supply.CreatedByUserID,
		ConfirmedByUserID: supply.ConfirmedByUserID,
		CancelledByUserID: supply.CancelledByUserID,
		Items:             items,
		CreatedAt:         supply.CreatedAt,
		UpdatedAt:         supply.UpdatedAt,
		ConfirmedAt:       supply.ConfirmedAt,
		CancelledAt:       supply.CancelledAt,
	}
}

func toSupplyResponses(supplies []domain.Supply) []SupplyResponse {
	responses := make([]SupplyResponse, 0, len(supplies))
	for _, supply := range supplies {
		responses = append(responses, toSupplyResponse(supply))
	}

	return responses
}

func toSupplyItemResponse(item domain.SupplyItem) SupplyItemResponse {
	return SupplyItemResponse{
		ID:             item.ID,
		SupplyID:       item.SupplyID,
		ProductID:      item.ProductID,
		Quantity:       item.Quantity,
		UnitPriceCents: item.UnitPriceCents,
	}
}

func toOutboundRequestResponse(request domain.OutboundRequest) OutboundRequestResponse {
	items := make([]OutboundRequestItemResponse, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, toOutboundRequestItemResponse(item))
	}

	return OutboundRequestResponse{
		ID:                request.ID,
		WarehouseID:       request.WarehouseID,
		Status:            request.Status,
		CreatedByUserID:   request.CreatedByUserID,
		ApprovedByUserID:  request.ApprovedByUserID,
		ExecutedByUserID:  request.ExecutedByUserID,
		CancelledByUserID: request.CancelledByUserID,
		Items:             items,
		CreatedAt:         request.CreatedAt,
		UpdatedAt:         request.UpdatedAt,
		ApprovedAt:        request.ApprovedAt,
		ExecutedAt:        request.ExecutedAt,
		CancelledAt:       request.CancelledAt,
	}
}

func toOutboundRequestResponses(requests []domain.OutboundRequest) []OutboundRequestResponse {
	responses := make([]OutboundRequestResponse, 0, len(requests))
	for _, request := range requests {
		responses = append(responses, toOutboundRequestResponse(request))
	}

	return responses
}

func toOutboundRequestItemResponse(item domain.OutboundRequestItem) OutboundRequestItemResponse {
	return OutboundRequestItemResponse{
		ID:                item.ID,
		OutboundRequestID: item.OutboundRequestID,
		ProductID:         item.ProductID,
		Quantity:          item.Quantity,
	}
}

func toTransferResponse(transfer domain.Transfer) TransferResponse {
	items := make([]TransferItemResponse, 0, len(transfer.Items))
	for _, item := range transfer.Items {
		items = append(items, toTransferItemResponse(item))
	}

	return TransferResponse{
		ID:                transfer.ID,
		SourceWarehouseID: transfer.SourceWarehouseID,
		TargetWarehouseID: transfer.TargetWarehouseID,
		Status:            transfer.Status,
		CreatedByUserID:   transfer.CreatedByUserID,
		ConfirmedByUserID: transfer.ConfirmedByUserID,
		CancelledByUserID: transfer.CancelledByUserID,
		Items:             items,
		CreatedAt:         transfer.CreatedAt,
		UpdatedAt:         transfer.UpdatedAt,
		ConfirmedAt:       transfer.ConfirmedAt,
		CancelledAt:       transfer.CancelledAt,
	}
}

func toTransferResponses(transfers []domain.Transfer) []TransferResponse {
	responses := make([]TransferResponse, 0, len(transfers))
	for _, transfer := range transfers {
		responses = append(responses, toTransferResponse(transfer))
	}

	return responses
}

func toTransferItemResponse(item domain.TransferItem) TransferItemResponse {
	return TransferItemResponse{
		ID:         item.ID,
		TransferID: item.TransferID,
		ProductID:  item.ProductID,
		Quantity:   item.Quantity,
	}
}

func toMovementResponse(movement application.MovementRecord) MovementResponse {
	return MovementResponse{
		Kind:               movement.Kind,
		OperationID:        movement.OperationID,
		OperationItemID:    movement.OperationItemID,
		ProductID:          movement.ProductID,
		WarehouseID:        movement.WarehouseID,
		RelatedWarehouseID: movement.RelatedWarehouseID,
		Quantity:           movement.Quantity,
		Status:             movement.Status,
		OccurredAt:         movement.OccurredAt,
	}
}

func toMovementResponses(movements []application.MovementRecord) []MovementResponse {
	responses := make([]MovementResponse, 0, len(movements))
	for _, movement := range movements {
		responses = append(responses, toMovementResponse(movement))
	}

	return responses
}
