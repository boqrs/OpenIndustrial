CREATE TABLE routings (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_routings_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_routings_product
        FOREIGN KEY (product_id)
        REFERENCES product_models(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_routings_resource_id
    ON routings(resource_id);

CREATE INDEX idx_routings_product_id
    ON routings(product_id);

CREATE INDEX idx_routings_status
    ON routings(status);


CREATE TABLE routing_operations (
    id BIGSERIAL PRIMARY KEY,
    routing_id BIGINT NOT NULL,
    sequence INTEGER NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    workstation_id BIGINT,
    standard_duration_seconds BIGINT NOT NULL DEFAULT 0,
    required BOOLEAN NOT NULL DEFAULT TRUE,
    parameters JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_routing_operations_routing
        FOREIGN KEY (routing_id)
        REFERENCES routings(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_routing_operations_routing_id
    ON routing_operations(routing_id);

CREATE INDEX idx_routing_operations_workstation_id
    ON routing_operations(workstation_id);


CREATE TABLE work_orders (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id)
    production_plan_id BIGINT NOT NULL,
    factory_id BIGINT NOT NULL,
    production_line_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    bom_id BIGINT NOT NULL,
    routing_id BIGINT NOT NULL,
    code VARCHAR(100) NOT NULL,
    planned_quantity BIGINT NOT NULL DEFAULT 0,
    completed_quantity BIGINT NOT NULL DEFAULT 0,
    priority INTEGER NOT NULL DEFAULT 0,
    due_date TIMESTAMPTZ,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_work_orders_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_work_orders_production_plan
        FOREIGN KEY (production_plan_id)
        REFERENCES production_plans(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_work_orders_factory
        FOREIGN KEY (factory_id)
        REFERENCES factories(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_work_orders_product
        FOREIGN KEY (product_id)
        REFERENCES product_models(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_work_orders_routing
        FOREIGN KEY (routing_id)
        REFERENCES routings(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_work_orders_resource_id
    ON work_orders(resource_id);

CREATE INDEX idx_work_orders_tenant_id
    ON work_orders(tenant_id);

CREATE INDEX idx_work_orders_production_plan_id
    ON work_orders(production_plan_id);

CREATE INDEX idx_work_orders_factory_id
    ON work_orders(factory_id);

CREATE INDEX idx_work_orders_production_line_id
    ON work_orders(production_line_id);

CREATE INDEX idx_work_orders_product_id
    ON work_orders(product_id);

CREATE INDEX idx_work_orders_bom_id
    ON work_orders(bom_id);

CREATE INDEX idx_work_orders_routing_id
    ON work_orders(routing_id);

CREATE INDEX idx_work_orders_status
    ON work_orders(status);


CREATE TABLE production_executions (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id)
    work_order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    routing_id BIGINT NOT NULL,
    routing_version INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_production_executions_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_production_executions_work_order
        FOREIGN KEY (work_order_id)
        REFERENCES work_orders(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_production_executions_product
        FOREIGN KEY (product_id)
        REFERENCES product_models(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_production_executions_routing
        FOREIGN KEY (routing_id)
        REFERENCES routings(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_production_executions_resource_id
    ON production_executions(resource_id);

CREATE INDEX idx_production_executions_tenant_id
    ON production_executions(tenant_id);

CREATE INDEX idx_production_executions_work_order_id
    ON production_executions(work_order_id);

CREATE INDEX idx_production_executions_product_id
    ON production_executions(product_id);

CREATE INDEX idx_production_executions_routing_id
    ON production_executions(routing_id);

CREATE INDEX idx_production_executions_status
    ON production_executions(status);


CREATE TABLE execution_operations (
    id BIGSERIAL PRIMARY KEY,
    execution_id BIGINT NOT NULL,
    routing_operation_id BIGINT,
    sequence INTEGER NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    workstation_id BIGINT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    parameters JSONB,
    result JSONB,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_execution_operations_execution
        FOREIGN KEY (execution_id)
        REFERENCES production_executions(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_execution_operations_routing_operation
        FOREIGN KEY (routing_operation_id)
        REFERENCES routing_operations(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_execution_operations_execution_id
    ON execution_operations(execution_id);

CREATE INDEX idx_execution_operations_routing_operation_id
    ON execution_operations(routing_operation_id);

CREATE INDEX idx_execution_operations_workstation_id
    ON execution_operations(workstation_id);


CREATE TABLE execution_results (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id)
    execution_id BIGINT NOT NULL UNIQUE,
    work_order_id BIGINT NOT NULL,
    produced_quantity BIGINT NOT NULL DEFAULT 0,
    qualified_quantity BIGINT NOT NULL DEFAULT 0,
    rejected_quantity BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_execution_results_execution
        FOREIGN KEY (execution_id)
        REFERENCES production_executions(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_execution_results_work_order
        FOREIGN KEY (work_order_id)
        REFERENCES work_orders(id)
        ON DELETE RESTRICT
);

CREATE INDEX idx_execution_results_tenant_id
    ON execution_results(tenant_id);

CREATE INDEX idx_execution_results_work_order_id
    ON execution_results(work_order_id);

CREATE INDEX idx_execution_results_status
    ON execution_results(status);