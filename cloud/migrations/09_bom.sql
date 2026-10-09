CREATE TABLE boms (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    product_id BIGINT NOT NULL,
    bom_no VARCHAR(100) NOT NULL,
    version INTEGER NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_boms_tenant_id
    ON boms(tenant_id);

CREATE INDEX idx_boms_product_id
    ON boms(product_id);

CREATE INDEX idx_boms_status
    ON boms(status);

CREATE INDEX idx_boms_deleted_at
    ON boms(deleted_at);


CREATE TABLE bom_items (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    bom_id BIGINT NOT NULL,
    material_id BIGINT NOT NULL,
    quantity NUMERIC(20,6) NOT NULL,
    unit VARCHAR(32) NOT NULL,
    sequence INTEGER NOT NULL DEFAULT 0,
    operation_code VARCHAR(100),
    is_optional BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_bom_items_tenant_id
    ON bom_items(tenant_id);

CREATE INDEX idx_bom_items_bom_id
    ON bom_items(bom_id);

CREATE INDEX idx_bom_items_material_id
    ON bom_items(material_id);

CREATE INDEX idx_bom_items_deleted_at
    ON bom_items(deleted_at);