Monorepo for Kadeem, a League of Legends match tracking application. Go API server + React/TypeScript frontend, deployed on Pop!_OS behind Cloudflare Tunnel and nginx.

## Stack

- **Backend:** Go 1.23, chi router, PostgreSQL
- **Frontend:** React 18, Vite 7, Tailwind CSS 4, shadcn/ui
- **Infra:** Ansible, systemd, nginx, Cloudflare Tunnel

## Repository Structure

```
packages/server/          Go API server
  cmd/daemon/main.go      Entry point
  cmd/migrate/main.go     Database migrations runner
  internal/               All Go packages (api, handler, service, store, model, riot)
  migrations/             SQL migration files
  go.mod                  Module: github.com/galchammat/kadeem

packages/web/             React frontend
  src/
    lib/api.ts            API client (fetch-based)
    lib/matchTransformer.ts  Match data transformation + DataDragon CDN
    types/index.ts         TypeScript types matching Go models
    hooks/                 React hooks (useLolAccounts, useLolMatches, useStreamer)
    contexts/              React contexts (streamerContext)
    pages/                 Page components
    components/            UI components (shadcn/ui based)

ansible/                  Deployment automation (roles: postgresql, nginx, server)
scripts/                  Operational scripts
```

## Development

- `make run` — Starts Go server + Vite dev server
- `make build` — Builds Go binary + frontend bundle
- `make test` — Runs Go tests
- `make migrate-up` / `make migrate-down` — Database migrations

## Guidelines

1. Follow Go idioms. Use `any` not `interface{}`.
2. Frontend types in `packages/web/src/types/index.ts` mirror Go models. Keep them in sync.
3. API client in `packages/web/src/lib/api.ts` matches routes in `packages/server/internal/api/routes.go`.
4. Use shadcn/ui components. Prefer editing existing files over creating new ones.
5. Minimize code. YAGNI.
