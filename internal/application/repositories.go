package application

import (
	"context"
	"time"

	"github.com/so1icitx/StockWise/internal/domain"
)

// ListOptions controls pagination for repository list operations.
type ListOptions struct {
	Limit  int
	Offset int
}

// UserFilter contains optional criteria for listing users.
type UserFilter struct {
	Role     *domain.UserRole
	IsActive *bool
	Search   string
	ListOptions
}

// WarehouseFilter contains optional criteria for listing warehouses.
type WarehouseFilter struct {
	IsActive *bool
	Code     string
	Search   string
	ListOptions
}

// CategoryFilter contains optional criteria for listing categories.
type CategoryFilter struct {
	IsActive *bool
	Search   string
	ListOptions
}

// ProductFilter contains optional criteria for listing products.
type ProductFilter struct {
	CategoryID *domain.ID
	IsActive   *bool
	SKU        string
	Name       string
	Search     string
	ListOptions
}

// SupplyFilter contains optional criteria for listing supplies.
type SupplyFilter struct {
	WarehouseID *domain.ID
	Status      *domain.SupplyStatus
	ListOptions
}

// OutboundRequestFilter contains optional criteria for listing outbound requests.
type OutboundRequestFilter struct {
	WarehouseID *domain.ID
	Status      *domain.OutboundRequestStatus
	ListOptions
}

// TransferFilter contains optional criteria for listing transfers.
type TransferFilter struct {
	WarehouseID       *domain.ID
	SourceWarehouseID *domain.ID
	TargetWarehouseID *domain.ID
	Status            *domain.TransferStatus
	ListOptions
}

// MovementKind identifies the business operation that produced a movement record.
type MovementKind string

const (
	// MovementKindSupply identifies stock added through a completed supply.
	MovementKindSupply MovementKind = "Supply"
	// MovementKindOutbound identifies stock removed through a completed outbound request.
	MovementKindOutbound MovementKind = "Outbound"
	// MovementKindTransferIn identifies stock received from another warehouse.
	MovementKindTransferIn MovementKind = "TransferIn"
	// MovementKindTransferOut identifies stock sent to another warehouse.
	MovementKindTransferOut MovementKind = "TransferOut"
)

// MovementRecord is a read model for product and warehouse movement history.
type MovementRecord struct {
	Kind               MovementKind
	OperationID        domain.ID
	OperationItemID    domain.ID
	ProductID          domain.ID
	WarehouseID        domain.ID
	RelatedWarehouseID *domain.ID
	Quantity           int64
	Status             string
	OccurredAt         time.Time
}

// UserRepository persists users.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id domain.ID) error
	GetByID(ctx context.Context, id domain.ID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, filter UserFilter) ([]domain.User, error)
}

// WarehouseRepository persists warehouses.
type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *domain.Warehouse) error
	Update(ctx context.Context, warehouse *domain.Warehouse) error
	Delete(ctx context.Context, id domain.ID) error
	GetByID(ctx context.Context, id domain.ID) (*domain.Warehouse, error)
	GetByCode(ctx context.Context, code string) (*domain.Warehouse, error)
	List(ctx context.Context, filter WarehouseFilter) ([]domain.Warehouse, error)
	SetActive(ctx context.Context, id domain.ID, isActive bool) error
	HasStock(ctx context.Context, id domain.ID) (bool, error)
	HasActiveOperations(ctx context.Context, id domain.ID) (bool, error)
}

// CategoryRepository persists categories.
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id domain.ID) error
	GetByID(ctx context.Context, id domain.ID) (*domain.Category, error)
	GetByName(ctx context.Context, name string) (*domain.Category, error)
	List(ctx context.Context, filter CategoryFilter) ([]domain.Category, error)
	HasActiveProducts(ctx context.Context, id domain.ID) (bool, error)
}

// ProductRepository persists products.
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id domain.ID) error
	GetByID(ctx context.Context, id domain.ID) (*domain.Product, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Product, error)
	List(ctx context.Context, filter ProductFilter) ([]domain.Product, error)
	SetActive(ctx context.Context, id domain.ID, isActive bool) error
	HasConfirmedMovements(ctx context.Context, id domain.ID) (bool, error)
}

// StockRepository persists stock item quantities.
type StockRepository interface {
	GetByWarehouse(ctx context.Context, warehouseID domain.ID) ([]domain.StockItem, error)
	GetByProductAcrossWarehouses(ctx context.Context, productID domain.ID) ([]domain.StockItem, error)
	GetByWarehouseAndProduct(ctx context.Context, warehouseID domain.ID, productID domain.ID) (*domain.StockItem, error)
	Upsert(ctx context.Context, stockItem *domain.StockItem) error
	Increment(ctx context.Context, warehouseID domain.ID, productID domain.ID, quantity int64) (*domain.StockItem, error)
	Decrement(ctx context.Context, warehouseID domain.ID, productID domain.ID, quantity int64) (*domain.StockItem, error)
	GetTotalQuantityForProduct(ctx context.Context, productID domain.ID) (int64, error)
	GetLowStock(ctx context.Context) ([]domain.StockItem, error)
}

// SupplyRepository persists inbound supply operations.
type SupplyRepository interface {
	Create(ctx context.Context, supply *domain.Supply) error
	Update(ctx context.Context, supply *domain.Supply) error
	Delete(ctx context.Context, id domain.ID) error
	GetByID(ctx context.Context, id domain.ID) (*domain.Supply, error)
	List(ctx context.Context, filter SupplyFilter) ([]domain.Supply, error)
	AddItem(ctx context.Context, item *domain.SupplyItem) error
	UpdateItem(ctx context.Context, item *domain.SupplyItem) error
	DeleteItem(ctx context.Context, id domain.ID) error
	ListItems(ctx context.Context, supplyID domain.ID) ([]domain.SupplyItem, error)
}

// OutboundRequestRepository persists outbound request operations.
type OutboundRequestRepository interface {
	Create(ctx context.Context, request *domain.OutboundRequest) error
	Update(ctx context.Context, request *domain.OutboundRequest) error
	Delete(ctx context.Context, id domain.ID) error
	GetByID(ctx context.Context, id domain.ID) (*domain.OutboundRequest, error)
	List(ctx context.Context, filter OutboundRequestFilter) ([]domain.OutboundRequest, error)
	AddItem(ctx context.Context, item *domain.OutboundRequestItem) error
	UpdateItem(ctx context.Context, item *domain.OutboundRequestItem) error
	DeleteItem(ctx context.Context, id domain.ID) error
	ListItems(ctx context.Context, outboundRequestID domain.ID) ([]domain.OutboundRequestItem, error)
}

// TransferRepository persists transfer operations.
type TransferRepository interface {
	Create(ctx context.Context, transfer *domain.Transfer) error
	Update(ctx context.Context, transfer *domain.Transfer) error
	Delete(ctx context.Context, id domain.ID) error
	GetByID(ctx context.Context, id domain.ID) (*domain.Transfer, error)
	List(ctx context.Context, filter TransferFilter) ([]domain.Transfer, error)
	AddItem(ctx context.Context, item *domain.TransferItem) error
	UpdateItem(ctx context.Context, item *domain.TransferItem) error
	DeleteItem(ctx context.Context, id domain.ID) error
	ListItems(ctx context.Context, transferID domain.ID) ([]domain.TransferItem, error)
}

// MovementRepository reads product and warehouse movement history.
type MovementRepository interface {
	ListByProduct(ctx context.Context, productID domain.ID, options ListOptions) ([]MovementRecord, error)
	ListByWarehouse(ctx context.Context, warehouseID domain.ID, options ListOptions) ([]MovementRecord, error)
}
