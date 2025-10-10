# Env Vars
ENV_FILE := $(PWD)/.env
BIN_DIR := $(PWD)/bin
export BIN_DIR ENV_FILE

run:
	wails dev -tags webkit2_41
build:
	wails build -tags webkit2_41

test-integration:
	@echo "Running integration tests..."
	RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./tests/...

migrate:
	go run cmd/migrate/main.go

# Usage: make migrate-create name=description_of_migration
migrate-create:
	migrate create -seq -digits 3 -dir ./migrations -ext sql ${name}