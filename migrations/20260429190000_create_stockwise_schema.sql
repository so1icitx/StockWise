-- +goose Up
CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	role TEXT NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	CONSTRAINT users_email_unique UNIQUE (email),
	CONSTRAINT users_role_check CHECK (role IN ('Admin', 'Manager', 'Operator'))
);

CREATE TABLE warehouses (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	code TEXT NOT NULL,
	location TEXT NOT NULL DEFAULT '',
	is_active BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	CONSTRAINT warehouses_code_unique UNIQUE (code)
);

CREATE TABLE categories (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	is_active BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	CONSTRAINT categories_name_unique UNIQUE (name)
);

CREATE TABLE products (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	sku TEXT NOT NULL,
	category_id BIGINT NOT NULL REFERENCES categories(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	unit_of_measure TEXT NOT NULL,
	min_stock_threshold BIGINT NOT NULL DEFAULT 0,
	is_active BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	CONSTRAINT products_sku_unique UNIQUE (sku),
	CONSTRAINT products_min_stock_threshold_non_negative CHECK (min_stock_threshold >= 0)
);

CREATE TABLE stock_items (
	id BIGSERIAL PRIMARY KEY,
	warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	product_id BIGINT NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	quantity BIGINT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	CONSTRAINT stock_items_warehouse_product_unique UNIQUE (warehouse_id, product_id),
	CONSTRAINT stock_items_quantity_non_negative CHECK (quantity >= 0)
);

CREATE TABLE supplies (
	id BIGSERIAL PRIMARY KEY,
	warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	status TEXT NOT NULL DEFAULT 'Draft',
	created_by_user_id BIGINT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	confirmed_by_user_id BIGINT REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
	cancelled_by_user_id BIGINT REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	confirmed_at TIMESTAMPTZ,
	cancelled_at TIMESTAMPTZ,
	CONSTRAINT supplies_status_check CHECK (status IN ('Draft', 'Completed', 'Cancelled'))
);

CREATE TABLE supply_items (
	id BIGSERIAL PRIMARY KEY,
	supply_id BIGINT NOT NULL REFERENCES supplies(id) ON UPDATE CASCADE ON DELETE CASCADE,
	product_id BIGINT NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	quantity BIGINT NOT NULL,
	unit_price_cents BIGINT NOT NULL DEFAULT 0,
	CONSTRAINT supply_items_supply_product_unique UNIQUE (supply_id, product_id),
	CONSTRAINT supply_items_quantity_positive CHECK (quantity > 0),
	CONSTRAINT supply_items_unit_price_non_negative CHECK (unit_price_cents >= 0)
);

CREATE TABLE outbound_requests (
	id BIGSERIAL PRIMARY KEY,
	warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	status TEXT NOT NULL DEFAULT 'Draft',
	created_by_user_id BIGINT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	approved_by_user_id BIGINT REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
	executed_by_user_id BIGINT REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
	cancelled_by_user_id BIGINT REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	approved_at TIMESTAMPTZ,
	executed_at TIMESTAMPTZ,
	cancelled_at TIMESTAMPTZ,
	CONSTRAINT outbound_requests_status_check CHECK (status IN ('Draft', 'Approved', 'Completed', 'Cancelled'))
);

CREATE TABLE outbound_request_items (
	id BIGSERIAL PRIMARY KEY,
	outbound_request_id BIGINT NOT NULL REFERENCES outbound_requests(id) ON UPDATE CASCADE ON DELETE CASCADE,
	product_id BIGINT NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	quantity BIGINT NOT NULL,
	CONSTRAINT outbound_request_items_request_product_unique UNIQUE (outbound_request_id, product_id),
	CONSTRAINT outbound_request_items_quantity_positive CHECK (quantity > 0)
);

CREATE TABLE transfers (
	id BIGSERIAL PRIMARY KEY,
	source_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	target_warehouse_id BIGINT NOT NULL REFERENCES warehouses(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	status TEXT NOT NULL DEFAULT 'Draft',
	created_by_user_id BIGINT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	confirmed_by_user_id BIGINT REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
	cancelled_by_user_id BIGINT REFERENCES users(id) ON UPDATE CASCADE ON DELETE SET NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	confirmed_at TIMESTAMPTZ,
	cancelled_at TIMESTAMPTZ,
	CONSTRAINT transfers_status_check CHECK (status IN ('Draft', 'Completed', 'Cancelled')),
	CONSTRAINT transfers_different_warehouses CHECK (source_warehouse_id <> target_warehouse_id)
);

CREATE TABLE transfer_items (
	id BIGSERIAL PRIMARY KEY,
	transfer_id BIGINT NOT NULL REFERENCES transfers(id) ON UPDATE CASCADE ON DELETE CASCADE,
	product_id BIGINT NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
	quantity BIGINT NOT NULL,
	CONSTRAINT transfer_items_transfer_product_unique UNIQUE (transfer_id, product_id),
	CONSTRAINT transfer_items_quantity_positive CHECK (quantity > 0)
);

CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_is_active ON products(is_active);
CREATE INDEX idx_stock_items_warehouse_id ON stock_items(warehouse_id);
CREATE INDEX idx_stock_items_product_id ON stock_items(product_id);
CREATE INDEX idx_supplies_warehouse_id ON supplies(warehouse_id);
CREATE INDEX idx_supplies_status ON supplies(status);
CREATE INDEX idx_supply_items_product_id ON supply_items(product_id);
CREATE INDEX idx_outbound_requests_warehouse_id ON outbound_requests(warehouse_id);
CREATE INDEX idx_outbound_requests_status ON outbound_requests(status);
CREATE INDEX idx_outbound_request_items_product_id ON outbound_request_items(product_id);
CREATE INDEX idx_transfers_source_warehouse_id ON transfers(source_warehouse_id);
CREATE INDEX idx_transfers_target_warehouse_id ON transfers(target_warehouse_id);
CREATE INDEX idx_transfers_status ON transfers(status);
CREATE INDEX idx_transfer_items_product_id ON transfer_items(product_id);

-- +goose Down
DROP TABLE IF EXISTS transfer_items;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS outbound_request_items;
DROP TABLE IF EXISTS outbound_requests;
DROP TABLE IF EXISTS supply_items;
DROP TABLE IF EXISTS supplies;
DROP TABLE IF EXISTS stock_items;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS warehouses;
DROP TABLE IF EXISTS users;
