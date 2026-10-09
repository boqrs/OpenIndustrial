CREATE TABLE resources (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),  
    resource_type VARCHAR(100) NOT NULL,
    resource_name VARCHAR(255) NOT NULL,
    resource_status VARCHAR(50) NOT NULL DEFAULT 'active',
    code VARCHAR(100),
    metadata JSONB,
    record_version INTEGER NOT NULL DEFAULT 1,
    parent_id BIGINT,
    owner_group_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_resources_code
    ON resources(code)
    WHERE code IS NOT NULL;

CREATE INDEX idx_resources_tenant_id
    ON resources(tenant_id);

CREATE INDEX idx_resources_resource_type
    ON resources(resource_type);

CREATE INDEX idx_resources_parent_id
    ON resources(parent_id);

CREATE INDEX idx_resources_owner_group_id
    ON resources(owner_group_id);

CREATE INDEX idx_resources_deleted_at
    ON resources(deleted_at);


CREATE TABLE resource_connections (
    id BIGSERIAL PRIMARY KEY,
    source_resource_id BIGINT NOT NULL,
    target_resource_id BIGINT NOT NULL,
    connection_type VARCHAR(100) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT fk_resource_connections_source
        FOREIGN KEY (source_resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_resource_connections_target
        FOREIGN KEY (target_resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_resource_connections
        UNIQUE (
            source_resource_id,
            target_resource_id,
            connection_type
        )
);

CREATE INDEX idx_resource_connections_source
    ON resource_connections(source_resource_id);

CREATE INDEX idx_resource_connections_target
    ON resource_connections(target_resource_id);

CREATE INDEX idx_resource_connections_type
    ON resource_connections(connection_type);

CREATE INDEX idx_resource_connections_deleted_at
    ON resource_connections(deleted_at);


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

CREATE TABLE resource_credentials (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    secret_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    consumed_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_resource_credentials_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_resource_credentials_resource_id
    ON resource_credentials(resource_id);


CREATE TABLE resource_identities (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    hardware_id VARCHAR(255),
    serial_number VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_resource_identities_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_resource_identities_resource_id
    ON resource_identities(resource_id);

CREATE INDEX idx_resource_identities_hardware_id
    ON resource_identities(hardware_id);

CREATE INDEX idx_resource_identities_serial_number
    ON resource_identities(serial_number);


CREATE TABLE resource_certificates (
    id BIGSERIAL PRIMARY KEY,
    resource_id BIGINT NOT NULL,
    certificate_id VARCHAR(512) NOT NULL,
    certificate_serial_number VARCHAR(255),
    fingerprint VARCHAR(255) NOT NULL UNIQUE,
    subject TEXT,
    issuer TEXT,
    status VARCHAR(50) NOT NULL,
    not_before TIMESTAMPTZ NOT NULL,
    not_after TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    activated_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,

    CONSTRAINT fk_resource_certificates_resource
        FOREIGN KEY (resource_id)
        REFERENCES resources(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_resource_certificates_resource_id
    ON resource_certificates(resource_id);

CREATE INDEX idx_resource_certificates_certificate_serial_number
    ON resource_certificates(certificate_serial_number);