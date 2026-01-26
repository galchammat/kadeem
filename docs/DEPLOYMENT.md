# Kadeem Deployment Guide

Complete guide for deploying Kadeem with PostgreSQL, automated backups, and CI/CD.

## Architecture Overview

```
┌─────────────────┐         ┌──────────────────┐         ┌─────────────┐
│  GitHub Actions │────────>│   Machine B      │<───────>│  Machine A  │
│  (CI/CD)        │  Deploy │   (Laptop/Server)│  Backup │  (2TB HDD)  │
└─────────────────┘         └──────────────────┘         └─────────────┘
                                     │
                                     ├─ PostgreSQL DB
                                     ├─ Kadeem Daemon
                                     └─ Systemd Services
```

## Prerequisites

- **Machine B** (Server): Linux machine with systemd
- **Machine A** (Backup): Linux machine with 2TB storage
- **GitHub Account**: For CI/CD
- **Discord Webhook**: For notifications (optional)
- **Riot API Key**: From Riot Developer Portal

---

## Part 1: Initial Setup on Machine B

### 1.1 Clone Repository

```bash
git clone https://github.com/yourusername/kadeem.git
cd kadeem
```

### 1.2 Install Tailscale

For secure networking without port forwarding:

```bash
# Install Tailscale
curl -fsSL https://tailscale.com/install.sh | sh

# Start and authenticate
sudo tailscale up
sudo tailscale set --hostname kadeem-server

# Get your Tailscale IP
tailscale ip -4
```

See [TAILSCALE_SETUP.md](./TAILSCALE_SETUP.md) for detailed instructions.

### 1.3 Setup PostgreSQL

Run the zero-touch setup script:

```bash
sudo make db-setup
```

This will:
- Install PostgreSQL 16
- Create `kadeem` database and user
- Generate secure random password
- Configure network access
- Create `/etc/kadeem/.env` with credentials

**Important**: Save the password shown in the output!

### 1.4 Configure Environment

Edit `/etc/kadeem/.env`:

```bash
sudo nano /etc/kadeem/.env
```

Add your configuration:

```env
DATABASE_URL=postgres://kadeem:YOUR_PASSWORD@localhost:5432/kadeem?sslmode=disable
RIOT_API_KEY=RGAPI-your-key-here
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
BACKUP_HOST=kadeem-backup  # Machine A's Tailscale hostname
BACKUP_USER=your_username
BACKUP_PATH=/mnt/2tb/kadeem-backups
```

### 1.5 Run Migrations

```bash
make migrate-up
```

### 1.6 Build and Deploy Daemon

```bash
make build-daemon
sudo make deploy-daemon
```

### 1.7 Install Systemd Services

```bash
make install-services
sudo systemctl enable kadeem-daemon
sudo systemctl enable kadeem-backup.timer
sudo systemctl start kadeem-daemon
sudo systemctl start kadeem-backup.timer
```

### 1.8 Verify Installation

```bash
# Check daemon status
make daemon-status

# Run health checks
make healthcheck

# View logs
make daemon-logs
```

---

## Part 2: Setup Machine A (Backup Server)

### 2.1 Install Tailscale

```bash
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up
sudo tailscale set --hostname kadeem-backup
```

### 2.2 Create Backup Directory

```bash
# Mount your 2TB drive (if not already mounted)
sudo mkdir -p /mnt/2tb
# Add to /etc/fstab for automatic mounting

# Create backup directory
mkdir -p /mnt/2tb/kadeem-backups
```

### 2.3 Setup SSH Access

On Machine B, generate SSH key:

```bash
ssh-keygen -t ed25519 -f ~/.ssh/kadeem_backup -N ""
```

Copy public key to Machine A:

```bash
ssh-copy-id -i ~/.ssh/kadeem_backup.pub user@kadeem-backup
```

Test connection:

```bash
ssh -i ~/.ssh/kadeem_backup user@kadeem-backup
```

### 2.4 Test Backup

On Machine B:

```bash
sudo make backup
```

Verify backup on Machine A:

```bash
ls -lh /mnt/2tb/kadeem-backups/
```

---

## Part 3: GitHub Actions CI/CD Setup

### 3.1 Generate Deploy SSH Key

On Machine B:

```bash
ssh-keygen -t ed25519 -f ~/.ssh/github_deploy -N ""
cat ~/.ssh/github_deploy.pub >> ~/.ssh/authorized_keys
```

Copy the private key:

```bash
cat ~/.ssh/github_deploy
```

### 3.2 Configure GitHub Secrets

Go to your repository: **Settings** → **Secrets and variables** → **Actions** → **New repository secret**

Add these secrets:

| Secret Name | Value | Example |
|------------|-------|---------|
| `DEPLOY_HOST` | Machine B's Tailscale hostname or IP | `kadeem-server` or `100.x.x.x` |
| `DEPLOY_USER` | SSH username on Machine B | `your_username` |
| `DEPLOY_SSH_KEY` | Private key content from above | `-----BEGIN OPENSSH PRIVATE KEY-----...` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://kadeem:pass@localhost:5432/kadeem` |
| `RIOT_API_KEY` | Your Riot API key | `RGAPI-...` |
| `DISCORD_WEBHOOK_URL` | Discord webhook for notifications | `https://discord.com/api/webhooks/...` |

### 3.3 Setup Repository on Machine B

Create a permanent repository location for migrations:

```bash
sudo mkdir -p /opt/kadeem/repo
sudo chown $USER:$USER /opt/kadeem/repo
cd /opt/kadeem/repo
git clone https://github.com/yourusername/kadeem.git .
```

### 3.4 Test Deployment

Push to `main` branch to trigger deployment:

```bash
git push origin main
```

Watch the workflow at: `https://github.com/yourusername/kadeem/actions`

---

## Part 4: Discord Notifications

### 4.1 Create Discord Webhook

1. Open Discord server settings
2. Go to **Integrations** → **Webhooks**
3. Click **New Webhook**
4. Name it "Kadeem Notifications"
5. Select channel
6. Copy webhook URL

### 4.2 Add to Environment

On Machine B, edit `/etc/kadeem/.env`:

```bash
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/123456/abcdef...
```

### 4.3 Test Notification

```bash
sudo /opt/kadeem/scripts/notify-discord.sh
```

You should receive a Discord message!

### 4.4 Notifications Configured For

- ✅ Successful deployments
- ❌ Failed deployments
- ⚠️  Daemon service failures
- 💾 Database backup completions

---

## Part 5: Daily Operations

### View Daemon Logs

```bash
make daemon-logs
```

### Manual Backup

```bash
sudo make backup
```

### Restart Daemon

```bash
make daemon-restart
```

### Run Health Checks

```bash
make healthcheck
```

### Manual Rollback

If deployment goes wrong:

1. Go to GitHub Actions
2. Click "Run workflow"
3. Select "Deploy Daemon to Machine B"
4. Check "Rollback to previous version"
5. Click "Run workflow"

---

## Part 6: Monitoring

### Check Backup Timer

```bash
sudo systemctl status kadeem-backup.timer
sudo systemctl list-timers kadeem-backup.timer
```

### Database Size

```bash
sudo -u postgres psql -d kadeem -c "SELECT pg_size_pretty(pg_database_size('kadeem'));"
```

### Recent Backups

On Machine A:

```bash
ls -lh /mnt/2tb/kadeem-backups/ | tail -10
```

### Disk Usage

Machine B:
```bash
df -h
```

Machine A:
```bash
df -h /mnt/2tb
```

---

## Part 7: Troubleshooting

### Daemon Won't Start

```bash
# Check logs
sudo journalctl -u kadeem-daemon -n 50

# Check environment
cat /etc/kadeem/.env

# Test database connection
psql $(grep DATABASE_URL /etc/kadeem/.env | cut -d= -f2-)
```

### Backup Fails

```bash
# Check SSH connection
ssh $BACKUP_USER@$BACKUP_HOST

# Check disk space on Machine A
ssh $BACKUP_USER@$BACKUP_HOST "df -h /mnt/2tb"

# Run backup manually with verbose output
sudo bash -x /opt/kadeem/scripts/backup-database.sh
```

### GitHub Actions Fails

1. Check secrets are set correctly
2. Verify SSH key is valid
3. Test SSH connection manually:
   ```bash
   ssh -i ~/.ssh/github_deploy $DEPLOY_USER@$DEPLOY_HOST
   ```

### PostgreSQL Connection Issues

```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Check pg_hba.conf
sudo cat /etc/postgresql/*/main/pg_hba.conf | grep kadeem

# Test connection
psql -h localhost -U kadeem -d kadeem
```

---

## Part 8: Maintenance

### Update Riot API Key

```bash
sudo nano /etc/kadeem/.env
# Update RIOT_API_KEY
sudo systemctl restart kadeem-daemon
```

### Restore from Backup

On Machine B:

```bash
# Stop daemon
sudo systemctl stop kadeem-daemon

# Get backup from Machine A
scp $BACKUP_USER@kadeem-backup:/mnt/2tb/kadeem-backups/kadeem_backup_YYYYMMDD_HHMMSS.sql.gz /tmp/

# Restore
gunzip < /tmp/kadeem_backup_*.sql.gz | psql $(grep DATABASE_URL /etc/kadeem/.env | cut -d= -f2- | sed 's/?.*$//')

# Start daemon
sudo systemctl start kadeem-daemon
```

### Clean Old Backups

The backup script automatically keeps last 30 days. To change retention:

```bash
sudo nano /etc/kadeem/.env
# Add: RETENTION_DAYS=60
```

---

## Part 9: AWS Migration Path (Future)

When ready to migrate to AWS:

### Option 1: RDS PostgreSQL

1. Create RDS instance
2. Update `DATABASE_URL` in GitHub secrets
3. Run migrations: `make migrate-up`
4. Deploy daemon to EC2 instance

### Option 2: Keep Local with Remote Access

1. Keep Machine B as-is
2. Use Tailscale for secure access
3. Access from anywhere via Tailscale VPN

See architecture decision in main README - current setup is AWS-compatible!

---

## Quick Reference

### Makefile Commands

```bash
make help              # Show all commands
make daemon-status     # Check daemon
make daemon-logs       # View live logs
make healthcheck       # Run health checks
make backup            # Manual backup
make migrate-up        # Run migrations
make deploy-daemon     # Deploy daemon
```

### Systemd Services

```bash
# Daemon
sudo systemctl status kadeem-daemon
sudo systemctl restart kadeem-daemon

# Backup
sudo systemctl status kadeem-backup.timer
sudo systemctl list-timers

# PostgreSQL
sudo systemctl status postgresql
```

### Logs

```bash
# Daemon
sudo journalctl -u kadeem-daemon -f

# Backup
sudo journalctl -u kadeem-backup -f

# PostgreSQL
sudo journalctl -u postgresql -f
```

---

## Support

- **GitHub Issues**: https://github.com/yourusername/kadeem/issues
- **Discord**: (your server link)
- **Docs**: See `docs/` directory

---

## Security Checklist

- [ ] PostgreSQL password is strong and unique
- [ ] `/etc/kadeem/.env` has correct permissions (600)
- [ ] SSH keys are secured
- [ ] GitHub secrets are set
- [ ] Tailscale MFA enabled
- [ ] Backups are encrypted (optional: add GPG encryption)
- [ ] Firewall rules are configured
- [ ] PostgreSQL only listens on Tailscale network

---

**Congratulations!** 🎉 Your Kadeem deployment is complete with:

- ✅ Zero-touch PostgreSQL setup
- ✅ Automated daemon deployment
- ✅ Daily backups to Machine A
- ✅ CI/CD via GitHub Actions
- ✅ Discord notifications
- ✅ Secure networking with Tailscale
- ✅ One-command operations via Makefile
