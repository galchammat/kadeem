# Env Vars
ENV_FILE := $(PWD)/.env
BIN_DIR := $(PWD)/bin
export BIN_DIR ENV_FILE

.PHONY: deps run build tests test migrate migrate-create migrate-force migrate-reset
.PHONY: db-setup build-daemon deploy-daemon migrate-up migrate-down
.PHONY: daemon-start daemon-stop daemon-restart daemon-status daemon-logs
.PHONY: healthcheck backup install-services

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

# Database Management (Legacy SQLite)
migrate:
	mkdir -p ./bin
	go run cmd/migrate/main.go

migrate-create:
	migrate create -seq -digits 3 -dir ./migrations -ext sql ${name}

migrate-force:
	migrate -database "sqlite3://kadeem.db" -path ./migrations force ${version}

migrate-reset:
	rm -f ./bin/kadeem.db && $(MAKE) migrate

# PostgreSQL Setup
db-setup:
	@echo "Setting up PostgreSQL..."
	sudo ./scripts/setup-postgres.sh

migrate-up:
	@echo "Running PostgreSQL migrations..."
	go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back PostgreSQL migrations..."
	go run cmd/migrate/main.go down

# Daemon Management
build-daemon:
	@echo "Building daemon..."
	go build -o bin/daemon cmd/daemon/main.go
	chmod +x bin/daemon

deploy-daemon:
	@echo "Deploying daemon..."
	sudo ./scripts/deploy-daemon.sh

daemon-start:
	sudo systemctl start kadeem-daemon

daemon-stop:
	sudo systemctl stop kadeem-daemon

daemon-restart:
	sudo systemctl restart kadeem-daemon

daemon-status:
	sudo systemctl status kadeem-daemon

daemon-logs:
	sudo journalctl -u kadeem-daemon -f

# Monitoring & Maintenance
healthcheck:
	@echo "Running health checks..."
	sudo ./scripts/healthcheck.sh

backup:
	@echo "Running database backup..."
	sudo ./scripts/backup-database.sh

# Installation
install-services:
	@echo "Installing systemd services..."
	sudo cp systemd/kadeem-daemon.service /etc/systemd/system/
	sudo cp systemd/kadeem-notify-failure.service /etc/systemd/system/
	sudo cp systemd/kadeem-backup.service /etc/systemd/system/
	sudo cp systemd/kadeem-backup.timer /etc/systemd/system/
	sudo mkdir -p /opt/kadeem/scripts
	sudo cp scripts/*.sh /opt/kadeem/scripts/
	sudo chmod +x /opt/kadeem/scripts/*.sh
	sudo systemctl daemon-reload
	@echo "Services installed. Enable with:"
	@echo "  sudo systemctl enable kadeem-daemon"
	@echo "  sudo systemctl enable kadeem-backup.timer"

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
	@echo "Database:"
	@echo "  make db-setup          - Install and configure PostgreSQL"
	@echo "  make migrate-up        - Run database migrations"
	@echo "  make migrate-down      - Rollback migrations"
	@echo "  make migrate-create    - Create new migration (name=...)"
	@echo ""
	@echo "Daemon:"
	@echo "  make build-daemon      - Build daemon binary"
	@echo "  make deploy-daemon     - Deploy daemon to system"
	@echo "  make daemon-start      - Start daemon service"
	@echo "  make daemon-stop       - Stop daemon service"
	@echo "  make daemon-restart    - Restart daemon service"
	@echo "  make daemon-status     - Check daemon status"
	@echo "  make daemon-logs       - View daemon logs (live)"
	@echo ""
	@echo "Maintenance:"
	@echo "  make healthcheck       - Run system health checks"
	@echo "  make backup            - Backup database to Machine A"
	@echo "  make install-services  - Install all systemd services"
