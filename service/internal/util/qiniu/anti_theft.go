package qiniu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"yanny-service/internal/config"
)

// SignURL 对七牛私有空间 URL 签名，支持多媒体处理参数
// rawURL: 原始文件 URL
// process: 处理指令（如 "imageView2/1/w/300/h/300"），空字符串表示不处理
// ttl: 签名有效期
func SignURL(rawURL, process string, ttl time.Duration) string {
	cfg := config.AppConfig.Qiniu
	if cfg.SecretKey == "" || cfg.AccessKey == "" {
		if process != "" {
			return rawURL + "?" + process
		}
		return rawURL
	}

	// 构造待签名 URL：<base>?<process>&e=<deadline>
	url := rawURL
	if process != "" {
		// 处理参数之间用 | 分隔，最终拼成一个 ?process&e=deadline 的格式
		url = rawURL + "?" + process
	}

	deadline := time.Now().Add(ttl).Unix()
	urlToSign := fmt.Sprintf("%s&e=%d", url, deadline)

	mac := hmac.New(sha1.New, []byte(cfg.SecretKey))
	mac.Write([]byte(urlToSign))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)

	return fmt.Sprintf("%s&token=%s:%s", urlToSign, cfg.AccessKey, encodedSign)
}

// 通用图片处理规则：自动旋转 + 限宽 750px + 质量 85%
const ImageProcessRule = "imageMogr2/auto-orient/thumbnail/750x/quality/85"

// SignImageURL 对图片 URL 签名并附加处理规则（logo、icon、封面等通用）
func SignImageURL(rawURL string) string {
	rawURL = EnsureProtocol(rawURL)
	ttl := time.Duration(config.AppConfig.Qiniu.GetMediaTTL()) * time.Hour
	return SignURL(rawURL, ImageProcessRule, ttl)
}

// EnsureProtocol 补全缺少的 https:// 前缀
func EnsureProtocol(rawURL string) string {
	if rawURL == "" {
		return rawURL
	}
	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		return rawURL
	}
	return "https://" + rawURL
}

// BuildProcessRule 拼接多个处理指令（用 | 连接）
func BuildProcessRule(ops ...string) string {
	var parts []string
	for _, op := range ops {
		if op != "" {
			parts = append(parts, op)
		}
	}
	return strings.Join(parts, "|")
}
