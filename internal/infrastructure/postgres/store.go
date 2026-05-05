package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"gorm.io/gorm"
)

// Store owns the PostgreSQL-backed repository implementations.
type Store struct {
	db *gorm.DB
}

// NewStore creates a repository store backed by GORM.
func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

// Repositories returns repository implementations outside an explicit transaction.
func (store *Store) Repositories() application.Repositories {
	return repositoriesForDB(store.db)
}

// WithinTransaction executes repository work inside one database transaction.
func (store *Store) WithinTransaction(ctx context.Context, fn func(ctx context.Context, repositories application.Repositories) error) error {
	return store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(ctx, repositoriesForDB(tx))
	})
}

func repositoriesForDB(db *gorm.DB) application.Repositories {
	return application.Repositories{
		Users:            &userRepository{db: db},
		Warehouses:       &warehouseRepository{db: db},
		Categories:       &categoryRepository{db: db},
		Products:         &productRepository{db: db},
		Stock:            &stockRepository{db: db},
		Supplies:         &supplyRepository{db: db},
		OutboundRequests: &outboundRequestRepository{db: db},
		Transfers:        &transferRepository{db: db},
		Movements:        &movementRepository{db: db},
	}
}
