BEGIN;

CREATE TABLE products (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price       DOUBLE PRECISION NOT NULL,
    quantity    INTEGER NOT NULL DEFAULT 0,
    sku         TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_products_sku ON products (sku);
CREATE INDEX idx_products_deleted_at ON products (deleted_at);

COMMIT;
