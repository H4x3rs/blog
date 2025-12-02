-- 添加ICP备案号设置
-- 如果不存在则插入，如果存在则更新
INSERT INTO settings (`key`, `value`, `type`, `desc`, `created_at`, `updated_at`)
VALUES ('icpNumber', '京ICP备2021032569号', 'string', 'ICP备案号', NOW(), NOW())
ON DUPLICATE KEY UPDATE 
  `value` = '京ICP备2021032569号',
  `updated_at` = NOW();

