#!/usr/bin/env bash
set -e

echo "Starting StockWise development environment..."

CONTAINER_NAME="stockwise-postgres"

if ! docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
  echo "Creating PostgreSQL container..."
  docker run --name "$CONTAINER_NAME" \
    -e POSTGRES_USER=stockwise \
    -e POSTGRES_PASSWORD=stockwise \
    -e POSTGRES_DB=stockwise \
    -p 5432:5432 \
    -d postgres:16
else
  echo "PostgreSQL container already exists."
  docker start "$CONTAINER_NAME" >/dev/null || true
fi

echo "Waiting for PostgreSQL..."
sleep 4

export DATABASE_URL="postgres://stockwise:stockwise@localhost:5432/stockwise?sslmode=disable"

echo "Running migrations..."
go run ./cmd/migrate up

echo "Running seed data..."
go run ./cmd/seed

echo "Starting API server..."
go run ./cmd/api