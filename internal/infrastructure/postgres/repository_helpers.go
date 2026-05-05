package postgres

import (
	"errors"

	"github.com/so1icitx/StockWise/internal/application"
	"gorm.io/gorm"
)

func mapError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return application.ErrNotFound
	}

	return err
}

func notFoundOnZeroRows(result *gorm.DB) error {
	if result.Error != nil {
		return mapError(result.Error)
	}

	if result.RowsAffected == 0 {
		return application.ErrNotFound
	}

	return nil
}

func applyListOptions(query *gorm.DB, options application.ListOptions) *gorm.DB {
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	return query
}
