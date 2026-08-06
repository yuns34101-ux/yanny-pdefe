package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"yanny-service/internal/config"
	"yanny-service/internal/database"
	"yanny-service/internal/dto"
	"yanny-service/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 载荷
type Claims struct {
	AdminID     uint64 `json:"admin_id,omitempty"`
	UserID      uint64 `json:"user_id,omitempty"`
	MpAccountID uint64 `json:"mp_account_id,omitempty"`
	Username    string `json:"username,omitempty"`
	Type        string `json:"type"` // "admin" / "mp_user"
	jwt.RegisteredClaims
}

// GenerateAdminToken 生成管理后台 JWT
func GenerateAdminToken(adminID uint64, username string) (string, error) {
	cfg := config.AppConfig.JWT
	claims := Claims{
		AdminID:  adminID,
		Username: username,
		Type:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// GenerateMpToken 生成小程序端 JWT
func GenerateMpToken(userID, mpAccountID uint64) (string, error) {
	cfg := config.AppConfig.JWT
	claims := Claims{
		UserID:      userID,
		MpAccountID: mpAccountID,
		Type:        "mp_user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.MpExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// parseToken 解析 JWT
func parseToken(tokenStr string) (*Claims, error) {
	cfg := config.AppConfig.JWT
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名方法: %v", t.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("token 无效")
	}
	return claims, nil
}

// AdminAuthMiddleware 管理后台鉴权中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			dto.ErrorWithStatus(c, http.StatusUnauthorized, dto.ErrCodeUnauthorized, "请先登录")
			c.Abort()
			return
		}

		claims, err := parseToken(tokenStr)
		if err != nil || claims.Type != "admin" {
			dto.ErrorWithStatus(c, http.StatusUnauthorized, dto.ErrCodeUnauthorized, "token 无效或已过期")
			c.Abort()
			return
		}

		// 检查 token 是否在登出黑名单中
		if IsTokenBlacklisted(tokenStr) {
			dto.ErrorWithStatus(c, http.StatusUnauthorized, dto.ErrCodeUnauthorized, "token 已失效，请重新登录")
			c.Abort()
			return
		}

		// 校验管理员是否存在且状态正常
		var count int64
		if config.AppConfig != nil {
			// database.DB 在 main.go 中初始化后才可用，测试中可能为 nil
			if database.DB != nil {
				database.DB.Model(&model.Admin{}).Where("id = ? AND status = 1", claims.AdminID).Count(&count)
				if count == 0 {
					dto.ErrorWithStatus(c, http.StatusUnauthorized, dto.ErrCodeUnauthorized, "账号已被禁用或不存在")
					c.Abort()
					return
				}
			}
		}

		c.Set("admin_id", claims.AdminID)
		c.Set("username", claims.Username)
		c.Set("token_type", "admin")
		c.Next()
	}
}

// MpAuthMiddleware 小程序用户鉴权中间件（必须登录）
func MpAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			dto.ErrorWithStatus(c, http.StatusUnauthorized, dto.ErrCodeUnauthorized, "请先登录")
			c.Abort()
			return
		}

		claims, err := parseToken(tokenStr)
		if err != nil || claims.Type != "mp_user" {
			dto.ErrorWithStatus(c, http.StatusUnauthorized, dto.ErrCodeUnauthorized, "token 无效或已过期")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("mp_account_id", claims.MpAccountID)
		c.Set("token_type", "mp_user")
		c.Next()
	}
}

// MpAuthOptional 小程序用户可选鉴权（游客可访问，登录用户设置 user_id）
func MpAuthOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := extractToken(c)
		if tokenStr == "" {
			c.Next()
			return
		}

		claims, err := parseToken(tokenStr)
		if err != nil || claims.Type != "mp_user" {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("mp_account_id", claims.MpAccountID)
		c.Set("token_type", "mp_user")
		c.Next()
	}
}

// GetAdminID 从上下文获取管理员 ID
func GetAdminID(c *gin.Context) uint64 {
	id, _ := c.Get("admin_id")
	if v, ok := id.(uint64); ok {
		return v
	}
	return 0
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) uint64 {
	id, _ := c.Get("user_id")
	if v, ok := id.(uint64); ok {
		return v
	}
	return 0
}

// GetMpAccountID 从上下文获取小程序账号 ID
func GetMpAccountID(c *gin.Context) uint64 {
	id, _ := c.Get("mp_account_id")
	if v, ok := id.(uint64); ok {
		return v
	}
	return 0
}

func extractToken(c *gin.Context) string {
	// 优先从 Authorization Header 取
	auth := c.GetHeader("Authorization")
	if auth != "" {
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
		return auth
	}
	// 从 query 参数取（兼容 SSE/WebSocket）
	return c.Query("token")
}

// FormatAdminID 格式化管理员 ID 为字符串
func FormatAdminID(id uint64) string {
	return strconv.FormatUint(id, 10)
}
