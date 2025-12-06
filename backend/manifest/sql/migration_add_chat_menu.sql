-- =================================================================================
-- 添加AI对话菜单和权限
-- 执行时间: 2024
-- =================================================================================

-- 1. 添加AI对话权限
INSERT INTO `permission` (`parent_id`, `name`, `code`, `type`, `path`, `icon`, `sort_order`, `status`) VALUES
(0, 'AI对话', 'chat:view', 'menu', '/admin/chat', 'ChatLineRound', 100, 1),
(0, '发送消息', 'chat:message:send', 'btn', '', '', 101, 1),
(0, '清除会话', 'chat:session:clear', 'btn', '', '', 102, 1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`);

-- 2. 添加AI对话菜单
-- 作为一级菜单添加到系统管理目录下
SET @system_menu_id = (SELECT id FROM menu WHERE path = '/admin/system' AND type = 'M' LIMIT 1);
INSERT INTO `menu` (`parent_id`, `type`, `title`, `icon`, `path`, `component`, `permission`, `order`, `status`) VALUES
(@system_menu_id, 'C', 'AI 对话', 'ChatLineRound', '/admin/chat', 'views/admin/Chat', 'chat:view', 6, 'active')
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`);

-- 3. 将AI对话权限分配给admin角色
SET @admin_role_id = (SELECT id FROM role WHERE `key` = 'admin' LIMIT 1);
INSERT INTO `role_permission` (`role_id`, `permission_id`)
SELECT @admin_role_id, id FROM permission 
WHERE code IN ('chat:view', 'chat:message:send', 'chat:session:clear')
ON DUPLICATE KEY UPDATE `role_id`=VALUES(`role_id`);

-- 4. 将AI对话权限分配给editor角色（可选，如果需要编辑也能使用）
SET @editor_role_id = (SELECT id FROM role WHERE `key` = 'editor' LIMIT 1);
INSERT INTO `role_permission` (`role_id`, `permission_id`)
SELECT @editor_role_id, id FROM permission 
WHERE code IN ('chat:view', 'chat:message:send', 'chat:session:clear')
ON DUPLICATE KEY UPDATE `role_id`=VALUES(`role_id`);


