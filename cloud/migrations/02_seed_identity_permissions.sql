-- 002_seed_identity_permissions.sql

-- I. Define all permissions for the identity kernel
-- Format: <resource>:<action>
INSERT INTO permissions (name, description) VALUES
('identity.users:create', 'Create a new user'),
('identity.users:read', 'Read user information and list users'),
('identity.users:update', 'Update a user''s information'),
('identity.users:delete', 'Delete a user'),
('identity.roles:assign', 'Assign a role to a user');

-- II. Define which permissions each system role gets
-- We use subqueries to make this script idempotent (runnable multiple times without error)

-- Admin (Boss) gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'Admin' AND is_system_role = true),
    p.id
FROM permissions p
WHERE p.id NOT IN (SELECT permission_id FROM role_permissions WHERE role_id = (SELECT id FROM roles WHERE name = 'Admin' AND is_system_role = true));

-- Manager (厂长) gets permissions to manage users but not delete them or assign all roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'Manager' AND is_system_role = true),
    p.id
FROM permissions p
WHERE p.name IN ('identity.users:create', 'identity.users:read', 'identity.users:update', 'identity.roles:assign')
  AND p.id NOT IN (SELECT permission_id FROM role_permissions WHERE role_id = (SELECT id FROM roles WHERE name = 'Manager' AND is_system_role = true));

-- Employee (工程师) can only read user information (e.g., to see team members)
INSERT INTO role_permissions (role_id, permission_id)
SELECT
    (SELECT id FROM roles WHERE name = 'Employee' AND is_system_role = true),
    p.id
FROM permissions p
WHERE p.name = 'identity.users:read'
  AND p.id NOT IN (SELECT permission_id FROM role_permissions WHERE role_id = (SELECT id FROM roles WHERE name = 'Employee' AND is_system_role = true));
