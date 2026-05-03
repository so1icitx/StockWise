package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/so1icitx/StockWise/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := seed(ctx, db); err != nil {
		log.Fatalf("seed database: %v", err)
	}

	log.Println("seed data applied")
}

func seed(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for _, statement := range seedStatements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Printf("rollback failed: %v", rollbackErr)
			}
			return err
		}
	}

	return tx.Commit()
}

var seedStatements = []string{
	`INSERT INTO users (name, email, role, is_active)
VALUES
	('Alice Admin', 'alice.admin@stockwise.local', 'Admin', true),
	('Martin Manager', 'martin.manager@stockwise.local', 'Manager', true),
	('Olivia Operator', 'olivia.operator@stockwise.local', 'Operator', true)
ON CONFLICT (email) DO UPDATE
SET name = EXCLUDED.name,
	role = EXCLUDED.role,
	is_active = EXCLUDED.is_active,
	updated_at = now();`,

	`INSERT INTO warehouses (name, code, location, is_active)
VALUES
	('Central Warehouse', 'WH-CENTRAL', 'Sofia, Bulgaria', true),
	('Retail Warehouse', 'WH-RETAIL', 'Plovdiv, Bulgaria', true),
	('Reserve Warehouse', 'WH-RESERVE', 'Varna, Bulgaria', true)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
	location = EXCLUDED.location,
	is_active = EXCLUDED.is_active,
	updated_at = now();`,

	`INSERT INTO categories (name, description, is_active)
VALUES
	('Electronics', 'Barcode scanners, handheld terminals, and warehouse electronics.', true),
	('Packaging', 'Boxes, labels, stretch film, and other packaging materials.', true),
	('Food', 'Shelf-stable food products for inventory demos.', true)
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
	is_active = EXCLUDED.is_active,
	updated_at = now();`,

	`INSERT INTO products (name, sku, category_id, unit_of_measure, min_stock_threshold, is_active)
SELECT 'Wireless Barcode Scanner', 'ELEC-001', id, 'pcs', 5, true FROM categories WHERE name = 'Electronics'
ON CONFLICT (sku) DO UPDATE
SET name = EXCLUDED.name,
	category_id = EXCLUDED.category_id,
	unit_of_measure = EXCLUDED.unit_of_measure,
	min_stock_threshold = EXCLUDED.min_stock_threshold,
	is_active = EXCLUDED.is_active,
	updated_at = now();`,

	`INSERT INTO products (name, sku, category_id, unit_of_measure, min_stock_threshold, is_active)
SELECT 'Thermal Label Printer', 'ELEC-002', id, 'pcs', 3, true FROM categories WHERE name = 'Electronics'
ON CONFLICT (sku) DO UPDATE
SET name = EXCLUDED.name,
	category_id = EXCLUDED.category_id,
	unit_of_measure = EXCLUDED.unit_of_measure,
	min_stock_threshold = EXCLUDED.min_stock_threshold,
	is_active = EXCLUDED.is_active,
	updated_at = now();`,

	`INSERT INTO products (name, sku, category_id, unit_of_measure, min_stock_threshold, is_active)
SELECT 'Cardboard Box Medium', 'PACK-001', id, 'pcs', 100, true FROM categories WHERE name = 'Packaging'
ON CONFLICT (sku) DO UPDATE
SET name = EXCLUDED.name,
	category_id = EXCLUDED.category_id,
	unit_of_measure = EXCLUDED.unit_of_measure,
	min_stock_threshold = EXCLUDED.min_stock_threshold,
	is_active = EXCLUDED.is_active,
	updated_at = now();`,

	`INSERT INTO products (name, sku, category_id, unit_of_measure, min_stock_threshold, is_active)
SELECT 'Pasta Pack 500g', 'FOOD-001', id, 'pcs', 25, true FROM categories WHERE name = 'Food'
ON CONFLICT (sku) DO UPDATE
SET name = EXCLUDED.name,
	category_id = EXCLUDED.category_id,
	unit_of_measure = EXCLUDED.unit_of_measure,
	min_stock_threshold = EXCLUDED.min_stock_threshold,
	is_active = EXCLUDED.is_active,
	updated_at = now();`,

	`INSERT INTO stock_items (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 12 FROM warehouses w, products p
WHERE w.code = 'WH-CENTRAL' AND p.sku = 'ELEC-001'
ON CONFLICT (warehouse_id, product_id) DO UPDATE
SET quantity = EXCLUDED.quantity,
	updated_at = now();`,

	`INSERT INTO stock_items (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 2 FROM warehouses w, products p
WHERE w.code = 'WH-CENTRAL' AND p.sku = 'ELEC-002'
ON CONFLICT (warehouse_id, product_id) DO UPDATE
SET quantity = EXCLUDED.quantity,
	updated_at = now();`,

	`INSERT INTO stock_items (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 250 FROM warehouses w, products p
WHERE w.code = 'WH-RETAIL' AND p.sku = 'PACK-001'
ON CONFLICT (warehouse_id, product_id) DO UPDATE
SET quantity = EXCLUDED.quantity,
	updated_at = now();`,

	`INSERT INTO stock_items (warehouse_id, product_id, quantity)
SELECT w.id, p.id, 0 FROM warehouses w, products p
WHERE w.code = 'WH-RESERVE' AND p.sku = 'FOOD-001'
ON CONFLICT (warehouse_id, product_id) DO UPDATE
SET quantity = EXCLUDED.quantity,
	updated_at = now();`,
}
