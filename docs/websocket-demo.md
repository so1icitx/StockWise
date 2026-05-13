# WebSocket Notifications Demo

StockWise exposes real-time notifications at:

```text
ws://localhost:8080/ws/notifications
```

The WebSocket stream is backed by the application service layer. Events are emitted by real service actions such as confirming a supply, approving or executing an outbound request, confirming a transfer, and deactivating a warehouse. There is no debug endpoint for manually emitting events.

## Connection Steps

Start PostgreSQL, run migrations and seed data, then start the API:

```bash
go run ./cmd/migrate up
go run ./cmd/seed
go run ./cmd/api
```

Connect a WebSocket client in a second terminal:

```bash
npx wscat -c ws://localhost:8080/ws/notifications
```

Alternative with `websocat`:

```bash
websocat ws://localhost:8080/ws/notifications
```

Keep that terminal open, then run real REST or GraphQL workflow commands in another terminal. The WebSocket client receives JSON messages automatically after the service action succeeds.

## Real Service Actions That Trigger Events

### Supply Confirmed

Confirming a real supply triggers `supply.confirmed`.

```bash
curl -X POST http://localhost:8080/api/v1/supplies/1/confirm \
  -H "X-User-ID: 2"
```

Example message:

```json
{
  "event": "supply.confirmed",
  "timestamp": "2026-04-30T12:00:00Z",
  "data": {
    "supply_id": 1,
    "warehouse_id": 1,
    "status": "Completed",
    "confirmed_by_user_id": 2,
    "item_count": 2,
    "items": [
      {
        "supply_item_id": 1,
        "product_id": 1,
        "quantity": 25,
        "unit_price_cents": 119900
      }
    ]
  }
}
```

### Outbound Approved

Approving a real outbound request triggers `outbound.approved`.

```bash
curl -X POST http://localhost:8080/api/v1/outbound-requests/1/approve \
  -H "X-User-ID: 2"
```

Example message:

```json
{
  "event": "outbound.approved",
  "timestamp": "2026-04-30T12:01:00Z",
  "data": {
    "outbound_request_id": 1,
    "warehouse_id": 1,
    "status": "Approved",
    "approved_by_user_id": 2,
    "item_count": 1,
    "items": [
      {
        "outbound_request_item_id": 1,
        "product_id": 1,
        "quantity": 5
      }
    ]
  }
}
```

### Outbound Completed

Executing an approved outbound request triggers `outbound.completed`. If the execution drops stock to the minimum threshold or to zero, the same real service action also triggers `stock.low` or `stock.out`.

```bash
curl -X POST http://localhost:8080/api/v1/outbound-requests/1/execute \
  -H "X-User-ID: 3"
```

Example messages:

```json
{
  "event": "outbound.completed",
  "timestamp": "2026-04-30T12:02:00Z",
  "data": {
    "outbound_request_id": 1,
    "warehouse_id": 1,
    "status": "Completed",
    "executed_by_user_id": 3,
    "item_count": 1,
    "items": [
      {
        "outbound_request_item_id": 1,
        "product_id": 1,
        "quantity": 5
      }
    ]
  }
}
```

```json
{
  "event": "stock.low",
  "timestamp": "2026-04-30T12:02:00Z",
  "data": {
    "warehouse_id": 1,
    "product_id": 1,
    "product_sku": "ELEC-001",
    "product_name": "Wireless Scanner",
    "quantity": 4,
    "min_stock_threshold": 5,
    "state": "LowStock",
    "triggered_by_event": "outbound.completed",
    "unit_of_measure": "pcs",
    "product_category_id": 1,
    "stock_item_id": 1
  }
}
```

### Transfer Confirmed

Confirming a real transfer triggers `transfer.confirmed`. If stock in the source warehouse reaches a low or zero state, the same service action also emits `stock.low` or `stock.out`.

```bash
curl -X POST http://localhost:8080/api/v1/transfers/1/confirm \
  -H "X-User-ID: 2"
```

Example message:

```json
{
  "event": "transfer.confirmed",
  "timestamp": "2026-04-30T12:03:00Z",
  "data": {
    "transfer_id": 1,
    "source_warehouse_id": 1,
    "target_warehouse_id": 2,
    "status": "Completed",
    "confirmed_by_user_id": 2,
    "item_count": 1,
    "items": [
      {
        "transfer_item_id": 1,
        "product_id": 1,
        "quantity": 3
      }
    ]
  }
}
```

### Stock Out

Any successful outbound execution or transfer confirmation that reduces a stock row to zero triggers `stock.out`.

```json
{
  "event": "stock.out",
  "timestamp": "2026-04-30T12:04:00Z",
  "data": {
    "warehouse_id": 1,
    "product_id": 1,
    "product_sku": "ELEC-001",
    "product_name": "Wireless Scanner",
    "quantity": 0,
    "min_stock_threshold": 5,
    "state": "OutOfStock",
    "triggered_by_event": "transfer.confirmed",
    "unit_of_measure": "pcs",
    "product_category_id": 1,
    "stock_item_id": 1
  }
}
```

### Warehouse Deactivated

Deactivating a real warehouse triggers `warehouse.deactivated`.

```bash
curl -X PATCH http://localhost:8080/api/v1/warehouses/2/deactivate \
  -H "X-User-ID: 1"
```

Example message:

```json
{
  "event": "warehouse.deactivated",
  "timestamp": "2026-04-30T12:05:00Z",
  "data": {
    "warehouse_id": 2,
    "name": "Retail Warehouse",
    "code": "RTL",
    "location": "Sofia",
    "is_active": false
  }
}
```

## Proof During Demo

Use this two-terminal flow:

1. Terminal A:

   ```bash
   npx wscat -c ws://localhost:8080/ws/notifications
   ```

2. Terminal B:

   ```bash
   curl -X POST http://localhost:8080/api/v1/outbound-requests/1/approve -H "X-User-ID: 2"
   curl -X POST http://localhost:8080/api/v1/outbound-requests/1/execute -H "X-User-ID: 3"
   ```

3. Terminal A receives messages like:

   ```text
   < {"event":"outbound.approved","timestamp":"2026-04-30T12:01:00Z","data":{"outbound_request_id":1,"warehouse_id":1,"status":"Approved","approved_by_user_id":2,"item_count":1,"items":[{"outbound_request_item_id":1,"product_id":1,"quantity":5}]}}
   < {"event":"outbound.completed","timestamp":"2026-04-30T12:02:00Z","data":{"outbound_request_id":1,"warehouse_id":1,"status":"Completed","executed_by_user_id":3,"item_count":1,"items":[{"outbound_request_item_id":1,"product_id":1,"quantity":5}]}}
   ```

Because the messages appear only after real REST or GraphQL operations complete, this proves the notification system is integrated with the business services rather than being simulated.
