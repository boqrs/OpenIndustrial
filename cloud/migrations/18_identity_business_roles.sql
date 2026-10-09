
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