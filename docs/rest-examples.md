# StockWise REST Examples

All business workflow requests that create, approve, confirm, execute, or cancel operations require an `X-User-ID` header. The seeded users include IDs for an admin, manager, and operator after `go run ./cmd/seed`.

Base URL:

```text
http://localhost:8080/api/v1
```

## Create Warehouse

```bash
curl -X POST http://localhost:8080/api/v1/warehouses \
  -H "Content-Type: application/json" \
  -d '{
    "name": "North Warehouse",
    "code": "WH-NORTH",
    "location": "Ruse, Bulgaria"
  }'
```

## Create Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Industrial Scale",
    "sku": "EQUIP-001",
    "category_id": 1,
    "unit_of_measure": "pcs",
    "min_stock_threshold": 2
  }'
```

## Confirm Supply

Create the supply:

```bash
curl -X POST http://localhost:8080/api/v1/supplies \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 3" \
  -d '{
    "warehouse_id": 1
  }'
```

Add an item:

```bash
curl -X POST http://localhost:8080/api/v1/supplies/1/items \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 10,
    "unit_price_cents": 12999
  }'
```

Confirm the supply:

```bash
curl -X POST http://localhost:8080/api/v1/supplies/1/confirm \
  -H "X-User-ID: 2"
```

## Execute Outbound Request

Create the outbound request:

```bash
curl -X POST http://localhost:8080/api/v1/outbound-requests \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 3" \
  -d '{
    "warehouse_id": 1
  }'
```

Add an item:

```bash
curl -X POST http://localhost:8080/api/v1/outbound-requests/1/items \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 2
  }'
```

Approve it:

```bash
curl -X POST http://localhost:8080/api/v1/outbound-requests/1/approve \
  -H "X-User-ID: 2"
```

Execute it:

```bash
curl -X POST http://localhost:8080/api/v1/outbound-requests/1/execute \
  -H "X-User-ID: 3"
```

## Confirm Transfer

Create the transfer:

```bash
curl -X POST http://localhost:8080/api/v1/transfers \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 3" \
  -d '{
    "source_warehouse_id": 1,
    "target_warehouse_id": 2
  }'
```

Add an item:

```bash
curl -X POST http://localhost:8080/api/v1/transfers/1/items \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "quantity": 1
  }'
```

Confirm the transfer:

```bash
curl -X POST http://localhost:8080/api/v1/transfers/1/confirm \
  -H "X-User-ID: 2"
```

## Get Low-Stock Products

```bash
curl http://localhost:8080/api/v1/products/low-stock
```
