package postgres

import (
	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

func toUserModel(user domain.User) userModel {
	return userModel{
		ID:        uint64(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toDomainUser(model userModel) domain.User {
	return domain.User{
		ID:        domain.ID(model.ID),
		Name:      model.Name,
		Email:     model.Email,
		Role:      domain.UserRole(model.Role),
		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func toWarehouseModel(warehouse domain.Warehouse) warehouseModel {
	return warehouseModel{
		ID:        uint64(warehouse.ID),
		Name:      warehouse.Name,
		Code:      warehouse.Code,
		Location:  warehouse.Location,
		IsActive:  warehouse.IsActive,
		CreatedAt: warehouse.CreatedAt,
		UpdatedAt: warehouse.UpdatedAt,
	}
}

func toDomainWarehouse(model warehouseModel) domain.Warehouse {
	return domain.Warehouse{
		ID:        domain.ID(model.ID),
		Name:      model.Name,
		Code:      model.Code,
		Location:  model.Location,
		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func toCategoryModel(category domain.Category) categoryModel {
	return categoryModel{
		ID:          uint64(category.ID),
		Name:        category.Name,
		Description: category.Description,
		IsActive:    category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func toDomainCategory(model categoryModel) domain.Category {
	return domain.Category{
		ID:          domain.ID(model.ID),
		Name:        model.Name,
		Description: model.Description,
		IsActive:    model.IsActive,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func toProductModel(product domain.Product) productModel {
	return productModel{
		ID:                uint64(product.ID),
		Name:              product.Name,
		SKU:               product.SKU,
		CategoryID:        uint64(product.CategoryID),
		UnitOfMeasure:     product.UnitOfMeasure,
		MinStockThreshold: product.MinStockThreshold,
		IsActive:          product.IsActive,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
	}
}

func toDomainProduct(model productModel) domain.Product {
	return domain.Product{
		ID:                domain.ID(model.ID),
		Name:              model.Name,
		SKU:               model.SKU,
		CategoryID:        domain.ID(model.CategoryID),
		UnitOfMeasure:     model.UnitOfMeasure,
		MinStockThreshold: model.MinStockThreshold,
		IsActive:          model.IsActive,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
	}
}

func toStockItemModel(stockItem domain.StockItem) stockItemModel {
	return stockItemModel{
		ID:          uint64(stockItem.ID),
		WarehouseID: uint64(stockItem.WarehouseID),
		ProductID:   uint64(stockItem.ProductID),
		Quantity:    stockItem.Quantity,
		CreatedAt:   stockItem.CreatedAt,
		UpdatedAt:   stockItem.UpdatedAt,
	}
}

func toDomainStockItem(model stockItemModel) domain.StockItem {
	return domain.StockItem{
		ID:          domain.ID(model.ID),
		WarehouseID: domain.ID(model.WarehouseID),
		ProductID:   domain.ID(model.ProductID),
		Quantity:    model.Quantity,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func toSupplyModel(supply domain.Supply) supplyModel {
	model := supplyModel{
		ID:                uint64(supply.ID),
		WarehouseID:       uint64(supply.WarehouseID),
		Status:            string(supply.Status),
		CreatedByUserID:   uint64(supply.CreatedByUserID),
		ConfirmedByUserID: domainIDPtrToUint64(supply.ConfirmedByUserID),
		CancelledByUserID: domainIDPtrToUint64(supply.CancelledByUserID),
		CreatedAt:         supply.CreatedAt,
		UpdatedAt:         supply.UpdatedAt,
		ConfirmedAt:       supply.ConfirmedAt,
		CancelledAt:       supply.CancelledAt,
	}

	if len(supply.Items) > 0 {
		model.Items = make([]supplyItemModel, 0, len(supply.Items))
		for _, item := range supply.Items {
			model.Items = append(model.Items, toSupplyItemModel(item))
		}
	}

	return model
}

func toDomainSupply(model supplyModel) domain.Supply {
	supply := domain.Supply{
		ID:                domain.ID(model.ID),
		WarehouseID:       domain.ID(model.WarehouseID),
		Status:            domain.SupplyStatus(model.Status),
		CreatedByUserID:   domain.ID(model.CreatedByUserID),
		ConfirmedByUserID: uint64PtrToDomainID(model.ConfirmedByUserID),
		CancelledByUserID: uint64PtrToDomainID(model.CancelledByUserID),
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		ConfirmedAt:       model.ConfirmedAt,
		CancelledAt:       model.CancelledAt,
	}

	if len(model.Items) > 0 {
		supply.Items = make([]domain.SupplyItem, 0, len(model.Items))
		for _, item := range model.Items {
			supply.Items = append(supply.Items, toDomainSupplyItem(item))
		}
	}

	return supply
}

func toSupplyItemModel(item domain.SupplyItem) supplyItemModel {
	return supplyItemModel{
		ID:             uint64(item.ID),
		SupplyID:       uint64(item.SupplyID),
		ProductID:      uint64(item.ProductID),
		Quantity:       item.Quantity,
		UnitPriceCents: item.UnitPriceCents,
	}
}

func toDomainSupplyItem(model supplyItemModel) domain.SupplyItem {
	return domain.SupplyItem{
		ID:             domain.ID(model.ID),
		SupplyID:       domain.ID(model.SupplyID),
		ProductID:      domain.ID(model.ProductID),
		Quantity:       model.Quantity,
		UnitPriceCents: model.UnitPriceCents,
	}
}

func toOutboundRequestModel(request domain.OutboundRequest) outboundRequestModel {
	model := outboundRequestModel{
		ID:                uint64(request.ID),
		WarehouseID:       uint64(request.WarehouseID),
		Status:            string(request.Status),
		CreatedByUserID:   uint64(request.CreatedByUserID),
		ApprovedByUserID:  domainIDPtrToUint64(request.ApprovedByUserID),
		ExecutedByUserID:  domainIDPtrToUint64(request.ExecutedByUserID),
		CancelledByUserID: domainIDPtrToUint64(request.CancelledByUserID),
		CreatedAt:         request.CreatedAt,
		UpdatedAt:         request.UpdatedAt,
		ApprovedAt:        request.ApprovedAt,
		ExecutedAt:        request.ExecutedAt,
		CancelledAt:       request.CancelledAt,
	}

	if len(request.Items) > 0 {
		model.Items = make([]outboundRequestItemModel, 0, len(request.Items))
		for _, item := range request.Items {
			model.Items = append(model.Items, toOutboundRequestItemModel(item))
		}
	}

	return model
}

func toDomainOutboundRequest(model outboundRequestModel) domain.OutboundRequest {
	request := domain.OutboundRequest{
		ID:                domain.ID(model.ID),
		WarehouseID:       domain.ID(model.WarehouseID),
		Status:            domain.OutboundRequestStatus(model.Status),
		CreatedByUserID:   domain.ID(model.CreatedByUserID),
		ApprovedByUserID:  uint64PtrToDomainID(model.ApprovedByUserID),
		ExecutedByUserID:  uint64PtrToDomainID(model.ExecutedByUserID),
		CancelledByUserID: uint64PtrToDomainID(model.CancelledByUserID),
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		ApprovedAt:        model.ApprovedAt,
		ExecutedAt:        model.ExecutedAt,
		CancelledAt:       model.CancelledAt,
	}

	if len(model.Items) > 0 {
		request.Items = make([]domain.OutboundRequestItem, 0, len(model.Items))
		for _, item := range model.Items {
			request.Items = append(request.Items, toDomainOutboundRequestItem(item))
		}
	}

	return request
}

func toOutboundRequestItemModel(item domain.OutboundRequestItem) outboundRequestItemModel {
	return outboundRequestItemModel{
		ID:                uint64(item.ID),
		OutboundRequestID: uint64(item.OutboundRequestID),
		ProductID:         uint64(item.ProductID),
		Quantity:          item.Quantity,
	}
}

func toDomainOutboundRequestItem(model outboundRequestItemModel) domain.OutboundRequestItem {
	return domain.OutboundRequestItem{
		ID:                domain.ID(model.ID),
		OutboundRequestID: domain.ID(model.OutboundRequestID),
		ProductID:         domain.ID(model.ProductID),
		Quantity:          model.Quantity,
	}
}

func toTransferModel(transfer domain.Transfer) transferModel {
	model := transferModel{
		ID:                uint64(transfer.ID),
		SourceWarehouseID: uint64(transfer.SourceWarehouseID),
		TargetWarehouseID: uint64(transfer.TargetWarehouseID),
		Status:            string(transfer.Status),
		CreatedByUserID:   uint64(transfer.CreatedByUserID),
		ConfirmedByUserID: domainIDPtrToUint64(transfer.ConfirmedByUserID),
		CancelledByUserID: domainIDPtrToUint64(transfer.CancelledByUserID),
		CreatedAt:         transfer.CreatedAt,
		UpdatedAt:         transfer.UpdatedAt,
		ConfirmedAt:       transfer.ConfirmedAt,
		CancelledAt:       transfer.CancelledAt,
	}

	if len(transfer.Items) > 0 {
		model.Items = make([]transferItemModel, 0, len(transfer.Items))
		for _, item := range transfer.Items {
			model.Items = append(model.Items, toTransferItemModel(item))
		}
	}

	return model
}

func toDomainTransfer(model transferModel) domain.Transfer {
	transfer := domain.Transfer{
		ID:                domain.ID(model.ID),
		SourceWarehouseID: domain.ID(model.SourceWarehouseID),
		TargetWarehouseID: domain.ID(model.TargetWarehouseID),
		Status:            domain.TransferStatus(model.Status),
		CreatedByUserID:   domain.ID(model.CreatedByUserID),
		ConfirmedByUserID: uint64PtrToDomainID(model.ConfirmedByUserID),
		CancelledByUserID: uint64PtrToDomainID(model.CancelledByUserID),
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		ConfirmedAt:       model.ConfirmedAt,
		CancelledAt:       model.CancelledAt,
	}

	if len(model.Items) > 0 {
		transfer.Items = make([]domain.TransferItem, 0, len(model.Items))
		for _, item := range model.Items {
			transfer.Items = append(transfer.Items, toDomainTransferItem(item))
		}
	}

	return transfer
}

func toTransferItemModel(item domain.TransferItem) transferItemModel {
	return transferItemModel{
		ID:         uint64(item.ID),
		TransferID: uint64(item.TransferID),
		ProductID:  uint64(item.ProductID),
		Quantity:   item.Quantity,
	}
}

func toDomainTransferItem(model transferItemModel) domain.TransferItem {
	return domain.TransferItem{
		ID:         domain.ID(model.ID),
		TransferID: domain.ID(model.TransferID),
		ProductID:  domain.ID(model.ProductID),
		Quantity:   model.Quantity,
	}
}

func toDomainMovementRecord(model movementRecordModel) application.MovementRecord {
	return application.MovementRecord{
		Kind:               application.MovementKind(model.Kind),
		OperationID:        domain.ID(model.OperationID),
		OperationItemID:    domain.ID(model.OperationItemID),
		ProductID:          domain.ID(model.ProductID),
		WarehouseID:        domain.ID(model.WarehouseID),
		RelatedWarehouseID: uint64PtrToDomainID(model.RelatedWarehouseID),
		Quantity:           model.Quantity,
		Status:             model.Status,
		OccurredAt:         model.OccurredAt,
	}
}

func domainIDPtrToUint64(id *domain.ID) *uint64 {
	if id == nil {
		return nil
	}

	value := uint64(*id)
	return &value
}

func uint64PtrToDomainID(id *uint64) *domain.ID {
	if id == nil {
		return nil
	}

	value := domain.ID(*id)
	return &value
}
