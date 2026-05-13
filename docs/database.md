# StockWise Database Design

The database is PostgreSQL and is versioned through Goose SQL migrations. The schema mirrors the required Bulgarian assignment entities and keeps business-critical invariants close to the data where possible.

## ER Diagram

```mermaid
erDiagram
    USERS ||--o{ SUPPLIES : creates
    USERS ||--o{ OUTBOUND_REQUESTS : creates
    USERS ||--o{ TRANSFERS : creates
    WAREHOUSES ||--o{ STOCK_ITEMS : stores
    PRODUCTS ||--o{ STOCK_ITEMS : appears_in
    CATEGORIES ||--o{ PRODUCTS : groups
    WAREHOUSES ||--o{ SUPPLIES : receives
    SUPPLIES ||--o{ SUPPLY_ITEMS : contains
    PRODUCTS ||--o{ SUPPLY_ITEMS : supplied
    WAREHOUSES ||--o{ OUTBOUND_REQUESTS : fulfills
    OUTBOUND_REQUESTS ||--o{ OUTBOUND_REQUEST_ITEMS : contains
    PRODUCTS ||--o{ OUTBOUND_REQUEST_ITEMS : requested
    WAREHOUSES ||--o{ TRANSFERS : source
    WAREHOUSES ||--o{ TRANSFERS : target
    TRANSFERS ||--o{ TRANSFER_ITEMS : contains
    PRODUCTS ||--o{ TRANSFER_ITEMS : transferred

    USERS {
        bigint id PK
        text name
        text email UK
        text role
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    WAREHOUSES {
        bigint id PK
        text name
        text code UK
        text location
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    CATEGORIES {
        bigint id PK
        text name UK
        text description
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    PRODUCTS {
        bigint id PK
        text name
        text sku UK
        bigint category_id FK
        text unit_of_measure
        bigint min_stock_threshold
        boolean is_active
        timestamptz created_at
        timestamptz updated_at
    }

    STOCK_ITEMS {
        bigint id PK
        bigint warehouse_id FK
        bigint product_id FK
        bigint quantity
        timestamptz created_at
        timestamptz updated_at
    }

    SUPPLIES {
        bigint id PK
        bigint warehouse_id FK
        text status
        bigint created_by_user_id FK
        bigint confirmed_by_user_id FK
        bigint cancelled_by_user_id FK
        timestamptz created_at
        timestamptz updated_at
        timestamptz confirmed_at
        timestamptz cancelled_at
    }

    SUPPLY_ITEMS {
        bigint id PK
        bigint supply_id FK
        bigint product_id FK
        bigint quantity
        bigint unit_price_cents
    }

    OUTBOUND_REQUESTS {
        bigint id PK
        bigint warehouse_id FK
        text status
        bigint created_by_user_id FK
        bigint approved_by_user_id FK
        bigint executed_by_user_id FK
        bigint cancelled_by_user_id FK
        timestamptz created_at
        timestamptz updated_at
        timestamptz approved_at
        timestamptz executed_at
        timestamptz cancelled_at
    }

    OUTBOUND_REQUEST_ITEMS {
        bigint id PK
        bigint outbound_request_id FK
        bigint product_id FK
        bigint quantity
    }

    TRANSFERS {
        bigint id PK
        bigint source_warehouse_id FK
        bigint target_warehouse_id FK
        text status
        bigint created_by_user_id FK
        bigint confirmed_by_user_id FK
        bigint cancelled_by_user_id FK
        timestamptz created_at
        timestamptz updated_at
        timestamptz confirmed_at
        timestamptz cancelled_at
    }

    TRANSFER_ITEMS {
        bigint id PK
        bigint transfer_id FK
        bigint product_id FK
        bigint quantity
    }
```

## Required Constraints

- `products.sku` is unique.
- `warehouses.code` is unique.
- `stock_items(warehouse_id, product_id)` is unique, so each product has one stock row per warehouse.
- `stock_items.quantity` cannot be negative.
- `supply_items.quantity`, `outbound_request_items.quantity`, and `transfer_items.quantity` must be positive.
- `supply_items.unit_price_cents` cannot be negative.
- `products.min_stock_threshold` cannot be negative.
- `transfers.source_warehouse_id` and `transfers.target_warehouse_id` must be different.
- Operation statuses are limited to the domain-supported lifecycle values.

## Deletion Behavior

The schema uses restrictive foreign keys for products, warehouses, and categories. This supports the assignment rules:

- Products that participate in movements cannot be hard-deleted because operation item rows reference them.
- Warehouses with stock or operations cannot be hard-deleted because stock and operation rows reference them.
- Categories with products cannot be hard-deleted. Service logic additionally blocks deletion while active products remain.
- Products with confirmed movements are soft-deactivated by the product service instead of being hard-deleted.
- Categories with inactive products can be deactivated by service logic when hard deletion would break references.

## Seed Data

The seed command is idempotent. It uses natural keys such as user email, warehouse code, category name, and product SKU with `ON CONFLICT` upserts.

Run after migrations:

```bash
go run ./cmd/seed
```
