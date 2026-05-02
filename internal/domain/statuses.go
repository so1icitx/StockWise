package domain

// UserRole describes the responsibility level of a StockWise user.
type UserRole string

const (
	// UserRoleAdmin can manage master data and operational workflows.
	UserRoleAdmin UserRole = "Admin"
	// UserRoleManager can approve and supervise warehouse operations.
	UserRoleManager UserRole = "Manager"
	// UserRoleOperator can create and execute warehouse operations.
	UserRoleOperator UserRole = "Operator"
)

// SupplyStatus describes the lifecycle state of a supply.
type SupplyStatus string

const (
	// SupplyStatusDraft means the supply can still be edited.
	SupplyStatusDraft SupplyStatus = "Draft"
	// SupplyStatusCompleted means the supply has been confirmed and applied to stock.
	SupplyStatusCompleted SupplyStatus = "Completed"
	// SupplyStatusCancelled means the supply was cancelled without changing stock.
	SupplyStatusCancelled SupplyStatus = "Cancelled"
)

// IsFinal reports whether the supply status prevents further editing.
func (status SupplyStatus) IsFinal() bool {
	return status == SupplyStatusCompleted || status == SupplyStatusCancelled
}

// CanEdit reports whether a supply in this status can be changed.
func (status SupplyStatus) CanEdit() bool {
	return !status.IsFinal()
}

// CanConfirm reports whether a supply in this status can be confirmed.
func (status SupplyStatus) CanConfirm() bool {
	return status == SupplyStatusDraft
}

// CanCancel reports whether a supply in this status can be cancelled.
func (status SupplyStatus) CanCancel() bool {
	return status == SupplyStatusDraft
}

// OutboundRequestStatus describes the lifecycle state of an outbound request.
type OutboundRequestStatus string

const (
	// OutboundRequestStatusDraft means the outbound request can still be edited.
	OutboundRequestStatusDraft OutboundRequestStatus = "Draft"
	// OutboundRequestStatusApproved means the outbound request can be executed.
	OutboundRequestStatusApproved OutboundRequestStatus = "Approved"
	// OutboundRequestStatusCompleted means the outbound request has decreased stock.
	OutboundRequestStatusCompleted OutboundRequestStatus = "Completed"
	// OutboundRequestStatusCancelled means the outbound request was cancelled.
	OutboundRequestStatusCancelled OutboundRequestStatus = "Cancelled"
)

// IsFinal reports whether the outbound request status prevents further editing.
func (status OutboundRequestStatus) IsFinal() bool {
	return status == OutboundRequestStatusCompleted || status == OutboundRequestStatusCancelled
}

// CanEdit reports whether an outbound request in this status can be changed.
func (status OutboundRequestStatus) CanEdit() bool {
	return status == OutboundRequestStatusDraft
}

// CanApprove reports whether an outbound request in this status can be approved.
func (status OutboundRequestStatus) CanApprove() bool {
	return status == OutboundRequestStatusDraft
}

// CanExecute reports whether an outbound request in this status can be executed.
func (status OutboundRequestStatus) CanExecute() bool {
	return status == OutboundRequestStatusApproved
}

// CanCancel reports whether an outbound request in this status can be cancelled.
func (status OutboundRequestStatus) CanCancel() bool {
	return status == OutboundRequestStatusDraft || status == OutboundRequestStatusApproved
}

// TransferStatus describes the lifecycle state of a transfer.
type TransferStatus string

const (
	// TransferStatusDraft means the transfer can still be edited.
	TransferStatusDraft TransferStatus = "Draft"
	// TransferStatusCompleted means the transfer has moved stock between warehouses.
	TransferStatusCompleted TransferStatus = "Completed"
	// TransferStatusCancelled means the transfer was cancelled without moving stock.
	TransferStatusCancelled TransferStatus = "Cancelled"
)

// IsFinal reports whether the transfer status prevents further editing.
func (status TransferStatus) IsFinal() bool {
	return status == TransferStatusCompleted || status == TransferStatusCancelled
}

// CanEdit reports whether a transfer in this status can be changed.
func (status TransferStatus) CanEdit() bool {
	return !status.IsFinal()
}

// CanConfirm reports whether a transfer in this status can be confirmed.
func (status TransferStatus) CanConfirm() bool {
	return status == TransferStatusDraft
}

// CanCancel reports whether a transfer in this status can be cancelled.
func (status TransferStatus) CanCancel() bool {
	return status == TransferStatusDraft
}

// StockState describes how a stock item compares to its minimum threshold.
type StockState string

const (
	// StockStateInStock means the stock quantity is above the configured threshold.
	StockStateInStock StockState = "InStock"
	// StockStateLow means the stock quantity is positive but at or below the threshold.
	StockStateLow StockState = "LowStock"
	// StockStateOut means the stock quantity is zero or below.
	StockStateOut StockState = "OutOfStock"
)
