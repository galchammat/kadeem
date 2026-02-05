# Env Vars
ENV_FILE := $(PWD)/.env
BIN_DIR := $(PWD)/bin
export BIN_DIR ENV_FILE

.PHONY: deps run build tests test migrate migrate-create migrate-force migrate-reset
.PHONY: ansible migrate-up migrate-down migrate-version migrate-force-reset

# Development
deps:
	sudo apt update
	sudo add-apt-repository universe
	sudo apt install -y pkg-config libgtk-3-dev libwebkit2gtk-4.1-dev
	cd frontend && npm install

run:
	wails dev -tags webkit2_41

build:
	wails build -tags webkit2_41

tests:
	@echo "Running integration tests..."
	go test -v ./tests/integration

test:
	@echo "Running integration test $(filter-out $@,$(MAKECMDGOALS))..."
	go test -v ./tests/integration -run "$(filter-out $@,$(MAKECMDGOALS))"

%: 
	@:

# Infrastructure
ansible:
	@echo "Running Ansible PostgreSQL setup..."
	@if [ ! -f $(ENV_FILE) ]; then \
		echo "Error: .env file not found at $(ENV_FILE)"; \
		exit 1; \
	fi
	@ANSIBLE_ARGS=""; \
	if echo "$(filter-out ansible,$(MAKECMDGOALS))" | grep -q "check"; then \
		ANSIBLE_ARGS="--check --diff"; \
		echo "Running in check mode (dry run)..."; \
	fi; \
	export $$(cat $(ENV_FILE) | grep -v '^#' | xargs) && \
	ansible-playbook -i ansible/inventory/production.yml ansible/playbook.yml $$ANSIBLE_ARGS

# Database
migrate-create:
	migrate create -seq -digits 3 -dir ./migrations -ext sql ${name}

migrate-up:
	@echo "Running PostgreSQL migrations..."
	go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back PostgreSQL migrations..."
	go run cmd/migrate/main.go down

migrate-version:
	@echo "Checking migration version..."
	go run cmd/migrate/main.go version

migrate-force:
	@echo "Forcing migration to version $(filter-out $@,$(MAKECMDGOALS))..."
	@if [ -z "$(filter-out $@,$(MAKECMDGOALS))" ]; then \
		echo "Error: version is required. Use: make migrate-force 1"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go force $(filter-out $@,$(MAKECMDGOALS))

migrate-force-reset:
	@echo "Dropping all tables and resetting database..."
	go run cmd/migrate/main.go drop-all
	@echo "Running PostgreSQL migrations..."
	go run cmd/migrate/main.go up

# Help
help:
	@echo "Kadeem Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make deps              - Install dependencies"
	@echo "  make run               - Run Wails app in dev mode"
	@echo "  make build             - Build Wails app"
	@echo "  make tests             - Run all tests"
	@echo ""
	@echo "Infrastructure:"
	@echo "  make ansible           - Setup PostgreSQL with Ansible"
	@echo "  make ansible check     - Dry run (show what would change)"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-create       - Create new migration (name=...)"
	@echo "  make migrate-up           - Run database migrations"
	@echo "  make migrate-down         - Rollback migrations"
	@echo "  make migrate-version      - Show current migration version"
	@echo "  make migrate-force        - Force migration version (make migrate-force 0)"
	@echo "  make migrate-force-reset  - Force to version 0 and re-run all migrations"
