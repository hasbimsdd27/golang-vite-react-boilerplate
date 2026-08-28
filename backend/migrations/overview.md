# Migrations — Domain Overview

> Status: Sqitch active (products migration shipped) · last updated 2026-08-28
> Back to [project overview](../../docs/overview.md)

## Purpose

Versioned database schema management with [Sqitch](https://sqitch.org/), replacing GORM `AutoMigrate`. Every schema change is a named, ordered, reversible SQL change. The Go app never mutates the schema — it assumes the schema already exists and only runs queries.

## Why Sqitch (and not AutoMigrate)

`AutoMigrate` creates/alters tables from Go structs at startup. Problems with that approach:

- No history: there is no record of what the schema looked like at any point in time.
- No reversibility: dropping a column is a one-way door; there is no rollback.
- Implicit DDL: schema changes happen silently on every server start, across every environment.
- Not reviewable: a struct-tag tweak can produce a different table than intended.

Sqitch makes every change explicit: a deploy script (forward), a revert script (backward), and a verify script (proof), all reviewed in a PR.

## Layout

```
backend/migrations/
├── sqitch.conf           # Engine config (postgres); plan + top_dir defaults
├── sqitch.plan           # Ordered list of changes — the single source of truth
├── deploy/<change>.sql   # Forward SQL (applies the change)
├── revert/<change>.sql   # Back-out SQL (undoes the change)
├── verify/<change>.sql   # Assertions that the change is correctly applied
└── overview.md           # This file
```

## How Sqitch works

- The `sqitch.plan` file lists changes in dependency order. Deploying runs the deploy scripts in order; reverting runs revert scripts in reverse order.
- Sqitch records what it applied in a registry schema called `sqitch` (created automatically in the database). It contains `changes`, `events`, and `dependencies` tables.
- Each change must have `deploy`, `revert`, and `verify` scripts. `verify` is optional by default but mandatory in this project's CI (`sqitch verify` runs after every deploy).
- **Changes are immutable once deployed.** Never edit an already-deployed script — add a new change, or use `sqitch rework` if a change must be rewritten.
- Multi-environment handling: a single `sqitch.plan` + separate targets (dev, staging, prod). Sqitch tracks per-target state via the registry, so each database advances independently.

## The three-script model

| Script | Runs when | Contract |
| --- | --- | --- |
| `deploy/<change>.sql` | `sqitch deploy` | Applies the change idempotently per sqitch's registry (never re-run for an applied change) |
| `revert/<change>.sql` | `sqitch revert` | Exactly undoes the deploy script |
| `verify/<change>.sql` | `sqitch verify` | Fails (RAISE) if the schema does not match expectations |

## Change lifecycle

```bash
# 1. Create the change skeleton (plan entry + empty deploy/revert/verify scripts)
sqitch add <change_name> --note "what and why"

# 2. Edit deploy/<change>.sql, revert/<change>.sql, verify/<change>.sql

# 3. Apply to your database
sqitch deploy <target>

# 4. Prove it
sqitch verify <target>

# 5. Undo (dev only — see CI below)
sqitch revert <target>

# 6. Rewrite an already-deployed change (rare)
sqitch rework <change_name> --note "why"
```

## Targets

The plan has no hardcoded connection — targets are passed per-invocation:

```bash
sqitch deploy db:pg://USER:PASS@HOST:PORT/DBNAME
# or via env var:
export SQITCH_TARGET=db:pg://USER:PASS@HOST:PORT/DBNAME
sqitch deploy
```

## Running Sqitch locally

Sqitch is a Perl app. The project standardizes on the official Docker image so no host install is required:

```bash
# One-shot deploy against the dev database (values from backend/.env)
docker run --rm \
  -v "$(pwd)/backend/migrations:/migrations" -w /migrations \
  sqitch/sqitch deploy db:pg://postgres:postgres@localhost:5432/inventory

# Verify + revert (dev reset)
docker run --rm -v "$(pwd)/backend/migrations:/migrations" -w /migrations \
  sqitch/sqitch verify db:pg://postgres:postgres@localhost:5432/inventory
docker run --rm -v "$(pwd)/backend/migrations:/migrations" -w /migrations \
  sqitch/sqitch revert -y db:pg://postgres:postgres@localhost:5432/inventory
```

## How migrations run in CI (deployment)

Migrations are applied by CI, **not** by the app container. See `.github/workflows/deploy.yml`:

- **PR validation** (`validate` job): spins up an ephemeral postgres service container, then runs `deploy → verify → revert → deploy`. Bad SQL fails the PR before it reaches main.
- **Deploy** (`migrate` job, push to main): runs `sqitch deploy` against the production database using GitHub secrets, followed by `sqitch verify`. If verify fails, the job fails and the deployment pipeline stops — the app image is not promoted.

This keeps the runtime image free of the sqitch toolchain and guarantees the schema is correct before new app code ships.

## Coupling with GORM models

The Go app still reads/writes through GORM, and the schema must match what the models expect:

- `backend/internal/models/...` structs are the query-side contract; `backend/migrations/deploy/*.sql` is the schema.
- **Rule: any change to a model (new field, type change, index, constraint) MUST ship with a paired Sqitch change in the same commit.** The products table is defined in `deploy/products.sql` to match `backend/internal/products/product.go`.
- Column/table naming follows GORM's conventions (snake_case, plural table names) so no explicit mapping is needed.

## Current changes

| Change | Deploy | Revert | Verify | Schema produced |
| --- | --- | --- | --- | --- |
| `products` | `deploy/products.sql` | `revert/products.sql` | `verify/products.sql` | `products` table + `idx_products_sku` (unique) + `idx_products_deleted_at` |

`products` columns match the `Product` model: `id BIGSERIAL PK`, `name TEXT NOT NULL`, `description TEXT`, `price DOUBLE PRECISION NOT NULL`, `quantity INTEGER NOT NULL DEFAULT 0`, `sku TEXT NOT NULL` (unique), `created_at`/`updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`, `deleted_at TIMESTAMPTZ` (GORM soft delete).

## Key Files

| File | Responsibility |
| --- | --- |
| `sqitch.conf` | Engine + plan configuration |
| `sqitch.plan` | Ordered change list (source of truth) |
| `deploy/*.sql`, `revert/*.sql`, `verify/*.sql` | The actual changes |
| `.github/workflows/deploy.yml` | CI deploy + validation |
| `docker-compose.yml` | Dev postgres database |
| `backend/cmd/server/main.go` | App entry point (no schema logic — connects only) |

## Dependencies

- Sqitch (postgres engine), run via the `sqitch/sqitch` Docker image or a native install.
- PostgreSQL 15+ (project uses the `postgres` image in dev/CI).
- GORM models define the expected schema shape (see coupling rule above).

## Decisions & Conventions

- **Schema is code, reviewed in PRs** alongside app code.
- **Immutable changes** — deployed changes are never edited; `sqitch rework` is the escape hatch.
- **Verify is mandatory** — every change ships a `verify/` script and CI runs it.
- **No schema logic in the app** — `database.Connect` only opens a connection.
- **Targets are passed explicitly** (CLI arg or `SQITCH_TARGET`), never hardcoded in `sqitch.conf`.
- Registry schema kept at the default name `sqitch`.

## Tests

Sqitch itself is exercised in CI, not by `go test`:

- `.github/workflows/deploy.yml` → `validate` job runs the full `deploy → verify → revert → deploy` cycle against a fresh postgres on every PR.
- Manual verification loop (dev): the deploy/verify/revert commands in this doc.
- Gap: no test yet asserts the migrated schema against the GORM model programmatically (e.g. a Go integration test reading `information_schema`).
