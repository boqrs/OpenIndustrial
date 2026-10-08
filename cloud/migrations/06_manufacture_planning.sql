CREATE TABLE product_models (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    code VARCHAR(100) NOT NULL,
    version VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_product_models_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_product_models_resource_id
    ON product_models(resource_id);

CREATE INDEX idx_product_models_category
    ON product_models(category);


CREATE TABLE factories (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    code VARCHAR(100) NOT NULL UNIQUE,
    address TEXT,
    timezone VARCHAR(100) NOT NULL DEFAULT 'UTC',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_factories_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_factories_resource_id
    ON factories(resource_id);


CREATE TABLE production_plans (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL,
    plan_no VARCHAR(100) NOT NULL,
    product_id BIGINT NOT NULL,
    factory_id BIGINT NOT NULL,
    planned_quantity BIGINT NOT NULL,
    planned_start_at TIMESTAMPTZ NOT NULL,
    planned_end_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_production_plans_product
        FOREIGN KEY (product_id)
        REFERENCES product_models(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_production_plans_factory
        FOREIGN KEY (factory_id)
        REFERENCES factories(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_production_plans_tenant_id
    ON production_plans(tenant_id);

CREATE INDEX idx_production_plans_product_id
    ON production_plans(product_id);

CREATE INDEX idx_production_plans_factory_id
    ON production_plans(factory_id);

CREATE INDEX idx_production_plans_planned_start_at
    ON production_plans(planned_start_at);

CREATE INDEX idx_production_plans_status
    ON production_plans(status);