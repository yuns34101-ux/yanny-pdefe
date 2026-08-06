package qiniu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"

	"yanny-service/internal/config"
)

// GeneratePrivateURL 七牛私有空间签名 URL
// rawURL: 原始文件 URL（如 https://cdn.example.com/videos/xxx.mp4）
// deadline: 过期时间戳（Unix 秒）
func GeneratePrivateURL(rawURL string, deadline int64) string {
	cfg := config.AppConfig.Qiniu
	if cfg.SecretKey == "" || cfg.AccessKey == "" {
		return rawURL
	}

	// 1. 构造待签名字符串：<rawURL>?e=<deadline>
	urlToSign := fmt.Sprintf("%s?e=%d", rawURL, deadline)

	// 2. HMAC-SHA1 签名
	mac := hmac.New(sha1.New, []byte(cfg.SecretKey))
	mac.Write([]byte(urlToSign))
	sign := mac.Sum(nil)

	// 3. URL-safe Base64（带填充）
	encodedSign := base64.URLEncoding.EncodeToString(sign)

	return fmt.Sprintf("%s&token=%s:%s", urlToSign, cfg.AccessKey, encodedSign)
}

// GenerateVideoURL 生成视频播放签名 URL（默认 2 小时有效）
func GenerateVideoURL(rawURL string) string {
	deadline := time.Now().Add(2 * time.Hour).Unix()
	return GeneratePrivateURL(rawURL, deadline)
}

// GenerateCoverURL 生成封面图签名 URL（默认 24 小时有效）
func GenerateCoverURL(rawURL string) string {
	deadline := time.Now().Add(24 * time.Hour).Unix()
	return GeneratePrivateURL(rawURL, deadline)
}
