# AGENTS.md - Project Documentation

## Overview

Full-stack inventory management system with React frontend and Go backend.

## Documentation Protocol (MANDATORY)

Every feature/domain has its own `overview.md` placed **inline inside its own directory**, next to the code it documents. The root index is `docs/overview.md`. Documentation is not optional — it is part of the definition of done.

### Before implementing (always)

1. Read `docs/overview.md` first — it is the single entry point: architecture, domain index, and this protocol.
2. Read the relevant inline `overview.md` for any domain the change touches (e.g. `backend/internal/products/overview.md`, `frontend/src/features/products/overview.md`).
3. If a domain overview does not exist yet, read `docs/overview.md` for context and create the domain directory + overview during implementation.

### After implementing (always)

4. **New feature/domain** → create its directory with an inline `overview.md` using the Domain Overview Template below, then add a row to the domain index table in `docs/overview.md`.
5. **Changes to an existing domain** → update its inline `overview.md` in the same commit/change: endpoints, model fields, key files, decisions, tests.
6. Update the `Last updated` date and status in the overview header and the index table.

### Rules

- **One domain per file.** Never merge domains into a single overview.
- **Co-location.** Each `overview.md` sits inside the directory it documents (backend feature packages, frontend `features/` dirs).
- **Document reality, not intention.** Overviews describe the code as it exists. If code diverges from a documented plan, fix the doc.
- **No drift.** A feature change is not complete until its documentation is updated.
- Root-level project docs live in `docs/overview.md`, not in AGENTS.md. AGENTS.md stays focused on conventions.

### Domain Overview Template

```markdown
# <Domain> — Domain Overview

> Status: <phase> · last updated YYYY-MM-DD
> Back to [project overview](../overview.md)  <!-- relative path varies; point it at docs/overview.md -->

## Purpose

<What this domain does, in 1-3 sentences.>

## Data Model

<Entities/fields table: name, type, constraints, JSON keys. Omit if stateless.>

## API Endpoints

| Method | Path | Handler | Success | Errors |
| --- | --- | --- | --- | --- |

## Key Files

| File | Responsibility |
| --- | --- |

## Dependencies

<Libraries, services, and cross-domain couplings.>

## Decisions & Conventions

<Non-obvious choices, patterns, and deviations from AGENTS.md.>

## Tests

<Test files and what they cover. Note gaps explicitly.>
```

## Tech Stack

### Frontend
- **Framework:** React 19 with Vite
- **Styling:** Tailwind CSS v4 + shadcn/ui
- **State Management:** TanStack React Query
- **Routing:** React Router v8 (declarative mode)
- **HTTP Client:** Axios
- **Testing:** Vitest + React Testing Library
- **Package Manager:** pnpm

### Backend
- **Language:** Go 1.24+
- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL
- **Authentication:** JWT (golang-jwt)
- **Testing:** Go testing package + testify

## Project Structure

```
inventory-management/
├── docs/                # Living documentation
│   └── overview.md      # Project overview + domain index + protocol
│
├── frontend/              # React application
│   ├── src/
│   │   ├── features/      # One directory per domain, each with its own
│   │   │   └── products/  #   overview.md inline next to the code
│   │   ├── components/   # Reusable UI components
│   │   │   ├── layout/   # Layout components (MainLayout)
│   │   │   └── ui/       # shadcn/ui components
│   │   ├── pages/        # Route page components
│   │   ├── lib/          # Utilities (api.js, queryClient.js)
│   │   ├── hooks/        # Custom React hooks
│   │   └── test/         # Test setup
│   ├── components.json   # shadcn/ui configuration
│   └── vite.config.js
│
├── backend/              # Go application
│   ├── cmd/server/       # Main entry point
│   ├── migrations/       # Sqitch schema migrations (sqitch.conf, plan, deploy/revert/verify, overview.md)
│   └── internal/
│       ├── products/     # Products domain: model, handlers, tests, overview.md
│       ├── auth/         # JWT authentication
│       ├── config/       # Configuration loading
│       └── database/     # Database connection
│   └── .env.example      # Environment template
│
├── .github/workflows/    # CI: validates + deploys Sqitch migrations before rollout
├── Dockerfile            # Multi-stage build
├── docker-compose.yml    # Development database (postgres)
├── package.json          # Root scripts (concurrently)
└── .tool-versions        # mise tool versions
```

## Development Commands

### Prerequisites
```bash
mise install  # Install node, go, pnpm from .tool-versions
```

### Frontend
```bash
cd frontend
pnpm install
pnpm dev          # Start dev server (port 5173)
pnpm build        # Production build
pnpm test         # Run tests
pnpm test:run     # Run tests once
pnpm lint         # Run oxlint
```

### Backend
```bash
cd backend
go mod download
go run cmd/server/main.go  # Start server (port 8080) — requires schema already deployed
go test ./...              # Run all tests
go build ./...             # Build all packages
```

### Database (Sqitch migrations)
```bash
docker compose up -d db    # Start postgres for development

# Run sqitch from backend/migrations via the official image:
cd backend/migrations
docker run --rm -v "$(pwd):/migrations" -w /migrations sqitch/sqitch deploy db:pg://postgres:postgres@localhost:5432/inventory
docker run --rm -v "$(pwd):/migrations" -w /migrations sqitch/sqitch verify db:pg://postgres:postgres@localhost:5432/inventory
docker run --rm -v "$(pwd):/migrations" -w /migrations sqitch/sqitch revert -y db:pg://postgres:postgres@localhost:5432/inventory

# Create a new migration (skeleton + plan entry):
docker run --rm -v "$(pwd):/migrations" -w /migrations sqitch/sqitch add <change_name> --note "what and why"
```
Migrations run in CI before app rollout (`.github/workflows/deploy.yml`); the app itself never mutates schema. See `backend/migrations/overview.md`.

### Full Stack
```bash
pnpm install      # Install root dependencies (concurrently)
pnpm dev          # Run both frontend and backend concurrently
```

### Docker
```bash
docker build -t inventory-app .
docker run -p 8080:8080 --env-file backend/.env inventory-app
```

## Architecture Decisions

### Frontend
- **Vite:** Fast dev server and optimized builds
- **React Router (declarative):** `<BrowserRouter>` + `<Routes>` pattern
- **React Query:** Server state management with caching
- **shadcn/ui:** Copy-paste component library with Tailwind
- **Code Splitting:** Lazy-loaded pages for better performance

### Backend
- **Gin:** High-performance HTTP framework
- **GORM:** Type-safe ORM with migrations
- **JWT:** Stateless authentication
- **Clean Architecture:** Separated concerns (handlers → services → repositories)

### Database
- **PostgreSQL:** Production database
- **Sqitch:** Versioned SQL migrations (deploy/revert/verify), applied by CI before app rollout — never by the app
- **Soft Deletes:** Using `gorm.DeletedAt` (requires `deleted_at` column from migrations)

## API Endpoints

### Products
- `GET /api/products` - List all products
- `GET /api/products/:id` - Get product by ID
- `POST /api/products` - Create product
- `PUT /api/products/:id` - Update product
- `DELETE /api/products/:id` - Delete product

### Health
- `GET /health` - Health check endpoint

## Environment Variables

### Backend (.env)
```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=inventory
DB_PORT=5432
PORT=8080
JWT_SECRET=your-secret-key
STATIC_DIR=./frontend/dist
```

### Frontend
- No environment variables needed
- API proxy configured in `vite.config.js` for development

## Code Conventions

### Frontend
- **Components:** PascalCase (e.g., `MainLayout.jsx`)
- **Utilities:** camelCase (e.g., `queryClient.js`)
- **Tests:** Co-located with components (e.g., `Home.test.jsx`)
- **Imports:** Use `@/` alias for `src/` directory

### Backend
- **Packages:** lowercase (e.g., `handlers`, `models`)
- **Files:** snake_case (e.g., `product.go`)
- **Tests:** `_test.go` suffix (e.g., `product_test.go`)
- **Error Handling:** Return errors, don't panic

## Testing Strategy

### Frontend
- **Unit Tests:** Vitest for components and utilities
- **Integration Tests:** React Testing Library for user interactions
- **Coverage:** Focus on business logic and critical paths

### Backend
- **Unit Tests:** Go testing package
- **Table-Driven Tests:** For multiple scenarios
- **Mocking:** Use interfaces for database mocking

## Deployment

### Docker Multi-Stage Build
1. **Stage 1:** Build frontend with Node
2. **Stage 2:** Build backend with Go
3. **Stage 3:** Minimal Alpine image with binary + static files

### Production
- Backend serves both API and frontend static files
- SPA routing handled by Go (fallback to index.html)
- CORS configured for cross-origin requests
- Migrations deployed by CI (`.github/workflows/deploy.yml`): `sqitch deploy` + `sqitch verify` against prod before the app image is promoted; PRs run the full deploy/verify/revert cycle against an ephemeral postgres

## Performance Optimizations

### Frontend
- **Code Splitting:** Lazy-loaded pages
- **Vendor Chunking:** Separate chunks for React, Router, Query
- **Tree Shaking:** Unused code removed in production
- **Compression:** Gzip enabled in production

### Backend
- **Connection Pooling:** GORM connection pool
- **Caching:** React Query cache for API responses
- **Static Files:** Served directly by Go (no reverse proxy needed)

## Security Considerations

- **JWT:** Stateless authentication with expiration
- **CORS:** Configured for specific origins in production
- **Input Validation:** Gin binding + custom validators
- **SQL Injection:** GORM parameterized queries
- **XSS:** React escapes output by default

## Troubleshooting

### Frontend Issues
- **Port 5173 in use:** Change port in `vite.config.js`
- **API 404:** Check backend is running on port 8080
- **Build errors:** Clear `node_modules` and reinstall

### Backend Issues
- **Database connection:** Verify PostgreSQL is running
- **Migration errors:** Check database credentials in `.env`
- **Port 8080 in use:** Change `PORT` in `.env`

### Docker Issues
- **Build fails:** Ensure Docker has enough memory (4GB+)
- **Container exits:** Check logs with `docker logs <container>`
- **Volume issues:** Use named volumes for database persistence
