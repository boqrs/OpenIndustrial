CREATE TABLE attribute_definitions (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    data_type VARCHAR(50) NOT NULL,
    unit VARCHAR(50),
    label VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_attribute_definitions_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_attribute_definitions_resource_id
    ON attribute_definitions(resource_id);

CREATE INDEX idx_attribute_definitions_deleted_at
    ON attribute_definitions(deleted_at);

CREATE UNIQUE INDEX idx_attr_def_model_name
    ON attribute_definitions(resource_id, name);


CREATE TABLE resource_attributes (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    attribute_definition_id BIGINT NOT NULL,
    value JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_resource_attributes_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_resource_attributes_definition
        FOREIGN KEY (attribute_definition_id)
        REFERENCES attribute_definitions(id)
        ON DELETE RESTRICT,

    CONSTRAINT uq_resource_attributes
        UNIQUE (resource_id, attribute_definition_id)
);

CREATE INDEX idx_resource_attributes_resource_id
    ON resource_attributes(resource_id);

CREATE INDEX idx_resource_attributes_definition_id
    ON resource_attributes(attribute_definition_id);