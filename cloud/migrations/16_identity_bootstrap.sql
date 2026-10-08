-- ============================================================
-- OpenIndustrial initial tenants
-- ============================================================

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

-- ============================================================
-- Admin roles
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

-- ============================================================
-- Initial administrators
--
-- IMPORTANT:
-- These are development/bootstrap credentials.
-- Change them immediately after first login.
-- ============================================================

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

-- ============================================================
-- Initial password principals
--
-- Development password:
-- OpenIndustrial@123
--
-- MUST be changed after first login.
-- ============================================================

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

-- ============================================================
-- Grant all currently defined permissions to tenant Admin roles
-- ============================================================

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