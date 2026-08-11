-- Run this manually against the orders Postgres DB before starting the
-- service (or wire it into your EC2/EKS setup script). We're not using an
-- ORM's auto-migration here on purpose — keeping schema changes explicit
-- and reviewable is the production-grade habit; tools like golang-migrate
-- are the next step up from running this by hand (worth learning next).

CREATE TABLE IF NOT EXISTS orders (
    id             UUID PRIMARY KEY,
    customer_email TEXT NOT NULL,
    item           TEXT NOT NULL,
    quantity       INTEGER NOT NULL CHECK (quantity > 0),
    status         TEXT NOT NULL DEFAULT 'PENDING',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders (created_at DESC);
