-- =============================================================================
-- WMS <-> Sales Order Fulfillment
-- =============================================================================
--
-- SalesOrder is a group-level business entity.
-- Shipment is a tenant-owned WMS entity.
--
-- This migration establishes the business bridge:
--
-- SalesOrder
--     |
-- SalesOrderItem
--     |
-- Shipment
--     |
-- ShipmentItem
--     |
-- Device
--
-- Quantity lifecycle:
--
-- OrderedQuantity
--     |
--     +-- ReservedQuantity
--     |       |
--     |       +-- shipment cancelled -> release reservation
--     |       |
--     |       +-- shipment delivered -> convert to fulfilled
--     |
--     +-- FulfilledQuantity
--
-- Constraint:
--
-- FulfilledQuantity + ReservedQuantity <= OrderedQuantity
-- =============================================================================

BEGIN;

-- =============================================================================
-- Sales Order Item quantities
-- =============================================================================

ALTER TABLE sales_order_items
    ADD COLUMN IF NOT EXISTS reserved_quantity BIGINT
        NOT NULL DEFAULT 0;

ALTER TABLE sales_order_items
    ADD COLUMN IF NOT EXISTS fulfilled_quantity BIGINT
        NOT NULL DEFAULT 0;

ALTER TABLE sales_order_items
    DROP CONSTRAINT IF EXISTS chk_sales_order_items_fulfillment_quantity;

ALTER TABLE sales_order_items
    ADD CONSTRAINT chk_sales_order_items_fulfillment_quantity
    CHECK (
        ordered_quantity > 0
        AND reserved_quantity >= 0
        AND fulfilled_quantity >= 0
        AND reserved_quantity + fulfilled_quantity <= ordered_quantity
    );

-- =============================================================================
-- Shipment -> SalesOrder
-- =============================================================================

ALTER TABLE shipments
    ADD COLUMN IF NOT EXISTS sales_order_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_shipments_sales_order_id
    ON shipments(sales_order_id);

ALTER TABLE shipments
    DROP CONSTRAINT IF EXISTS fk_shipments_sales_order;

ALTER TABLE shipments
    ADD CONSTRAINT fk_shipments_sales_order
    FOREIGN KEY (sales_order_id)
    REFERENCES sales_orders(id)
    ON DELETE RESTRICT;

-- =============================================================================
-- ShipmentItem -> SalesOrderItem
-- =============================================================================

ALTER TABLE shipment_items
    ADD COLUMN IF NOT EXISTS sales_order_item_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_shipment_items_sales_order_item_id
    ON shipment_items(sales_order_item_id);

ALTER TABLE shipment_items
    DROP CONSTRAINT IF EXISTS fk_shipment_items_sales_order_item;

ALTER TABLE shipment_items
    ADD CONSTRAINT fk_shipment_items_sales_order_item
    FOREIGN KEY (sales_order_item_id)
    REFERENCES sales_order_items(id)
    ON DELETE RESTRICT;

-- =============================================================================
-- Shipment item uniqueness
-- =============================================================================
--
-- The same Device must not appear twice inside the same Shipment.
--
-- A Device may appear in different shipments over its lifetime because
-- shipped devices can later be returned and reshipped.
-- =============================================================================

CREATE UNIQUE INDEX IF NOT EXISTS
    uq_shipment_items_shipment_device
ON shipment_items(shipment_id, device_id);

-- =============================================================================
-- Tracking event idempotency
-- =============================================================================
--
-- ExternalEventID is the carrier's idempotency key.
-- It must exist and must be unique within a shipment.
-- =============================================================================

UPDATE shipment_tracking_events
SET external_event_id = 'legacy-' || id::text
WHERE external_event_id IS NULL
   OR btrim(external_event_id) = '';

ALTER TABLE shipment_tracking_events
    ALTER COLUMN external_event_id SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS
    uq_shipment_tracking_events_external
ON shipment_tracking_events(
    shipment_id,
    external_event_id
);

COMMIT;