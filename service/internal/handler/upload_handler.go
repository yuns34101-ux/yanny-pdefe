package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"yanny-service/internal/config"
	"yanny-service/internal/dto"
	"yanny-service/internal/model"
	"yanny-service/internal/repository"
	"yanny-service/internal/util/qiniu"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	gormErrors "gorm.io/gorm"
)

// qiniuPutPolicy 七牛上传策略
type qiniuPutPolicy struct {
	Scope      string `json:"scope"`
	Deadline   int64  `json:"deadline"`
	ReturnBody string `json:"returnBody"`
}

// mysqlDuplicateEntryErrNo MySQL 唯一键冲突错误码
const mysqlDuplicateEntryErrNo = 1062

// buildQiniuUploadToken 生成七牛上传凭证（admin/mp 共用）
func buildQiniuUploadToken() (gin.H, error) {
	cfg := config.AppConfig.Qiniu
	if cfg.AccessKey == "" || cfg.SecretKey == "" || cfg.Bucket == "" {
		return nil, errors.New("七牛云配置未设置（access_key/secret_key/bucket）")
	}

	// 上传凭证 1 小时有效
	deadline := time.Now().Add(1 * time.Hour).Unix()

	policy := qiniuPutPolicy{
		Scope:      cfg.Bucket,
		Deadline:   deadline,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize)}`,
	}
	policyJSON, _ := json.Marshal(policy)

	// 1. 对 PutPolicy JSON 做 URL-safe Base64（带填充，与官方 SDK 一致）
	encodedPolicy := base64.URLEncoding.EncodeToString(policyJSON)

	// 2. HMAC-SHA1 签名
	mac := hmac.New(sha1.New, []byte(cfg.SecretKey))
	mac.Write([]byte(encodedPolicy))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)

	// 3. 组装 Token：<AK>:<encodedSign>:<encodedPolicy>
	token := fmt.Sprintf("%s:%s:%s", cfg.AccessKey, encodedSign, encodedPolicy)

	return gin.H{
		"token":       token,
		"domain":      cfg.Domain,
		"bucket":      cfg.Bucket,
		"region":      cfg.Region,
		"upload_host": qiniuUploadHost(cfg.Region),
		"deadline":    deadline,
	}, nil
}

// GetUploadToken 获取七牛上传 Token（管理后台）
func GetUploadToken(c *gin.Context) {
	var req struct {
		FileType string `json:"file_type" binding:"required"` // video / image
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "请指定 file_type: video 或 image")
		return
	}

	result, err := buildQiniuUploadToken()
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, err.Error())
		return
	}
	dto.Success(c, result)
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

// CheckMediaAsset 上传前查重（基于前端预计算的哈希，仅用于加速判断，不作为最终权威依据）
func CheckMediaAsset(c *gin.Context) {
	var req struct {
		MpAccountID uint64 `json:"mp_account_id"`
		ClientHash  string `json:"client_hash" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "请提供 client_hash")
		return
	}

	asset, err := repository.FindMediaAssetByClientHash(req.MpAccountID, req.ClientHash)
	if err != nil {
		if errors.Is(err, gormErrors.ErrRecordNotFound) {
			dto.Success(c, gin.H{"exists": false})
			return
		}
		dto.Error(c, dto.ErrCodeInternal, "查重失败："+err.Error())
		return
	}

	dto.Success(c, gin.H{"exists": true, "url": asset.URL})
}

// ConfirmMediaAsset 上传成功后确认落库（content_hash 取自七牛 returnBody 的 etag，服务端权威校验）
func ConfirmMediaAsset(c *gin.Context) {
	var req struct {
		MpAccountID uint64 `json:"mp_account_id"`
		DirType     string `json:"dir_type" binding:"required"`
		ObjectKey   string `json:"object_key" binding:"required"`
		ContentHash string `json:"content_hash" binding:"required"`
		ClientHash  string `json:"client_hash"`
		FileSize    uint64 `json:"file_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "参数无效："+err.Error())
		return
	}

	cfg := config.AppConfig.Qiniu
	baseURL := qiniu.EnsureProtocol(cfg.Domain)
	url := fmt.Sprintf("%s/%s", baseURL, req.ObjectKey)

	asset := &model.MediaAsset{
		MpAccountID: req.MpAccountID,
		DirType:     req.DirType,
		ContentHash: req.ContentHash,
		ClientHash:  req.ClientHash,
		ObjectKey:   req.ObjectKey,
		URL:         url,
		FileSize:    req.FileSize,
	}

	if err := repository.CreateMediaAsset(asset); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryErrNo {
			existing, findErr := repository.FindMediaAssetByContentHash(req.MpAccountID, req.ContentHash)
			if findErr != nil {
				dto.Error(c, dto.ErrCodeInternal, "查询已存在记录失败："+findErr.Error())
				return
			}
			dto.Success(c, gin.H{"url": existing.URL, "reused": true})
			return
		}
		dto.Error(c, dto.ErrCodeInternal, "保存媒体资源失败："+err.Error())
		return
	}

	dto.Success(c, gin.H{"url": asset.URL, "reused": false})
}

