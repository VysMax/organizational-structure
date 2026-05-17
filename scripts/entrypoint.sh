#!/bin/sh

export GOOSE_DRIVER=postgres
export DB_DSN="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable"

echo "Checking migration status..."
goose -dir ./migrations status

echo "Applying migrations..."
goose -dir ./migrations up

./service