package config

import (
	"os"
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	yamlContent := `
server:
  port: 8088
  mode: debug
mysql:
  host: 127.0.0.1
  port: 3306
  user: root
  password: ""
  database: yanny
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0
  pool_size: 10
jwt:
  secret: test-secret
  expire_hours: 2
  mp_expire_hours: 720
wechat:
  app_id: wx123
  app_secret: secret
qiniu:
  access_key: ak
  secret_key: sk
  bucket: test
  domain: test.example.com
  region: z0
`
	path := t.TempDir() + "/config.yaml"
	os.WriteFile(path, []byte(yamlContent), 0644)

	err := Load(path)
	if err != nil {
		t.Fatalf("Load() 失败: %v", err)
	}

	if AppConfig.Server.Port != 8088 {
		t.Errorf("Server.Port = %d, want 8088", AppConfig.Server.Port)
	}
	if AppConfig.MySQL.Database != "yanny" {
		t.Errorf("MySQL.Database = %s, want yanny", AppConfig.MySQL.Database)
	}
	if AppConfig.JWT.Secret != "test-secret" {
		t.Errorf("JWT.Secret = %s, want test-secret", AppConfig.JWT.Secret)
	}
	if AppConfig.Qiniu.Bucket != "test" {
		t.Errorf("Qiniu.Bucket = %s, want test", AppConfig.Qiniu.Bucket)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	err := Load("/non/existent/path.yaml")
	if err == nil {
		t.Error("Load() 应返回错误，文件不存在")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := t.TempDir() + "/invalid.yaml"
	os.WriteFile(path, []byte("{{{invalid yaml!!!}}"), 0644)

	err := Load(path)
	if err == nil {
		t.Error("Load() 应返回错误，YAML 格式非法")
	}
}

func TestMySQLConfig_DSN(t *testing.T) {
	cfg := MySQLConfig{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "pwd",
		Database: "testdb",
		Charset:  "utf8mb4",
	}
	dsn := cfg.DSN()
	expected := "root:pwd@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	if dsn != expected {
		t.Errorf("DSN = %s, want %s", dsn, expected)
	}
}

func TestRedisConfig_Addr(t *testing.T) {
	cfg := RedisConfig{Host: "127.0.0.1", Port: 6379}
	if cfg.Addr() != "127.0.0.1:6379" {
		t.Errorf("Addr = %s, want 127.0.0.1:6379", cfg.Addr())
	}
}
