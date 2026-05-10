package httpapi

import (
	"time"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
)

// CreateUserRequest is the payload for creating a user.
type CreateUserRequest struct {
	Name  string          `json:"name" binding:"required"`
	Email string          `json:"email" binding:"required,email"`
	Role  domain.UserRole `json:"role" binding:"required"`
}

// UpdateUserRequest is the payload for updating a user.
type UpdateUserRequest struct {
	Name     string          `json:"name" binding:"required"`
	Email    string          `json:"email" binding:"required,email"`
	Role     domain.UserRole `json:"role" binding:"required"`
	IsActive bool            `json:"is_active"`
}

// UserResponse is the REST representation of a user.
type UserResponse struct {
	ID        domain.ID       `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// WarehouseRequest is the payload for creating or updating a warehouse.
type WarehouseRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Location string `json:"location"`
	IsActive bool   `json:"is_active"`
}

// WarehouseResponse is the REST representation of a warehouse.
type WarehouseResponse struct {
	ID        domain.ID `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Location  string    `json:"location"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CategoryRequest is the payload for creating or updating a category.
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

// CategoryResponse is the REST representation of a category.
type CategoryResponse struct {
	ID          domain.ID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductRequest is the payload for creating or updating a product.
type ProductRequest struct {
	Name              string    `json:"name" binding:"required"`
	SKU               string    `json:"sku" binding:"required"`
	CategoryID        domain.ID `json:"category_id" binding:"required"`
	UnitOfMeasure     string    `json:"unit_of_measure" binding:"required"`
	MinStockThreshold int64     `json:"min_stock_threshold" binding:"min=0"`
	IsActive          bool      `json:"is_active"`
}

// ProductResponse is the REST representation of a product.
type ProductResponse struct {
	ID                domain.ID `json:"id"`
	Name              string    `json:"name"`
	SKU               string    `json:"sku"`
	CategoryID        domain.ID `json:"category_id"`
	UnitOfMeasure     string    `json:"unit_of_measure"`
	MinStockThreshold int64     `json:"min_stock_threshold"`
	IsActive          bool      `json:"is_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// SupplyCreateRequest is the payload for creating a supply.
type SupplyCreateRequest struct {
	WarehouseID domain.ID `json:"warehouse_id" binding:"required"`
}

// SupplyItemRequest is the payload for creating or updating a supply item.
type SupplyItemRequest struct {
	ProductID      domain.ID `json:"product_id" binding:"required"`
	Quantity       int64     `json:"quantity" binding:"required,gt=0"`
	UnitPriceCents int64     `json:"unit_price_cents" binding:"min=0"`
}

// SupplyResponse is the REST representation of a supply.
type SupplyResponse struct {
	ID                domain.ID            `json:"id"`
	WarehouseID       domain.ID            `json:"warehouse_id"`
	Status            domain.SupplyStatus  `json:"status"`
	CreatedByUserID   domain.ID            `json:"created_by_user_id"`
	ConfirmedByUserID *domain.ID           `json:"confirmed_by_user_id,omitempty"`
	CancelledByUserID *domain.ID           `json:"cancelled_by_user_id,omitempty"`
	Items             []SupplyItemResponse `json:"items"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
	ConfirmedAt       *time.Time           `json:"confirmed_at,omitempty"`
	CancelledAt       *time.Time           `json:"cancelled_at,omitempty"`
}

// SupplyItemResponse is the REST representation of a supply item.
type SupplyItemResponse struct {
	ID             domain.ID `json:"id"`
	SupplyID       domain.ID `json:"supply_id"`
	ProductID      domain.ID `json:"product_id"`
	Quantity       int64     `json:"quantity"`
	UnitPriceCents int64     `json:"unit_price_cents"`
}

// OutboundRequestCreateRequest is the payload for creating an outbound request.
type OutboundRequestCreateRequest struct {
	WarehouseID domain.ID `json:"warehouse_id" binding:"required"`
}

// OutboundRequestItemRequest is the payload for creating or updating an outbound request item.
type OutboundRequestItemRequest struct {
	ProductID domain.ID `json:"product_id" binding:"required"`
	Quantity  int64     `json:"quantity" binding:"required,gt=0"`
}

// OutboundRequestResponse is the REST representation of an outbound request.
type OutboundRequestResponse struct {
	ID                domain.ID                     `json:"id"`
	WarehouseID       domain.ID                     `json:"warehouse_id"`
	Status            domain.OutboundRequestStatus  `json:"status"`
	CreatedByUserID   domain.ID                     `json:"created_by_user_id"`
	ApprovedByUserID  *domain.ID                    `json:"approved_by_user_id,omitempty"`
	ExecutedByUserID  *domain.ID                    `json:"executed_by_user_id,omitempty"`
	CancelledByUserID *domain.ID                    `json:"cancelled_by_user_id,omitempty"`
	Items             []OutboundRequestItemResponse `json:"items"`
	CreatedAt         time.Time                     `json:"created_at"`
	UpdatedAt         time.Time                     `json:"updated_at"`
	ApprovedAt        *time.Time                    `json:"approved_at,omitempty"`
	ExecutedAt        *time.Time                    `json:"executed_at,omitempty"`
	CancelledAt       *time.Time                    `json:"cancelled_at,omitempty"`
}

// OutboundRequestItemResponse is the REST representation of an outbound request item.
type OutboundRequestItemResponse struct {
	ID                domain.ID `json:"id"`
	OutboundRequestID domain.ID `json:"outbound_request_id"`
	ProductID         domain.ID `json:"product_id"`
	Quantity          int64     `json:"quantity"`
}

// TransferCreateRequest is the payload for creating a transfer.
type TransferCreateRequest struct {
	SourceWarehouseID domain.ID `json:"source_warehouse_id" binding:"required"`
	TargetWarehouseID domain.ID `json:"target_warehouse_id" binding:"required"`
}

// TransferItemRequest is the payload for creating or updating a transfer item.
type TransferItemRequest struct {
	ProductID domain.ID `json:"product_id" binding:"required"`
	Quantity  int64     `json:"quantity" binding:"required,gt=0"`
}

// TransferResponse is the REST representation of a transfer.
type TransferResponse struct {
	ID                domain.ID              `json:"id"`
	SourceWarehouseID domain.ID              `json:"source_warehouse_id"`
	TargetWarehouseID domain.ID              `json:"target_warehouse_id"`
	Status            domain.TransferStatus  `json:"status"`
	CreatedByUserID   domain.ID              `json:"created_by_user_id"`
	ConfirmedByUserID *domain.ID             `json:"confirmed_by_user_id,omitempty"`
	CancelledByUserID *domain.ID             `json:"cancelled_by_user_id,omitempty"`
	Items             []TransferItemResponse `json:"items"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	ConfirmedAt       *time.Time             `json:"confirmed_at,omitempty"`
	CancelledAt       *time.Time             `json:"cancelled_at,omitempty"`
}

// TransferItemResponse is the REST representation of a transfer item.
type TransferItemResponse struct {
	ID         domain.ID `json:"id"`
	TransferID domain.ID `json:"transfer_id"`
	ProductID  domain.ID `json:"product_id"`
	Quantity   int64     `json:"quantity"`
}

// StockItemResponse is the REST representation of a stock item.
type StockItemResponse struct {
	ID          domain.ID `json:"id"`
	WarehouseID domain.ID `json:"warehouse_id"`
	ProductID   domain.ID `json:"product_id"`
	Quantity    int64     `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockStatusResponse describes stock with low-stock state information.
type StockStatusResponse struct {
	StockItem StockItemResponse `json:"stock_item"`
	Product   ProductResponse   `json:"product"`
	State     domain.StockState `json:"state"`
}

// TotalStockResponse contains total stock for one product.
type TotalStockResponse struct {
	ProductID domain.ID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
}

// MovementResponse is the REST representation of a stock movement record.
type MovementResponse struct {
	Kind               application.MovementKind `json:"kind"`
	OperationID        domain.ID                `json:"operation_id"`
	OperationItemID    domain.ID                `json:"operation_item_id"`
	ProductID          domain.ID                `json:"product_id"`
	WarehouseID        domain.ID                `json:"warehouse_id"`
	RelatedWarehouseID *domain.ID               `json:"related_warehouse_id,omitempty"`
	Quantity           int64                    `json:"quantity"`
	Status             string                   `json:"status"`
	OccurredAt         time.Time                `json:"occurred_at"`
}
