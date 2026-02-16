# Infrastructure

Machine B (`mac`) runs PostgreSQL, nginx, and the `kadeem-daemon` service. Machine A (`pc3080`) stores daily database backups. CI runs on GitHub-hosted runners; CD runs on a repo-level self-hosted runner installed on Machine B.

## Provisioning

Ansible provisions:

- PostgreSQL
- nginx
- daemon systemd service
- backup timer/service
- GitHub self-hosted runner service

Run:

```bash
make ansible
make ansible check
```

`--check` is supported for dry runs; runner install/config operations are skipped in check mode.

## Required Environment Variables

Set these variables before running Ansible:

```bash
HOST_B=mac
USER_B=yog404
BECOME_PASSWORD=...

DATABASE_URL=postgres://kadeem:pass@mac:5432/kadeem?sslmode=require
RIOT_API_KEY=RGAPI-...
SUPABASE_JWKS_URL=https://xxx.supabase.co/auth/v1/.well-known/jwks.json

HOST_A=pc3080
USER_A=yog404
BACKUP_PATH=/mnt/d/kadeem-backups
BACKUP_RETENTION_DAYS=30

DISCORD_WEBHOOK_URL=https://...
LETSENCRYPT_EMAIL=you@example.com

# One-time token for runner registration during ansible apply
GITHUB_RUNNER_REGISTRATION_TOKEN=...
```

## CI/CD Workflows

- `/.github/workflows/ci.yml`
  - PR + `main` push checks
  - component-aware test/lint/build
  - publishes `daemon-binary` artifact on `push` to `main`
  - `CI Gate` is the required status check

- `/.github/workflows/cd.yml`
  - triggered by successful `CI` workflow (`workflow_run`) on `main`
  - runs on self-hosted runner labels: `self-hosted,linux,x64,machine-b,deploy`
  - downloads `daemon-binary` from the CI run and deploys locally
  - supports manual rollback via `workflow_dispatch`

## Deployment Flow

1. CI passes on `main`
2. CD downloads daemon artifact
3. Backup current daemon binary
4. Replace daemon binary and run migrations
5. Restart daemon and run health check
6. Roll back on failure

## Self-Hosted Runner Notes

- Runner is installed as `github-runner` user via Ansible.
- Runner uses outbound polling to GitHub (no inbound static IP required).
- Deploy commands are allowed via least-privilege sudoers entry.

## Useful Checks

```bash
# Runner status
ssh yog404@mac "sudo systemctl status github-runner"

# Daemon status
ssh yog404@mac "sudo systemctl status kadeem-daemon"
ssh yog404@mac "sudo journalctl -u kadeem-daemon -n 100 --no-pager"

# Backup timer
ssh yog404@mac "sudo systemctl list-timers kadeem-backup.timer"
```
