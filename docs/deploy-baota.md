# Yanny 宝塔面板部署流程

> 适用环境：CentOS 7+ / Ubuntu 18+，宝塔面板 8.x  
> 更新日期：2026-08-06

---

## 一、环境概览

| 组件 | 版本 | 端口 | 用途 |
|------|------|------|------|
| MySQL | 8.0 | 3306 | 主数据库 |
| Redis | 7.x | 6379 | 缓存 / 限流 / 黑名单 |
| Go | 1.21+ | — | API 编译 |
| Nginx | 1.24+ | 80/443 | 反向代理 + 静态资源 |
| Node.js | 18+ | — | 管理后台构建 |

---

## 二、宝塔面板安装基础环境

### 2.1 安装宝塔

```bash
# CentOS
yum install -y wget && wget -O install.sh https://download.bt.cn/install/install_6.0.sh && bash install.sh ed8484bec

# Ubuntu / Debian
wget -O install.sh https://download.bt.cn/install/install-ubuntu_6.0.sh && bash install.sh ed8484bec
```

安装完成后记录面板地址、用户名、密码。

### 2.2 软件商店安装

在宝塔面板「软件商店」中一键安装：

```
✅ MySQL 8.0
✅ Redis 7.x
✅ Nginx 1.24+
✅ PM2 管理器（Node 进程守护）
```

### 2.3 安装 Go（SSH 终端）

```bash
# 下载 Go 1.22
cd /usr/local
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
tar -xzf go1.22.5.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
echo 'export GOPATH=/home/go' >> /etc/profile
echo 'export GOPROXY=https://goproxy.cn,direct' >> /etc/profile
source /etc/profile

# 验证
go version
```

---

## 三、数据库初始化

### 3.1 创建数据库

宝塔面板 → 数据库 → 添加数据库：

```
数据库名: yanny
用户名:   yanny
密码:     <生成随机密码，记录下来>
字符集:   utf8mb4
```

### 3.2 导入 DDL

方式一：宝塔面板 → 数据库 → phpMyAdmin → 导入 `db/v1.0.0_ddl.sql`

方式二：SSH 命令行：

```bash
mysql -u yanny -p yanny < /www/wwwroot/yanny/db/v1.0.0_ddl.sql
```

### 3.3 修改管理员密码

```bash
mysql -u yanny -p yanny

# 生成新密码的 bcrypt hash
# 在本地 Go 环境中运行：
# go run -exec '' <<EOF
# package main; import ("fmt"; "golang.org/x/crypto/bcrypt")
# func main() { h,_:=bcrypt.GenerateFromPassword([]byte("your_password"), bcrypt.DefaultCost); fmt.Println(string(h)) }
# EOF

UPDATE admins SET password = '<bcrypt_hash>' WHERE username = 'admin';
```

---

## 四、项目部署

### 4.1 创建站点目录

宝塔面板 → 网站 → 添加站点：

```
域名:     yanny-api.example.com（API）/ yanny-admin.example.com（后台）
根目录:   /www/wwwroot/yanny
FTP:      不创建
数据库:   不创建（已手动创建）
```

### 4.2 上传代码

```bash
# 在服务器上
cd /www/wwwroot
git clone https://github.com/yuns34101-ux/yanny-pdefe.git yanny
```

或者通过宝塔面板「文件」上传 zip 包解压。

### 4.3 配置 API 服务

```bash
cd /www/wwwroot/yanny/service
cp config.yaml.example config.yaml
vi config.yaml
```

修改配置：

```yaml
mysql:
  host: 127.0.0.1
  user: yanny
  password: "<数据库密码>"
  database: yanny

redis:
  host: 127.0.0.1
  password: ""  # 宝塔默认无密码，可在 Redis 设置中配置

jwt:
  secret: "<生成随机 32 位字符串>"

dynamic_key: "<生成随机 16 位字符串>"
```

### 4.4 编译并启动 API

```bash
cd /www/wwwroot/yanny/service

# 编译
go build -o yanny-service ./cmd/

# 方式一：直接运行（测试）
./yanny-service

# 方式二：PM2 守护（推荐）
pm2 start yanny-service --name yanny-api
pm2 save
pm2 startup
```

---

## 五、Nginx 反向代理

### 5.1 API 反向代理

宝塔面板 → 网站 → yanny-api 站点 → 配置文件：

```nginx
server {
    listen 80;
    server_name yanny-api.example.com;

    # API 反向代理
    location / {
        proxy_pass http://127.0.0.1:8088;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超时设置
        proxy_connect_timeout 30s;
        proxy_read_timeout 60s;
        proxy_send_timeout 30s;

        # 文件上传大小限制
        client_max_body_size 100m;
    }
}
```

### 5.2 管理后台静态站点

```bash
# 在服务器上构建管理后台
cd /www/wwwroot/yanny/admin
npm install
npm run build
# 产出在 dist/ 目录
```

Nginx 配置：

```nginx
server {
    listen 80;
    server_name yanny-admin.example.com;
    root /www/wwwroot/yanny/admin/dist;
    index index.html;

    # SPA 路由回退
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 代理（管理后台 /api 请求转发到后端）
    location /api/ {
        proxy_pass http://127.0.0.1:8088;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 静态资源缓存
    location /assets/ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

### 5.3 SSL 证书

宝塔面板 → 网站 → SSL → Let's Encrypt 一键申请。

两个站点都需要开启「强制 HTTPS」。

---

## 六、Redis 安全配置（重要）

宝塔面板 → 软件商店 → Redis → 设置：

```
# 修改 requirepass
requirepass <生成随机密码>

# 绑定本地 IP
bind 127.0.0.1

# 禁用危险命令
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command KEYS ""
```

修改后 `config.yaml` 中的 Redis 密码需要同步更新。

---

## 七、防火墙配置

宝塔面板 → 安全 → 放行端口：

```
✅ 80   (HTTP)
✅ 443  (HTTPS)
✅ 3306 (MySQL - 仅本地 127.0.0.1)
✅ 6379 (Redis - 仅本地 127.0.0.1)
❌ 8088 (API - 不对外开放，通过 Nginx 代理)
```

> **重要**：MySQL 和 Redis 不要开放公网端口。API 只通过 Nginx 反向代理访问，不直接暴露 8088。

---

## 八、管理后台首次登录

### 8.1 访问后台

```
https://yanny-admin.example.com
```

### 8.2 登录凭据

```
用户名: admin
密码:   <静态密码> + <当前时分HHMM>

示例: 当前时间 15:06 → 输入 "admin1234561506"
```

静态密码为数据库初始化时 bcrypt 加密的密码，登录后可在后台修改。

### 8.3 首次配置顺序

1. **主体管理** → 创建运营主体
2. **小程序账号** → 添加微信小程序 AppID
3. **绑定管理** → 绑定主体与小程序
4. **分类管理** → 创建视频分类（萌宠/健身/知识等）
5. **CDN 配置** → 配置七牛云

---

## 九、小程序发布

### 9.1 构建

```bash
cd /www/wwwroot/yanny/mp
npm install

# 微信小程序
npm run build:mp-weixin
# 产出在 dist/build/mp-weixin/

# 或使用 HBuilderX 打开项目直接发行
```

### 9.2 修改 API 地址

发布前将 `mp/src/utils/request.js` 中的 `BASE_URL` 改为生产环境地址：

```js
const BASE_URL = 'https://yanny-api.example.com/api/v1/mp'
```

### 9.3 命令行上传体验版（推荐，脱离微信开发者工具）

基于官方 `miniprogram-ci`，可在服务器/CI 环境直接完成上传，无需打开微信开发者工具。

**首次准备**：

1. 微信公众平台 → 开发管理 → 开发设置 → 小程序代码上传密钥 → 生成并下载私钥文件，保存为 `mp/upload-key.pem`（该文件已在 `.gitignore` 中忽略，禁止提交）。
2. 同一页面把执行上传的服务器 IP 加入 IP 白名单。
3. `mp/src/manifest.json` 中的 `versionName` 是本次上传的版本号，发布前先手动改好。

**执行上传**：

```bash
# 先出预览二维码（不写入体验版，用于快速验证）
node scripts/upload-mp.js --preview --desc "预览说明"

# 正式上传体验版
node scripts/upload-mp.js --desc "版本说明"
```

上传成功后，微信公众平台后台「版本管理 → 体验版本」会出现新版本，扫码即可体验。

参数说明：`--key` 指定私钥路径（默认 `mp/upload-key.pem`），`--robot` 指定小程序 CI 机器人编号（默认 1）。

### 9.4 微信开发者工具（备选方式）

1. 打开微信开发者工具
2. 导入 `dist/build/mp-weixin/`
3. AppID 配置为正式小程序 AppID
4. 上传 → 提交审核

---

## 十、定时任务

宝塔面板 → 计划任务 → 添加：

| 任务 | 周期 | 脚本 |
|------|------|------|
| 日统计聚合 | 每天 01:00 | `curl -X POST http://127.0.0.1:8088/api/v1/admin/stats/aggregate` (待实现) |
| Redis 播放量同步 | 每 5 分钟 | 内置于 API 服务，无需额外配置 |
| 日志清理 | 每天 03:00 | `find /www/wwwroot/yanny/logs -mtime +30 -delete` |

---

## 十一、监控与运维

### 11.1 PM2 进程监控

```bash
pm2 list          # 查看进程状态
pm2 logs yanny-api # 查看日志
pm2 restart yanny-api # 重启服务
```

### 11.2 数据库备份

宝塔面板 → 计划任务 → 添加：

```
任务类型: 备份数据库
数据库:   yanny
周期:     每天 04:00
保留份数: 7
```

### 11.3 健康检查

```bash
curl https://yanny-api.example.com/health
# 返回 {"status":"ok"} 表示正常
```

---

## 十二、故障排查

| 问题 | 检查项 |
|------|--------|
| API 502 | `pm2 logs yanny-api` 查看错误日志，检查 MySQL/Redis 连接 |
| 登录失败(动态口令) | 检查服务器时间 `date`，确保与北京时间一致；`ntpdate ntp.aliyun.com` |
| IP 被锁定 | `redis-cli KEYS "yanny:login:locked:*"` → `DEL` 对应 key |
| 数据库连接失败 | 检查 `config.yaml` 中 MySQL 用户名密码，确认 `bind-address = 127.0.0.1` |
| 视频上传失败 | 检查七牛云 AccessKey/Bucket/Domain 配置 |

---

## 十三、常用命令速查

```bash
# API 服务
pm2 restart yanny-api       # 重启
pm2 logs yanny-api --lines 50 # 查看最近 50 行日志
cd /www/wwwroot/yanny/service && go build -o yanny-service ./cmd/ && pm2 restart yanny-api  # 更新代码后重新编译部署

# 数据库
mysql -u yanny -p yanny     # 登录数据库

# Redis
redis-cli -a <密码>          # 登录 Redis
redis-cli -a <密码> KEYS "yanny:*" # 查看所有缓存 key

# 管理后台
cd /www/wwwroot/yanny/admin && npm run build  # 重新构建

# 证书续期
/usr/bin/certbot renew --quiet
```

---

**版本记录**

| 版本 | 日期 | 变更 |
|------|------|------|
| 1.0 | 2026-08-06 | 初版部署文档 |
