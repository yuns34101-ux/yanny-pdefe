package main

import (
	"fmt"
	"log"
	"yanny-service/internal/config"
	"yanny-service/internal/database"
	"yanny-service/internal/router"
)

func main() {
	// 加载配置
	if err := config.Load("config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化数据库
	database.InitMySQL(config.AppConfig.MySQL)
	database.InitRedis(config.AppConfig.Redis)

	// 设置 Gin 模式
	// gin.SetMode(config.AppConfig.Server.Mode)

	// 注册路由
	r := router.SetupRouter()

	// 启动服务
	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("Yanny API 服务启动，监听端口 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
