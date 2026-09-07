# Products — Domain Overview (Frontend)

> Status: no UI yet · last updated 2026-09-07
> Back to [project overview](../../../docs/overview.md) · Backend side: [`backend/internal/products/overview.md`](../../../backend/internal/products/overview.md)

## Purpose

The frontend side of the products domain: pages, components, and React Query hooks for listing and managing products. No code exists yet — this directory is the feature's home.

## Data Model

None yet. API contract (shared with backend): `{ id, name, description, price, quantity, sku, created_at, updated_at }`.

## API Endpoints

Consumed from the backend API described in [`backend/internal/products/overview.md`](../../../backend/internal/products/overview.md):

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/api/products` | List products |
| GET | `/api/products/:id` | Get product by ID |
| POST | `/api/products` | Create product |
| PUT | `/api/products/:id` | Update product |
| DELETE | `/api/products/:id` | Delete product |

## Key Files

| File | Responsibility |
| --- | --- |
| (none yet) | — |

## Dependencies

- `frontend/src/lib/api.ts` — Axios client (`baseURL: http://localhost:3000/api`), ready for product requests.
- TanStack React Query (`frontend/src/lib/queryClient.ts`, 5-min staleTime) for server state.
- Intended React Query keys (when implemented): `["products"]`, `["products", id]`.

## Decisions & Conventions

- Feature directory mirrors the backend domain package so docs and code are co-located per domain.
- This overview stays a stub until the first UI lands; update it in the same change that adds code.

## Tests

- None yet — no code exists. When pages/hooks are added, co-locate tests here.
