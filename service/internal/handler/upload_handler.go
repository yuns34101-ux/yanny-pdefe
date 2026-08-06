package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"yanny-service/internal/config"
	"yanny-service/internal/dto"

	"github.com/gin-gonic/gin"
)

// qiniuPutPolicy 七牛上传策略
type qiniuPutPolicy struct {
	Scope    string `json:"scope"`
	Deadline int64  `json:"deadline"`
}

// GetUploadToken 获取七牛上传 Token
func GetUploadToken(c *gin.Context) {
	var req struct {
		FileType string `json:"file_type" binding:"required"` // video / image
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "请指定 file_type: video 或 image")
		return
	}

	cfg := config.AppConfig.Qiniu
	if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Bucket == "" {
		dto.Error(c, dto.ErrCodeInternal, "七牛云配置未设置（access_key/secret_key/bucket）")
		return
	}

	// 上传凭证 1 小时有效
	deadline := time.Now().Add(1 * time.Hour).Unix()

	policy := qiniuPutPolicy{
		Scope:    cfg.Bucket,
		Deadline: deadline,
	}
	policyJSON, _ := json.Marshal(policy)

	// 1. 对 PutPolicy JSON 做 URL-safe Base64（无填充）
	encodedPolicy := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(policyJSON)

	// 2. HMAC-SHA1 签名
	mac := hmac.New(sha1.New, []byte(cfg.SecretKey))
	mac.Write([]byte(encodedPolicy))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(sign)

	// 3. 组装 Token：<AK>:<encodedSign>:<encodedPolicy>
	token := fmt.Sprintf("%s:%s:%s", cfg.AccessKey, encodedSign, encodedPolicy)

	dto.Success(c, gin.H{
		"token":       token,
		"domain":      cfg.Domain,
		"bucket":      cfg.Bucket,
		"region":      cfg.Region,
		"upload_host": qiniuUploadHost(cfg.Region),
		"deadline":    deadline,
	})
}

// qiniuUploadHost 根据区域返回上传域名
func qiniuUploadHost(region string) string {
	switch region {
	case "z0":
		return "https://up.qiniup.com"
	case "z1":
		return "https://up-z1.qiniup.com"
	case "z2":
		return "https://up-z2.qiniup.com"
	case "na0":
		return "https://up-na0.qiniup.com"
	case "as0":
		return "https://up-as0.qiniup.com"
	default:
		return "https://up.qiniup.com"
	}
}

