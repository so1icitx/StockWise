package postgres

import "time"

type userModel struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string
	Email     string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (userModel) TableName() string {
	return "users"
}

type warehouseModel struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string
	Code      string
	Location  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (warehouseModel) TableName() string {
	return "warehouses"
}

type categoryModel struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (categoryModel) TableName() string {
	return "categories"
}

type productModel struct {
	ID                uint64 `gorm:"primaryKey"`
	Name              string
	SKU               string
	CategoryID        uint64
	UnitOfMeasure     string
	MinStockThreshold int64
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (productModel) TableName() string {
	return "products"
}

type stockItemModel struct {
	ID          uint64 `gorm:"primaryKey"`
	WarehouseID uint64
	ProductID   uint64
	Quantity    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (stockItemModel) TableName() string {
	return "stock_items"
}

type supplyModel struct {
	ID                uint64 `gorm:"primaryKey"`
	WarehouseID       uint64
	Status            string
	CreatedByUserID   uint64
	ConfirmedByUserID *uint64
	CancelledByUserID *uint64
	Items             []supplyItemModel `gorm:"foreignKey:SupplyID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ConfirmedAt       *time.Time
	CancelledAt       *time.Time
}

func (supplyModel) TableName() string {
	return "supplies"
}

type supplyItemModel struct {
	ID             uint64 `gorm:"primaryKey"`
	SupplyID       uint64
	ProductID      uint64
	Quantity       int64
	UnitPriceCents int64
}

func (supplyItemModel) TableName() string {
	return "supply_items"
}

type outboundRequestModel struct {
	ID                uint64 `gorm:"primaryKey"`
	WarehouseID       uint64
	Status            string
	CreatedByUserID   uint64
	ApprovedByUserID  *uint64
	ExecutedByUserID  *uint64
	CancelledByUserID *uint64
	Items             []outboundRequestItemModel `gorm:"foreignKey:OutboundRequestID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ApprovedAt        *time.Time
	ExecutedAt        *time.Time
	CancelledAt       *time.Time
}

func (outboundRequestModel) TableName() string {
	return "outbound_requests"
}

type outboundRequestItemModel struct {
	ID                uint64 `gorm:"primaryKey"`
	OutboundRequestID uint64
	ProductID         uint64
	Quantity          int64
}

func (outboundRequestItemModel) TableName() string {
	return "outbound_request_items"
}

type transferModel struct {
	ID                uint64 `gorm:"primaryKey"`
	SourceWarehouseID uint64
	TargetWarehouseID uint64
	Status            string
	CreatedByUserID   uint64
	ConfirmedByUserID *uint64
	CancelledByUserID *uint64
	Items             []transferItemModel `gorm:"foreignKey:TransferID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ConfirmedAt       *time.Time
	CancelledAt       *time.Time
}

func (transferModel) TableName() string {
	return "transfers"
}

type transferItemModel struct {
	ID         uint64 `gorm:"primaryKey"`
	TransferID uint64
	ProductID  uint64
	Quantity   int64
}

func (transferItemModel) TableName() string {
	return "transfer_items"
}

type movementRecordModel struct {
	Kind               string
	OperationID        uint64
	OperationItemID    uint64
	ProductID          uint64
	WarehouseID        uint64
	RelatedWarehouseID *uint64
	Quantity           int64
	Status             string
	OccurredAt         time.Time
}
