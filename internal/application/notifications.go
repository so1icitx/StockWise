package application

import (
	"context"
	"time"
)

// NotificationEventName identifies a real-time event published by StockWise services.
type NotificationEventName string

const (
	// NotificationSupplyConfirmed is emitted after a supply is confirmed and applied to stock.
	NotificationSupplyConfirmed NotificationEventName = "supply.confirmed"
	// NotificationOutboundApproved is emitted after an outbound request is approved.
	NotificationOutboundApproved NotificationEventName = "outbound.approved"
	// NotificationOutboundCompleted is emitted after an outbound request is executed.
	NotificationOutboundCompleted NotificationEventName = "outbound.completed"
	// NotificationTransferConfirmed is emitted after a transfer is confirmed and stock is moved.
	NotificationTransferConfirmed NotificationEventName = "transfer.confirmed"
	// NotificationStockLow is emitted when a stock row is positive and at or below its minimum threshold.
	NotificationStockLow NotificationEventName = "stock.low"
	// NotificationStockOut is emitted when a stock row reaches zero.
	NotificationStockOut NotificationEventName = "stock.out"
	// NotificationWarehouseDeactivated is emitted after a warehouse is deactivated.
	NotificationWarehouseDeactivated NotificationEventName = "warehouse.deactivated"
)

// NotificationEvent is a WebSocket-ready application event with frontend-friendly context.
type NotificationEvent struct {
	Event     NotificationEventName `json:"event"`
	Timestamp time.Time             `json:"timestamp"`
	Data      map[string]any        `json:"data"`
}

// NotificationPublisher publishes service events to interested real-time subscribers.
type NotificationPublisher interface {
	Publish(ctx context.Context, event NotificationEvent)
}

// NoopNotificationPublisher ignores events when no real-time transport is configured.
type NoopNotificationPublisher struct{}

// Publish accepts an event without side effects.
func (NoopNotificationPublisher) Publish(context.Context, NotificationEvent) {}

// NewNotificationEvent creates a timestamped notification event.
func NewNotificationEvent(name NotificationEventName, data map[string]any) NotificationEvent {
	if data == nil {
		data = map[string]any{}
	}

	return NotificationEvent{
		Event:     name,
		Timestamp: nowUTC(),
		Data:      data,
	}
}

func notificationPublisherFrom(publishers ...NotificationPublisher) NotificationPublisher {
	if len(publishers) > 0 && publishers[0] != nil {
		return publishers[0]
	}

	return NoopNotificationPublisher{}
}
