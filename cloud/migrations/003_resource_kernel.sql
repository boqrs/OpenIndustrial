-- 003_resource_kernel.sql

-- I. Core Resource Kernel Tables
-- These two tables are the heart of our Digital Twin / Digital Thread platform.

-- The 'resources' table is a generic container for ANY entity in our ecosystem.
-- This could be a physical device, a virtual product, a supplier, a customer, etc.
-- The 'type' column gives it semantic meaning.
-- The 'properties' JSONB column allows for flexible, type-specific attributes.
CREATE TABLE IF NOT EXISTS resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- e.g., 'product', 'device', 'supplier', 'customer', 'production_line'
    name VARCHAR(255) NOT NULL,
    properties JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    version INT NOT NULL DEFAULT 1,
    parent_id UUID REFERENCES resources(id) ON DELETE SET NULL,
    owner_group_id UUID REFERENCES groups(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_resources_tenant_id_type ON resources (tenant_id, type);
CREATE INDEX IF NOT EXISTS idx_resources_name ON resources (name);
CREATE INDEX IF NOT EXISTS idx_resources_parent_id ON resources (parent_id);
CREATE INDEX IF NOT EXISTS idx_resources_owner_group_id ON resources (owner_group_id);


-- The 'resource_relations' table defines the relationships between resources.
-- This forms the "Digital Thread" connecting all entities.
-- e.g., (product_A, 'IS_PART_OF', device_B), (device_B, 'LOCATED_IN', factory_C)
CREATE TABLE IF NOT EXISTS resource_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    source_resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    target_resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    relation_type VARCHAR(50) NOT NULL, -- e.g., 'CONTAINS', 'PRODUCED_ON', 'IS_PART_OF', 'OWNS'
    properties JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_resource_relations_source ON resource_relations (source_resource_id);
CREATE INDEX IF NOT EXISTS idx_resource_relations_target ON resource_relations (target_resource_id);
CREATE INDEX IF NOT EXISTS idx_resource_relations_type ON resource_relations (relation_type);


-- II. ABAC (Attribute-Based Access Control) Support Tables
-- These tables allow for instance-level permissions.

-- 'groups' define a collection of users or resources.
-- e.g., 'Factory-A-Staff', 'Customer-B-Users', 'High-Priority-Devices'
CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, name)
);

-- 'user_groups' is the many-to-many link between users and groups.
CREATE TABLE IF NOT EXISTS user_groups (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);

-- 'resource_groups' is the many-to-many link between resources and groups.
-- This is the core of our ABAC implementation, controlling which groups can see/manage which resources.
CREATE TABLE IF NOT EXISTS resource_groups (
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    PRIMARY KEY (resource_id, group_id)
);