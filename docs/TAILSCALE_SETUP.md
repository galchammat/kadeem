# Tailscale Setup for Kadeem

Tailscale provides secure, zero-config networking between your machines without port forwarding.

## Why Tailscale?

- **No Port Forwarding**: Works behind NAT/firewalls
- **Encrypted**: WireGuard-based encryption
- **Static Hostnames**: `machine-b.tail-scale.ts.net` instead of changing IPs
- **Free**: Up to 100 devices on free tier
- **Easy**: 5-minute setup

## Installation

### Machine B (Daemon Server)

```bash
# Install Tailscale
curl -fsSL https://tailscale.com/install.sh | sh

# Start and authenticate
sudo tailscale up

# Copy the URL and open in browser to authenticate

# Set hostname (optional but recommended)
sudo tailscale set --hostname kadeem-server

# Enable IP forwarding (if needed)
echo 'net.ipv4.ip_forward = 1' | sudo tee -a /etc/sysctl.conf
sudo sysctl -p
```

### Machine A (Backup Server)

```bash
# Same installation steps
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up
sudo tailscale set --hostname kadeem-backup
```

### Your Development Machine (optional)

```bash
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up
```

## Get Tailscale IP Addresses

```bash
# On each machine, get its Tailscale IP
tailscale ip -4

# Or list all devices
tailscale status
```

## Configure Kadeem for Tailscale

### Update Environment Variables

On Machine B, edit `/etc/kadeem/.env`:

```bash
# Use Tailscale hostname instead of localhost
DATABASE_URL=postgres://kadeem:password@100.x.x.x:5432/kadeem?sslmode=disable

# Or use Tailscale hostname (recommended)
# DATABASE_URL=postgres://kadeem:password@kadeem-server:5432/kadeem?sslmode=disable
```

### Update Backup Configuration

On Machine B, edit `/etc/kadeem/.env`:

```bash
# Set Machine A's Tailscale IP or hostname
BACKUP_HOST=kadeem-backup  # or 100.x.x.x
BACKUP_USER=your_username
BACKUP_PATH=/mnt/2tb/kadeem-backups
```

### Update GitHub Secrets

In GitHub Settings → Secrets → Actions:

```
DEPLOY_HOST=kadeem-server  # or Tailscale IP
```

## Allow PostgreSQL Connections

Edit `/etc/postgresql/*/main/pg_hba.conf`:

```conf
# Add Tailscale network (100.64.0.0/10)
host    kadeem    kadeem    100.64.0.0/10    md5
```

Restart PostgreSQL:
```bash
sudo systemctl restart postgresql
```

## Test Connection

### From Machine A to Machine B

```bash
# Test SSH
ssh user@kadeem-server

# Test PostgreSQL
psql postgres://kadeem:password@kadeem-server:5432/kadeem
```

### From GitHub Actions

The deploy workflow will automatically use the Tailscale hostname from `DEPLOY_HOST` secret.

## Advantages

1. **No Dynamic DNS needed**: Tailscale hostnames are stable
2. **No firewall rules**: Works through NAT
3. **Encrypted**: All traffic is encrypted end-to-end
4. **Fast**: Peer-to-peer when possible
5. **Access from anywhere**: Connect to your network from phone/laptop

## Optional: Subnet Routing

If you want to access other devices on Machine B's network:

```bash
# On Machine B, enable subnet routing
sudo tailscale up --advertise-routes=192.168.1.0/24

# In Tailscale admin console, approve the subnet route
```

## Monitoring

Check Tailscale status:
```bash
# Show connected peers
tailscale status

# Show your IP
tailscale ip

# Check connectivity to specific peer
tailscale ping kadeem-backup
```

## Troubleshooting

### Connection issues

```bash
# Check Tailscale status
sudo systemctl status tailscaled

# Restart Tailscale
sudo systemctl restart tailscaled

# Check logs
sudo journalctl -u tailscaled -n 50
```

### Can't reach other machine

```bash
# Verify both machines are online
tailscale status

# Ping the other machine
tailscale ping machine-name

# Check firewall isn't blocking
sudo ufw status  # If using ufw
```

## Security Notes

- Tailscale uses your Google/Microsoft/GitHub account for auth
- Each device gets a unique key
- You can revoke access from the admin console
- Enable MFA on your Tailscale account
