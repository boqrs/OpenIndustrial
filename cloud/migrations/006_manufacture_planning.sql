-- =============================================================================
-- Production Plans
-- =============================================================================

CREATE TABLE IF NOT EXISTS production_plans (
    id BIGSERIAL PRIMARY KEY,

    resource_id BIGINT NOT NULL,

    tenant_id UUID NOT NULL,

    plan_no VARCHAR(100) NOT NULL,

    product_id BIGINT NOT NULL,

    factory_id BIGINT NOT NULL,

    planned_quantity BIGINT NOT NULL,

    planned_start_at TIMESTAMPTZ NOT NULL,

    planned_end_at TIMESTAMPTZ NOT NULL,

    status VARCHAR(32) NOT NULL DEFAULT 'draft',

    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,

    CONSTRAINT uq_production_plans_resource_uuid
        UNIQUE (resource_uuid),

    CONSTRAINT uq_production_plans_plan_no
        UNIQUE (tenant_id, plan_no),

    CONSTRAINT fk_production_plans_resource
        FOREIGN KEY (resource_uuid)
        REFERENCES resources(uuid),

    CONSTRAINT fk_production_plans_product
        FOREIGN KEY (product_id)
        REFERENCES products(id),

    CONSTRAINT fk_production_plans_factory
        FOREIGN KEY (factory_id)
        REFERENCES factories(id)
);

CREATE INDEX IF NOT EXISTS idx_production_plans_tenant_id
    ON production_plans(tenant_id);

CREATE INDEX IF NOT EXISTS idx_production_plans_product_id
    ON production_plans(product_id);

CREATE INDEX IF NOT EXISTS idx_production_plans_factory_id
    ON production_plans(factory_id);

CREATE INDEX IF NOT EXISTS idx_production_plans_status
    ON production_plans(status);

CREATE INDEX IF NOT EXISTS idx_production_plans_planned_start_at
    ON production_plans(planned_start_at);
