# Env Vars
BIN_DIR := $(PWD)/bin
SERVER_DIR := packages/server
WEB_DIR := packages/web
export BIN_DIR

.PHONY: deps run build tests test migrate migrate-create migrate-force migrate-reset
.PHONY: ansible migrate-up migrate-down migrate-version migrate-force-reset

# Development
deps:
	cd $(WEB_DIR) && npm install

run:
	@echo "Starting API server and Vite dev server..."
	@trap 'kill 0' EXIT; \
	cd $(SERVER_DIR) && go run cmd/daemon/main.go & \
	cd $(WEB_DIR) && npm run dev & \
	wait

build:
	@echo "Building Go daemon..."
	cd $(SERVER_DIR) && CGO_ENABLED=0 go build -o ../../bin/daemon cmd/daemon/main.go
	@echo "Building frontend..."
	cd $(WEB_DIR) && npm run build

tests:
	@echo "Running integration tests..."
	cd $(SERVER_DIR) && go test -v ./tests/integration

test:
	@echo "Running integration test $(filter-out $@,$(MAKECMDGOALS))..."
	cd $(SERVER_DIR) && go test -v ./tests/integration -run "$(filter-out $@,$(MAKECMDGOALS))"

%: 
	@:

# Infrastructure
ansible:
	@echo "Running Ansible PostgreSQL setup..."
	@ANSIBLE_ARGS=""; \
	if echo "$(filter-out ansible,$(MAKECMDGOALS))" | grep -q "check"; then \
		ANSIBLE_ARGS="$$ANSIBLE_ARGS --check --diff"; \
		echo "Running in check mode (dry run)..."; \
	fi; \
	ansible-playbook -i ansible/inventory/production.yml ansible/playbook.yml $$ANSIBLE_ARGS

# Database
migrate-create:
	migrate create -seq -digits 3 -dir ./$(SERVER_DIR)/migrations -ext sql ${name}

migrate-up:
	@echo "Running PostgreSQL migrations..."
	cd $(SERVER_DIR) && go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back PostgreSQL migrations..."
	cd $(SERVER_DIR) && go run cmd/migrate/main.go down

migrate-version:
	@echo "Checking migration version..."
	cd $(SERVER_DIR) && go run cmd/migrate/main.go version

migrate-force:
	@echo "Forcing migration to version $(filter-out $@,$(MAKECMDGOALS))..."
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Error: version is required. Use: make migrate-force 1"; \
		exit 1; \
	fi
	cd $(SERVER_DIR) && go run cmd/migrate/main.go force $(filter-out $@,$(MAKECMDGOALS))

migrate-force-reset:
	@echo "Dropping all tables and resetting database..."
	cd $(SERVER_DIR) && go run cmd/migrate/main.go drop-all
	@echo "Running PostgreSQL migrations..."
	cd $(SERVER_DIR) && go run cmd/migrate/main.go up

# Help
help:
	@echo "Kadeem Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make deps              - Install frontend dependencies"
	@echo "  make run               - Run API server + Vite dev server"
	@echo "  make build             - Build Go daemon + frontend"
	@echo "  make tests             - Run all integration tests"
	@echo ""
	@echo "Infrastructure:"
	@echo "  make ansible           - Run Ansible playbook"
	@echo "  make ansible check     - Dry run (show what would change)"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-create       - Create new migration (name=...)"
	@echo "  make migrate-up           - Run database migrations"
	@echo "  make migrate-down         - Rollback migrations"
	@echo "  make migrate-version      - Show current migration version"
	@echo "  make migrate-force        - Force migration version (make migrate-force 0)"
	@echo "  make migrate-force-reset  - Drop all and re-run migrations"
