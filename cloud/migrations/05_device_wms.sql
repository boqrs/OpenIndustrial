CREATE TABLE warehouses (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_warehouses_tenant_id
    ON warehouses(tenant_id);

CREATE INDEX idx_warehouses_code
    ON warehouses(code);


CREATE TABLE warehouse_locations (
    id BIGSERIAL PRIMARY KEY,
    warehouse_id BIGINT NOT NULL,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_warehouse_locations_warehouse
        FOREIGN KEY (warehouse_id)
        REFERENCES warehouses(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_warehouse_locations_warehouse_id
    ON warehouse_locations(warehouse_id);

CREATE INDEX idx_warehouse_locations_code
    ON warehouse_locations(code);


CREATE TABLE device_inventories (
    id BIGSERIAL PRIMARY KEY,
    device_id BIGINT NOT NULL UNIQUE,
    warehouse_id BIGINT NOT NULL,
    location_id BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL,
    inbound_at TIMESTAMPTZ NOT NULL,
    outbound_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_device_inventories_warehouse_id
    ON device_inventories(warehouse_id);

CREATE INDEX idx_device_inventories_location_id
    ON device_inventories(location_id);

CREATE INDEX idx_device_inventories_status
    ON device_inventories(status);


CREATE TABLE shipments (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    sales_order_id BIGINT,
    external_order_id VARCHAR(255),
    carrier VARCHAR(100) NOT NULL,
    tracking_number VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shipments_tenant_id
    ON shipments(tenant_id);

CREATE INDEX idx_shipments_sales_order_id
    ON shipments(sales_order_id);

CREATE INDEX idx_shipments_external_order_id
    ON shipments(external_order_id);

CREATE INDEX idx_shipments_tracking_number
    ON shipments(tracking_number);

CREATE INDEX idx_shipments_status
    ON shipments(status);


CREATE TABLE shipment_items (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL,
    device_id BIGINT NOT NULL,
    sales_order_item_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shipment_items_shipment_id
    ON shipment_items(shipment_id);

CREATE INDEX idx_shipment_items_device_id
    ON shipment_items(device_id);

CREATE INDEX idx_shipment_items_sales_order_item_id
    ON shipment_items(sales_order_item_id);


CREATE TABLE shipment_tracking_events (
    id BIGSERIAL PRIMARY KEY,
    shipment_id BIGINT NOT NULL,
    external_event_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    location VARCHAR(255),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shipment_tracking_events_shipment_id
    ON shipment_tracking_events(shipment_id);

CREATE INDEX idx_shipment_tracking_events_external_event_id
    ON shipment_tracking_events(external_event_id);

CREATE INDEX idx_shipment_tracking_events_status
    ON shipment_tracking_events(status);

CREATE INDEX idx_shipment_tracking_events_occurred_at
    ON shipment_tracking_events(occurred_at);

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

ALTER TABLE sales_order_items
    ADD COLUMN IF NOT EXISTS reserved_quantity BIGINT NOT NULL DEFAULT 0;

ALTER TABLE sales_order_items
    ADD COLUMN IF NOT EXISTS fulfilled_quantity BIGINT NOT NULL DEFAULT 0;

ALTER TABLE sales_order_items
    ADD CONSTRAINT chk_sales_order_items_fulfillment_quantity
    CHECK (
        ordered_quantity > 0
        AND reserved_quantity >= 0
        AND fulfilled_quantity >= 0
        AND reserved_quantity + fulfilled_quantity <= ordered_quantity
    );


ALTER TABLE shipments
    ADD CONSTRAINT fk_shipments_sales_order
    FOREIGN KEY (sales_order_id)
    REFERENCES sales_orders(id)
    ON DELETE RESTRICT;


ALTER TABLE shipment_items
    ADD CONSTRAINT fk_shipment_items_device
    FOREIGN KEY (device_id)
    REFERENCES devices(id)
    ON DELETE RESTRICT;

ALTER TABLE shipment_items
    ADD CONSTRAINT fk_shipment_items_sales_order_item
    FOREIGN KEY (sales_order_item_id)
    REFERENCES sales_order_items(id)
    ON DELETE RESTRICT;


ALTER TABLE shipment_tracking_events
    ADD CONSTRAINT fk_shipment_tracking_events_shipment
    FOREIGN KEY (shipment_id)
    REFERENCES shipments(id)
    ON DELETE CASCADE;