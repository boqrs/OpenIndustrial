-- ============================================================================
-- 015_wms_tenant.sql
--
-- Add tenant isolation to WMS.
--
-- Tenant-owned root entities:
--   warehouses
--   shipments
--
-- Child entities inherit tenant ownership through their parent:
--   warehouse_locations -> warehouses
--   shipment_items      -> shipments
--   shipment_tracking_events -> shipments
--
-- DeviceInventory inherits tenant ownership through:
--   device_inventories -> devices -> resources -> tenant
-- ============================================================================

BEGIN;

-- ============================================================================
-- Warehouses
-- ============================================================================

ALTER TABLE warehouses
    ADD COLUMN IF NOT EXISTS tenant_id UUID;

-- The old schema used a globally unique warehouse code.
-- Warehouse codes only need to be unique inside a tenant.
DROP INDEX IF EXISTS warehouses_code_key;

CREATE UNIQUE INDEX IF NOT EXISTS uq_warehouses_tenant_code
    ON warehouses (tenant_id, code);

CREATE INDEX IF NOT EXISTS idx_warehouses_tenant_id
    ON warehouses (tenant_id);

-- ============================================================================
-- Shipments
-- ============================================================================

ALTER TABLE shipments
    ADD COLUMN IF NOT EXISTS tenant_id UUID;

CREATE INDEX IF NOT EXISTS idx_shipments_tenant_id
    ON shipments (tenant_id);

-- ============================================================================
-- Existing data guard
-- ============================================================================
--
-- We intentionally do NOT invent a tenant for existing WMS records.
--
-- This project is a new system and historical data migration must explicitly
-- assign the correct tenant before this migration can be completed against a
-- database containing old WMS records.
--
-- On an empty/new database the following checks pass immediately.
-- ============================================================================

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM warehouses
        WHERE tenant_id IS NULL
    ) THEN
        RAISE EXCEPTION
            '015_wms_tenant: warehouses contains rows without tenant_id; backfill tenant_id before applying this migration';
    END IF;

    IF EXISTS (
        SELECT 1
        FROM shipments
        WHERE tenant_id IS NULL
    ) THEN
        RAISE EXCEPTION
            '015_wms_tenant: shipments contains rows without tenant_id; backfill tenant_id before applying this migration';
    END IF;
END
$$;

ALTER TABLE warehouses
    ALTER COLUMN tenant_id SET NOT NULL;

ALTER TABLE shipments
    ALTER COLUMN tenant_id SET NOT NULL;

COMMIT;