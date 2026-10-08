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