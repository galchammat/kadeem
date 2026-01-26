#!/bin/bash
set -e

# Configuration - override with environment variables
BACKUP_HOST="${BACKUP_HOST:-}"
BACKUP_USER="${BACKUP_USER:-}"
BACKUP_PATH="${BACKUP_PATH:-/mnt/2tb/kadeem-backups}"
RETENTION_DAYS="${RETENTION_DAYS:-30}"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Load environment
if [ -f "/etc/kadeem/.env" ]; then
    source /etc/kadeem/.env
else
    log_error "/etc/kadeem/.env not found"
    exit 1
fi

# Validate configuration
if [ -z "$BACKUP_HOST" ]; then
    log_error "BACKUP_HOST not set. Set it in /etc/kadeem/.env or as environment variable"
    exit 1
fi

# Extract database details
DB_USER=$(echo "$DATABASE_URL" | sed -n 's|.*://\([^:]*\):.*|\1|p')
DB_PASS=$(echo "$DATABASE_URL" | sed -n 's|.*://[^:]*:\([^@]*\)@.*|\1|p')
DB_HOST=$(echo "$DATABASE_URL" | sed -n 's|.*@\([^:]*\):.*|\1|p')
DB_NAME=$(echo "$DATABASE_URL" | sed -n 's|.*/\([^?]*\).*|\1|p')

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="kadeem_backup_${TIMESTAMP}.sql.gz"
LOCAL_BACKUP="/tmp/$BACKUP_FILE"

log_info "Starting PostgreSQL backup..."

# Create backup
PGPASSWORD="$DB_PASS" pg_dump -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" \
    --format=plain \
    --no-owner \
    --no-privileges | gzip > "$LOCAL_BACKUP"

if [ $? -eq 0 ]; then
    BACKUP_SIZE=$(du -h "$LOCAL_BACKUP" | cut -f1)
    log_info "Backup created: $BACKUP_FILE ($BACKUP_SIZE)"
else
    log_error "Backup failed!"
    exit 1
fi

# Transfer to Machine A
log_info "Transferring backup to $BACKUP_HOST:$BACKUP_PATH..."

# Create remote directory if needed
ssh "${BACKUP_USER}@${BACKUP_HOST}" "mkdir -p $BACKUP_PATH"

# Copy backup
if scp "$LOCAL_BACKUP" "${BACKUP_USER}@${BACKUP_HOST}:${BACKUP_PATH}/"; then
    log_info "Backup transferred successfully"
else
    log_error "Failed to transfer backup"
    rm "$LOCAL_BACKUP"
    exit 1
fi

# Clean up local backup
rm "$LOCAL_BACKUP"

# Clean old backups on remote (keep last N days)
log_info "Cleaning old backups (keeping last $RETENTION_DAYS days)..."
ssh "${BACKUP_USER}@${BACKUP_HOST}" \
    "find $BACKUP_PATH -name 'kadeem_backup_*.sql.gz' -mtime +$RETENTION_DAYS -delete" || \
    log_warn "Failed to clean old backups"

# List recent backups
log_info "Recent backups on $BACKUP_HOST:"
ssh "${BACKUP_USER}@${BACKUP_HOST}" \
    "ls -lh $BACKUP_PATH/kadeem_backup_*.sql.gz | tail -5" || true

log_info "Backup complete!"

# Optional: Send Discord notification
if [ -n "$DISCORD_WEBHOOK_URL" ]; then
    curl -X POST "$DISCORD_WEBHOOK_URL" \
        -H "Content-Type: application/json" \
        -d "{\"content\": \"✅ Kadeem database backup completed successfully\n**File:** $BACKUP_FILE\n**Size:** $BACKUP_SIZE\"}" \
        &> /dev/null || log_warn "Failed to send Discord notification"
fi
