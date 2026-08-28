# Inventory Management — Project Overview

Full-stack inventory management system: React 19 frontend (Vite + Tailwind CSS v4 + shadcn/ui) and Go backend (Gin + GORM + PostgreSQL). The Go backend serves the built frontend as static files in production, making it a single deployable binary.

## Stack

### Frontend (`frontend/`)
- React 19, Vite, Tailwind CSS v4, shadcn/ui (base-ui)
- TanStack React Query for server state, Axios for HTTP
- React Router v8 (declarative `<BrowserRouter>` + `<Routes>`), lazy-loaded pages
- Vitest + React Testing Library, oxlint

### Backend (`backend/`)
- Go 1.24+, Gin, GORM (PostgreSQL driver), golang-jwt/v5
- Schema migrations via Sqitch (versioned SQL in `backend/migrations/`), applied by CI — not by the app
- Soft deletes via `gorm.DeletedAt` (requires `deleted_at` column from migrations)
- Config loaded from environment / `.env` via `godotenv`
- Tests: Go `testing` package + testify (planned)

## Directory Layout

```
inventory-management/
├── docs/                      # Root project overview + domain index
│   └── overview.md            # This file
├── frontend/                  # React application
│   ├── src/
│   │   ├── components/        # layout/ + ui/ (shadcn) components
│   │   ├── features/          # Per-domain feature dirs (own overview.md inline)
│   │   │   └── products/      # Products feature (overview.md inside)
│   │   ├── pages/             # Route pages (lazy-loaded)
│   │   ├── lib/               # api.js, queryClient.js, utils.js
│   │   ├── hooks/             # Custom React hooks (future)
│   │   └── test/              # Test setup
├── backend/
│   ├── cmd/server/main.go     # Entry point: config, DB connect, routes, static files
│   ├── migrations/            # Sqitch project: sqitch.conf, sqitch.plan, deploy/revert/verify/ + overview.md
│   └── internal/
│       ├── products/          # Products domain: model + handler + tests + overview.md
│       ├── auth/              # JWT issue/validate (not yet wired into routes)
│       ├── config/            # Env-based configuration
│       └── database/          # GORM connection + DSN
├── Dockerfile                 # Multi-stage: frontend build -> Go build -> alpine
├── docker-compose.yml         # Development database (postgres) — migrations run via CI, not compose
├── .github/workflows/         # CI: validates + deploys Sqitch migrations
└── AGENTS.md                  # Agent instructions (project conventions)
```

## Architecture

### Schema management (Sqitch)
Schema is defined entirely in versioned SQL under `backend/migrations/` and applied by CI before app rollout (`.github/workflows/deploy.yml`). The app never runs DDL — `database.Connect` only opens a connection. See the [Migrations domain overview](../backend/migrations/overview.md) for the full workflow.

### Backend flow
`config.Load()` → `database.Connect()` (schema already applied by Sqitch) → Gin router with CORS middleware → routes under `/api/*` → `products.ProductHandler` (direct `gorm.DB` access, no service/repository layer yet) → PostgreSQL.

### Serving
- Dev: frontend on :5173 (Vite proxy for `/api`), backend on :8080.
- Prod: backend serves `frontend/dist` from `STATIC_DIR`; SPA fallback to `index.html`; `/api/*` 404s return JSON.
- CORS: permissive (`Access-Control-Allow-Origin: *`) in `main.go`.

### Request lifecycle (products)
`GET /api/products` → handler `List` → `db.Find(&products)` → JSON array. Errors map to HTTP 4xx/5xx with `{"error": "..."}` shape.

## Domains

Each domain directory carries its own inline `overview.md` next to the code it documents. This table is the index.

| Domain | Overview(s) | Status | Last updated |
| --- | --- | --- | --- |
| Migrations | `backend/migrations/overview.md` | Sqitch active (products migration shipped) | 2026-08-28 |
| Products (backend) | `backend/internal/products/overview.md` | CRUD API live | 2026-08-28 |
| Products (frontend) | `frontend/src/features/products/overview.md` | no UI yet | 2026-08-28 |

## Development Commands

```bash
mise install                        # Install toolchain from .tool-versions
pnpm dev                            # Root: frontend + backend concurrently
cd frontend && pnpm dev|build|test|lint
cd backend && go run cmd/server/main.go
cd backend && go test ./... && go build ./...

# Database (dev)
docker compose up -d db             # Start postgres on :5432 (see docker-compose.yml)
cd backend/migrations
docker run --rm -v "$(pwd):/migrations" -w /migrations \
  sqitch/sqitch deploy db:pg://postgres:postgres@localhost:5432/inventory
docker run --rm -v "$(pwd):/migrations" -w /migrations \
  sqitch/sqitch verify db:pg://postgres:postgres@localhost:5432/inventory
```

## Agent Documentation Protocol

This file is the single entry point to project documentation. Follow it on every session.

1. **Read before implementing.** Always read `docs/overview.md` and the relevant domain `overview.md`(s) before writing or changing code. Overviews live inline inside each domain directory. If a domain has no overview yet, note it and create it (step 3).
2. **Scope the domain.** Each feature/domain owns its `overview.md` inside its own directory. Never merge two domains into one file, and keep overviews next to the code they describe.
3. **New feature/domain → create.** When implementing a new feature or domain, create its directory and a `overview.md` inside it after implementation, using the domain template (see AGENTS.md), then add a row to the domain index table above.
4. **Existing feature → update.** When modifying an existing domain, update its inline `overview.md` in the same change (endpoints, model fields, key files, decisions, tests).
5. **Keep the index current.** Update the domain table and `Last updated` dates whenever docs change.
6. **Document reality.** Overviews describe the actual code, not intentions. If code diverges from a documented plan, fix the doc.
7. **No drift.** A feature change is not complete until its documentation is updated.
