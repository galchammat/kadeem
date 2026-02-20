#!/bin/bash

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

EXIT_CODE=0

check_service() {
    local service=$1
    if systemctl is-active --quiet "$service"; then
        echo -e "${GREEN}✓${NC} $service is running"
        return 0
    else
        echo -e "${RED}✗${NC} $service is NOT running"
        return 1
    fi
}

check_database() {
    if [ -f "/etc/kadeem/kadeem.env" ]; then
        source /etc/kadeem/kadeem.env
    elif [ -f "/etc/kadeem/.env" ]; then
        source /etc/kadeem/.env
    fi
    
    if [ -z "$DATABASE_URL" ]; then
        echo -e "${RED}✗${NC} DATABASE_URL not set"
        return 1
    fi
    
    # Extract connection details from DATABASE_URL
    # Format: postgres://user:pass@host:port/dbname
    DB_CHECK=$(echo "$DATABASE_URL" | sed -n 's|.*://\([^:]*\):\([^@]*\)@\([^:]*\):\([^/]*\)/\(.*\)|\1 \3 \5|p')
    read -r DB_USER DB_HOST DB_NAME <<< "$DB_CHECK"
    
    if PGPASSWORD=$(echo "$DATABASE_URL" | sed -n 's|.*://[^:]*:\([^@]*\)@.*|\1|p') \
       psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1" &> /dev/null; then
        echo -e "${GREEN}✓${NC} PostgreSQL connection OK"
        return 0
    else
        echo -e "${RED}✗${NC} PostgreSQL connection FAILED"
        return 1
    fi
}

check_daemon_logs() {
    local errors=$(journalctl -u kadeem-daemon --since "5 minutes ago" -p err -q | wc -l)
    if [ "$errors" -eq 0 ]; then
        echo -e "${GREEN}✓${NC} No errors in daemon logs (last 5 min)"
        return 0
    else
        echo -e "${YELLOW}⚠${NC} Found $errors error(s) in daemon logs (last 5 min)"
        journalctl -u kadeem-daemon --since "5 minutes ago" -p err -q --no-pager | tail -n 3
        return 1
    fi
}

echo "=== Kadeem Health Check ==="
echo ""

check_service "postgresql" || EXIT_CODE=1
check_service "kadeem-daemon" || EXIT_CODE=1
echo ""
check_database || EXIT_CODE=1
echo ""
check_daemon_logs || EXIT_CODE=1

echo ""
if [ $EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}All checks passed!${NC}"
else
    echo -e "${RED}Some checks failed!${NC}"
fi

exit $EXIT_CODE
