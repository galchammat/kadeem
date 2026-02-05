#!/usr/bin/env bash
set -euo pipefail

# Script to sync environment variables from .env to GitHub repository secrets
# Requires: gh CLI (GitHub CLI) to be installed and authenticated

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="$PROJECT_ROOT/.env"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# List of environment variables to sync to GitHub secrets
SECRETS=(
    "TAILSCALE_CLIENT_ID"
    "TAILSCALE_AUDIENCE"
    "DATABASE_URL"
    "RIOT_API_KEY"
    "TWITCH_CLIENT_ID"
    "TWITCH_CLIENT_SECRET"
    "DISCORD_WEBHOOK_URL"
    "BACKUP_PATH"
    "BACKUP_RETENTION_DAYS"
)

echo -e "${GREEN}GitHub Secrets Sync Tool${NC}"
echo "================================"
echo

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo -e "${RED}Error: GitHub CLI (gh) is not installed${NC}"
    echo "Install it from: https://cli.github.com/"
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo -e "${RED}Error: Not authenticated with GitHub CLI${NC}"
    echo "Run: gh auth login"
    exit 1
fi

# Check if .env file exists
if [[ ! -f "$ENV_FILE" ]]; then
    echo -e "${RED}Error: .env file not found at $ENV_FILE${NC}"
    exit 1
fi

# Load .env file
echo "Loading environment variables from .env..."
set -a
source "$ENV_FILE"
set +a

echo -e "${GREEN}✓${NC} Loaded .env file"
echo

# Sync secrets to GitHub
echo "Syncing secrets to GitHub repository..."
echo

SUCCESS_COUNT=0
SKIP_COUNT=0
ERROR_COUNT=0

for secret_name in "${SECRETS[@]}"; do
    # Get the value from environment
    secret_value="${!secret_name:-}"
    
    if [[ -z "$secret_value" ]]; then
        echo -e "${YELLOW}⊘${NC} Skipping $secret_name (not set in .env)"
        ((SKIP_COUNT++))
        continue
    fi
    
    # Set the secret using gh CLI
    if echo "$secret_value" | gh secret set "$secret_name" 2>/dev/null; then
        echo -e "${GREEN}✓${NC} Set $secret_name"
        ((SUCCESS_COUNT++))
    else
        echo -e "${RED}✗${NC} Failed to set $secret_name"
        ((ERROR_COUNT++))
    fi
done

echo
echo "================================"
echo -e "Summary:"
echo -e "  ${GREEN}Success:${NC} $SUCCESS_COUNT secrets"
echo -e "  ${YELLOW}Skipped:${NC} $SKIP_COUNT secrets"
if [[ $ERROR_COUNT -gt 0 ]]; then
    echo -e "  ${RED}Failed:${NC} $ERROR_COUNT secrets"
fi
echo

if [[ $ERROR_COUNT -eq 0 ]]; then
    echo -e "${GREEN}✓ All secrets synced successfully!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some secrets failed to sync${NC}"
    exit 1
fi
