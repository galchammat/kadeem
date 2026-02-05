#!/bin/bash

# Discord notification script for Kadeem daemon monitoring
# Can be called manually or via systemd OnFailure

# Load environment
if [ -f "/etc/kadeem/.env" ]; then
    source /etc/kadeem/.env
fi

if [ -z "$DISCORD_WEBHOOK_URL" ]; then
    echo "DISCORD_WEBHOOK_URL not set"
    exit 0
fi

# Get service status
SERVICE_STATUS=$(systemctl is-active kadeem-daemon)
SERVICE_FAILED=$(systemctl is-failed kadeem-daemon 2>/dev/null || echo "unknown")

# Get recent logs
RECENT_LOGS=$(journalctl -u kadeem-daemon -n 10 --no-pager 2>/dev/null | tail -5)

# Determine color based on status
if [ "$SERVICE_STATUS" = "active" ]; then
    COLOR=3066993  # Green
    TITLE="✅ Kadeem Daemon - Healthy"
elif [ "$SERVICE_STATUS" = "inactive" ]; then
    COLOR=16776960  # Yellow
    TITLE="⚠️  Kadeem Daemon - Stopped"
else
    COLOR=15158332  # Red
    TITLE="❌ Kadeem Daemon - Failed"
fi

# Escape logs for JSON
LOGS_ESCAPED=$(echo "$RECENT_LOGS" | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')

# Send to Discord
curl -X POST "$DISCORD_WEBHOOK_URL" \
    -H "Content-Type: application/json" \
    -d "{
        \"embeds\": [{
            \"title\": \"$TITLE\",
            \"description\": \"Status change detected on Machine B\",
            \"color\": $COLOR,
            \"fields\": [
                {\"name\": \"Status\", \"value\": \"\`$SERVICE_STATUS\`\", \"inline\": true},
                {\"name\": \"Host\", \"value\": \"\`$(hostname)\`\", \"inline\": true},
                {\"name\": \"Recent Logs\", \"value\": \"\`\`\`\\n$LOGS_ESCAPED\\n\`\`\`\", \"inline\": false}
            ],
            \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
        }]
    }" 2>&1 | logger -t kadeem-discord
