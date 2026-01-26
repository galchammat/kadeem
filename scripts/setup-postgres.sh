#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    log_error "Please run as root or with sudo"
    exit 1
fi

log_info "Starting PostgreSQL setup for Kadeem..."

# Install PostgreSQL if not already installed
if ! command -v psql &> /dev/null; then
    log_info "PostgreSQL not found. Installing..."
    apt-get update
    apt-get install -y postgresql postgresql-contrib
    log_info "PostgreSQL installed successfully"
else
    log_info "PostgreSQL already installed"
fi

# Start and enable PostgreSQL
log_info "Starting PostgreSQL service..."
systemctl start postgresql
systemctl enable postgresql

# Generate a strong random password if not provided
if [ -z "$DB_PASSWORD" ]; then
    DB_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
    log_info "Generated random database password"
fi

DB_NAME="${DB_NAME:-kadeem}"
DB_USER="${DB_USER:-kadeem}"

log_info "Creating database and user..."

# Switch to postgres user and create database/user
sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
    sudo -u postgres psql -c "CREATE DATABASE $DB_NAME;"

sudo -u postgres psql -tc "SELECT 1 FROM pg_user WHERE usename = '$DB_USER'" | grep -q 1 || \
    sudo -u postgres psql -c "CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASSWORD';"

sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"
sudo -u postgres psql -c "ALTER DATABASE $DB_NAME OWNER TO $DB_USER;"

log_info "Database and user created successfully"

# Configure PostgreSQL for local and network access
PG_VERSION=$(psql --version | awk '{print $3}' | cut -d. -f1)
PG_HBA="/etc/postgresql/$PG_VERSION/main/pg_hba.conf"
PG_CONF="/etc/postgresql/$PG_VERSION/main/postgresql.conf"

log_info "Configuring PostgreSQL access..."

# Backup original config
cp "$PG_HBA" "$PG_HBA.backup.$(date +%s)"

# Add local and network access rules (if not already present)
grep -q "kadeem" "$PG_HBA" || cat >> "$PG_HBA" << EOF

# Kadeem application access
local   $DB_NAME        $DB_USER                                md5
host    $DB_NAME        $DB_USER        127.0.0.1/32            md5
host    $DB_NAME        $DB_USER        ::1/128                 md5
EOF

# Enable network listening (for Tailscale access)
if ! grep -q "^listen_addresses = '\*'" "$PG_CONF"; then
    log_info "Enabling network access on all interfaces..."
    sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" "$PG_CONF"
fi

# Restart PostgreSQL to apply changes
log_info "Restarting PostgreSQL..."
systemctl restart postgresql

# Create config directory for Kadeem
mkdir -p /etc/kadeem
chmod 755 /etc/kadeem

# Generate .env file
ENV_FILE="/etc/kadeem/.env"
log_info "Creating environment file at $ENV_FILE..."

cat > "$ENV_FILE" << EOF
DATABASE_URL=postgres://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME?sslmode=disable
RIOT_API_KEY=${RIOT_API_KEY:-your_riot_api_key_here}
DISCORD_WEBHOOK_URL=${DISCORD_WEBHOOK_URL:-}
EOF

chmod 600 "$ENV_FILE"

log_info "PostgreSQL setup complete!"
echo ""
log_info "Database details:"
echo "  Database: $DB_NAME"
echo "  User: $DB_USER"
echo "  Password: $DB_PASSWORD"
echo ""
log_info "Environment file created at: $ENV_FILE"
log_info "Connection string: postgres://$DB_USER:$DB_PASSWORD@localhost:5432/$DB_NAME"
echo ""
log_warn "Next steps:"
echo "  1. Update RIOT_API_KEY in $ENV_FILE"
echo "  2. Run migrations: make migrate-up"
echo "  3. Install daemon: make install-daemon"
