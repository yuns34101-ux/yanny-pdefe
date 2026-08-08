package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

	deadline := time.Now().Add(1 * time.Hour).Unix()

	policy := qiniuPutPolicy{
		Scope:      cfg.Bucket,
		Deadline:   deadline,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize)}`,
	}
	policyJSON, _ := json.Marshal(policy)

	encodedPolicy := base64.URLEncoding.EncodeToString(policyJSON)
	mac := hmac.New(sha1.New, []byte(cfg.SecretKey))
	mac.Write([]byte(encodedPolicy))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)

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

// buildQiniuUploadTokenForKey 生成限定 key 的上传凭证（服务端中转用，防 key 篡改）
func buildQiniuUploadTokenForKey(key string) string {
	cfg := config.AppConfig.Qiniu
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return ""
	}
	deadline := time.Now().Add(10 * time.Minute).Unix()
	policy := qiniuPutPolicy{
		Scope:      cfg.Bucket + ":" + key,
		Deadline:   deadline,
		ReturnBody: `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize)}`,
	}
	policyJSON, _ := json.Marshal(policy)
	encodedPolicy := base64.URLEncoding.EncodeToString(policyJSON)
	mac := hmac.New(sha1.New, []byte(cfg.SecretKey))
	mac.Write([]byte(encodedPolicy))
	encodedSign := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("%s:%s:%s", cfg.AccessKey, encodedSign, encodedPolicy)
}

// uploadFileToQiniu 服务端将文件流直传七牛，返回完整 URL
func uploadFileToQiniu(token, key string, file io.Reader, fileSize int64, mimeType string) (string, error) {
	cfg := config.AppConfig.Qiniu
	uploadHost := qiniuUploadHost(cfg.Region)

	// 构建 multipart form
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("token", token)
	w.WriteField("key", key)
	part, err := w.CreateFormFile("file", key)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	w.Close()

	// POST 到七牛
	req, err := http.NewRequest("POST", uploadHost, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("七牛返回 %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析七牛响应
	var result struct {
		Key  string `json:"key"`
		Hash string `json:"hash"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析七牛响应失败: %v", err)
	}

	baseURL := qiniu.EnsureProtocol(cfg.Domain)
	return fmt.Sprintf("%s/%s", baseURL, result.Key), nil
}

// isValidImageBytes 校验文件头魔数，防止 Content-Type 伪造
func isValidImageBytes(data []byte) bool {
	n := len(data)
	if n < 4 {
		return false
	}
	// JPEG: FF D8 FF
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return true
	}
	// PNG: 89 50 4E 47
	if data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return true
	}
	// GIF: 47 49 46 38
	if data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x38 {
		return true
	}
	// WEBP: 52 49 46 46 ... 57 45 42 50 (at offset 8)
	if n >= 12 && data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50 {
		return true
	}
	return false
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

// CheckMediaAsset 上传前查重
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

// ConfirmMediaAsset 上传成功后确认落库
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

// ========== 服务端中转上传（小程序不允许直传七牛，走后端代理） ==========

// MpUploadAvatar 小程序头像上传（服务端校验 → 七牛）
// 安全措施：需登录、仅允许图片、限制 2MB、校验文件魔数、服务端随机 key
func MpUploadAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		dto.Error(c, dto.ErrCodeParamInvalid, "请选择文件")
		return
	}
	defer file.Close()

	// 限制 2MB
	if header.Size > 2<<20 {
		dto.Error(c, dto.ErrCodeParamInvalid, "图片大小不能超过 2MB")
		return
	}

	// 校验 MIME
	contentType := header.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true, "image/png": true, "image/gif": true, "image/webp": true,
	}
	if !allowedTypes[contentType] {
		dto.Error(c, dto.ErrCodeParamInvalid, "仅支持 JPG/PNG/GIF/WEBP 格式")
		return
	}

	// 读入内存 buffer（multipart.File 不支持 Seek，先用 buffer 校验魔数再上传）
	fileBytes, err := io.ReadAll(io.LimitReader(file, 2<<20))
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, "读取文件失败")
		return
	}

	// 校验文件头魔数（防 Content-Type 伪造）
	if !isValidImageBytes(fileBytes) {
		dto.Error(c, dto.ErrCodeParamInvalid, "文件不是有效图片")
		return
	}

	// 服务端生成随机 key
	randomBytes := make([]byte, 6)
	for i := range randomBytes {
		randomBytes[i] = byte(time.Now().UnixNano()>>uint(i*8)) ^ 0xAA
	}
	ext := ".jpg"
	switch contentType {
	case "image/png":
		ext = ".png"
	case "image/gif":
		ext = ".gif"
	case "image/webp":
		ext = ".webp"
	}
	key := fmt.Sprintf("images/avatars/%d_%x%s", time.Now().Unix(), randomBytes, ext)

	// 生成限定 key 的上传 token
	uploadToken := buildQiniuUploadTokenForKey(key)
	if uploadToken == "" {
		dto.Error(c, dto.ErrCodeInternal, "生成上传凭证失败")
		return
	}

	// 上传到七牛
	qiniuURL, err := uploadFileToQiniu(uploadToken, key, bytes.NewReader(fileBytes), int64(len(fileBytes)), contentType)
	if err != nil {
		dto.Error(c, dto.ErrCodeInternal, "上传失败："+err.Error())
		return
	}

	dto.Success(c, gin.H{"url": qiniuURL})
}
