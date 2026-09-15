-- =============================================================================
-- WMS
-- =============================================================================

-- -----------------------------------------------------------------------------
-- Warehouses
-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS warehouses (
    id BIGSERIAL PRIMARY KEY,

    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,

    address TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- -----------------------------------------------------------------------------
-- Warehouse Locations
-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS warehouse_locations (
    id BIGSERIAL PRIMARY KEY,

    warehouse_id BIGINT NOT NULL,

    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_warehouse_locations_warehouse
        FOREIGN KEY (warehouse_id)
        REFERENCES warehouses(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_warehouse_locations_code
        UNIQUE (warehouse_id, code)
);

CREATE INDEX IF NOT EXISTS idx_warehouse_locations_warehouse_id
    ON warehouse_locations(warehouse_id);

-- -----------------------------------------------------------------------------
-- Device Inventory
-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS device_inventories (
    id BIGSERIAL PRIMARY KEY,

    device_id BIGINT NOT NULL UNIQUE,

    warehouse_id BIGINT NOT NULL,
    location_id BIGINT NOT NULL,

    status VARCHAR(50) NOT NULL DEFAULT 'in_stock',

    inbound_at TIMESTAMPTZ NOT NULL,
    outbound_at TIMESTAMPTZ NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_device_inventories_device
        FOREIGN KEY (device_id)
        REFERENCES devices(id),

    CONSTRAINT fk_device_inventories_warehouse
        FOREIGN KEY (warehouse_id)
        REFERENCES warehouses(id),

    CONSTRAINT fk_device_inventories_location
        FOREIGN KEY (location_id)
        REFERENCES warehouse_locations(id)
);

CREATE INDEX IF NOT EXISTS idx_device_inventories_warehouse_id
    ON device_inventories(warehouse_id);

CREATE INDEX IF NOT EXISTS idx_device_inventories_location_id
    ON device_inventories(location_id);

CREATE INDEX IF NOT EXISTS idx_device_inventories_status
    ON device_inventories(status);

-- -----------------------------------------------------------------------------
-- Shipments
-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS shipments (
    id BIGSERIAL PRIMARY KEY,

    external_order_id VARCHAR(255),

    carrier VARCHAR(100) NOT NULL,
    tracking_number VARCHAR(255) NOT NULL,

    status VARCHAR(50) NOT NULL DEFAULT 'created',

    shipped_at TIMESTAMPTZ NULL,
    delivered_at TIMESTAMPTZ NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_shipments_external_order_id
    ON shipments(external_order_id);

CREATE INDEX IF NOT EXISTS idx_shipments_tracking_number
    ON shipments(tracking_number);

CREATE INDEX IF NOT EXISTS idx_shipments_status
    ON shipments(status);

-- -----------------------------------------------------------------------------
-- Shipment Items
-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS shipment_items (
    id BIGSERIAL PRIMARY KEY,

    shipment_id BIGINT NOT NULL,
    device_id BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_shipment_items_shipment
        FOREIGN KEY (shipment_id)
        REFERENCES shipments(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_shipment_items_device
        FOREIGN KEY (device_id)
        REFERENCES devices(id)
);

CREATE INDEX IF NOT EXISTS idx_shipment_items_shipment_id
    ON shipment_items(shipment_id);

CREATE INDEX IF NOT EXISTS idx_shipment_items_device_id
    ON shipment_items(device_id);

-- -----------------------------------------------------------------------------
-- Shipment Tracking Events
-- -----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS shipment_tracking_events (
    id BIGSERIAL PRIMARY KEY,

    shipment_id BIGINT NOT NULL,

    external_event_id VARCHAR(255),

    status VARCHAR(50) NOT NULL,

    occurred_at TIMESTAMPTZ NOT NULL,

    location VARCHAR(255),
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_shipment_tracking_events_shipment
        FOREIGN KEY (shipment_id)
        REFERENCES shipments(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_shipment_tracking_events_shipment_id
    ON shipment_tracking_events(shipment_id);

CREATE INDEX IF NOT EXISTS idx_shipment_tracking_events_external_event_id
    ON shipment_tracking_events(external_event_id);

CREATE INDEX IF NOT EXISTS idx_shipment_tracking_events_status
    ON shipment_tracking_events(status);

CREATE INDEX IF NOT EXISTS idx_shipment_tracking_events_occurred_at
    ON shipment_tracking_events(occurred_at);

-- Same third-party event must not be processed twice for one shipment.
CREATE UNIQUE INDEX IF NOT EXISTS uq_shipment_tracking_events_external_event
    ON shipment_tracking_events(shipment_id, external_event_id)
    WHERE external_event_id IS NOT NULL;