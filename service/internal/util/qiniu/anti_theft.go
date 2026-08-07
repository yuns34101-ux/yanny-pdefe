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
// rawURL: 原始文件 URL（若已带签名参数会自动剥离，避免重复签名）
// process: 处理指令（如 "imageView2/1/w/300/h/300"），空字符串表示不处理
// ttl: 签名有效期
func SignURL(rawURL, process string, ttl time.Duration) string {
	// 防御：剥离已有的 query string（历史脏数据可能已经是签名过的完整 URL）
	baseURL := rawURL
	if idx := strings.Index(rawURL, "?"); idx >= 0 {
		baseURL = rawURL[:idx]
	}

	cfg := config.AppConfig.Qiniu
	if cfg.SecretKey == "" || cfg.AccessKey == "" {
		if process != "" {
			return baseURL + "?" + process
		}
		return baseURL
	}

	// 构造待签名 URL：<base>?<process>&e=<deadline>
	url := baseURL
	if process != "" {
		// 处理参数之间用 | 分隔，最终拼成一个 ?process&e=deadline 的格式
		url = baseURL + "?" + process
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

// StripQuery 剥离 URL 中的签名/处理参数，仅保留纯净路径部分
// 用于写库前清洗：避免把前端回显的已签名 URL 误存为原始值
func StripQuery(rawURL string) string {
	if idx := strings.Index(rawURL, "?"); idx >= 0 {
		return rawURL[:idx]
	}
	return rawURL
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
