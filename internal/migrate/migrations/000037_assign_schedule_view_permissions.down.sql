-- Remove schedule:view permission from branch_admin and super_admin roles
DELETE FROM role_permissions
WHERE role_id IN (SELECT role_id FROM roles WHERE name IN ('branch_admin', 'super_admin'))
  AND permission_id IN (SELECT permission_id FROM permissions WHERE name = 'schedule:view');
