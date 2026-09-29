CREATE TABLE IF NOT EXISTS materials (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    material_type VARCHAR(32) NOT NULL,
    unit VARCHAR(32) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS
idx_materials_tenant_code
ON materials (tenant_id, code)
WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS
idx_materials_tenant_type
ON materials (tenant_id, material_type)
WHERE deleted_at IS NULL;