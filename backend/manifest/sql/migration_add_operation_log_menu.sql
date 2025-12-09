-- =================================================================================
-- 添加操作日志管理菜单和权限
-- 执行时间: 2024
-- =================================================================================

-- 1. 添加操作日志管理权限
SET @system_manage_id = (SELECT id FROM permission WHERE code = 'system:manage' LIMIT 1);
-- 如果system:manage不存在，尝试查找其他系统管理权限
SET @system_manage_id = IFNULL(@system_manage_id, (SELECT id FROM permission WHERE code LIKE 'system:%' LIMIT 1));

INSERT INTO `permission` (`parent_id`, `name`, `code`, `type`, `path`, `icon`, `sort_order`, `status`) VALUES
(IFNULL(@system_manage_id, 0), '操作日志管理', 'system:log:manage', 'menu', '/admin/operation-logs', 'Document', 6, 1),
(IFNULL(@system_manage_id, 0), '操作日志列表', 'system:log:list', 'btn', '', '', 1, 1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`);

-- 2. 添加操作日志菜单到系统管理菜单下
SET @system_menu_id = (SELECT id FROM menu WHERE path = '/admin/system' AND type = 'M' LIMIT 1);
INSERT INTO `menu` (`parent_id`, `type`, `title`, `icon`, `path`, `component`, `permission`, `order`, `status`) VALUES
(@system_menu_id, 'C', '操作日志', 'Document', '/admin/operation-logs', 'views/admin/OperationLogManage', 'system:log:list', 6, 'active')
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`);

-- 3. 将操作日志权限分配给admin角色
SET @admin_role_id = (SELECT id FROM role WHERE `key` = 'admin' LIMIT 1);
INSERT INTO `role_permission` (`role_id`, `permission_id`)
SELECT @admin_role_id, id FROM permission 
WHERE code IN ('system:log:manage', 'system:log:list')
ON DUPLICATE KEY UPDATE `role_id`=VALUES(`role_id`);

