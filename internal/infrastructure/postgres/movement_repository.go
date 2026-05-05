package postgres

import (
	"context"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/domain"
	"gorm.io/gorm"
)

type movementRepository struct {
	db *gorm.DB
}

func (repository *movementRepository) ListByProduct(ctx context.Context, productID domain.ID, options application.ListOptions) ([]application.MovementRecord, error) {
	query := `
SELECT * FROM (
	SELECT
		'Supply' AS kind,
		s.id AS operation_id,
		si.id AS operation_item_id,
		si.product_id AS product_id,
		s.warehouse_id AS warehouse_id,
		NULL::BIGINT AS related_warehouse_id,
		si.quantity AS quantity,
		s.status AS status,
		COALESCE(s.confirmed_at, s.updated_at) AS occurred_at
	FROM supply_items si
	JOIN supplies s ON s.id = si.supply_id
	WHERE si.product_id = ? AND s.status = 'Completed'

	UNION ALL

	SELECT
		'Outbound' AS kind,
		ore.id AS operation_id,
		ori.id AS operation_item_id,
		ori.product_id AS product_id,
		ore.warehouse_id AS warehouse_id,
		NULL::BIGINT AS related_warehouse_id,
		ori.quantity AS quantity,
		ore.status AS status,
		COALESCE(ore.executed_at, ore.updated_at) AS occurred_at
	FROM outbound_request_items ori
	JOIN outbound_requests ore ON ore.id = ori.outbound_request_id
	WHERE ori.product_id = ? AND ore.status = 'Completed'

	UNION ALL

	SELECT
		'TransferOut' AS kind,
		t.id AS operation_id,
		ti.id AS operation_item_id,
		ti.product_id AS product_id,
		t.source_warehouse_id AS warehouse_id,
		t.target_warehouse_id AS related_warehouse_id,
		ti.quantity AS quantity,
		t.status AS status,
		COALESCE(t.confirmed_at, t.updated_at) AS occurred_at
	FROM transfer_items ti
	JOIN transfers t ON t.id = ti.transfer_id
	WHERE ti.product_id = ? AND t.status = 'Completed'

	UNION ALL

	SELECT
		'TransferIn' AS kind,
		t.id AS operation_id,
		ti.id AS operation_item_id,
		ti.product_id AS product_id,
		t.target_warehouse_id AS warehouse_id,
		t.source_warehouse_id AS related_warehouse_id,
		ti.quantity AS quantity,
		t.status AS status,
		COALESCE(t.confirmed_at, t.updated_at) AS occurred_at
	FROM transfer_items ti
	JOIN transfers t ON t.id = ti.transfer_id
	WHERE ti.product_id = ? AND t.status = 'Completed'
) movements
ORDER BY occurred_at DESC, operation_id DESC`

	return repository.queryMovements(ctx, query, []any{productID, productID, productID, productID}, options)
}

func (repository *movementRepository) ListByWarehouse(ctx context.Context, warehouseID domain.ID, options application.ListOptions) ([]application.MovementRecord, error) {
	query := `
SELECT * FROM (
	SELECT
		'Supply' AS kind,
		s.id AS operation_id,
		si.id AS operation_item_id,
		si.product_id AS product_id,
		s.warehouse_id AS warehouse_id,
		NULL::BIGINT AS related_warehouse_id,
		si.quantity AS quantity,
		s.status AS status,
		COALESCE(s.confirmed_at, s.updated_at) AS occurred_at
	FROM supply_items si
	JOIN supplies s ON s.id = si.supply_id
	WHERE s.warehouse_id = ? AND s.status = 'Completed'

	UNION ALL

	SELECT
		'Outbound' AS kind,
		ore.id AS operation_id,
		ori.id AS operation_item_id,
		ori.product_id AS product_id,
		ore.warehouse_id AS warehouse_id,
		NULL::BIGINT AS related_warehouse_id,
		ori.quantity AS quantity,
		ore.status AS status,
		COALESCE(ore.executed_at, ore.updated_at) AS occurred_at
	FROM outbound_request_items ori
	JOIN outbound_requests ore ON ore.id = ori.outbound_request_id
	WHERE ore.warehouse_id = ? AND ore.status = 'Completed'

	UNION ALL

	SELECT
		'TransferOut' AS kind,
		t.id AS operation_id,
		ti.id AS operation_item_id,
		ti.product_id AS product_id,
		t.source_warehouse_id AS warehouse_id,
		t.target_warehouse_id AS related_warehouse_id,
		ti.quantity AS quantity,
		t.status AS status,
		COALESCE(t.confirmed_at, t.updated_at) AS occurred_at
	FROM transfer_items ti
	JOIN transfers t ON t.id = ti.transfer_id
	WHERE t.source_warehouse_id = ? AND t.status = 'Completed'

	UNION ALL

	SELECT
		'TransferIn' AS kind,
		t.id AS operation_id,
		ti.id AS operation_item_id,
		ti.product_id AS product_id,
		t.target_warehouse_id AS warehouse_id,
		t.source_warehouse_id AS related_warehouse_id,
		ti.quantity AS quantity,
		t.status AS status,
		COALESCE(t.confirmed_at, t.updated_at) AS occurred_at
	FROM transfer_items ti
	JOIN transfers t ON t.id = ti.transfer_id
	WHERE t.target_warehouse_id = ? AND t.status = 'Completed'
) movements
ORDER BY occurred_at DESC, operation_id DESC`

	return repository.queryMovements(ctx, query, []any{warehouseID, warehouseID, warehouseID, warehouseID}, options)
}

func (repository *movementRepository) queryMovements(ctx context.Context, query string, args []any, options application.ListOptions) ([]application.MovementRecord, error) {
	if options.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, options.Limit)
	}
	if options.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, options.Offset)
	}

	var models []movementRecordModel
	if err := repository.db.WithContext(ctx).Raw(query, args...).Scan(&models).Error; err != nil {
		return nil, err
	}

	movements := make([]application.MovementRecord, 0, len(models))
	for _, model := range models {
		movements = append(movements, toDomainMovementRecord(model))
	}

	return movements, nil
}
