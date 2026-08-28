# Products — Domain Overview (Backend)

> Status: CRUD API live · last updated 2026-08-28
> Back to [project overview](../../../docs/overview.md) · Frontend side: [`frontend/src/features/products/overview.md`](../../../frontend/src/features/products/overview.md)

## Purpose

The backend side of the products domain: a REST CRUD API backed by PostgreSQL. Self-contained feature package — model, handlers, and tests all live in this directory.

## Data Model

`Product` — `backend/internal/products/product.go`

| Field | Type | Constraints | JSON |
| --- | --- | --- | --- |
| `ID` | `uint` | primary key | `id` |
| `Name` | `string` | not null | `name` |
| `Description` | `string` | — | `description` |
| `Price` | `float64` | not null | `price` |
| `Quantity` | `int` | not null, default 0 | `quantity` |
| `SKU` | `string` | not null, unique index | `sku` |
| `CreatedAt` / `UpdatedAt` | `time.Time` | GORM auto | `created_at` / `updated_at` |
| `DeletedAt` | `gorm.DeletedAt` | index, soft delete | omitted (`json:"-"`) |

`Validate()` enforces: `name` and `sku` required, `price` and `quantity` non-negative.

## API Endpoints

All under `/api`, registered in `backend/cmd/server/main.go` via `products.NewProductHandler(db)`.

| Method | Path | Handler | Success | Errors |
| --- | --- | --- | --- | --- |
| GET | `/api/products` | `List` | 200 JSON array | 500 |
| GET | `/api/products/:id` | `Get` | 200 | 400 invalid id, 404 |
| POST | `/api/products` | `Create` | 201 | 400 bind/validate, 500 |
| PUT | `/api/products/:id` | `Update` | 200 | 400, 404 |
| DELETE | `/api/products/:id` | `Delete` | 200 `{"message": ...}` | 400, 500 |

- Request/response shape is the `Product` model JSON above (bind + validate on write).
- Update is a partial-overwrite of `Name`, `Description`, `Price`, `Quantity`, `SKU`.
- Error shape is `{"error": "..."}`.

## Key Files

| File | Responsibility |
| --- | --- |
| `product.go` | Model, table mapping, validation |
| `handler.go` | HTTP handlers (direct `gorm.DB`, no service layer) |
| `product_test.go` | Model construction + `Validate()` cases |
| `handler_test.go` | Handler tests: invalid input → 400, invalid id → 400 |
| `backend/cmd/server/main.go` | Route registration only (schema managed by Sqitch) |
| `backend/migrations/deploy/products.sql` | Products table DDL (Sqitch change) |

## Dependencies

- Gin (HTTP), GORM + postgres driver (persistence).
- Schema is owned by the [Migrations domain](../../migrations/overview.md) — `deploy/products.sql` creates the table this model queries; GORM `AutoMigrate` was removed (2026-08-28). Any model change must ship with a paired Sqitch change.

## Decisions & Conventions

- Feature-first layout: each domain is one package containing model, handler, tests, and this overview.
- Handlers operate directly on `gorm.DB` (AGENTS.md mentions services/repositories, but none exist — documented reality wins).
- Soft deletes enabled; `DeletedAt` hidden from API responses.
- Validation lives on the model (`Validate()`), called by both create and update handlers.
- **Schema coupling:** the `products` table DDL lives in `backend/migrations/deploy/products.sql` (Sqitch) and must match this model — struct changes require a paired migration.

## Tests

- `product_test.go`: field defaults + all `Validate()` error paths and the happy path.
- `handler_test.go`: POST with empty name → 400; GET `/products/invalid` → 400.
- Gap: no tests for 200/201/404/500 responses with a real DB — handler tests use `nil` DB.
