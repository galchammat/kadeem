# Kadeem Ansible

Infrastructure playbook for Machine B (`mac`): PostgreSQL, nginx, daemon service, backups, and GitHub self-hosted runner.

## Usage

```bash
make ansible
make ansible check
```

## Prerequisites

- Ansible installed
- SSH access to target host
- `ansible/.env` configured (see `.env.template`)
- `GITHUB_RUNNER_REGISTRATION_TOKEN` set when first installing runner

## Roles

- `postgresql`
- `nginx`
- `server`
- `github_runner`

## Tags

```bash
make ansible ARGS="--tags postgresql"
make ansible ARGS="--tags nginx"
make ansible ARGS="--tags server"
make ansible ARGS="--tags runner"
```

## Notes

- The playbook supports dry-run mode (`--check`).
- Runner install/config tasks are skipped in check mode.
- Runner is configured at repo-level with labels `self-hosted,linux,x64,machine-b,deploy`.
