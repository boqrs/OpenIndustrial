CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE tenants (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    name VARCHAR(128) NOT NULL,
    code VARCHAR(64) NOT NULL UNIQUE,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_tenants_deleted_at
    ON tenants(deleted_at);


CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_roles_tenant_id
    ON roles(tenant_id);

CREATE INDEX idx_roles_deleted_at
    ON roles(deleted_at);


CREATE TABLE permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_permissions_deleted_at
    ON permissions(deleted_at);


CREATE TABLE role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);


CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    role_id BIGINT NOT NULL REFERENCES roles(id),
    email VARCHAR(255) NOT NULL,
    name VARCHAR(128) NOT NULL,
    user_type VARCHAR(32) NOT NULL DEFAULT 'employee',
    status VARCHAR(32) NOT NULL DEFAULT 'invited',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_tenant_id
    ON users(tenant_id);

CREATE INDEX idx_users_role_id
    ON users(role_id);

CREATE INDEX idx_users_deleted_at
    ON users(deleted_at);


CREATE TABLE principals (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(32) NOT NULL,
    identifier VARCHAR(255) NOT NULL,
    secret_hash TEXT,
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT uq_principals_identity
        UNIQUE (tenant_id, provider, identifier)
);

CREATE INDEX idx_principals_tenant_id
    ON principals(tenant_id);

CREATE INDEX idx_principals_user_id
    ON principals(user_id);

CREATE INDEX idx_principals_deleted_at
    ON principals(deleted_at);


CREATE TABLE user_invitations (
    id BIGSERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_user_invitations_tenant_id
    ON user_invitations(tenant_id);

CREATE INDEX idx_user_invitations_user_id
    ON user_invitations(user_id);

CREATE INDEX idx_user_invitations_token_hash
    ON user_invitations(token_hash);

CREATE INDEX idx_user_invitations_deleted_at
    ON user_invitations(deleted_at);