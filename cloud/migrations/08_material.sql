CREATE TABLE materials (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id)
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    material_type VARCHAR(32) NOT NULL,
    unit VARCHAR(32) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_materials_tenant_id
    ON materials(tenant_id);

CREATE INDEX idx_materials_code
    ON materials(code);

CREATE INDEX idx_materials_material_type
    ON materials(material_type);

CREATE INDEX idx_materials_deleted_at
    ON materials(deleted_at);