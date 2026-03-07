include .env
export

# Database URL for golang-migrate
DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

# Migrate binary path
MIGRATE=$(shell go env GOPATH)/bin/migrate

.PHONY: help migrate-create migrate-up migrate-down migrate-drop seed db-reset

help:
	@echo "Available commands:"
	@echo "  make migrate-create name=<name>  Create a new migration pair"
	@echo "  make migrate-up                  Apply all pending migrations"
	@echo "  make migrate-down                Roll back the last migration"
	@echo "  make migrate-drop                Drop all tables and recreate public schema (DANGER)"
	@echo "  make seed                         Run the Go seeder script"
	@echo "  make db-reset                     Drop, Up, and Seed (DANGER)"

migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	$(MIGRATE) create -ext sql -dir migrations -seq $(name)

migrate-up:
	$(MIGRATE) -path migrations -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path migrations -database "$(DB_URL)" down 1

migrate-drop:
	$(MIGRATE) -path migrations -database "$(DB_URL)" drop -f

seed:
	go run cmd/seeder/main.go

db-reset: migrate-drop migrate-up seed
