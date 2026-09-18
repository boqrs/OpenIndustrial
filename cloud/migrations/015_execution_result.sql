-- =============================================================================
-- Execution Result Constraints
-- =============================================================================

-- One production execution represents exactly one physical product.
-- Therefore an execution can have at most one execution result.
CREATE TABLE IF NOT EXISTS execution_results (
    id BIGSERIAL PRIMARY KEY,

    tenant_id UUID NOT NULL,

    execution_id BIGINT NOT NULL,
    work_order_id BIGINT NOT NULL,

    produced_quantity BIGINT NOT NULL DEFAULT 0,
    qualified_quantity BIGINT NOT NULL DEFAULT 0,
    rejected_quantity BIGINT NOT NULL DEFAULT 0,

    status VARCHAR(50) NOT NULL DEFAULT 'draft',

    confirmed_at TIMESTAMPTZ NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_execution_results_execution_id
        UNIQUE (execution_id),

    CONSTRAINT fk_execution_results_execution
        FOREIGN KEY (execution_id)
        REFERENCES production_executions(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_execution_results_work_order
        FOREIGN KEY (work_order_id)
        REFERENCES work_orders(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_execution_results_quantities
        CHECK (
            produced_quantity >= 0
            AND qualified_quantity >= 0
            AND rejected_quantity >= 0
            AND qualified_quantity + rejected_quantity = produced_quantity
        )
);

CREATE INDEX IF NOT EXISTS idx_execution_results_tenant_id
    ON execution_results(tenant_id);

CREATE INDEX IF NOT EXISTS idx_execution_results_work_order_id
    ON execution_results(work_order_id);

CREATE INDEX IF NOT EXISTS idx_execution_results_execution_id
    ON execution_results(execution_id);

CREATE INDEX IF NOT EXISTS idx_execution_results_status
    ON execution_results(status);
