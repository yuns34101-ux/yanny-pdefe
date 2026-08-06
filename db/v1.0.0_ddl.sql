-- ============================================================
-- Yanny 小程序短视频平台 - 数据库初始化 DDL
-- Version: v1.0.0
-- Date: 2026-08-06
-- DB: MySQL 8.0
-- Charset: utf8mb4 / utf8mb4_unicode_ci
-- ============================================================

CREATE DATABASE IF NOT EXISTS `yanny` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `yanny`;

-- ============================================================
-- 1. 主体表 (entities)
--    运营主体/品牌方，一个主体可绑定多个小程序账号
-- ============================================================
DROP TABLE IF EXISTS `entities`;
CREATE TABLE `entities` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name`          VARCHAR(100)     NOT NULL                COMMENT '主体名称（品牌名）',
    `logo_url`      VARCHAR(500)     NOT NULL DEFAULT ''     COMMENT '主体 Logo URL',
    `description`   VARCHAR(500)     NOT NULL DEFAULT ''     COMMENT '主体简介',
    `contact_phone` VARCHAR(20)      NOT NULL DEFAULT ''     COMMENT '联系电话',
    `contact_email` VARCHAR(100)     NOT NULL DEFAULT ''     COMMENT '联系邮箱',
    `address`       VARCHAR(300)     NOT NULL DEFAULT ''     COMMENT '地址',
    `extra`         JSON             NULL                    COMMENT '扩展信息（自定义字段）',
    `sort_order`    INT              NOT NULL DEFAULT 0      COMMENT '排序（降序）',
    `status`        TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=启用 0=禁用',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`    DATETIME         NULL                    COMMENT '软删除',
    PRIMARY KEY (`id`),
    INDEX `idx_status_sort` (`status`, `sort_order` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='运营主体表';


-- ============================================================
-- 2. 小程序账号表 (mp_accounts)
--    每个小程序账号有独立的 AppID / AppSecret
-- ============================================================
DROP TABLE IF EXISTS `mp_accounts`;
CREATE TABLE `mp_accounts` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `app_id`        VARCHAR(64)      NOT NULL                COMMENT '小程序 AppID（微信分配）',
    `app_secret`    VARCHAR(128)     NOT NULL                COMMENT '小程序 AppSecret（加密存储）',
    `app_name`      VARCHAR(100)     NOT NULL                COMMENT '小程序名称',
    `app_icon`      VARCHAR(500)     NOT NULL DEFAULT ''     COMMENT '小程序图标 URL',
    `description`   VARCHAR(300)     NOT NULL DEFAULT ''     COMMENT '备注说明',
    `status`        TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=启用 0=禁用',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`    DATETIME         NULL                    COMMENT '软删除',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_app_id` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='小程序账号表';


-- ============================================================
-- 3. 主体-小程序绑定表 (entity_mp_bindings)
--    多对多：一个主体可绑定多个小程序，一个小程序也可绑定多个主体
-- ============================================================
DROP TABLE IF EXISTS `entity_mp_bindings`;
CREATE TABLE `entity_mp_bindings` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `entity_id`     BIGINT UNSIGNED  NOT NULL                COMMENT '主体 ID',
    `mp_account_id` BIGINT UNSIGNED  NOT NULL                COMMENT '小程序账号 ID',
    `is_default`    TINYINT          NOT NULL DEFAULT 0      COMMENT '是否默认绑定：1=是 0=否',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_entity_mp` (`entity_id`, `mp_account_id`),
    INDEX `idx_mp_account` (`mp_account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='主体-小程序账号绑定表（多对多）';


-- ============================================================
-- 4. CDN 配置表 (cdn_configs)
--    七牛云等 CDN 配置，按小程序账号维度管理
-- ============================================================
DROP TABLE IF EXISTS `cdn_configs`;
CREATE TABLE `cdn_configs` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id` BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `provider`      VARCHAR(30)      NOT NULL DEFAULT 'qiniu' COMMENT 'CDN 厂商：qiniu/aliyun/tencent',
    `access_key`    VARCHAR(200)     NOT NULL                COMMENT 'Access Key（加密存储）',
    `secret_key`    VARCHAR(200)     NOT NULL                COMMENT 'Secret Key（加密存储）',
    `bucket`        VARCHAR(100)     NOT NULL                COMMENT '存储空间名称',
    `domain`        VARCHAR(200)     NOT NULL                COMMENT 'CDN 加速域名',
    `region`        VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '存储区域（如 z0/z1/z2/na0）',
    `callback_url`  VARCHAR(300)     NOT NULL DEFAULT ''     COMMENT '回调 URL（上传完成通知）',
    `status`        TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=启用 0=禁用',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_mp_account` (`mp_account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='CDN 配置表';


-- ============================================================
-- 5. 视频分类表 (video_categories)
--    分类属于主体（entity），同一主体在不同小程序中共享同一套分类体系
--    mp_account_id 用于标记该分类在哪个小程序中生效展示
-- ============================================================
DROP TABLE IF EXISTS `video_categories`;
CREATE TABLE `video_categories` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `entity_id`     BIGINT UNSIGNED  NOT NULL                COMMENT '关联主体 ID（分类所有权归主体）',
    `mp_account_id` BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID（生效范围）',
    `name`          VARCHAR(50)      NOT NULL                COMMENT '分类名称（如：精选、萌宠、健身）',
    `icon_url`      VARCHAR(500)     NOT NULL DEFAULT ''     COMMENT '分类图标 URL',
    `sort_order`    INT              NOT NULL DEFAULT 0      COMMENT '排序（降序）',
    `status`        TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=启用 0=禁用',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`    DATETIME         NULL                    COMMENT '软删除',
    PRIMARY KEY (`id`),
    INDEX `idx_entity_mp_sort` (`entity_id`, `mp_account_id`, `sort_order` DESC),
    INDEX `idx_mp` (`mp_account_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频分类表（归属主体）';


-- ============================================================
-- 6. 视频表 (videos)
--    核心内容表，关联主体、小程序、分类
-- ============================================================
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos` (
    `id`              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id`   BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `entity_id`       BIGINT UNSIGNED  NOT NULL                COMMENT '关联主体 ID',
    `category_id`     BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '关联分类 ID',
    `title`           VARCHAR(200)     NOT NULL DEFAULT ''     COMMENT '视频标题',
    `description`     VARCHAR(1000)    NOT NULL DEFAULT ''     COMMENT '视频描述',
    `cover_url`       VARCHAR(500)     NOT NULL                COMMENT '封面图 URL',
    `video_url`       VARCHAR(500)     NOT NULL                COMMENT '视频文件 URL',
    `duration`        INT UNSIGNED     NOT NULL DEFAULT 0      COMMENT '视频时长（秒）',
    `width`           INT UNSIGNED     NOT NULL DEFAULT 0      COMMENT '视频宽度（px）',
    `height`          INT UNSIGNED     NOT NULL DEFAULT 0      COMMENT '视频高度（px）',
    `file_size`       BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '文件大小（字节）',
    `tags`            VARCHAR(500)     NOT NULL DEFAULT ''     COMMENT '标签（JSON数组，如["搞笑","萌宠"]）',
    `status`          TINYINT          NOT NULL DEFAULT 0      COMMENT '状态：0=待审核 1=已发布 2=已下架',
    `is_recommended`  TINYINT          NOT NULL DEFAULT 0      COMMENT '是否推荐：1=是 0=否',
    `view_count`      BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '播放量（冗余计数）',
    `like_count`      BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '点赞数（冗余计数）',
    `collect_count`   BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '收藏数（冗余计数）',
    `share_count`     BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '分享数（冗余计数）',
    `comment_count`   BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '评论数（冗余计数）',
    `published_at`    DATETIME         NULL                    COMMENT '发布时间',
    `created_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at`      DATETIME         NULL                    COMMENT '软删除',
    PRIMARY KEY (`id`),
    INDEX `idx_mp_status` (`mp_account_id`, `status`, `published_at` DESC),
    INDEX `idx_mp_category` (`mp_account_id`, `category_id`, `status`),
    INDEX `idx_entity` (`entity_id`),
    INDEX `idx_published` (`status`, `published_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频表';


-- ============================================================
-- 7. 用户表 (users)
--    按小程序维度隔离，一个用户在同一小程序下唯一
-- ============================================================
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `id`              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id`   BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `openid`          VARCHAR(64)      NOT NULL                COMMENT '微信 OpenID',
    `unionid`         VARCHAR(64)      NOT NULL DEFAULT ''     COMMENT '微信 UnionID（跨应用唯一）',
    `session_key`     VARCHAR(128)     NOT NULL DEFAULT ''     COMMENT '微信 SessionKey（加密存储）',
    `nickname`        VARCHAR(100)     NOT NULL DEFAULT ''     COMMENT '用户昵称',
    `avatar_url`      VARCHAR(500)     NOT NULL DEFAULT ''     COMMENT '头像 URL',
    `phone`           VARCHAR(20)      NOT NULL DEFAULT ''     COMMENT '手机号（加密存储）',
    `gender`          TINYINT          NOT NULL DEFAULT 0      COMMENT '性别：0=未知 1=男 2=女',
    `province`        VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '省份',
    `city`            VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '城市',
    `country`         VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '国家',
    `status`          TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=正常 0=禁用',
    `last_login_at`   DATETIME         NULL                    COMMENT '最近登录时间',
    `last_login_ip`   VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '最近登录 IP',
    `created_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
    `updated_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_mp_openid` (`mp_account_id`, `openid`),
    INDEX `idx_mp_phone` (`mp_account_id`, `phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';


-- ============================================================
-- 8. 评论表 (comments)
--    二级评论：parent_id=NULL 为一级评论，非 NULL 为回复
-- ============================================================
DROP TABLE IF EXISTS `comments`;
CREATE TABLE `comments` (
    `id`                BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id`     BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `video_id`          BIGINT UNSIGNED  NOT NULL                COMMENT '关联视频 ID',
    `user_id`           BIGINT UNSIGNED  NOT NULL                COMMENT '评论用户 ID',
    `parent_id`         BIGINT UNSIGNED  NULL                    COMMENT '父评论 ID（NULL=一级评论，非NULL=二级回复）',
    `root_id`           BIGINT UNSIGNED  NULL                    COMMENT '根评论 ID（一级评论ID，加速查询）',
    `reply_to_user_id`  BIGINT UNSIGNED  NULL                    COMMENT '被回复用户 ID',
    `content`           VARCHAR(1000)    NOT NULL                COMMENT '评论内容',
    `like_count`        INT UNSIGNED     NOT NULL DEFAULT 0      COMMENT '点赞数（冗余计数）',
    `reply_count`       INT UNSIGNED     NOT NULL DEFAULT 0      COMMENT '回复数（仅一级评论有效）',
    `status`            TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=正常 0=已删除',
    `is_top`            TINYINT          NOT NULL DEFAULT 0      COMMENT '是否置顶（仅一级评论）：1=是 0=否',
    `created_at`        DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论时间',
    `updated_at`        DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    INDEX `idx_video_root` (`mp_account_id`, `video_id`, `root_id`, `created_at`),
    INDEX `idx_video_parent` (`video_id`, `parent_id`, `created_at`),
    INDEX `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='评论表（支持二级回复）';


-- ============================================================
-- 9. 点赞表 (likes)
--    支持视频点赞 + 评论点赞，统一存储
-- ============================================================
DROP TABLE IF EXISTS `likes`;
CREATE TABLE `likes` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id` BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `user_id`       BIGINT UNSIGNED  NOT NULL                COMMENT '用户 ID',
    `target_type`   VARCHAR(20)      NOT NULL                COMMENT '点赞目标类型：video / comment',
    `target_id`     BIGINT UNSIGNED  NOT NULL                COMMENT '点赞目标 ID',
    `status`        TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=已点赞 0=已取消',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '点赞时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_user_target` (`mp_account_id`, `user_id`, `target_type`, `target_id`),
    INDEX `idx_target` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='点赞表';


-- ============================================================
-- 10. 收藏表 (favorites)
-- ============================================================
DROP TABLE IF EXISTS `favorites`;
CREATE TABLE `favorites` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id` BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `user_id`       BIGINT UNSIGNED  NOT NULL                COMMENT '用户 ID',
    `video_id`      BIGINT UNSIGNED  NOT NULL                COMMENT '视频 ID',
    `status`        TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=已收藏 0=已取消',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_user_video` (`mp_account_id`, `user_id`, `video_id`),
    INDEX `idx_user` (`mp_account_id`, `user_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收藏表';


-- ============================================================
-- 11. 分享记录表 (shares)
-- ============================================================
DROP TABLE IF EXISTS `shares`;
CREATE TABLE `shares` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id` BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `user_id`       BIGINT UNSIGNED  NOT NULL                COMMENT '分享用户 ID',
    `video_id`      BIGINT UNSIGNED  NOT NULL                COMMENT '视频 ID',
    `share_type`    VARCHAR(20)      NOT NULL DEFAULT ''     COMMENT '分享渠道：wechat_friend / wechat_moments / copy_link',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '分享时间',
    PRIMARY KEY (`id`),
    INDEX `idx_video` (`video_id`, `created_at`),
    INDEX `idx_user` (`user_id`, `created_at`),
    INDEX `idx_mp_date` (`mp_account_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分享记录表';


-- ============================================================
-- 12. 视频播放记录表 (view_logs)
--    颗粒度：每次播放行为一条记录（最细粒度原始数据）
--    数据量预估：日活用户 × 人均观看视频数，高并发场景考虑异步批量写入
-- ============================================================
DROP TABLE IF EXISTS `view_logs`;
CREATE TABLE `view_logs` (
    `id`              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id`   BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `entity_id`       BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '关联主体 ID（冗余，加速按主体统计）',
    `user_id`         BIGINT UNSIGNED  NULL                    COMMENT '用户 ID（游客为 NULL）',
    `video_id`        BIGINT UNSIGNED  NOT NULL                COMMENT '视频 ID',
    `category_id`     BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '视频所属分类 ID（快照，避免JOIN）',
    `watch_duration`  INT UNSIGNED     NOT NULL DEFAULT 0      COMMENT '观看时长（秒）',
    `is_complete`     TINYINT          NOT NULL DEFAULT 0      COMMENT '是否完整播放：1=是 0=否',
    `source`          VARCHAR(30)      NOT NULL DEFAULT ''     COMMENT '来源：recommend / category / search / share',
    `ip`              VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '客户端 IP',
    `province`        VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT 'IP 解析省份',
    `city`            VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT 'IP 解析城市',
    `device`          VARCHAR(100)     NOT NULL DEFAULT ''     COMMENT '设备型号（如 iPhone 15 Pro）',
    `os`              VARCHAR(30)      NOT NULL DEFAULT ''     COMMENT '操作系统：ios / android / ohos',
    `created_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '播放时间',
    PRIMARY KEY (`id`),
    INDEX `idx_mp_entity_date` (`mp_account_id`, `entity_id`, `created_at`),
    INDEX `idx_video_date` (`video_id`, `created_at`),
    INDEX `idx_user_date` (`user_id`, `created_at`),
    INDEX `idx_category_date` (`category_id`, `created_at`),
    INDEX `idx_province_date` (`province`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频播放记录表（最细粒度：每次播放一条）';


-- ============================================================
-- 13. 用户行为事件表 (action_logs)
--    颗粒度：每次用户行为一条记录（点赞/收藏/分享/页面浏览等）
--    注意：点赞/收藏状态变更 likes/favorites 表已有完整记录，
--    此表侧重埋点分析（页面来源、停留时长、点击路径），
--    避免与业务表重复，重点覆盖页面级行为
-- ============================================================
DROP TABLE IF EXISTS `action_logs`;
CREATE TABLE `action_logs` (
    `id`              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id`   BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `entity_id`       BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '关联主体 ID（冗余）',
    `user_id`         BIGINT UNSIGNED  NULL                    COMMENT '用户 ID（游客为 NULL）',
    `event_type`      VARCHAR(30)      NOT NULL                COMMENT '事件类型：page_view / click / search / share_trigger / follow',
    `target_type`     VARCHAR(30)      NOT NULL DEFAULT ''     COMMENT '目标类型：video / category / button / banner / page',
    `target_id`       BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '目标 ID',
    `page_path`       VARCHAR(200)     NOT NULL DEFAULT ''     COMMENT '触发页面路径（如 /pages/player/index）',
    `extra_data`      JSON             NULL                    COMMENT '附加数据（按钮名称、搜索关键词、停留时长等）',
    `ip`              VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '客户端 IP',
    `province`        VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '省份',
    `city`            VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '城市',
    `device`          VARCHAR(100)     NOT NULL DEFAULT ''     COMMENT '设备型号',
    `os`              VARCHAR(30)      NOT NULL DEFAULT ''     COMMENT '操作系统',
    `created_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '事件时间',
    PRIMARY KEY (`id`),
    INDEX `idx_mp_type_date` (`mp_account_id`, `event_type`, `created_at`),
    INDEX `idx_user_date` (`user_id`, `created_at`),
    INDEX `idx_page_date` (`page_path`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户行为事件表（侧重页面级行为埋点）';


-- ============================================================
-- 14. 视频每日统计表 (stats_video_daily)
--    颗粒度：每个视频 + 每天一条记录
--    数据来源：定时任务从 view_logs / likes / favorites / comments / shares 聚合
--    查询场景：视频排行榜、单视频趋势图、分类维度汇总
-- ============================================================
DROP TABLE IF EXISTS `stats_video_daily`;
CREATE TABLE `stats_video_daily` (
    `id`                    BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id`         BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `entity_id`             BIGINT UNSIGNED  NOT NULL                COMMENT '关联主体 ID',
    `video_id`              BIGINT UNSIGNED  NOT NULL                COMMENT '视频 ID',
    `category_id`           BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '分类 ID',
    `stat_date`             DATE             NOT NULL                COMMENT '统计日期',
    `view_count`            BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '播放次数',
    `view_users`            BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '播放人数（去重）',
    `avg_watch_duration`    DECIMAL(8,2)     NOT NULL DEFAULT 0.00   COMMENT '平均观看时长（秒）',
    `complete_count`        BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '完整播放次数',
    `like_count`            BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '点赞数（当日新增）',
    `collect_count`         BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '收藏数（当日新增）',
    `share_count`           BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '分享数（当日新增）',
    `comment_count`         BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '评论数（当日新增）',
    `created_at`            DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`            DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_video_date` (`video_id`, `stat_date`),
    INDEX `idx_mp_entity_date` (`mp_account_id`, `entity_id`, `stat_date`),
    INDEX `idx_category_date` (`category_id`, `stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='视频每日统计表（每视频每天一条）';


-- ============================================================
-- 15. 平台每日统计表 (stats_platform_daily)
--    颗粒度：每个小程序 + 每天一条记录
--    数据来源：定时任务聚合
--    查询场景：管理后台首页看板（总播放量、活跃用户、新增用户趋势）
-- ============================================================
DROP TABLE IF EXISTS `stats_platform_daily`;
CREATE TABLE `stats_platform_daily` (
    `id`              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id`   BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `entity_id`       BIGINT UNSIGNED  NOT NULL                COMMENT '关联主体 ID',
    `stat_date`       DATE             NOT NULL                COMMENT '统计日期',
    `total_views`     BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '总播放次数',
    `total_view_users` BIGINT UNSIGNED NOT NULL DEFAULT 0      COMMENT '播放用户数（去重）',
    `active_users`    BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '日活跃用户数（有任意行为的去重用户）',
    `new_users`       BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '日新增注册用户数',
    `total_users`     BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '累计注册用户数（截至当日）',
    `total_likes`     BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '总点赞数（当日新增）',
    `total_collects`  BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '总收藏数（当日新增）',
    `total_shares`    BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '总分享数（当日新增）',
    `total_comments`  BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '总评论数（当日新增）',
    `avg_online_minutes` DECIMAL(8,2)  NOT NULL DEFAULT 0.00   COMMENT '人均在线时长（分钟）',
    `created_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_mp_entity_date` (`mp_account_id`, `entity_id`, `stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='平台每日统计表（每主体每天一条）';


-- ============================================================
-- 16. 地域每日统计表 (stats_region_daily)
--    颗粒度：每个省份 + 每天一条记录
--    数据来源：定时任务从 view_logs / users 聚合
--    查询场景：用户地域分布热力图、各省播放量排行
-- ============================================================
DROP TABLE IF EXISTS `stats_region_daily`;
CREATE TABLE `stats_region_daily` (
    `id`              BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `mp_account_id`   BIGINT UNSIGNED  NOT NULL                COMMENT '关联小程序账号 ID',
    `entity_id`       BIGINT UNSIGNED  NOT NULL                COMMENT '关联主体 ID',
    `stat_date`       DATE             NOT NULL                COMMENT '统计日期',
    `province`        VARCHAR(50)      NOT NULL                COMMENT '省份',
    `view_count`      BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '播放次数',
    `view_users`      BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '播放用户数（去重）',
    `active_users`    BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '活跃用户数',
    `new_users`       BIGINT UNSIGNED  NOT NULL DEFAULT 0      COMMENT '新增用户数',
    `created_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`      DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_mp_province_date` (`mp_account_id`, `province`, `stat_date`),
    INDEX `idx_entity_date` (`entity_id`, `stat_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='地域每日统计表（每省份每天一条）';


-- ============================================================
-- 17. 管理员表 (admins)
--    管理后台账号（角色通过 RBAC 关联，不在此表存 role 字段）
-- ============================================================
DROP TABLE IF EXISTS `admins`;
CREATE TABLE `admins` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `username`      VARCHAR(64)      NOT NULL                COMMENT '登录用户名',
    `password`      VARCHAR(200)     NOT NULL                COMMENT '密码（bcrypt 加密）',
    `real_name`     VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '真实姓名',
    `avatar_url`    VARCHAR(500)     NOT NULL DEFAULT ''     COMMENT '头像 URL',
    `status`        TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=正常 0=禁用',
    `last_login_at` DATETIME         NULL                    COMMENT '最近登录时间',
    `last_login_ip` VARCHAR(50)      NOT NULL DEFAULT ''     COMMENT '最近登录 IP',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';


-- ============================================================
-- 18. 角色表 (roles)
--    预置角色：super_admin / admin / editor
--    支持自定义角色（按需扩展）
-- ============================================================
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name`          VARCHAR(50)      NOT NULL                COMMENT '角色名称（如：超级管理员、内容编辑）',
    `code`          VARCHAR(50)      NOT NULL                COMMENT '角色编码（如：super_admin、editor）',
    `description`   VARCHAR(200)     NOT NULL DEFAULT ''     COMMENT '角色描述',
    `is_system`     TINYINT          NOT NULL DEFAULT 0      COMMENT '是否系统预置角色：1=是（不可删除） 0=否',
    `status`        TINYINT          NOT NULL DEFAULT 1      COMMENT '状态：1=启用 0=禁用',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';


-- ============================================================
-- 19. 权限表 (permissions)
--    采用资源:动作的粒度，前端路由守卫 + 后端接口鉴权双重校验
--    code 格式: {module}:{action}  如 entity:view / video:delete
-- ============================================================
DROP TABLE IF EXISTS `permissions`;
CREATE TABLE `permissions` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name`          VARCHAR(100)     NOT NULL                COMMENT '权限名称（如：查看主体列表）',
    `code`          VARCHAR(100)     NOT NULL                COMMENT '权限编码（如：entity:view）',
    `module`        VARCHAR(50)      NOT NULL                COMMENT '所属模块：entity / mp_account / cdn / video / user / analytics / admin / role',
    `action`        VARCHAR(50)      NOT NULL                COMMENT '操作类型：view / create / edit / delete / audit',
    `description`   VARCHAR(200)     NOT NULL DEFAULT ''     COMMENT '权限描述',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_code` (`code`),
    INDEX `idx_module` (`module`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';


-- ============================================================
-- 20. 角色-权限关联表 (role_permissions)
--    多对多：一个角色拥有多个权限
-- ============================================================
DROP TABLE IF EXISTS `role_permissions`;
CREATE TABLE `role_permissions` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `role_id`       BIGINT UNSIGNED  NOT NULL                COMMENT '角色 ID',
    `permission_id` BIGINT UNSIGNED  NOT NULL                COMMENT '权限 ID',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_role_perm` (`role_id`, `permission_id`),
    INDEX `idx_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色-权限关联表';


-- ============================================================
-- 21. 管理员-角色关联表 (admin_roles)
--    多对多：一个管理员可拥有多个角色（权限叠加）
-- ============================================================
DROP TABLE IF EXISTS `admin_roles`;
CREATE TABLE `admin_roles` (
    `id`            BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
    `admin_id`      BIGINT UNSIGNED  NOT NULL                COMMENT '管理员 ID',
    `role_id`       BIGINT UNSIGNED  NOT NULL                COMMENT '角色 ID',
    `created_at`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE INDEX `uk_admin_role` (`admin_id`, `role_id`),
    INDEX `idx_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员-角色关联表';


-- ============================================================
-- 初始化数据
-- ============================================================

-- 默认管理员（密码：yanny2024，bcrypt 加密，生产环境务必更换）
INSERT INTO `admins` (`username`, `password`, `real_name`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '超级管理员');

-- 预置角色
INSERT INTO `roles` (`name`, `code`, `description`, `is_system`) VALUES
('超级管理员', 'super_admin', '拥有全部权限，不可删除', 1),
('管理员',     'admin',       '日常运营管理，除系统配置外的全部权限', 1),
('内容编辑',   'editor',      '仅视频内容管理与数据查看', 1);

-- 预置权限（8 个模块 × 标准 CRUD + 特殊操作）
INSERT INTO `permissions` (`name`, `code`, `module`, `action`) VALUES
-- 主体管理
('查看主体列表',   'entity:view',   'entity',      'view'),
('创建主体',       'entity:create', 'entity',      'create'),
('编辑主体',       'entity:edit',   'entity',      'edit'),
('删除主体',       'entity:delete', 'entity',      'delete'),
-- 小程序账号管理
('查看小程序账号', 'mp_account:view',   'mp_account', 'view'),
('创建小程序账号', 'mp_account:create', 'mp_account', 'create'),
('编辑小程序账号', 'mp_account:edit',   'mp_account', 'edit'),
('删除小程序账号', 'mp_account:delete', 'mp_account', 'delete'),
-- CDN 配置
('查看 CDN 配置',  'cdn:view',   'cdn', 'view'),
('创建 CDN 配置',  'cdn:create', 'cdn', 'create'),
('编辑 CDN 配置',  'cdn:edit',   'cdn', 'edit'),
('删除 CDN 配置',  'cdn:delete', 'cdn', 'delete'),
-- 视频管理
('查看视频列表',   'video:view',   'video', 'view'),
('创建/上传视频',  'video:create', 'video', 'create'),
('编辑视频信息',   'video:edit',   'video', 'edit'),
('删除视频',       'video:delete', 'video', 'delete'),
('审核视频',       'video:audit', 'video', 'audit'),
-- 用户管理
('查看用户列表',   'user:view',   'user', 'view'),
('禁用/启用用户',  'user:edit',   'user', 'edit'),
-- 数据看板
('查看数据看板',   'analytics:view', 'analytics', 'view'),
('导出数据报表',   'analytics:export', 'analytics', 'export'),
-- 管理员管理
('查看管理员列表', 'admin:view',   'admin', 'view'),
('创建管理员',     'admin:create', 'admin', 'create'),
('编辑管理员',     'admin:edit',   'admin', 'edit'),
('删除管理员',     'admin:delete', 'admin', 'delete'),
-- 角色管理
('查看角色列表',   'role:view',   'role', 'view'),
('创建角色',       'role:create', 'role', 'create'),
('编辑角色',       'role:edit',   'role', 'edit'),
('删除角色',       'role:delete', 'role', 'delete');

-- 超级管理员：拥有全部权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 1, id FROM `permissions`;

-- 管理员：除 admin/role 模块外的全部权限
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 2, id FROM `permissions` WHERE `module` NOT IN ('admin', 'role');

-- 内容编辑：仅视频管理 + 数据看板查看
INSERT INTO `role_permissions` (`role_id`, `permission_id`)
SELECT 3, id FROM `permissions` WHERE `module` IN ('video', 'analytics') AND `action` != 'delete';

-- 默认管理员关联超级管理员角色
INSERT INTO `admin_roles` (`admin_id`, `role_id`) VALUES (1, 1);
