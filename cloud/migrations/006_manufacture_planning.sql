CREATE TABLE production_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    plan_no VARCHAR(100) NOT NULL,
    product_id UUID NOT NULL,
    factory_id UUID NOT NULL,
    planned_quantity INTEGER NOT NULL,
    planned_start_at TIMESTAMPTZ NOT NULL,
    planned_end_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT production_plans_quantity_positive
        CHECK (planned_quantity > 0),
    CONSTRAINT production_plans_time_valid
        CHECK (planned_end_at > planned_start_at),
    CONSTRAINT production_plans_plan_no_unique
        UNIQUE (tenant_id, plan_no)
);

CREATE INDEX idx_production_plans_tenant_id ON production_plans (tenant_id);
CREATE INDEX idx_production_plans_product_id ON production_plans (product_id);
CREATE INDEX idx_production_plans_factory_id ON production_plans (factory_id);
CREATE INDEX idx_production_plans_status ON production_plans (tenant_id, status);
CREATE INDEX idx_production_plans_schedule ON production_plans (tenant_id, planned_start_at);