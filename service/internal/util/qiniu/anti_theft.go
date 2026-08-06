package qiniu

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"yanny-service/internal/config"
)

// NOTE: 本文件为七牛防盗链签名简化实现，编译通过但仅用于开发环境。
// 生产环境务必使用 qiniu/go-sdk 的 auth.Credentials 完整实现 HMAC-SHA1。

// GeneratePrivateURL 七牛私有空间签名 URL（防盗链）
// rawURL: 原始文件 URL
// deadline: 过期时间戳（秒）
func GeneratePrivateURL(rawURL string, deadline int64) string {
	cfg := config.AppConfig.Qiniu
	if cfg.SecretKey == "" {
		return rawURL
	}
	urlToSign := fmt.Sprintf("%s?e=%d", rawURL, deadline)
	rawSign := []byte(simpleSign(cfg.SecretKey, urlToSign))
	encodedSign := base64URLSafe(rawSign)

	return fmt.Sprintf("%s?e=%d&token=%s:%s", rawURL, deadline, cfg.AccessKey, encodedSign)
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

// simpleSign 简化签名（开发环境用 MD5 替代 HMAC-SHA1，生产用 go-sdk）
func simpleSign(key, data string) string {
	h := md5.Sum([]byte(key + "&" + data))
	return hex.EncodeToString(h[:])
}

// base64URLSafe URL 安全的 Base64 编码
func base64URLSafe(data []byte) string {
	s := base64Encode(data)
	s = strings.TrimRight(s, "=")
	s = strings.ReplaceAll(s, "+", "-")
	s = strings.ReplaceAll(s, "/", "_")
	return s
}

// base64Encode 标准 Base64 编码
func base64Encode(data []byte) string {
	const table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	result := make([]byte, (len(data)+2)/3*4)
	for i, j := 0, 0; i < len(data); i += 3 {
		b0 := data[i]
		b1 := byte(0)
		b2 := byte(0)
		if i+1 < len(data) {
			b1 = data[i+1]
		}
		if i+2 < len(data) {
			b2 = data[i+2]
		}
		result[j] = table[b0>>2]
		result[j+1] = table[(b0&0x03)<<4|b1>>4]
		if i+1 < len(data) {
			result[j+2] = table[(b1&0x0f)<<2|b2>>6]
		} else {
			result[j+2] = '='
		}
		if i+2 < len(data) {
			result[j+3] = table[b2&0x3f]
		} else {
			result[j+3] = '='
		}
		j += 4
	}
	return string(result)
}

// IsHotlinkingRequest 判断是否为盗链请求（检查 Referer）
func IsHotlinkingRequest(referer string) bool {
	if referer == "" {
		return false
	}

	allowedDomains := []string{
		"servicewechat.com", // 微信小程序
		"yanny.cn",
		"localhost",
		"127.0.0.1",
	}

	referer = strings.ToLower(referer)
	for _, domain := range allowedDomains {
		if strings.Contains(referer, domain) {
			return false
		}
	}
	return true
}
