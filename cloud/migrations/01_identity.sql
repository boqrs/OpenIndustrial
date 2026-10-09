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

-- ============================================================
-- Identity permissions
-- ============================================================

INSERT INTO permissions (name, description)
VALUES
    ('identity.users:create', 'Create a user'),
    ('identity.users:read', 'Read users'),
    ('identity.users:update', 'Update users'),
    ('identity.users:delete', 'Delete users'),
    ('identity.roles:assign', 'Assign user roles')
ON CONFLICT (name) DO NOTHING;

INSERT INTO tenants (
    uuid,
    name,
    code,
    status
)
VALUES
    (
        gen_random_uuid(),
        '苏州工厂',
        'suzhou',
        'active'
    ),
    (
        gen_random_uuid(),
        '东莞工厂',
        'dongguan',
        'active'
    ),
    (
        gen_random_uuid(),
        '成都工厂',
        'chengdu',
        'active'
    )
ON CONFLICT (code) DO NOTHING;


INSERT INTO roles (
    uuid,
    tenant_id,
    name,
    description,
    is_system
)
SELECT
    gen_random_uuid(),
    t.id,
    'Admin',
    'Tenant administrator',
    TRUE
FROM tenants t
WHERE t.code IN ('suzhou', 'dongguan', 'chengdu')
  AND NOT EXISTS (
      SELECT 1
      FROM roles r
      WHERE r.tenant_id = t.id
        AND r.name = 'Admin'
        AND r.deleted_at IS NULL
  );


INSERT INTO users (
    uuid,
    tenant_id,
    role_id,
    email,
    name,
    user_type,
    status
)
SELECT
    gen_random_uuid(),
    t.id,
    r.id,
    CASE t.code
        WHEN 'suzhou' THEN 'admin@suzhou.local'
        WHEN 'dongguan' THEN 'admin@dongguan.local'
        WHEN 'chengdu' THEN 'admin@chengdu.local'
    END,
    CASE t.code
        WHEN 'suzhou' THEN '苏州管理员'
        WHEN 'dongguan' THEN '东莞管理员'
        WHEN 'chengdu' THEN '成都管理员'
    END,
    'admin',
    'active'
FROM tenants t
JOIN roles r
    ON r.tenant_id = t.id
   AND r.name = 'Admin'
   AND r.deleted_at IS NULL
WHERE t.code IN ('suzhou', 'dongguan', 'chengdu')
  AND NOT EXISTS (
      SELECT 1
      FROM users u
      WHERE u.tenant_id = t.id
        AND u.user_type = 'admin'
        AND u.deleted_at IS NULL
  );


INSERT INTO principals (
    uuid,
    tenant_id,
    user_id,
    provider,
    identifier,
    secret_hash,
    status
)
SELECT
    gen_random_uuid(),
    u.tenant_id,
    u.id,
    'password',
    u.email,
    crypt(
        'OpenIndustrial@123',
        gen_salt('bf', 10)
    ),
    'active'
FROM users u
WHERE u.user_type = 'admin'
  AND u.status = 'active'
  AND u.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1
      FROM principals p
      WHERE p.tenant_id = u.tenant_id
        AND p.user_id = u.id
        AND p.provider = 'password'
        AND p.deleted_at IS NULL
  );


INSERT INTO role_permissions (
    role_id,
    permission_id
)
SELECT
    r.id,
    p.id
FROM roles r
JOIN permissions p ON TRUE
WHERE r.name = 'Admin'
  AND r.is_system = TRUE
  AND r.deleted_at IS NULL
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO roles (
    uuid,
    tenant_id,
    name,
    description,
    is_system
)
SELECT
    gen_random_uuid(),
    t.id,
    'Employee',
    'Tenant employee',
    FALSE
FROM tenants t
WHERE t.code IN ('suzhou', 'dongguan', 'chengdu')
  AND NOT EXISTS (
      SELECT 1
      FROM roles r
      WHERE r.tenant_id = t.id
        AND r.name = 'Employee'
        AND r.deleted_at IS NULL
  );


-- ============================================================
-- Identity business roles
-- Run after 16_identity_bootstrap.sql and
-- 17_identity_employee_role.sql.
-- ============================================================

INSERT INTO roles (
    uuid,
    tenant_id,
    name,
    description,
    is_system
)
SELECT
    gen_random_uuid(),
    t.id,
    role_data.name,
    role_data.description,
    TRUE
FROM tenants t
CROSS JOIN (
    VALUES
        (
            'Process Engineer',
            '工艺工程师：负责产品、BOM、工艺路线等生产配置'
        ),
        (
            'Production Planner',
            '生产计划员：负责生产计划和工单管理'
        ),
        (
            'Operator',
            '生产操作员：负责执行授权的生产工序'
        ),
        (
            'Quality Inspector',
            '质量检验员：负责授权的质量检验和结果录入'
        )
) AS role_data(name, description)
WHERE t.code IN ('suzhou', 'dongguan', 'chengdu')
  AND NOT EXISTS (
      SELECT 1
      FROM roles existing
      WHERE existing.tenant_id = t.id
        AND existing.name = role_data.name
        AND existing.deleted_at IS NULL
  );