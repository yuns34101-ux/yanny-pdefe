package middleware

import (
	"net/http"
	"strings"
	"yanny-service/internal/dto"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders 安全响应头中间件
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Next()
	}
}

// ContentTypeCheck 请求 Content-Type 校验
func ContentTypeCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			c.Next()
			return
		}

		ct := c.GetHeader("Content-Type")

		// 文件上传
		if strings.HasPrefix(ct, "multipart/form-data") {
			if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
				dto.ErrorWithStatus(c, http.StatusBadRequest, dto.ErrCodeParamInvalid, "请求格式错误")
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// JSON 请求
		if strings.HasPrefix(ct, "application/json") || ct == "" {
			c.Next()
			return
		}

		dto.ErrorWithStatus(c, http.StatusUnsupportedMediaType, dto.ErrCodeParamInvalid, "不支持的 Content-Type，仅接受 application/json")
		c.Abort()
	}
}

// RequestSizeLimit 请求体大小限制
func RequestSizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// SanitizeInput 基础输入清洗（XSS 简单防护）
func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 简单检测 SQL 注入特征（实际防护由 GORM 参数化查询保证）
		for _, values := range c.Request.URL.Query() {
			for _, v := range values {
				if containsSQLInjection(v) {
					dto.ErrorWithStatus(c, http.StatusBadRequest, dto.ErrCodeParamInvalid, "请求参数包含非法字符")
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}

func containsSQLInjection(s string) bool {
	dangerous := []string{
		"UNION SELECT", "DROP TABLE", "DROP DATABASE",
		"INSERT INTO", "DELETE FROM", "1=1", "1'='1",
		"--", ";--", "xp_cmdshell",
	}
	upper := strings.ToUpper(s)
	for _, d := range dangerous {
		if strings.Contains(upper, d) {
			return true
		}
	}
	return false
}
