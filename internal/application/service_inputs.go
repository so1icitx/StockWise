package application

import "github.com/so1icitx/StockWise/internal/domain"

// CreateUserInput contains values for creating a user.
type CreateUserInput struct {
	Name  string
	Email string
	Role  domain.UserRole
}

// UpdateUserInput contains values for updating a user.
type UpdateUserInput struct {
	Name     string
	Email    string
	Role     domain.UserRole
	IsActive bool
}

// CreateWarehouseInput contains values for creating a warehouse.
type CreateWarehouseInput struct {
	Name     string
	Code     string
	Location string
}

// UpdateWarehouseInput contains values for updating a warehouse.
type UpdateWarehouseInput struct {
	Name     string
	Code     string
	Location string
	IsActive bool
}

// CreateCategoryInput contains values for creating a category.
type CreateCategoryInput struct {
	Name        string
	Description string
}

// UpdateCategoryInput contains values for updating a category.
type UpdateCategoryInput struct {
	Name        string
	Description string
	IsActive    bool
}

// CreateProductInput contains values for creating a product.
type CreateProductInput struct {
	Name              string
	SKU               string
	CategoryID        domain.ID
	UnitOfMeasure     string
	MinStockThreshold int64
}

// UpdateProductInput contains values for updating a product.
type UpdateProductInput struct {
	Name              string
	SKU               string
	CategoryID        domain.ID
	UnitOfMeasure     string
	MinStockThreshold int64
	IsActive          bool
}

// CreateSupplyInput contains values for creating a supply.
type CreateSupplyInput struct {
	WarehouseID     domain.ID
	CreatedByUserID domain.ID
}

// SupplyItemInput contains values for creating or updating a supply item.
type SupplyItemInput struct {
	ProductID      domain.ID
	Quantity       int64
	UnitPriceCents int64
}

// CreateOutboundRequestInput contains values for creating an outbound request.
type CreateOutboundRequestInput struct {
	WarehouseID     domain.ID
	CreatedByUserID domain.ID
}

// OutboundRequestItemInput contains values for creating or updating an outbound request item.
type OutboundRequestItemInput struct {
	ProductID domain.ID
	Quantity  int64
}

// CreateTransferInput contains values for creating a transfer.
type CreateTransferInput struct {
	SourceWarehouseID domain.ID
	TargetWarehouseID domain.ID
	CreatedByUserID   domain.ID
}

// TransferItemInput contains values for creating or updating a transfer item.
type TransferItemInput struct {
	ProductID domain.ID
	Quantity  int64
}
