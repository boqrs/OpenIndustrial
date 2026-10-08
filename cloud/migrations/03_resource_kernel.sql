CREATE TABLE resources (
    id BIGSERIAL PRIMARY KEY,
    tenant_id UUID NOT NULL,
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