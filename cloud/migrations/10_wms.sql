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