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
// 签名算法：HMAC-SHA1(SecretKey, <rawURL>?e=<deadline>)，URL-safe Base64
func GeneratePrivateURL(rawURL string, deadline int64) string {
	cfg := config.AppConfig.Qiniu
	if cfg.SecretKey == "" || cfg.AccessKey == "" {
		return rawURL
	}

	urlToSign := fmt.Sprintf("%s?e=%d", rawURL, deadline)

	mac := hmac.New(sha1.New, []byte(cfg.SecretKey))
	mac.Write([]byte(urlToSign))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)

	return fmt.Sprintf("%s&token=%s:%s", urlToSign, cfg.AccessKey, encodedSign)
}

// GenerateVideoURL 生成视频播放签名 URL（默认 6 小时有效）
func GenerateVideoURL(rawURL string) string {
	deadline := time.Now().Add(6 * time.Hour).Unix()
	return GeneratePrivateURL(rawURL, deadline)
}

// GenerateCoverURL 生成封面图签名 URL（默认 6 小时有效）
func GenerateCoverURL(rawURL string) string {
	deadline := time.Now().Add(6 * time.Hour).Unix()
	return GeneratePrivateURL(rawURL, deadline)
}
