-- =============================================================================
-- Customer
-- =============================================================================

CREATE TABLE IF NOT EXISTS customers (
    id BIGSERIAL PRIMARY KEY,

    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(32) NOT NULL,

    contact_name VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(64),
    address TEXT,

    status VARCHAR(32) NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_customers_code
        UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS idx_customers_name
    ON customers(name);

CREATE INDEX IF NOT EXISTS idx_customers_type
    ON customers(type);

CREATE INDEX IF NOT EXISTS idx_customers_status
    ON customers(status);