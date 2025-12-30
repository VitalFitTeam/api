-- Assign schedule:view permission to branch_admin and super_admin roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.name IN ('branch_admin', 'super_admin')
  AND p.name = 'schedule:view'
ON CONFLICT DO NOTHING;
