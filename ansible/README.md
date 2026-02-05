# Kadeem Ansible

Infrastructure-as-code for Kadeem project.

## Usage

```bash
# Run full setup
make ansible

# Dry run (check what would change)
make ansible check
```

## Structure

```
ansible/
├── playbook.yml           # Main playbook (currently: PostgreSQL + backups)
├── inventory/
│   └── production.yml     # Target hosts (reads from .env)
└── roles/
    ├── postgresql/        # PostgreSQL installation & configuration
    └── postgres_backup/   # Automated backup system
```

## Prerequisites

- Ansible installed: `pipx install ansible`
- Tailscale running on both machines
- SSH key auth configured to target host
- `.env` file with required variables (see below)

## Environment Variables

Set in project root `.env`:

```bash
HOST_B=mac                              # PostgreSQL host
USER_B=yog404                           # SSH user
DATABASE_URL=postgres://user:pass@...  # Full connection string
HOST_A=pc3080                           # Backup destination
USER_A=yog404                           # Backup SSH user
BACKUP_PATH=/mnt/d/kadeem-backups      # Backup directory
BACKUP_RETENTION_DAYS=30               # Backup retention
DISCORD_WEBHOOK_URL=https://...        # Optional notifications
```

## What It Does

1. Installs PostgreSQL latest stable
2. Creates database and user from `DATABASE_URL`
3. Configures for production (optimized for small-scale apps)
4. Generates self-signed SSL certificate
5. Sets up daily backups at 2:00 AM
6. Pushes backups to remote host via Tailscale

## Tags

Run specific parts:

```bash
# Only PostgreSQL setup
make ansible ARGS="--tags postgresql"

# Only backup setup
make ansible ARGS="--tags backup"
```

## Adding New Roles

Create a new role in `roles/` and add it to `playbook.yml`:

```yaml
roles:
  - role: postgresql
  - role: postgres_backup
  - role: your_new_role  # Add here
```
