-- 评论表
CREATE TABLE IF NOT EXISTS `comment` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `article_id` int(11) NOT NULL COMMENT '文章ID',
  `user_id` int(11) DEFAULT NULL COMMENT '用户ID（可为空，支持匿名评论）',
  `parent_id` int(11) DEFAULT 0 COMMENT '父评论ID（0表示顶级评论）',
  `content` text NOT NULL COMMENT '评论内容',
  `status` varchar(20) DEFAULT 'approved' COMMENT '状态:approved已审核,pending待审核,rejected已拒绝',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_article_id` (`article_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表';

