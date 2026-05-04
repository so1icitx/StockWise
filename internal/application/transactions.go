package application

import "context"

// Repositories groups all repository contracts available to application services.
type Repositories struct {
	Users            UserRepository
	Warehouses       WarehouseRepository
	Categories       CategoryRepository
	Products         ProductRepository
	Stock            StockRepository
	Supplies         SupplyRepository
	OutboundRequests OutboundRequestRepository
	Transfers        TransferRepository
	Movements        MovementRepository
}

// RepositoryProvider exposes repository implementations outside an explicit transaction.
type RepositoryProvider interface {
	Repositories() Repositories
}

// TransactionManager runs repository work inside a single atomic transaction.
type TransactionManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context, repositories Repositories) error) error
}
