-- ============================================================
-- Yanny v1.0.2 - 媒体资源去重
-- 新增 media_assets 表，基于内容哈希（七牛 etag）去重
-- ============================================================

USE `yanny`;

-- 媒体资源去重记录表
DROP TABLE IF EXISTS `media_assets`;
CREATE TABLE `media_assets` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id` BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '所属小程序账号 ID，0=全局（entity_logo/mp_icon 等无归属场景）',
    `dir_type`      VARCHAR(20)      NOT NULL                COMMENT '上传目录类型：entity_logo/mp_icon/video_cover/video',
    `content_hash`  VARCHAR(64)      NOT NULL                COMMENT '七牛返回的内容哈希（etag），权威去重键',
    `client_hash`   VARCHAR(64)      NOT NULL DEFAULT ''     COMMENT '前端预计算的 SHA-256，仅用于上传前查重加速',
    `object_key`    VARCHAR(300)     NOT NULL                COMMENT '七牛存储对象 key',
    `url`           VARCHAR(500)     NOT NULL                COMMENT '文件 CDN URL（未签名，原始路径）',
    `file_size`     BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '文件大小（字节）',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_mp_hash` (`mp_account_id`, `content_hash`),
    INDEX `idx_client_hash` (`client_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='媒体资源去重记录表';
