-- ============================================================
-- Yanny v1.0.1 - 管理员数据级权限（主体隔离）
-- 新增 admin_entities 表，实现管理员→主体绑定
-- ============================================================

USE `yanny`;

-- 管理员-主体绑定表（数据级权限隔离）
DROP TABLE IF EXISTS `admin_entities`;
CREATE TABLE `admin_entities` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `admin_id`      BIGINT UNSIGNED  NOT NULL                COMMENT '管理员 ID',
    `entity_id`     BIGINT UNSIGNED  NOT NULL                COMMENT '主体 ID',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_admin_entity` (`admin_id`, `entity_id`),
    INDEX `idx_entity` (`entity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员-主体绑定表';

-- 超级管理员默认绑定所有已有主体
INSERT INTO `admin_entities` (`admin_id`, `entity_id`)
SELECT 1, id FROM `entities`;
