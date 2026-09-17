CREATE TABLE IF NOT EXISTS devices (
    id BIGSERIAL PRIMARY KEY,

    resource_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,

    -- Manufacturing provenance
    work_order_id BIGINT NOT NULL,
    execution_id BIGINT NOT NULL,
    execution_result_id BIGINT NOT NULL,

    -- Device identity
    serial_number VARCHAR(255) NOT NULL,
    hardware_id VARCHAR(255),

    -- Runtime state
    status VARCHAR(50) NOT NULL,

    activated_at TIMESTAMPTZ,
    last_online_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_devices_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_devices_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_devices_work_order
        FOREIGN KEY (work_order_id)
        REFERENCES work_orders(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_devices_execution
        FOREIGN KEY (execution_id)
        REFERENCES production_executions(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_devices_execution_result
        FOREIGN KEY (execution_result_id)
        REFERENCES execution_results(id)
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_devices_resource_id
    ON devices(resource_id);

CREATE INDEX IF NOT EXISTS idx_devices_product_id
    ON devices(product_id);

CREATE INDEX IF NOT EXISTS idx_devices_work_order_id
    ON devices(work_order_id);

CREATE INDEX IF NOT EXISTS idx_devices_execution_id
    ON devices(execution_id);

CREATE INDEX IF NOT EXISTS idx_devices_execution_result_id
    ON devices(execution_result_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_serial_number_unique
    ON devices(serial_number);

CREATE INDEX IF NOT EXISTS idx_devices_hardware_id
    ON devices(hardware_id);

CREATE INDEX IF NOT EXISTS idx_devices_status
    ON devices(status);
