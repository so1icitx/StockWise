package graphqlapi

import (
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"github.com/so1icitx/StockWise/internal/transport/graphql/model"
)

func toGraphQLUser(user domain.User) *model.User {
	return &model.User{
		ID:        idString(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: timeString(user.CreatedAt),
		UpdatedAt: timeString(user.UpdatedAt),
	}
}

func toGraphQLUsers(users []domain.User) []*model.User {
	result := make([]*model.User, 0, len(users))
	for _, user := range users {
		result = append(result, toGraphQLUser(user))
	}
	return result
}

func toGraphQLWarehouse(warehouse domain.Warehouse) *model.Warehouse {
	return &model.Warehouse{
		ID:        idString(warehouse.ID),
		Name:      warehouse.Name,
		Code:      warehouse.Code,
		Location:  warehouse.Location,
		IsActive:  warehouse.IsActive,
		CreatedAt: timeString(warehouse.CreatedAt),
		UpdatedAt: timeString(warehouse.UpdatedAt),
	}
}

func toGraphQLWarehouses(warehouses []domain.Warehouse) []*model.Warehouse {
	result := make([]*model.Warehouse, 0, len(warehouses))
	for _, warehouse := range warehouses {
		result = append(result, toGraphQLWarehouse(warehouse))
	}
	return result
}

func toGraphQLCategory(category domain.Category) *model.Category {
	return &model.Category{
		ID:          idString(category.ID),
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   timeString(category.CreatedAt),
		UpdatedAt:   timeString(category.UpdatedAt),
	}
}

func toGraphQLCategories(categories []domain.Category) []*model.Category {
	result := make([]*model.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, toGraphQLCategory(category))
	}
	return result
}

func toGraphQLProduct(product domain.Product) *model.Product {
	return &model.Product{
		ID:                idString(product.ID),
		Name:              product.Name,
		Sku:               product.SKU,
		CategoryID:        idString(product.CategoryID),
		UnitOfMeasure:     product.UnitOfMeasure,
		MinStockThreshold: int(product.MinStockThreshold),
		IsActive:          product.IsActive,
		CreatedAt:         timeString(product.CreatedAt),
		UpdatedAt:         timeString(product.UpdatedAt),
	}
}

func toGraphQLProducts(products []domain.Product) []*model.Product {
	result := make([]*model.Product, 0, len(products))
	for _, product := range products {
		result = append(result, toGraphQLProduct(product))
	}
	return result
}

func toGraphQLStockItem(stockItem domain.StockItem) *model.StockItem {
	return &model.StockItem{
		ID:          idString(stockItem.ID),
		WarehouseID: idString(stockItem.WarehouseID),
		ProductID:   idString(stockItem.ProductID),
		Quantity:    int(stockItem.Quantity),
		CreatedAt:   timeString(stockItem.CreatedAt),
		UpdatedAt:   timeString(stockItem.UpdatedAt),
	}
}

func toGraphQLStockItems(stockItems []domain.StockItem) []*model.StockItem {
	result := make([]*model.StockItem, 0, len(stockItems))
	for _, stockItem := range stockItems {
		result = append(result, toGraphQLStockItem(stockItem))
	}
	return result
}

func toGraphQLStockStatus(status application.StockStatus) *model.StockStatus {
	return &model.StockStatus{
		StockItem: toGraphQLStockItem(status.StockItem),
		Product:   toGraphQLProduct(status.Product),
		State:     string(status.State),
	}
}

func toGraphQLStockStatuses(statuses []application.StockStatus) []*model.StockStatus {
	result := make([]*model.StockStatus, 0, len(statuses))
	for _, status := range statuses {
		result = append(result, toGraphQLStockStatus(status))
	}
	return result
}

func toGraphQLSupply(supply domain.Supply) *model.Supply {
	items := make([]*model.SupplyItem, 0, len(supply.Items))
	for _, item := range supply.Items {
		items = append(items, toGraphQLSupplyItem(item))
	}

	return &model.Supply{
		ID:                idString(supply.ID),
		WarehouseID:       idString(supply.WarehouseID),
		Status:            string(supply.Status),
		CreatedByUserID:   idString(supply.CreatedByUserID),
		ConfirmedByUserID: idStringPtr(supply.ConfirmedByUserID),
		CancelledByUserID: idStringPtr(supply.CancelledByUserID),
		Items:             items,
		CreatedAt:         timeString(supply.CreatedAt),
		UpdatedAt:         timeString(supply.UpdatedAt),
		ConfirmedAt:       timeStringPtr(supply.ConfirmedAt),
		CancelledAt:       timeStringPtr(supply.CancelledAt),
	}
}

func toGraphQLSupplies(supplies []domain.Supply) []*model.Supply {
	result := make([]*model.Supply, 0, len(supplies))
	for _, supply := range supplies {
		result = append(result, toGraphQLSupply(supply))
	}
	return result
}

func toGraphQLSupplyItem(item domain.SupplyItem) *model.SupplyItem {
	return &model.SupplyItem{
		ID:             idString(item.ID),
		SupplyID:       idString(item.SupplyID),
		ProductID:      idString(item.ProductID),
		Quantity:       int(item.Quantity),
		UnitPriceCents: int(item.UnitPriceCents),
	}
}

func toGraphQLOutboundRequest(request domain.OutboundRequest) *model.OutboundRequest {
	items := make([]*model.OutboundRequestItem, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, toGraphQLOutboundRequestItem(item))
	}

	return &model.OutboundRequest{
		ID:                idString(request.ID),
		WarehouseID:       idString(request.WarehouseID),
		Status:            string(request.Status),
		CreatedByUserID:   idString(request.CreatedByUserID),
		ApprovedByUserID:  idStringPtr(request.ApprovedByUserID),
		ExecutedByUserID:  idStringPtr(request.ExecutedByUserID),
		CancelledByUserID: idStringPtr(request.CancelledByUserID),
		Items:             items,
		CreatedAt:         timeString(request.CreatedAt),
		UpdatedAt:         timeString(request.UpdatedAt),
		ApprovedAt:        timeStringPtr(request.ApprovedAt),
		ExecutedAt:        timeStringPtr(request.ExecutedAt),
		CancelledAt:       timeStringPtr(request.CancelledAt),
	}
}

func toGraphQLOutboundRequests(requests []domain.OutboundRequest) []*model.OutboundRequest {
	result := make([]*model.OutboundRequest, 0, len(requests))
	for _, request := range requests {
		result = append(result, toGraphQLOutboundRequest(request))
	}
	return result
}

func toGraphQLOutboundRequestItem(item domain.OutboundRequestItem) *model.OutboundRequestItem {
	return &model.OutboundRequestItem{
		ID:                idString(item.ID),
		OutboundRequestID: idString(item.OutboundRequestID),
		ProductID:         idString(item.ProductID),
		Quantity:          int(item.Quantity),
	}
}

func toGraphQLTransfer(transfer domain.Transfer) *model.Transfer {
	items := make([]*model.TransferItem, 0, len(transfer.Items))
	for _, item := range transfer.Items {
		items = append(items, toGraphQLTransferItem(item))
	}

	return &model.Transfer{
		ID:                idString(transfer.ID),
		SourceWarehouseID: idString(transfer.SourceWarehouseID),
		TargetWarehouseID: idString(transfer.TargetWarehouseID),
		Status:            string(transfer.Status),
		CreatedByUserID:   idString(transfer.CreatedByUserID),
		ConfirmedByUserID: idStringPtr(transfer.ConfirmedByUserID),
		CancelledByUserID: idStringPtr(transfer.CancelledByUserID),
		Items:             items,
		CreatedAt:         timeString(transfer.CreatedAt),
		UpdatedAt:         timeString(transfer.UpdatedAt),
		ConfirmedAt:       timeStringPtr(transfer.ConfirmedAt),
		CancelledAt:       timeStringPtr(transfer.CancelledAt),
	}
}

func toGraphQLTransfers(transfers []domain.Transfer) []*model.Transfer {
	result := make([]*model.Transfer, 0, len(transfers))
	for _, transfer := range transfers {
		result = append(result, toGraphQLTransfer(transfer))
	}
	return result
}

func toGraphQLTransferItem(item domain.TransferItem) *model.TransferItem {
	return &model.TransferItem{
		ID:         idString(item.ID),
		TransferID: idString(item.TransferID),
		ProductID:  idString(item.ProductID),
		Quantity:   int(item.Quantity),
	}
}

func toGraphQLMovement(movement application.MovementRecord) *model.Movement {
	return &model.Movement{
		Kind:               string(movement.Kind),
		OperationID:        idString(movement.OperationID),
		OperationItemID:    idString(movement.OperationItemID),
		ProductID:          idString(movement.ProductID),
		WarehouseID:        idString(movement.WarehouseID),
		RelatedWarehouseID: idStringPtr(movement.RelatedWarehouseID),
		Quantity:           int(movement.Quantity),
		Status:             movement.Status,
		OccurredAt:         timeString(movement.OccurredAt),
	}
}

func toGraphQLMovements(movements []application.MovementRecord) []*model.Movement {
	result := make([]*model.Movement, 0, len(movements))
	for _, movement := range movements {
		result = append(result, toGraphQLMovement(movement))
	}
	return result
}
