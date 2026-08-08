-- ============================================================
-- Yanny v1.0.4 - 主体经纬度 + 分享裂变邀请关系
-- entities 新增经纬度字段（用于小程序端地址导航）
-- users 新增邀请人字段（单层归属，用于分享裂变统计）
-- ============================================================

USE `yanny`;

ALTER TABLE `entities`
    ADD COLUMN `latitude`  DECIMAL(10,7) NULL COMMENT '纬度（管理后台手动录入，用于小程序端地址导航）' AFTER `address`,
    ADD COLUMN `longitude` DECIMAL(10,7) NULL COMMENT '经度（管理后台手动录入，用于小程序端地址导航）' AFTER `latitude`;

ALTER TABLE `users`
    ADD COLUMN `inviter_user_id` BIGINT UNSIGNED NULL COMMENT '邀请人用户 ID（分享裂变，单层归属，首次登录时绑定后不可变）' AFTER `status`,
    ADD INDEX `idx_inviter` (`inviter_user_id`);
