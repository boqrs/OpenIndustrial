QL

-- +migrate Up
-- This migration script establishes the core 'resources' table with best practices.
-- It uses an auto-incrementing integer as the primary key for performance
-- and a separate UUID for the public-facing business ID to enhance security.
-- Column names that could conflict with SQL keywords have been renamed (e.g., 'type' -> 'resource_type').

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS resources (
    -- Internal, auto-incrementing primary key for performance and joins
    id BIGSERIAL PRIMARY KEY,
    
    -- External-facing, unique business identifier to prevent enumeration attacks
    uuid uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4(),

    tenant_id uuid NOT NULL,
    
    -- Renamed to avoid SQL keyword conflicts
    resource_type VARCHAR(100) NOT NULL,
    resource_name VARCHAR(255) NOT NULL,
    code VARCHAR(100) UNIQUE,
    resource_status VARCHAR(50) NOT NULL DEFAULT 'active',
    
    metadata JSONB,
    record_version INTEGER NOT NULL DEFAULT 1,

    -- Foreign key now references the internal auto-incrementing 'id'
    parent_id uuid,

    -- Assuming owner_group_id references the 'uuid' of a group, so it remains uuid
    owner_group_id uuid,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_resources_parent
        FOREIGN KEY(parent_id) 
        REFERENCES resources(id)
);

-- Create indexes for frequently queried columns
CREATE INDEX IF NOT EXISTS idx_resources_tenant_id ON resources(tenant_id);
CREATE INDEX IF NOT EXISTS idx_resources_resource_type ON resources(resource_type);
CREATE INDEX IF NOT EXISTS idx_resources_parent_id ON resources(parent_id);
CREATE INDEX IF NOT EXISTS idx_resources_owner_group_id ON resources(owner_group_id);
CREATE INDEX IF NOT EXISTS idx_resources_deleted_at ON resources(deleted_at);

-- +migrate Down
DROP TABLE IF EXISTS resources;

-- The 'resource_relations' table defines the relationships between resources.
-- This forms the "Digital Thread" connecting all entities.
-- e.g., (product_A, 'IS_PART_OF', device_B), (device_B, 'LOCATED_IN', factory_C)
CREATE TABLE IF NOT EXISTS resource_relations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    from_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE, -- Renamed from source_resource_id
    to_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,   -- Renamed from target_resource_id
    relation_type VARCHAR(50) NOT NULL, -- e.g., 'CONTAINS', 'PRODUCED_ON', 'IS_PART_OF', 'OWNS'
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_resource_relations_from ON resource_relations (from_id);
CREATE INDEX IF NOT EXISTS idx_resource_relations_to ON resource_relations (to_id);
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

-- 'group_members' is the many-to-many link between users and groups.
-- Note: Renamed from user_groups for consistency with other potential member types.
CREATE TABLE IF NOT EXISTS group_members (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    PRIMARY KEY (tenant_id, user_id, group_id)
);

-- 'resource_groups' is the many-to-many link between resources and groups.
-- This is the core of our ABAC implementation, controlling which groups can see/manage which resources.
CREATE TABLE IF NOT EXISTS resource_groups (
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    PRIMARY KEY (tenant_id, resource_id, group_id)
);

-- +goose Down
DROP TABLE IF EXISTS resource_groups;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS resource_relations;
DROP TABLE IF EXISTS resources;