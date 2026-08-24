CREATE TABLE work_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    tenant_id UUID NOT NULL,

    order_no VARCHAR(100) NOT NULL,

    production_plan_id UUID NOT NULL,

    product_id UUID NOT NULL,

    factory_id UUID NOT NULL,

    planned_quantity INTEGER NOT NULL DEFAULT 0,

    completed_quantity INTEGER NOT NULL DEFAULT 0,

    planned_start_at TIMESTAMPTZ NOT NULL,

    planned_end_at TIMESTAMPTZ NOT NULL,

    status VARCHAR(32) NOT NULL DEFAULT 'draft',

    priority INTEGER NOT NULL DEFAULT 0,

    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT work_orders_planned_quantity_positive
        CHECK (planned_quantity > 0),

    CONSTRAINT work_orders_completed_quantity_valid
        CHECK (
            completed_quantity >= 0
            AND completed_quantity <= planned_quantity
        ),

    CONSTRAINT work_orders_time_valid
        CHECK (planned_end_at > planned_start_at),

    CONSTRAINT work_orders_order_no_unique
        UNIQUE (tenant_id, order_no)
);

CREATE INDEX idx_work_orders_tenant_id
    ON work_orders (tenant_id);

CREATE INDEX idx_work_orders_production_plan_id
    ON work_orders (production_plan_id);

CREATE INDEX idx_work_orders_product_id
    ON work_orders (product_id);

CREATE INDEX idx_work_orders_factory_id
    ON work_orders (factory_id);

CREATE INDEX idx_work_orders_status
    ON work_orders (tenant_id, status);

CREATE INDEX idx_work_orders_schedule
    ON work_orders (
        tenant_id,
        planned_start_at,
        planned_end_at
    );