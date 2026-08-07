package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server         ServerConfig         `yaml:"server"`
	MySQL          MySQLConfig          `yaml:"mysql"`
	Redis          RedisConfig          `yaml:"redis"`
	JWT            JWTConfig            `yaml:"jwt"`
	DynamicKey     string               `yaml:"dynamic_key"`
	AntiBruteForce AntiBruteForceConfig `yaml:"anti_brute_force"`
	Wechat         WechatConfig         `yaml:"wechat"`
	Qiniu          QiniuConfig          `yaml:"qiniu"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// MySQLConfig MySQL 配置
type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// DSN 生成 MySQL 连接串
func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// Addr Redis 地址
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret        string `yaml:"secret"`
	ExpireHours   int    `yaml:"expire_hours"`
	MpExpireHours int    `yaml:"mp_expire_hours"`
}

// WechatConfig 微信小程序配置
type WechatConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

// AntiBruteForceConfig 防爆破配置
type AntiBruteForceConfig struct {
	MaxAttempts   int `yaml:"max_attempts"`
	LockDuration  int `yaml:"lock_duration"`
	WindowMinutes int `yaml:"window_minutes"`
}

// QiniuConfig 七牛云配置
type QiniuConfig struct {
	AccessKey    string `yaml:"access_key"`
	SecretKey    string `yaml:"secret_key"`
	Bucket       string `yaml:"bucket"`
	Domain       string `yaml:"domain"`
	Region       string `yaml:"region"`
	CallbackURL  string `yaml:"callback_url"`
	AntiTheftKey string `yaml:"anti_theft_key"` // 时间戳防盗链 key
}

var AppConfig *Config

// Load 加载配置
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}
	AppConfig = cfg
	return nil
}
