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
    TRUE
FROM tenants t
WHERE t.code IN ('suzhou', 'dongguan', 'chengdu')
  AND NOT EXISTS (
      SELECT 1
      FROM roles r
      WHERE r.tenant_id = t.id
        AND r.name = 'Employee'
        AND r.deleted_at IS NULL
  );