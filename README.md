<div align="center">
<h1>StockWise</h1>

Inventory and warehouse API written in Go. Tracks warehouses, products,
stock levels, supplies, outbound requests, transfers, product movements, and
real-time stock events.

[![Go Version](https://img.shields.io/github/go-mod/go-version/so1icitx/StockWise)](go.mod)
[![License](https://img.shields.io/github/license/so1icitx/StockWise)](LICENSE)

</div>

## How it works

StockWise is split into domain, application, infrastructure, and transport
layers. HTTP, GraphQL, and WebSocket handlers do not update data directly;
they call application services. The services enforce business rules and use
repositories for persistence.

Incoming stock is handled through supplies. A supply starts as a draft, gets
line items, then confirmation increases stock in the target warehouse.

Outgoing stock is handled through outbound requests. A request is created,
filled with items, approved, then executed. Execution checks available stock
before subtracting anything.

Warehouse transfers move stock between two different warehouses in one
transaction. Confirmation subtracts from the source warehouse and adds to the
target warehouse.

Stock changes are evaluated after each confirmed supply, executed outbound
request, or confirmed transfer. Low-stock and out-of-stock states are reported
through WebSocket notifications.

## Stack

- Go
- Gin
- GORM
- PostgreSQL
- Goose migrations
- gqlgen
- gorilla/websocket
- Testcontainers for integration tests

## Layout

```text
cmd/
  api/        HTTP server entrypoint
  migrate/    Goose migration runner
  seed/       demo data command
internal/
  domain/     entities, statuses, domain helpers
  application/ services, repository contracts, business rules
  infrastructure/postgres/ GORM repositories and database setup
  transport/httpapi/ REST routes, DTOs, validation, Swagger
  transport/graphql/ gqlgen schema and resolvers
  transport/websocket/ notification hub
migrations/   PostgreSQL schema migrations
docs/         API examples, database notes, OpenAPI spec
tests/        unit and integration tests
```

## Business rules

- Product SKU is unique.
- Quantities must be positive.
- Supply item price cannot be negative.
- A supply cannot be confirmed without items.
- An outbound request cannot execute without enough stock.
- Transfers cannot target the same warehouse.
- A transfer cannot be confirmed without items.
- Confirmed transfers subtract from the source and add to the target.
- Inactive warehouses cannot be used in new supplies, outbound requests, or transfers.
- Completed and cancelled operations cannot be edited.
- Warehouses with stock or active operations cannot be deleted.
- Products with confirmed movements are deactivated instead of hard-deleted.
- Categories with active products cannot be deleted.

## Getting started

Start PostgreSQL:

```bash
docker run --name stockwise-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=stockwise \
  -p 5432:5432 \
  -d postgres:16
```

Set the connection string:

```bash
export DATABASE_URL='postgres://postgres:postgres@127.0.0.1:5432/stockwise?sslmode=disable'
```

Run migrations, seed data, and start the API:

```bash
go run ./cmd/migrate up
go run ./cmd/seed
go run ./cmd/api
```

The API listens on `http://localhost:8080`.

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/health
```

## REST

REST endpoints are served under `/api/v1`.

Main resources:

- `/api/v1/users`
- `/api/v1/warehouses`
- `/api/v1/categories`
- `/api/v1/products`
- `/api/v1/supplies`
- `/api/v1/outbound-requests`
- `/api/v1/transfers`
- `/api/v1/products/low-stock`
- `/api/v1/products/{id}/movements`
- `/api/v1/warehouses/{id}/movements`

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

OpenAPI document:

```text
http://localhost:8080/swagger/openapi.yaml
```

Example REST workflows are in [`docs/rest-examples.md`](docs/rest-examples.md).

## GraphQL

GraphQL is served at:

```text
http://localhost:8080/graphql
```

Playground:

```text
http://localhost:8080/graphql/playground
```

Example:

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query { products { id sku name } warehouses { id code name } }"}'
```

More examples are in [`docs/graphql-examples.md`](docs/graphql-examples.md).

## WebSocket notifications

Clients connect to:

```text
ws://localhost:8080/ws/notifications
```

Events are emitted by real service actions:

- `supply.confirmed`
- `outbound.approved`
- `outbound.completed`
- `transfer.confirmed`
- `stock.low`
- `stock.out`
- `warehouse.deactivated`

Example:

```bash
npx wscat -c ws://localhost:8080/ws/notifications
```

Then confirm a real supply or execute an outbound request from another
terminal. Demo steps and example payloads are in
[`docs/websocket-demo.md`](docs/websocket-demo.md).

## Database

Migrations are stored in `migrations/` and run through Goose.

```bash
go run ./cmd/migrate status
go run ./cmd/migrate up
go run ./cmd/migrate down
```

The schema enforces unique product SKUs, unique warehouse codes, one stock row
per warehouse/product pair, positive operation quantities, non-negative stock,
and non-negative supply prices.

Database notes and the ER diagram are in [`docs/database.md`](docs/database.md).

## Tests

Run everything:

```bash
go test ./...
```

Run unit tests:

```bash
go test ./tests/unit/...
```

Run integration tests:

```bash
go test ./tests/integration/... -count=1 -v
```

Integration tests use Testcontainers PostgreSQL, real Goose migrations, a real
Gin server, REST calls, GraphQL calls, and a WebSocket client.

Full local verification:

```bash
./scripts/test.sh
```

## Useful commands

```bash
go mod tidy
go test ./...
go vet ./...
go build ./cmd/api ./cmd/migrate ./cmd/seed
```
