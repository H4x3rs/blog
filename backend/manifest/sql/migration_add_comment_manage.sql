-- =================================================================================
-- 添加评论管理菜单和权限
-- 执行时间: 2024
-- =================================================================================

-- 1. 添加评论管理权限
SET @content_manage_id = (SELECT id FROM permission WHERE code = 'content:manage' LIMIT 1);
INSERT INTO `permission` (`parent_id`, `name`, `code`, `type`, `path`, `icon`, `sort_order`, `status`) VALUES
(@content_manage_id, '评论管理', 'content:comment:manage', 'menu', '/admin/comments', 'ChatLineRound', 5, 1),
(@content_manage_id, '评论列表', 'content:comment:list', 'btn', '', '', 1, 1),
(@content_manage_id, '审核评论', 'content:comment:approve', 'btn', '', '', 2, 1),
(@content_manage_id, '更新评论', 'content:comment:update', 'btn', '', '', 3, 1),
(@content_manage_id, '删除评论', 'content:comment:delete', 'btn', '', '', 4, 1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`);

-- 2. 添加评论管理菜单
SET @content_menu_id = (SELECT id FROM menu WHERE path = '/admin/content' AND type = 'M' LIMIT 1);
INSERT INTO `menu` (`parent_id`, `type`, `title`, `icon`, `path`, `component`, `permission`, `order`, `status`) VALUES
(@content_menu_id, 'C', '评论管理', 'ChatLineRound', '/admin/comments', 'views/admin/CommentManage', 'content:comment:list', 5, 'active')
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`);

-- 3. 将评论管理权限分配给admin角色
SET @admin_role_id = (SELECT id FROM role WHERE `key` = 'admin' LIMIT 1);
INSERT INTO `role_permission` (`role_id`, `permission_id`)
SELECT @admin_role_id, id FROM permission 
WHERE code IN ('content:comment:manage', 'content:comment:list', 'content:comment:approve', 'content:comment:update', 'content:comment:delete')
ON DUPLICATE KEY UPDATE `role_id`=VALUES(`role_id`);

-- 4. 将评论管理权限分配给editor角色（可选）
SET @editor_role_id = (SELECT id FROM role WHERE `key` = 'editor' LIMIT 1);
INSERT INTO `role_permission` (`role_id`, `permission_id`)
SELECT @editor_role_id, id FROM permission 
WHERE code IN ('content:comment:manage', 'content:comment:list', 'content:comment:approve', 'content:comment:update', 'content:comment:delete')
ON DUPLICATE KEY UPDATE `role_id`=VALUES(`role_id`);

