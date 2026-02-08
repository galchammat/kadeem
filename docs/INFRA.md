# Infrastructure

Two machines connected via Tailscale. Machine B (hostname `mac`) runs PostgreSQL, nginx, and the kadeem daemon. Machine A (hostname `pc3080`) stores daily database backups. GitHub Actions deploys on push to `main`.

## Tailscale Setup

Install on both machines and your dev machine:

```bash
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up
```

Verify connectivity:

```bash
tailscale status
ssh yog404@mac
ssh yog404@pc3080
```

## Ansible Setup

Ansible provisions everything on Machine B: PostgreSQL 16, nginx, the kadeem daemon systemd service, and the backup timer.

### Prerequisites

1. Install Ansible: `sudo apt install ansible`
2. SSH key auth configured: `ssh-copy-id yog404@mac`
3. Tailscale running on your machine and Machine B
4. `.env` file in project root (see `.env.template` for all vars). Key infra vars:

```bash
HOST_B=mac                          # Target host (Tailscale hostname)
USER_B=yog404                       # SSH user
BECOME_PASSWORD=your_sudo           # Sudo password for Ansible
DATABASE_URL=postgres://kadeem:pass@mac:5432/kadeem?sslmode=require
HOST_A=pc3080                       # Backup destination
USER_A=yog404
BACKUP_PATH=/mnt/d/kadeem-backups
BACKUP_RETENTION_DAYS=30
DISCORD_WEBHOOK_URL=https://...     # Optional
RIOT_API_KEY=RGAPI-...
SUPABASE_JWKS_URL=https://xxx.supabase.co/auth/v1/.well-known/jwks.json
```

Frontend environment variables (VITE_ prefixed) are in `packages/web/.env`:

```bash
VITE_API_URL=https://api.cyanlab.cc
VITE_SUPABASE_URL=https://xxx.supabase.co
VITE_SUPABASE_PUBLISHABLE_KEY=sb_publishable_...
DOMAIN=api.cyanlab.cc
FRONTEND_DOMAIN=cyanlab.cc
LETSENCRYPT_EMAIL=you@example.com
```

### Run

```bash
make ansible          # Full setup
make ansible check    # Dry run (no changes applied)
```

The playbook is idempotent -- safe to re-run.

### What It Provisions

- **PostgreSQL 16** -- install, config, SSL cert, create DB/user from `DATABASE_URL`
- **Nginx** -- reverse proxy for the API, Let's Encrypt SSL, rate limiting
- **Kadeem daemon** -- systemd service at `/opt/kadeem/bin/daemon`
- **Backup timer** -- daily at 2 AM, pg_dump to Machine A, 30-day retention, Discord notification

### Run Specific Roles

```bash
ansible-playbook -i ansible/inventory/production.yml ansible/playbook.yml --tags postgresql
ansible-playbook -i ansible/inventory/production.yml ansible/playbook.yml --tags nginx
ansible-playbook -i ansible/inventory/production.yml ansible/playbook.yml --tags server
```

### Directory Structure

```
ansible/
├── ansible.cfg
├── inventory/production.yml
├── playbook.yml
└── roles/
    ├── postgresql/          # DB install, config, SSL, backup
    │   ├── tasks/
    │   ├── templates/       # postgresql.conf.j2, pg_hba.conf.j2, backup-script.sh.j2, etc.
    │   ├── handlers/
    │   └── defaults/
     ├── nginx/               # Reverse proxy, Let's Encrypt
     │   ├── tasks/
     │   ├── templates/       # nginx.conf.j2, kadeem-api.conf.j2, etc.
    │   ├── handlers/
    │   └── defaults/
    └── server/              # Daemon user, dirs, systemd service
        ├── tasks/
        ├── templates/       # kadeem-daemon.service.j2
        ├── handlers/
        └── defaults/
```

## Deployment (GitHub Actions)

On push to `main` (when `packages/server/` changes):

1. Run tests
2. Build Linux/amd64 binary
3. Connect to Machine B via Tailscale SSH
4. Run database migrations
5. Deploy binary to `/opt/kadeem/bin/daemon`
6. Restart `kadeem-daemon` service
7. Health check -- auto-rollback on failure
8. Discord notification

Manual rollback: trigger the workflow with the `rollback` option.

### GitHub Secrets

Run `./scripts/set-secrets.sh` from a machine authenticated with `gh` CLI.

### Tailscale SSH for GitHub Actions (OIDC)

No SSH keys needed. GitHub Actions authenticates via Tailscale's federated identity.

**1. Enable Tailscale SSH on Machine B:**

```bash
sudo tailscale up --ssh
```

**2. Configure Tailscale ACLs** (Admin Console > Access Controls):

```json
{
  "tagOwners": {
    "tag:ci": ["autogroup:admin"]
  },
  "acls": [
    { "action": "accept", "src": ["tag:ci"], "dst": ["mac:*"] }
  ],
  "ssh": [
    { "action": "accept", "src": ["tag:ci"], "dst": ["mac"], "users": ["yog404"] }
  ]
}
```

**3. Create OAuth client** (Admin Console > Settings > OAuth Clients):

- Provider: GitHub Actions
- Audience: `https://github.com/YOUR_GITHUB_USERNAME`
- Tags: `tag:ci`
- Save the client ID

The `tailscale/github-action@v2` in the workflow handles the rest. Temporary access is granted per workflow run and auto-revoked.

## Backup & Restore

### Trigger a backup manually

```bash
ssh yog404@mac "sudo systemctl start kadeem-backup.service"
```

### Check backup status

```bash
ssh yog404@mac "sudo journalctl -u kadeem-backup.service -n 50"
ssh yog404@mac "sudo systemctl list-timers kadeem-backup.timer"
ssh yog404@pc3080 "ls -lh /mnt/d/kadeem-backups/"
```

### Restore from backup

```bash
# Stop daemon
ssh yog404@mac "sudo systemctl stop kadeem-daemon"

# List available backups
ssh yog404@pc3080 "ls -lh /mnt/d/kadeem-backups/"

# Copy backup to Machine B and restore
ssh yog404@mac "scp yog404@pc3080:/mnt/d/kadeem-backups/kadeem_backup_YYYYMMDD_HHMMSS.sql.gz /tmp/"
ssh yog404@mac "gunzip < /tmp/kadeem_backup_*.sql.gz | psql -U kadeem -d kadeem"

# Restart daemon
ssh yog404@mac "sudo systemctl start kadeem-daemon"
```

## Maintenance

**Update PostgreSQL config:** Edit `ansible/roles/postgresql/templates/postgresql.conf.j2`, run `make ansible`. PG restarts automatically if config changed.

**Change backup schedule:** Edit `ansible/roles/postgresql/defaults/main.yml` (`backup_schedule`, `backup_time`), run `make ansible`.

**Rotate database password:** Update `DATABASE_URL` in `.env`, run `make ansible`.

## Troubleshooting

### Ansible / SSH

```bash
ssh yog404@mac "echo ok"                  # Test SSH
ping mac                                   # Test Tailscale DNS
ansible-playbook ... -vvv                  # Verbose output
grep BECOME_PASSWORD .env                  # Verify sudo password is set
```

### PostgreSQL

```bash
ssh yog404@mac "sudo systemctl status postgresql"
ssh yog404@mac "sudo tail -f /var/log/postgresql/postgresql-*.log"
ssh yog404@mac "psql -U kadeem -d kadeem -c 'SELECT version();'"
```

### Daemon

```bash
ssh yog404@mac "sudo systemctl status kadeem-daemon"
ssh yog404@mac "sudo journalctl -u kadeem-daemon -f"
```

### Backups

```bash
ssh yog404@mac "sudo systemctl is-active kadeem-backup.timer"
ssh yog404@mac "sudo journalctl -u kadeem-backup.service --since '7 days ago'"
ssh yog404@pc3080 "ls -lh /mnt/d/kadeem-backups/"
```

### Tailscale

```bash
sudo systemctl status tailscaled
sudo systemctl restart tailscaled
tailscale status
tailscale ping mac
```
