-- =================================================================================
-- 迁移脚本：为专题表添加阅读数字段
-- Migration: Add views field to topic table
-- Date: 2024
-- =================================================================================

-- 为专题表添加阅读数字段（如果不存在）
ALTER TABLE `topic` 
ADD COLUMN IF NOT EXISTS `views` int(11) DEFAULT 0 COMMENT '阅读数' AFTER `sort_order`;

-- 如果数据库不支持 IF NOT EXISTS，可以使用以下方式：
-- ALTER TABLE `topic` ADD COLUMN `views` int(11) DEFAULT 0 COMMENT '阅读数' AFTER `sort_order`;


