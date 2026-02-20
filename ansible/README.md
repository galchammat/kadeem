# Kadeem Ansible

Infrastructure playbook for Machine B (`mac`): PostgreSQL, nginx, daemon service, backups, and GitHub self-hosted runner.

## Usage

```bash
# first time only
cp ansible/group_vars/db_server/vault.yml.example ansible/group_vars/db_server/vault.yml

make ansible
make ansible check
```

## Prerequisites

- Ansible installed
- SSH access to target host
- `ansible/group_vars/db_server/main.yml` configured
- `ansible/group_vars/db_server/vault.yml` populated
- `cloudflare_api_token` set in `ansible/group_vars/db_server/vault.yml`

## Roles

- `postgresql`
- `nginx`
- `server`
- `github_runner`

## Tags

```bash
ansible-playbook -i ansible/inventory/production.yml ansible/playbook.yml --tags postgresql
ansible-playbook -i ansible/inventory/production.yml ansible/playbook.yml --tags nginx
ansible-playbook -i ansible/inventory/production.yml ansible/playbook.yml --tags server
ansible-playbook -i ansible/inventory/production.yml ansible/playbook.yml --tags runner
```

## Notes

- The playbook supports dry-run mode (`--check`).
- Runner install/config tasks are skipped in check mode.
- Runner is configured at repo-level with labels `self-hosted,linux,x64,machine-b,deploy`.
- Keep `ansible/group_vars/db_server/vault.yml` local and out of git.
- TLS uses Let's Encrypt DNS-01 via Cloudflare API token.
