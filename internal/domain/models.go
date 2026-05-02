package domain

import "time"

// User represents a person who creates, approves, confirms, or executes stock operations.
type User struct {
	ID        ID
	Name      string
	Email     string
	Role      UserRole
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Warehouse represents a physical or logical storage location for inventory.
type Warehouse struct {
	ID        ID
	Name      string
	Code      string
	Location  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CanBeUsedForOperations reports whether the warehouse can participate in new stock operations.
func (warehouse Warehouse) CanBeUsedForOperations() bool {
	return warehouse.IsActive
}

// Category groups related products.
type Category struct {
	ID          ID
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CanBeUsedForProducts reports whether the category can be assigned to active products.
func (category Category) CanBeUsedForProducts() bool {
	return category.IsActive
}

// Product describes an item that can be stocked, supplied, requested, or transferred.
type Product struct {
	ID                ID
	Name              string
	SKU               string
	CategoryID        ID
	UnitOfMeasure     string
	MinStockThreshold int64
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// CanBeUsedForOperations reports whether the product can appear in new stock operations.
func (product Product) CanBeUsedForOperations() bool {
	return product.IsActive
}

// StockItem stores the current quantity of a product in a warehouse.
type StockItem struct {
	ID          ID
	WarehouseID ID
	ProductID   ID
	Quantity    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// StateForProduct classifies this stock item using the product minimum threshold.
func (stockItem StockItem) StateForProduct(product Product) StockState {
	return EvaluateStockState(stockItem.Quantity, product.MinStockThreshold)
}

// IsLowForProduct reports whether this stock item is low for the supplied product threshold.
func (stockItem StockItem) IsLowForProduct(product Product) bool {
	return IsLowStock(stockItem.Quantity, product.MinStockThreshold)
}

// IsOutOfStock reports whether this stock item is depleted.
func (stockItem StockItem) IsOutOfStock() bool {
	return IsOutOfStock(stockItem.Quantity)
}

// Supply represents an inbound operation that increases warehouse stock when confirmed.
type Supply struct {
	ID                ID
	WarehouseID       ID
	Status            SupplyStatus
	CreatedByUserID   ID
	ConfirmedByUserID *ID
	CancelledByUserID *ID
	Items             []SupplyItem
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ConfirmedAt       *time.Time
	CancelledAt       *time.Time
}

// CanBeEdited reports whether the supply can still be changed.
func (supply Supply) CanBeEdited() bool {
	return supply.Status.CanEdit()
}

// HasItems reports whether the supply contains at least one item row.
func (supply Supply) HasItems() bool {
	return len(supply.Items) > 0
}

// HasValidItems reports whether every supply item has valid quantity and price values.
func (supply Supply) HasValidItems() bool {
	if !supply.HasItems() {
		return false
	}

	for _, item := range supply.Items {
		if !item.IsValid() {
			return false
		}
	}

	return true
}

// SupplyItem represents one product row in a supply.
type SupplyItem struct {
	ID             ID
	SupplyID       ID
	ProductID      ID
	Quantity       int64
	UnitPriceCents int64
}

// IsValid reports whether the supply item satisfies quantity and price rules.
func (item SupplyItem) IsValid() bool {
	return IsPositiveQuantity(item.Quantity) && IsNonNegativePriceCents(item.UnitPriceCents)
}

// OutboundRequest represents a request to remove stock from one warehouse.
type OutboundRequest struct {
	ID                ID
	WarehouseID       ID
	Status            OutboundRequestStatus
	CreatedByUserID   ID
	ApprovedByUserID  *ID
	ExecutedByUserID  *ID
	CancelledByUserID *ID
	Items             []OutboundRequestItem
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ApprovedAt        *time.Time
	ExecutedAt        *time.Time
	CancelledAt       *time.Time
}

// CanBeEdited reports whether the outbound request can still be changed.
func (request OutboundRequest) CanBeEdited() bool {
	return request.Status.CanEdit()
}

// HasItems reports whether the outbound request contains at least one item row.
func (request OutboundRequest) HasItems() bool {
	return len(request.Items) > 0
}

// HasValidItems reports whether every outbound request item has a valid quantity.
func (request OutboundRequest) HasValidItems() bool {
	if !request.HasItems() {
		return false
	}

	for _, item := range request.Items {
		if !item.IsValid() {
			return false
		}
	}

	return true
}

// OutboundRequestItem represents one product row in an outbound request.
type OutboundRequestItem struct {
	ID                ID
	OutboundRequestID ID
	ProductID         ID
	Quantity          int64
}

// IsValid reports whether the outbound request item satisfies quantity rules.
func (item OutboundRequestItem) IsValid() bool {
	return IsPositiveQuantity(item.Quantity)
}

// Transfer represents a stock movement from one warehouse to another warehouse.
type Transfer struct {
	ID                ID
	SourceWarehouseID ID
	TargetWarehouseID ID
	Status            TransferStatus
	CreatedByUserID   ID
	ConfirmedByUserID *ID
	CancelledByUserID *ID
	Items             []TransferItem
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ConfirmedAt       *time.Time
	CancelledAt       *time.Time
}

// CanBeEdited reports whether the transfer can still be changed.
func (transfer Transfer) CanBeEdited() bool {
	return transfer.Status.CanEdit()
}

// HasItems reports whether the transfer contains at least one item row.
func (transfer Transfer) HasItems() bool {
	return len(transfer.Items) > 0
}

// UsesSameWarehouse reports whether the source and target warehouse are the same.
func (transfer Transfer) UsesSameWarehouse() bool {
	return !transfer.SourceWarehouseID.IsZero() && transfer.SourceWarehouseID == transfer.TargetWarehouseID
}

// HasValidItems reports whether every transfer item has a valid quantity.
func (transfer Transfer) HasValidItems() bool {
	if !transfer.HasItems() {
		return false
	}

	for _, item := range transfer.Items {
		if !item.IsValid() {
			return false
		}
	}

	return true
}

// TransferItem represents one product row in a transfer.
type TransferItem struct {
	ID         ID
	TransferID ID
	ProductID  ID
	Quantity   int64
}

// IsValid reports whether the transfer item satisfies quantity rules.
func (item TransferItem) IsValid() bool {
	return IsPositiveQuantity(item.Quantity)
}
