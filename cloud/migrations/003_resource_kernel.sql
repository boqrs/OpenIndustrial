-- +migrate Up
-- This migration establishes the core "Resource Kernel" of the OpenIndustrial platform.
-- It defines the fundamental tables for identity, hierarchy, attributes, and connections.

-- Enable UUID generation extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Resource Status Enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'resource_status') THEN
        CREATE TYPE resource_status AS ENUM (
            'active',
            'inactive',
            'archived',
            'pending'
        );
    END IF;
END$$;

-- 2. Attribute Value Type Enum
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'attribute_value_type') THEN
        CREATE TYPE attribute_value_type AS ENUM (
            'string',
            'text',
            'integer',
            'float',
            'boolean',
            'datetime',
            'json',
            'decimal',
            'resource_reference',
            'resource_reference_list'
        );
    END IF;
END$$;


-- 3. Resources Table
CREATE TABLE IF NOT EXISTS resources (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL,
    type VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    status resource_status NOT NULL,
    parent_id UUID, -- Nullable for root resources, references resources.uuid
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_resources_tenant_id ON resources(tenant_id);
CREATE INDEX IF NOT EXISTS idx_resources_type ON resources(type);
CREATE INDEX IF NOT EXISTS idx_resources_parent_id ON resources(parent_id);
CREATE INDEX IF NOT EXISTS idx_resources_deleted_at ON resources(deleted_at);


-- 4. Attribute Definitions Table
CREATE TABLE IF NOT EXISTS attribute_definitions (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL, -- Points to a "model" or "template" resource
    name VARCHAR(255) NOT NULL,
    label VARCHAR(255),
    description TEXT,
    data_type attribute_value_type NOT NULL,
    unit VARCHAR(50),
    required BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_attribute_definitions_resource
        FOREIGN KEY(resource_id) 
        REFERENCES resources(uuid)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_attribute_definitions_tenant_id ON attribute_definitions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_attribute_definitions_resource_id ON attribute_definitions(resource_id);
CREATE INDEX IF NOT EXISTS idx_attribute_definitions_deleted_at ON attribute_definitions(deleted_at);


-- 5. Resource Attributes Table
CREATE TABLE IF NOT EXISTS resource_attributes (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    attribute_definition_id BIGINT NOT NULL,
    value JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_resource_attributes_resource
        FOREIGN KEY(resource_id) 
        REFERENCES resources(uuid)
        ON DELETE CASCADE,
    
    CONSTRAINT fk_resource_attributes_definition
        FOREIGN KEY(attribute_definition_id) 
        REFERENCES attribute_definitions(id)
        ON DELETE RESTRICT,

    UNIQUE (resource_id, attribute_definition_id)
);

CREATE INDEX IF NOT EXISTS idx_resource_attributes_tenant_id ON resource_attributes(tenant_id);


-- 6. Resource Connections Table
CREATE TABLE IF NOT EXISTS resource_connections (
    id BIGSERIAL PRIMARY KEY,
    source_resource_id BIGINT NOT NULL,
    target_resource_id BIGINT NOT NULL,
    connection_type VARCHAR(100) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_resource_connections_source
        FOREIGN KEY(source_resource_id) 
        REFERENCES resources(uuid)
        ON DELETE CASCADE,

    CONSTRAINT fk_resource_connections_target
        FOREIGN KEY(target_resource_id) 
        REFERENCES resources(uuid)
        ON DELETE CASCADE,

    UNIQUE (source_resource_id, target_resource_id, connection_type)
);

CREATE INDEX IF NOT EXISTS idx_resource_connections_tenant_id ON resource_connections(tenant_id);
CREATE INDEX IF NOT EXISTS idx_resource_connections_deleted_at ON resource_connections(deleted_at);


-- +migrate Down
DROP TABLE IF EXISTS resource_connections;
DROP TABLE IF EXISTS resource_attributes;
DROP TABLE IF EXISTS attribute_definitions;
DROP TABLE IF EXISTS resources;

DROP TYPE IF EXISTS attribute_value_type;
DROP TYPE IF EXISTS resource_status;