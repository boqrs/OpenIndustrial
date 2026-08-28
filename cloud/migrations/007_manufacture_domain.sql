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
    CONSTRAINT uq_production_plans_resource_uuid UNIQUE (resource_uuid),
    CONSTRAINT uq_production_plans_plan_no UNIQUE (tenant_id, plan_no),
    CONSTRAINT fk_production_plans_resource FOREIGN KEY (resource_uuid) REFERENCES resources(uuid),
    CONSTRAINT fk_production_plans_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_production_plans_factory FOREIGN KEY (factory_id) REFERENCES factories(id)
);

CREATE INDEX IF NOT EXISTS idx_production_plans_tenant_id ON production_plans(tenant_id);
CREATE INDEX IF NOT EXISTS idx_production_plans_product_id ON production_plans(product_id);
CREATE INDEX IF NOT EXISTS idx_production_plans_factory_id ON production_plans(factory_id);
CREATE INDEX IF NOT EXISTS idx_production_plans_status ON production_plans(status);
CREATE INDEX IF NOT EXISTS idx_production_plans_planned_start_at ON production_plans(planned_start_at);

-- =============================================================================
-- Routings
-- =============================================================================

CREATE TABLE IF NOT EXISTS routings (
    id BIGSERIAL PRIMARY KEY,
    resource_id  BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_routings_resource_uuid UNIQUE (resource_uuid),
    CONSTRAINT fk_routings_resource FOREIGN KEY (resource_uuid) REFERENCES resources(uuid),
    CONSTRAINT fk_routings_product FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE INDEX IF NOT EXISTS idx_routings_product_id ON routings(product_id);
CREATE INDEX IF NOT EXISTS idx_routings_status ON routings(status);
CREATE INDEX IF NOT EXISTS idx_routings_product_version ON routings(product_id, version);

-- =============================================================================
-- Routing Operations
-- =============================================================================

CREATE TABLE IF NOT EXISTS routing_operations (
    id BIGSERIAL PRIMARY KEY,
    routing_id BIGINT NOT NULL,
    sequence INTEGER NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    workstation_id BIGINT NULL,
    standard_duration_seconds BIGINT NOT NULL DEFAULT 0,
    required BOOLEAN NOT NULL DEFAULT TRUE,
    parameters JSONB NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_routing_operations_uuid UNIQUE (uuid),
    CONSTRAINT fk_routing_operations_routing FOREIGN KEY (routing_id) REFERENCES routings(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_routing_operations_routing_id ON routing_operations(routing_id);
CREATE INDEX IF NOT EXISTS idx_routing_operations_workstation_id ON routing_operations(workstation_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_routing_operations_sequence ON routing_operations(routing_id, sequence) WHERE deleted_at IS NULL;

-- =============================================================================
-- Work Orders
-- =============================================================================

CREATE TABLE IF NOT EXISTS work_orders (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGSERIAL NOT NULL,
    tenant_id UUID NOT NULL,
    production_plan_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    routing_id BIGINT NOT NULL,
    code VARCHAR(100) NOT NULL,
    planned_quantity BIGINT NOT NULL DEFAULT 0,
    priority INTEGER NOT NULL DEFAULT 0,
    due_date TIMESTAMPTZ NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    started_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_work_orders_resource_uuid UNIQUE (resource_uuid),
    CONSTRAINT uq_work_orders_code UNIQUE (tenant_id, code),
    CONSTRAINT fk_work_orders_resource FOREIGN KEY (resource_uuid) REFERENCES resources(uuid),
    CONSTRAINT fk_work_orders_production_plan FOREIGN KEY (production_plan_id) REFERENCES production_plans(id),
    CONSTRAINT fk_work_orders_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_work_orders_routing FOREIGN KEY (routing_id) REFERENCES routings(id)
);

CREATE INDEX IF NOT EXISTS idx_work_orders_tenant_id ON work_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_production_plan_id ON work_orders(production_plan_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_product_id ON work_orders(product_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_routing_id ON work_orders(routing_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_status ON work_orders(status);
CREATE INDEX IF NOT EXISTS idx_work_orders_due_date ON work_orders(due_date);

-- =============================================================================
-- Production Executions
-- =============================================================================

CREATE TABLE IF NOT EXISTS production_executions (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGSERIAL NOT NULL,
    tenant_id UUID NOT NULL,
    work_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    routing_id BIGINT NOT NULL,
    routing_version INTEGER NOT NULL,
    device_id BIGINT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_production_executions_resource_uuid UNIQUE (resource_uuid),
    CONSTRAINT fk_production_executions_resource FOREIGN KEY (resource_uuid) REFERENCES resources(uuid),
    CONSTRAINT fk_production_executions_work_order FOREIGN KEY (work_order_id) REFERENCES work_orders(id),
    CONSTRAINT fk_production_executions_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_production_executions_routing FOREIGN KEY (routing_id) REFERENCES routings(id),
    CONSTRAINT fk_production_executions_device FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE INDEX IF NOT EXISTS idx_production_executions_tenant_id ON production_executions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_production_executions_work_order_id ON production_executions(work_order_id);
CREATE INDEX IF NOT EXISTS idx_production_executions_product_id ON production_executions(product_id);
CREATE INDEX IF NOT EXISTS idx_production_executions_routing_id ON production_executions(routing_id);
CREATE INDEX IF NOT EXISTS idx_production_executions_device_id ON production_executions(device_id);
CREATE INDEX IF NOT EXISTS idx_production_executions_status ON production_executions(status);

-- =============================================================================
-- Execution Operations
-- =============================================================================

CREATE TABLE IF NOT EXISTS execution_operations (
    id BIGSERIAL PRIMARY KEY,
    execution_id BIGINT NOT NULL,
    routing_operation_id BIGINT NULL,
    sequence INTEGER NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    workstation_id BIGINT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ NULL,
    CONSTRAINT uq_execution_operations_uuid UNIQUE (uuid),
    CONSTRAINT fk_execution_operations_execution FOREIGN KEY (execution_id) REFERENCES production_executions(id) ON DELETE CASCADE,
    CONSTRAINT fk_execution_operations_routing_operation FOREIGN KEY (routing_operation_id) REFERENCES routing_operations(id)
);

CREATE INDEX IF NOT EXISTS idx_execution_operations_execution_id ON execution_operations(execution_id);
CREATE INDEX IF NOT EXISTS idx_execution_operations_routing_operation_id ON execution_operations(routing_operation_id);
CREATE INDEX IF NOT EXISTS idx_execution_operations_workstation_id ON execution_operations(workstation_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_execution_operations_sequence ON execution_operations(execution_id, sequence) WHERE deleted_at IS NULL;