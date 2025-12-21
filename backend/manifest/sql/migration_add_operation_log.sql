-- 操作日志表（用于记录登录和管理操作）
CREATE TABLE IF NOT EXISTS `operation_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL COMMENT '操作用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '操作用户名',
  `operation_type` varchar(50) NOT NULL COMMENT '操作类型：login, create, update, delete等',
  `module` varchar(50) DEFAULT NULL COMMENT '操作模块：user, article, category等',
  `operation_desc` varchar(255) DEFAULT NULL COMMENT '操作描述',
  `request_method` varchar(10) DEFAULT NULL COMMENT '请求方法：GET, POST等',
  `request_path` varchar(255) DEFAULT NULL COMMENT '请求路径',
  `request_params` text DEFAULT NULL COMMENT '请求参数（JSON格式）',
  `ip_address` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` varchar(500) DEFAULT NULL COMMENT '用户代理',
  `status` tinyint(1) DEFAULT 1 COMMENT '操作状态：0失败，1成功',
  `error_message` text DEFAULT NULL COMMENT '错误信息',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_operation_type` (`operation_type`),
  KEY `idx_module` (`module`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';


