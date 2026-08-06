package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"
	"yanny-service/internal/config"
	"yanny-service/internal/database"
	"yanny-service/internal/dto"

	"github.com/gin-gonic/gin"
)

// RequirePermission 权限校验中间件
// permCode 格式: "module:action" 如 "video:delete"
func RequirePermission(permCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := GetAdminID(c)
		if adminID == 0 {
			dto.ErrorWithStatus(c, 401, dto.ErrCodeUnauthorized, "请先登录")
			c.Abort()
			return
		}

		if !hasPermission(c, adminID, permCode) {
			dto.ErrorWithStatus(c, 403, dto.ErrCodeForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasPermission 判断管理员是否拥有指定权限
func hasPermission(c context.Context, adminID uint64, permCode string) bool {
	cacheKey := fmt.Sprintf("yanny:admin:%d:perms", adminID)

	// 先查 Redis
	isMember, err := database.Redis.SIsMember(c, cacheKey, permCode).Result()
	if err == nil {
		return isMember
	}

	// Redis 未命中，查 MySQL
	var count int64
	database.DB.Raw(`
		SELECT COUNT(1)
		FROM admin_roles ar
		JOIN role_permissions rp ON ar.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		JOIN roles r ON ar.role_id = r.id
		WHERE ar.admin_id = ? AND p.code = ? AND r.status = 1
	`, adminID, permCode).Scan(&count)

	has := count > 0
	if !has {
		// 异步加载权限到缓存（下次命中）
		go cacheAdminPermissions(adminID)
	}
	return has
}

// cacheAdminPermissions 加载管理员权限到 Redis
func cacheAdminPermissions(adminID uint64) {
	cacheKey := fmt.Sprintf("yanny:admin:%d:perms", adminID)

	var permCodes []string
	database.DB.Raw(`
		SELECT DISTINCT p.code
		FROM admin_roles ar
		JOIN role_permissions rp ON ar.role_id = rp.role_id
		JOIN permissions p ON rp.permission_id = p.id
		JOIN roles r ON ar.role_id = r.id
		WHERE ar.admin_id = ? AND r.status = 1
	`, adminID).Scan(&permCodes)

	if len(permCodes) > 0 {
		members := make([]interface{}, len(permCodes))
		for i, code := range permCodes {
			members[i] = code
		}
		database.Redis.Del(context.Background(), cacheKey)
		database.Redis.SAdd(context.Background(), cacheKey, members...)
		database.Redis.Expire(context.Background(), cacheKey, time.Duration(config.AppConfig.JWT.ExpireHours)*time.Hour)
	}
}

// LoadAdminPermissions 登录时加载权限到 Redis
func LoadAdminPermissions(adminID uint64) {
	cacheAdminPermissions(adminID)
}

// extractModuleAction 从请求路径+方法推断权限码（辅助用）
func extractModuleAction(path, method string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return ""
	}
	module := parts[1] // /api/v1/entity/xxx → entity

	actionMap := map[string]string{
		"GET":    "view",
		"POST":   "create",
		"PUT":    "edit",
		"PATCH":  "edit",
		"DELETE": "delete",
	}
	action, ok := actionMap[method]
	if !ok {
		action = "view"
	}
	return module + ":" + action
}
