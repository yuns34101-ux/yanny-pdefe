package database

import (
	"context"
	"log"
	"yanny-service/internal/config"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg config.RedisConfig) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx := context.Background()
	if err := Redis.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	log.Println("Redis 连接成功")
}
