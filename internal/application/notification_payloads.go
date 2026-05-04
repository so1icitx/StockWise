package application

import (
	"context"

	"github.com/so1icitx/StockWise/internal/domain"
)

func publishAll(ctx context.Context, publisher NotificationPublisher, events []NotificationEvent) {
	for _, event := range events {
		publisher.Publish(ctx, event)
	}
}

func supplyConfirmedEvent(supply domain.Supply) NotificationEvent {
	return NewNotificationEvent(NotificationSupplyConfirmed, map[string]any{
		"supply_id":            supply.ID,
		"warehouse_id":         supply.WarehouseID,
		"status":               supply.Status,
		"confirmed_by_user_id": supply.ConfirmedByUserID,
		"item_count":           len(supply.Items),
		"items":                supplyItemsPayload(supply.Items),
	})
}

func outboundApprovedEvent(request domain.OutboundRequest) NotificationEvent {
	return NewNotificationEvent(NotificationOutboundApproved, map[string]any{
		"outbound_request_id": request.ID,
		"warehouse_id":        request.WarehouseID,
		"status":              request.Status,
		"approved_by_user_id": request.ApprovedByUserID,
		"item_count":          len(request.Items),
		"items":               outboundItemsPayload(request.Items),
	})
}

func outboundCompletedEvent(request domain.OutboundRequest) NotificationEvent {
	return NewNotificationEvent(NotificationOutboundCompleted, map[string]any{
		"outbound_request_id": request.ID,
		"warehouse_id":        request.WarehouseID,
		"status":              request.Status,
		"executed_by_user_id": request.ExecutedByUserID,
		"item_count":          len(request.Items),
		"items":               outboundItemsPayload(request.Items),
	})
}

func transferConfirmedEvent(transfer domain.Transfer) NotificationEvent {
	return NewNotificationEvent(NotificationTransferConfirmed, map[string]any{
		"transfer_id":          transfer.ID,
		"source_warehouse_id":  transfer.SourceWarehouseID,
		"target_warehouse_id":  transfer.TargetWarehouseID,
		"status":               transfer.Status,
		"confirmed_by_user_id": transfer.ConfirmedByUserID,
		"item_count":           len(transfer.Items),
		"items":                transferItemsPayload(transfer.Items),
	})
}

func warehouseDeactivatedEvent(warehouse domain.Warehouse) NotificationEvent {
	return NewNotificationEvent(NotificationWarehouseDeactivated, map[string]any{
		"warehouse_id": warehouse.ID,
		"name":         warehouse.Name,
		"code":         warehouse.Code,
		"location":     warehouse.Location,
		"is_active":    warehouse.IsActive,
	})
}

func appendStockStateEvent(events []NotificationEvent, stockItem domain.StockItem, product domain.Product, trigger NotificationEventName) []NotificationEvent {
	state := stockItem.StateForProduct(product)
	if state != domain.StockStateLow && state != domain.StockStateOut {
		return events
	}

	name := NotificationStockLow
	if state == domain.StockStateOut {
		name = NotificationStockOut
	}

	return append(events, NewNotificationEvent(name, map[string]any{
		"warehouse_id":        stockItem.WarehouseID,
		"product_id":          stockItem.ProductID,
		"product_sku":         product.SKU,
		"product_name":        product.Name,
		"quantity":            stockItem.Quantity,
		"min_stock_threshold": product.MinStockThreshold,
		"state":               state,
		"triggered_by_event":  trigger,
		"unit_of_measure":     product.UnitOfMeasure,
		"product_category_id": product.CategoryID,
		"stock_item_id":       stockItem.ID,
	}))
}

func supplyItemsPayload(items []domain.SupplyItem) []map[string]any {
	payload := make([]map[string]any, 0, len(items))
	for _, item := range items {
		payload = append(payload, map[string]any{
			"supply_item_id":   item.ID,
			"product_id":       item.ProductID,
			"quantity":         item.Quantity,
			"unit_price_cents": item.UnitPriceCents,
		})
	}

	return payload
}

func outboundItemsPayload(items []domain.OutboundRequestItem) []map[string]any {
	payload := make([]map[string]any, 0, len(items))
	for _, item := range items {
		payload = append(payload, map[string]any{
			"outbound_request_item_id": item.ID,
			"product_id":               item.ProductID,
			"quantity":                 item.Quantity,
		})
	}

	return payload
}

func transferItemsPayload(items []domain.TransferItem) []map[string]any {
	payload := make([]map[string]any, 0, len(items))
	for _, item := range items {
		payload = append(payload, map[string]any{
			"transfer_item_id": item.ID,
			"product_id":       item.ProductID,
			"quantity":         item.Quantity,
		})
	}

	return payload
}
