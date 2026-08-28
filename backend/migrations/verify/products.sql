BEGIN;

DO $$
BEGIN
    IF to_regclass('public.products') IS NULL THEN
        RAISE EXCEPTION 'verify failed: products table does not exist';
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = 'public'
          AND tablename = 'products'
          AND indexname = 'idx_products_sku'
    ) THEN
        RAISE EXCEPTION 'verify failed: unique index idx_products_sku is missing on products';
    END IF;

    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes
        WHERE schemaname = 'public'
          AND tablename = 'products'
          AND indexname = 'idx_products_deleted_at'
    ) THEN
        RAISE EXCEPTION 'verify failed: index idx_products_deleted_at is missing on products';
    END IF;
END $$;

ROLLBACK;
