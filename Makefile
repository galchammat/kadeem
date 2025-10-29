# Env Vars
ENV_FILE := $(PWD)/.env
BIN_DIR := $(PWD)/bin
export BIN_DIR ENV_FILE

deps:
	sudo apt update
	sudo add-apt-repository universe
	sudo apt install -y pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev
	cd frontend && npm install

run:
	wails dev -tags webkit2_41
build:
	wails build -tags webkit2_41

test-integration:
	@echo "Running integration tests..."
	RUN_INTEGRATION_TESTS=true go test -v ./tests/...

migrate:
	mkdir -p ./bin
	go run cmd/migrate/main.go

# Usage: make migrate-create name=description_of_migration
migrate-create:
	migrate create -seq -digits 3 -dir ./migrations -ext sql ${name}

migrate-force:
	migrate -database "sqlite3://kadeem.db" -path ./migrations force ${version}
