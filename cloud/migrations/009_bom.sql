CREATE TABLE IF NOT EXISTS boms (
    id BIGSERIAL PRIMARY KEY,

    tenant_id UUID NOT NULL,

    product_id UUID NOT NULL,

    bom_no VARCHAR(100) NOT NULL,

    version INTEGER NOT NULL,

    status VARCHAR(32) NOT NULL DEFAULT 'draft',

    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS
idx_boms_tenant_no_version
ON boms (
    tenant_id,
    bom_no,
    version
)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS
idx_boms_tenant_product
ON boms (
    tenant_id,
    product_id
)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS
idx_boms_tenant_status
ON boms (
    tenant_id,
    status
)
WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS bom_items (
    id BIGSERIAL PRIMARY KEY,

    tenant_id UUID NOT NULL,

    bom_id BIGINT NOT NULL,

    material_id BIGINT NOT NULL,

    quantity NUMERIC(20, 6) NOT NULL,

    unit VARCHAR(32) NOT NULL,

    sequence INTEGER NOT NULL DEFAULT 0,

    operation_code VARCHAR(100),

    is_optional BOOLEAN NOT NULL DEFAULT FALSE,

    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS
idx_bom_items_tenant_bom
ON bom_items (
    tenant_id,
    bom_id
)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS
idx_bom_items_tenant_material
ON bom_items (
    tenant_id,
    material_id
)
WHERE deleted_at IS NULL;