CREATE TABLE IF NOT EXISTS production_plan_allocations (
    id BIGSERIAL PRIMARY KEY,

    production_plan_id BIGINT NOT NULL,
    sales_order_item_id BIGINT NOT NULL,

    allocated_quantity BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_production_plan_allocations_plan
        FOREIGN KEY (production_plan_id)
        REFERENCES production_plans(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_production_plan_allocations_sales_order_item
        FOREIGN KEY (sales_order_item_id)
        REFERENCES sales_order_items(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_production_plan_allocations_quantity
        CHECK (allocated_quantity > 0)
);

CREATE INDEX IF NOT EXISTS idx_production_plan_allocations_plan_id
    ON production_plan_allocations(production_plan_id);

CREATE INDEX IF NOT EXISTS idx_production_plan_allocations_sales_order_item_id
    ON production_plan_allocations(sales_order_item_id);