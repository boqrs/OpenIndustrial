CREATE TABLE devices (
    id BIGSERIAL PRIMARY KEY,

    resource_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,

    work_order_id BIGINT NOT NULL,
    execution_id BIGINT NOT NULL,
    execution_result_id BIGINT NOT NULL,

    serial_number VARCHAR(255) NOT NULL UNIQUE,
    hardware_id VARCHAR(255),

    status VARCHAR(50) NOT NULL,

    customer_uuid UUID,
    activated_at TIMESTAMPTZ,

    connection_status VARCHAR(50) NOT NULL DEFAULT 'disconnected',
    client_id VARCHAR(255),

    last_online_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_devices_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_devices_product
        FOREIGN KEY (product_id)
        REFERENCES product_models(id)
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

CREATE INDEX idx_devices_resource_id
    ON devices(resource_id);

CREATE INDEX idx_devices_product_id
    ON devices(product_id);

CREATE INDEX idx_devices_work_order_id
    ON devices(work_order_id);

CREATE INDEX idx_devices_execution_id
    ON devices(execution_id);

CREATE INDEX idx_devices_execution_result_id
    ON devices(execution_result_id);

CREATE INDEX idx_devices_hardware_id
    ON devices(hardware_id);

CREATE INDEX idx_devices_status
    ON devices(status);

CREATE INDEX idx_devices_customer_uuid
    ON devices(customer_uuid);

CREATE INDEX idx_devices_connection_status
    ON devices(connection_status);

CREATE INDEX idx_devices_client_id
    ON devices(client_id);