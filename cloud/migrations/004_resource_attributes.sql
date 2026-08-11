-- +migrate Up
CREATE TYPE attribute_value_type AS ENUM (
    'string',
    'text',
    'integer',
    'float',
    'boolean',
    'datetime',
    'json'
);

CREATE TABLE attribute_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    key VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    value_type attribute_value_type NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, key)
);

CREATE TABLE resource_attributes (
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    attribute_id UUID NOT NULL REFERENCES attribute_definitions(id) ON DELETE CASCADE,
    value_string VARCHAR(255),
    value_text TEXT,
    value_integer BIGINT,
    value_float DOUBLE PRECISION,
    value_boolean BOOLEAN,
    value_datetime TIMESTAMPTZ,
    value_json JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (resource_id, attribute_id)
);

CREATE INDEX idx_resource_attributes_resource_id ON resource_attributes(resource_id);
CREATE INDEX idx_resource_attributes_attribute_id ON resource_attributes(attribute_id);

-- Indexes on value columns for faster queries
CREATE INDEX idx_resource_attributes_value_string ON resource_attributes(value_string);
CREATE INDEX idx_resource_attributes_value_integer ON resource_attributes(value_integer);
CREATE INDEX idx_resource_attributes_value_float ON resource_attributes(value_float);
CREATE INDEX idx_resource_attributes_value_boolean ON resource_attributes(value_boolean);
CREATE INDEX idx_resource_attributes_value_datetime ON resource_attributes(value_datetime);


-- +migrate Down
DROP TABLE IF EXISTS resource_attributes;
DROP TABLE IF EXISTS attribute_definitions;
DROP TYPE IF EXISTS attribute_value_type;