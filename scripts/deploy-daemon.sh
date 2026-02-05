#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    log_error "Please run as root or with sudo"
    exit 1
fi

INSTALL_DIR="/opt/kadeem"
BIN_DIR="$INSTALL_DIR/bin"
SERVICE_FILE="/etc/systemd/system/kadeem-daemon.service"

log_info "Starting Kadeem daemon deployment..."

# Create kadeem user if doesn't exist
if ! id -u kadeem > /dev/null 2>&1; then
    log_step "Creating kadeem user..."
    useradd -r -s /bin/bash -d "$INSTALL_DIR" kadeem
    log_info "User created"
else
    log_info "User kadeem already exists"
fi

# Create directories
log_step "Creating installation directories..."
mkdir -p "$BIN_DIR"
chown -R kadeem:kadeem "$INSTALL_DIR"

# Build the daemon
log_step "Building daemon binary..."
cd "$(dirname "$0")/.."
go build -o "$BIN_DIR/daemon" cmd/daemon/main.go
chmod +x "$BIN_DIR/daemon"
chown kadeem:kadeem "$BIN_DIR/daemon"
log_info "Daemon binary built and installed"

# Run migrations
log_step "Running database migrations..."
source /etc/kadeem/.env
export DATABASE_URL
go run cmd/migrate/main.go up
log_info "Migrations completed"

# Install systemd service
if [ -f "systemd/kadeem-daemon.service" ]; then
    log_step "Installing systemd service..."
    cp systemd/kadeem-daemon.service "$SERVICE_FILE"
    systemctl daemon-reload
    log_info "Systemd service installed"
else
    log_warn "Systemd service file not found, skipping"
fi

# Enable and restart service
log_step "Starting daemon service..."
systemctl enable kadeem-daemon
systemctl restart kadeem-daemon

# Wait a moment for service to start
sleep 2

# Check service status
if systemctl is-active --quiet kadeem-daemon; then
    log_info "✓ Daemon is running successfully"
    systemctl status kadeem-daemon --no-pager -l
else
    log_error "✗ Daemon failed to start"
    journalctl -u kadeem-daemon -n 20 --no-pager
    exit 1
fi

echo ""
log_info "Deployment complete!"
log_info "View logs with: journalctl -u kadeem-daemon -f"
