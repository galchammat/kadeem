#!/usr/bin/env bash
set -euo pipefail

# Sync package server .env to GitHub repository secrets.
# Requires: gh CLI installed and authenticated.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
SERVER_ENV_FILE="$PROJECT_ROOT/packages/server/.env"

echo "GitHub secrets sync (packages/server/.env)"
echo "================================"
echo

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed"
    echo "Install it from: https://cli.github.com/"
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo "Error: Not authenticated with GitHub CLI"
    echo "Run: gh auth login"
    exit 1
fi

# Check if .env file exists
if [[ ! -f "$SERVER_ENV_FILE" ]]; then
    echo "Error: .env file not found at $SERVER_ENV_FILE"
    exit 1
fi

echo "Syncing all keys from $SERVER_ENV_FILE"
gh secret set -f "$SERVER_ENV_FILE"

echo "Done."
echo
