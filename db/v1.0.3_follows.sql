-- ============================================================
-- Yanny v1.0.3 - 关注功能
-- 新增 follows 表，记录用户对主体（entity）的关注关系
-- ============================================================

USE `yanny`;

DROP TABLE IF EXISTS `follows`;
CREATE TABLE `follows` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id` BIGINT UNSIGNED  NOT NULL                COMMENT '小程序账号 ID',
    `user_id`       BIGINT UNSIGNED  NOT NULL                COMMENT '用户 ID',
    `entity_id`     BIGINT UNSIGNED  NOT NULL                COMMENT '被关注的主体 ID',
    `status`        TINYINT          NOT NULL DEFAULT 1      COMMENT '状态 1=已关注 0=已取消',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_user_entity` (`mp_account_id`, `user_id`, `entity_id`),
    INDEX `idx_entity` (`entity_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='关注关系表';
