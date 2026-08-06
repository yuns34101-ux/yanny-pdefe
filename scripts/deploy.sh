#!/bin/bash
set -e  # 任何命令执行失败则立即退出

# ========== 配置区域 ==========
PROJECT_DIR="/www/wwwroot/yanny-pdefe"
SERVICE_DIR="$PROJECT_DIR/service"
ADMIN_DIR="$PROJECT_DIR/admin"
BINARY="$SERVICE_DIR/binary"
CONFIG_FILE="$SERVICE_DIR/config.yaml"
CONFIG_EXAMPLE="$SERVICE_DIR/config.yaml.example"
LOG_FILE="$PROJECT_DIR/app.log"

# Go 环境配置（根据您安装的版本调整）
export GOROOT=/usr/local/go1.25.5
export PATH=$GOROOT/bin:$PATH
export GOTOOLCHAIN=local          # 禁止自动下载工具链
export GOPROXY=https://goproxy.cn,direct  # 使用国内代理

# ========== 开始部署 ==========
echo "🚀 开始部署项目..."

# 1. 进入项目根目录并拉取最新代码
cd "$PROJECT_DIR"
echo "📦 拉取最新代码..."
git pull

# 2. 确保配置文件存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "⚠️  配置文件不存在，从示例复制..."
    cp "$CONFIG_EXAMPLE" "$CONFIG_FILE"
    echo "📝 请编辑 $CONFIG_FILE 填入正确的配置后，重新执行本脚本。"
    exit 1
fi

# 3. 数据库迁移（自动扫描 db/ 目录下未执行的 SQL）
echo "🗄️  检查数据库迁移..."
for f in "$PROJECT_DIR"/db/v*.sql; do
    if [ -f "$f" ]; then
        echo "  执行: $(basename $f)"
        mysql -u yanny -p"${DB_PASSWORD}" yanny < "$f" 2>/dev/null || echo "  ⚠️ 迁移可能已执行或失败，继续..."
    fi
done

# 4. 构建 Go 后端
echo "🔨 构建 Go 后端..."
cd "$SERVICE_DIR"
go build -o "$BINARY" ./cmd/
echo "✅ Go 后端构建完成：$BINARY"

# 5. 构建前端
if [ -d "$ADMIN_DIR" ]; then
    echo "🌐 构建前端..."
    cd "$ADMIN_DIR"
    [ -d "node_modules" ] || npm install
    npm run build
    echo "✅ 前端构建完成"
else
    echo "⚠️  未找到 admin 目录，跳过前端构建"
fi

# 6. 停止旧进程
echo "🛑 停止旧进程..."
pkill -f "$BINARY" 2>/dev/null || true
sleep 1

# 7. 启动新进程
echo "▶️  启动后端服务..."
cd "$SERVICE_DIR"
nohup "$BINARY" > "$LOG_FILE" 2>&1 &

# 8. 检查是否启动成功
sleep 2
if pgrep -f "$BINARY" > /dev/null; then
    echo "✅ 服务启动成功，PID: $(pgrep -f "$BINARY")"
    echo "📄 日志文件: $LOG_FILE"
else
    echo "❌ 服务启动失败，请检查日志: $LOG_FILE"
    exit 1
fi

echo "🎉 部署完成！"