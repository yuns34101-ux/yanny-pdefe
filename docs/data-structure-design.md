# Yanny 数据结构设计文档

> Version: v1.0.0 | Date: 2026-08-06 | DB: MySQL 8.0 + Redis

---

## 一、ER 关系图

```
┌──────────────┐         ┌──────────────────┐         ┌──────────────┐
│   entities   │         │ entity_mp_bindings│         │  mp_accounts │
│   (主体表)    │ 1     N │   (主体-小程序绑定) │ N     1 │  (小程序账号)  │
│──────────────│─────────│──────────────────│─────────│──────────────│
│ id           │         │ id                │         │ id           │
│ name         │         │ entity_id  ───────┘         │ app_id       │
│ logo_url     │         │ mp_account_id ──────────────│ app_secret   │
│ description  │         │ is_default       │         │ app_name     │
│ status       │         └──────────────────┘         │ status       │
└──┬────┬──────┘                                      └──────┬───────┘
   │    │ 1                                                   │ 1
   │    │                                                     │
   │    │ N                                                   │ N
   │    └──────────────────────┐                              │
   │                           │                              │
   │    ┌─────────────────┐    │                              │
   │    │ video_categories│    │                              │
   │    │  (视频分类)      │    │                              │
   │    │────────────────│    │                              │
   │    │ entity_id ──────┘    │                              │
   │    │ mp_account_id ──────────────────────────────────────│
   │    │ name           │    │                              │
   │    └────────┬───────┘    │                              │
   │             │            │                              │
   │  N          │            │                              │
   │             │            │                              │
   │  ┌──────────┴──────────┐ │                              │
   │  │       videos        │ │                              │
   │  │      (视频表)        │ │                              │
   │  │─────────────────────│ │                              │
   │  │ entity_id ──────────┘ │                              │
   │  │ mp_account_id ────────┘                              │
   │  │ category_id ─────────┘                              │
   │  │ title               │                               │
   │  │ cover_url           │          ┌──────────────────┐ │
   │  │ video_url           │          │   cdn_configs    │ │
   │  │ *_count (冗余计数)   │          │   (CDN配置)       │ │
   │  └──┬──┬──┬──┬─────────┘          │──────────────────│ │
   │     │  │  │  │                    │ mp_account_id ───┘ │
   │     │  │  │  │                    │ provider          │
   │     │  │  │  │                    │ bucket/domain     │
   │     │  │  │  │                    └──────────────────┘
   │     │  │  │  │
   │     │  │  │  └──────────────────────────────────────┐
   │     │  │  │                                         │
   │     │  │  │  ┌─────────────────┐                    │
   │     │  │  │  │    comments     │                    │
   │     │  │  │  │    (评论表)      │                    │
   │     │  │  │  │────────────────│                    │
   │     │  │  └──│ video_id       │                    │
   │     │  │     │ user_id ───────┼──┐                 │
   │     │  │     │ parent_id      │  │                 │
   │     │  │     │ root_id        │  │                 │
   │     │  │     │ content        │  │                 │
   │     │  │     └─────────────────┘  │                 │
   │     │  │                          │                 │
   │     │  │  ┌─────────────────┐     │                 │
   │     │  │  │     likes       │     │                 │
   │     │  │  │    (点赞表)      │     │                 │
   │     │  └──│ target_type     │     │                 │
   │     │     │ target_id       │     │                 │
   │     │     │ user_id ────────┼─────┤                 │
   │     │     └─────────────────┘     │                 │
   │     │                             │                 │
   │     │  ┌─────────────────┐        │                 │
   │     │  │   favorites     │        │                 │
   │     │  │   (收藏表)       │        │                 │
   │     └──│ video_id        │        │                 │
   │        │ user_id ────────┼────────┤                 │
   │        └─────────────────┘        │                 │
   │                                   │                 │
   │  ┌─────────────────┐              │                 │
   │  │    shares       │              │                 │
   │  │   (分享记录)     │              │                 │
   │  └──│ video_id      │              │                 │
   │     │ user_id ──────┼──────────────┘                 │
   │     └─────────────────┘                              │
   │                                                      │
   │  ┌──────────────────────────────────────────┐        │
   │  │              埋点统计层                     │        │
   │  │──────────────────────────────────────────│        │
   │  │  view_logs (每次播放一条，最细粒度)          │        │
   │  │    └── video_id / user_id / category_id   │        │
   │  │    └── province / city / ip              │        │
   │  │                                          │        │
   │  │  action_logs (页面级行为，侧重路径分析)     │        │
   │  │    └── event_type / page_path / target    │        │
   │  │    └── province / city / ip              │        │
   │  │                                          │        │
   │  │  ┌─ stats_video_daily (每视频·每天)        │        │
   │  │  ├─ stats_platform_daily (每主体·每天)     │        │
   │  │  └─ stats_region_daily (每省份·每天)       │        │
   │  └──────────────────────────────────────────┘        │
   │                                                      │
   │  ┌─────────────────┐                                 │
   │  │     users       │                                 │
   │  │    (用户表)      │                                 │
   │  │─────────────────│                                 │
   │  │ openid          │                                 │
   │  │ mp_account_id ──┘                                 │
   │  │ nickname        │      ┌─────────────────┐        │
   │  │ phone           │      │     admins      │        │
   │  │ avatar_url      │      │   (管理员表)     │        │
   │  └─────────────────┘      │ username        │        │
   │                           │ password(bcrypt) │        │
   │                           │ role            │        │
   │                           └─────────────────┘        │
```
```

---

## 二、核心业务表关系说明

### 2.1 主体 ↔ 小程序账号（多对多）

```
entities ──< entity_mp_bindings >── mp_accounts
```

- 一个运营主体（如"扬昵文化"）可同时运营多个小程序
- 一个小程序也可关联多个主体（联合运营场景）
- `is_default` 标记默认绑定关系，小程序首页默认展示默认主体的信息

### 2.2 视频归属链

```
video_categories.entity_id  → entities.id        （分类属于主体，所有权）
video_categories.mp_account_id → mp_accounts.id  （分类在哪个小程序中展示，生效范围）
videos.entity_id            → entities.id        （视频归属哪个主体）
videos.mp_account_id        → mp_accounts.id     （视频在哪个小程序展示）
videos.category_id          → video_categories.id（视频分类）
```

- 分类是主体级数据：同一主体在不同小程序中共享同一套分类体系
- 视频同时持有 `entity_id` + `category_id`，可直接 `WHERE entity_id = ? AND category_id = ?` 查某主体某分类下的视频，无需 JOIN

### 2.3 用户隔离

```
users.mp_account_id → mp_accounts.id
```

- **用户数据按小程序物理隔离**：同一微信用户在不同小程序下是两条独立记录
- `openid` 在小程序间不同，`unionid` 可用于跨应用识别（需微信开放平台绑定）

### 2.4 评论二级结构

```
comments.parent_id → comments.id  （自引用）
comments.root_id   → comments.id  （冗余加速）
```

| 层级 | parent_id | root_id | reply_to_user_id |
|------|-----------|---------|------------------|
| 一级评论 | NULL | NULL（写入后回填自己的 id） | NULL |
| 二级回复 | 父评论 id | 根评论 id（一级评论 id） | 被回复用户 id |

- `root_id` 是冗余字段，避免递归查询：查某个视频的所有一级评论 + 其下所有回复，只需 `WHERE root_id IN (一级评论ID列表)`
- 评论列表分页：先分页查一级评论（`parent_id IS NULL`），再根据一级评论 ID 列表批量拉回复

### 2.7 管理后台 RBAC 权限模型

```
admins ──< admin_roles >── roles ──< role_permissions >── permissions
```

| 表 | 职责 | 关联关系 |
|----|------|----------|
| `admins` | 管理员账号（仅存登录信息，不存角色） | `admin_roles` 多对多关联 `roles` |
| `roles` | 角色定义（支持预置 + 自定义） | `role_permissions` 多对多关联 `permissions` |
| `permissions` | 权限原子（`module:action` 格式） | 最小粒度，不可再分 |
| `admin_roles` | 管理员-角色绑定 | UNIQUE(admin_id, role_id) |
| `role_permissions` | 角色-权限绑定 | UNIQUE(role_id, permission_id) |

**权限编码格式**：`{module}:{action}`，如 `video:delete`、`analytics:view`。

**8 个权限模块**：

| 模块 | 权限 |
|------|------|
| `entity` | view / create / edit / delete |
| `mp_account` | view / create / edit / delete |
| `cdn` | view / create / edit / delete |
| `video` | view / create / edit / delete / audit |
| `user` | view / edit（禁用/启用） |
| `analytics` | view / export |
| `admin` | view / create / edit / delete |
| `role` | view / create / edit / delete |

**3 个预置角色**：

| 角色 | 权限范围 |
|------|----------|
| `super_admin` | 全部 29 项权限 |
| `admin` | 除 admin/role 模块外的全部权限（日常运营） |
| `editor` | 仅 video + analytics 模块的非删除权限 |

**鉴权流程**：
```
请求 → JWT 中间件（解析 admin_id）
     → 查 admin_roles → roles → role_permissions → permissions
     → 判断当前请求 {module}:{action} 是否在权限集合中
     → 前端路由守卫 + 后端接口中间件双重校验
```

**缓存策略**：管理员登录后，将其权限集合缓存到 Redis：
- Key: `yanny:admin:{admin_id}:perms`
- 类型: Set（members = 权限 code 列表）
- TTL: 2h（与 JWT 同步过期）

videos 表的 `view_count / like_count / collect_count / share_count / comment_count` 是**冗余计数器**，目的是避免高频 COUNT 查询。更新策略：

| 计数器 | 写入触发 | 一致性策略 |
|--------|----------|------------|
| `view_count` | view_logs 插入后异步 +1 | Redis INCR → 定时批量回写 MySQL |
| `like_count` | likes 状态变更时同步 +1/-1 | 事务内原子更新 |
| `collect_count` | favorites 状态变更时同步 +1/-1 | 事务内原子更新 |
| `share_count` | shares 插入后 +1 | 事务内原子更新 |
| `comment_count` | comments 插入后 +1 | 事务内原子更新 |

### 2.6 埋点统计分层设计

埋点数据按颗粒度分为三层：

```
Layer 0 — 业务表（最粗）
  likes / favorites / shares / comments
  粒度：每条记录 = 一次明确的操作
  用途：业务查询（我点赞了哪些视频）+ 实时计数

Layer 1 — 事件流水表（最细）
  view_logs（每次播放一条）/ action_logs（每次页面行为一条）
  粒度：每条记录 = 一次事件
  用途：明细查询、数据回刷、adhoc 分析
  注意：点赞/收藏/分享/评论的计数走 Layer 0 业务表，
        action_logs 侧重页面级行为（page_view / click / search），
        避免与业务表重复采集

Layer 2 — 聚合统计表（天级汇总）
  stats_video_daily / stats_platform_daily / stats_region_daily
  粒度：每条记录 = 一个维度 × 一天
  用途：管理后台报表、趋势图、排行榜
  生成：定时任务（每日凌晨）从 Layer 0 + Layer 1 聚合
```

**三张聚合表的查询场景：**

| 表 | 唯一键 | 典型查询 | 数据量估算 |
|----|--------|----------|------------|
| `stats_video_daily` | video_id + stat_date | 单视频 30 天趋势、分类播放排行、热门视频 Top100 | 视频数 × 天数 |
| `stats_platform_daily` | mp_account_id + entity_id + stat_date | 管理后台首页看板：日活/新增/播放量趋势 | 主体数 × 天数 |
| `stats_region_daily` | mp_account_id + province + stat_date | 全国省份热力图、某省用户播放量变化 | 34 省 × 天数 |

---

## 三、Redis 缓存设计

### 3.1 Key 命名规范

```
{namespace}:{entity}:{id}:{sub...}
yanny:{mp_id}:{resource}:{id}
```

### 3.2 缓存策略表

| 缓存对象 | Key 模式 | 类型 | TTL | 说明 |
|----------|----------|------|-----|------|
| 主体信息 | `yanny:entity:{entity_id}` | Hash | 1h | 全字段 |
| 小程序配置 | `yanny:mp:{mp_id}:config` | Hash | 1h | AppID/名称/状态 |
| 主体-小程序绑定 | `yanny:mp:{mp_id}:entities` | Set | 1h | 该小程序绑定的主体 ID 集合 |
| 视频分类列表 | `yanny:entity:{entity_id}:mp:{mp_id}:categories` | ZSet | 30min | score=sort_order |
| 视频详情 | `yanny:video:{video_id}` | Hash | 30min | 全字段 |
| 视频列表（按分类） | `yanny:mp:{mp_id}:videos:{category_id}:{page}` | ZSet | 10min | score=published_at 时间戳 |
| 视频列表（推荐） | `yanny:mp:{mp_id}:videos:recommend:{page}` | ZSet | 10min | score=published_at |
| 视频播放量 | `yanny:video:{video_id}:views` | String | 永久 | INCR 实时计数 |
| 用户信息 | `yanny:mp:{mp_id}:user:{user_id}` | Hash | 30min | |
| 用户 Token | `yanny:mp:{mp_id}:token:{user_id}` | String | 2h | JWT |
| 评论列表（一级） | `yanny:video:{video_id}:comments:p{page}` | List | 10min | 存 JSON 序列化 |
| 评论回复 | `yanny:comment:{comment_id}:replies` | List | 10min | |
| 用户点赞集合 | `yanny:mp:{mp_id}:user:{user_id}:likes` | Set | 1h | target_type:target_id |
| 用户收藏集合 | `yanny:mp:{mp_id}:user:{user_id}:favorites` | ZSet | 1h | member=video_id, score=时间戳 |
| 每日 UV（HyperLogLog） | `yanny:mp:{mp_id}:uv:{date}` | HyperLogLog | 30d | 日活去重 |
| 每日 PV | `yanny:mp:{mp_id}:pv:{date}` | String | 30d | INCR |
| IP 地域统计 | `yanny:mp:{mp_id}:region:{date}` | Hash | 30d | field=province, value=count |
| 平台看板缓存 | `yanny:mp:{mp_id}:dashboard:30d` | String | 5min | 30 天趋势 JSON（预热缓存，避免每次查聚合表） |

### 3.3 缓存更新策略

```
写操作：
  1. 更新 MySQL（主）
  2. 删除对应 Redis Key（Cache-Aside）
  3. 下次读取时自动重建缓存

高并发计数（播放量）：
  1. Redis INCR (实时)
  2. 定时任务（每 5 分钟）：批量读取 Redis 计数 → UPDATE MySQL
  3. 读取时：优先 Redis；降级 MySQL

用户互动状态（点赞/收藏）：
  1. 写入 MySQL + 更新 Redis Set/ZSet（同步双写）
  2. 前端展示"已点赞/已收藏"状态走 Redis Set SISMEMBER
```

---

## 四、索引设计要点

| 表 | 核心索引 | 覆盖查询场景 |
|----|----------|-------------|
| `videos` | `idx_mp_status (mp_account_id, status, published_at DESC)` | 小程序首页视频列表分页 |
| `videos` | `idx_mp_category (mp_account_id, category_id, status)` | 按分类筛选视频 |
| `video_categories` | `idx_entity_mp_sort (entity_id, mp_account_id, sort_order DESC)` | 查某主体在某小程序的分类列表（有序） |
| `comments` | `idx_video_root (mp_account_id, video_id, root_id, created_at)` | 视频评论列表（含二级回复） |
| `comments` | `idx_video_parent (video_id, parent_id, created_at)` | 查某条评论的回复 |
| `likes` | `uk_user_target (mp_account_id, user_id, target_type, target_id)` | 唯一约束 + 快速查用户是否点赞 |
| `favorites` | `uk_user_video (mp_account_id, user_id, video_id)` | 唯一约束 + 用户收藏列表 |
| `view_logs` | `idx_mp_entity_date (mp_account_id, entity_id, created_at)` | 管理后台按主体统计播放量 |
| `view_logs` | `idx_province_date (province, created_at)` | 按省份+时间范围统计 |
| `action_logs` | `idx_mp_type_date (mp_account_id, event_type, created_at)` | 按事件类型+时间范围统计 |
| `stats_video_daily` | `uk_video_date (video_id, stat_date)` | 单视频趋势图、排行榜 |
| `stats_platform_daily` | `uk_mp_entity_date (mp_account_id, entity_id, stat_date)` | 平台看板趋势 |
| `stats_region_daily` | `uk_mp_province_date (mp_account_id, province, stat_date)` | 地域分布热力图 |
| `users` | `uk_mp_openid (mp_account_id, openid)` | 登录时按 OpenID 查用户 |

---

## 五、扩展性预留

| 预留点 | 设计 | 说明 |
|--------|------|------|
| `entities.extra` | JSON 类型 | 运营方自定义字段，无需 DDL 变更 |
| `videos.tags` | JSON 数组字符串 | 灵活打标签，后续可建倒排索引 |
| `action_logs.extra_data` | JSON 类型 | 不同事件类型携带不同的附加信息 |
| `mp_accounts` 支持多平台 | `app_id/app_secret` 字段 | 当前仅微信，预留字节/支付宝/抖音小程序 |
| `cdn_configs.provider` | 枚举值 | 当前仅七牛，预留阿里云/腾讯云 |
