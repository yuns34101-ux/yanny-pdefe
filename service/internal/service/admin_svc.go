package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"yanny-service/internal/config"
	"yanny-service/internal/database"
	"yanny-service/internal/middleware"
	"yanny-service/internal/model"
	"yanny-service/internal/repository"
	"yanny-service/internal/util"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AdminLogin 管理员登录（静态密码 + 动态口令 + 防爆破）
func AdminLogin(username, passwordInput, ip string) (string, *model.Admin, error) {
	// 检查 IP 是否被锁定
	if isLoginLocked("ip", ip) {
		return "", nil, errors.New("登录失败次数过多，IP 已被暂时锁定，请稍后再试")
	}

	// 检查用户名是否被锁定
	lockKey := fmt.Sprintf("yanny:login:locked:user:%s", username)
	if isLoginLocked("user", username) {
		return "", nil, errors.New("该账号登录失败次数过多，已被暂时锁定，请稍后再试")
	}

	// 分离静态密码和动态口令
	staticPassword, dynamicCode := util.SplitPasswordAndCode(passwordInput)

	// 校验动态口令
	cfg := config.AppConfig
	if cfg.DynamicKey != "" && !util.ValidateDynamicCode(cfg.DynamicKey, dynamicCode, time.Now()) {
		recordLoginFailure("user", username, lockKey)
		recordLoginFailure("ip", ip, fmt.Sprintf("yanny:login:locked:ip:%s", ip))
		return "", nil, errors.New("动态口令错误")
	}

	// 查找管理员
	admin, err := repository.FindAdminByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			recordLoginFailure("user", username, lockKey)
			return "", nil, errors.New("用户名或密码错误")
		}
		return "", nil, err
	}

	if admin.Status == 0 {
		return "", nil, errors.New("账号已被禁用")
	}

	// 校验静态密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(staticPassword)); err != nil {
		recordLoginFailure("user", username, lockKey)
		recordLoginFailure("ip", ip, fmt.Sprintf("yanny:login:locked:ip:%s", ip))
		return "", nil, errors.New("用户名或密码错误")
	}

	// 登录成功，清除失败计数
	database.Redis.Del(context.Background(),
		fmt.Sprintf("yanny:login:fail:user:%s", username),
		fmt.Sprintf("yanny:login:fail:ip:%s", ip),
		lockKey,
		fmt.Sprintf("yanny:login:locked:ip:%s", ip),
	)

	// 生成 JWT
	token, err := middleware.GenerateAdminToken(admin.ID, admin.Username)
	if err != nil {
		return "", nil, err
	}

	// 更新登录信息
	_ = repository.UpdateAdminLoginInfo(admin.ID, ip)

	// 加载权限到 Redis
	middleware.LoadAdminPermissions(admin.ID)

	return token, admin, nil
}

// recordLoginFailure 记录登录失败并检查是否需要锁定
func recordLoginFailure(idType, id, lockKey string) {
	cfg := config.AppConfig
	if cfg.AntiBruteForce.MaxAttempts <= 0 {
		return // 防爆破未启用
	}

	ctx := context.Background()
	failKey := fmt.Sprintf("yanny:login:fail:%s:%s", idType, id)
	count, _ := database.Redis.Incr(ctx, failKey).Result()
	if count == 1 {
		database.Redis.Expire(ctx, failKey, time.Duration(cfg.AntiBruteForce.WindowMinutes)*time.Minute)
	}

	if int(count) >= cfg.AntiBruteForce.MaxAttempts {
		database.Redis.Set(ctx, lockKey, "1", time.Duration(cfg.AntiBruteForce.LockDuration)*time.Second)
		database.Redis.Del(ctx, failKey)
	}
}

// isLoginLocked 检查是否被锁定
func isLoginLocked(idType, id string) bool {
	lockKey := fmt.Sprintf("yanny:login:locked:%s:%s", idType, id)
	exists, _ := database.Redis.Exists(context.Background(), lockKey).Result()
	return exists > 0
}

// GetAdminInfo 获取管理员详情（含角色和权限）
func GetAdminInfo(adminID uint64) (map[string]interface{}, error) {
	admin, err := repository.FindAdminWithRoles(adminID)
	if err != nil {
		return nil, err
	}

	perms, _ := repository.FindAdminPermissions(adminID)

	result := map[string]interface{}{
		"id":         admin.ID,
		"username":   admin.Username,
		"real_name":  admin.RealName,
		"avatar_url": admin.AvatarURL,
		"status":     admin.Status,
		"roles":      admin.Roles,
		"perms":      perms,
	}
	return result, nil
}

// ========== 管理员 CRUD ==========

// CreateAdmin 创建管理员
func CreateAdmin(username, password, realName string, roleIDs []uint64) (*model.Admin, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	admin := &model.Admin{
		Username: username,
		Password: string(hashed),
		RealName: realName,
		Status:   1,
	}
	if err := repository.CreateAdmin(admin); err != nil {
		return nil, err
	}
	if len(roleIDs) > 0 {
		repository.AssignRoles(admin.ID, roleIDs)
	}
	return admin, nil
}

// UpdateAdmin 更新管理员
func UpdateAdmin(adminID uint64, realName string, status *int8, roleIDs []uint64) error {
	updates := map[string]interface{}{}
	if realName != "" {
		updates["real_name"] = realName
	}
	if status != nil {
		updates["status"] = *status
	}
	if len(updates) > 0 {
		if err := repository.UpdateAdmin(adminID, updates); err != nil {
			return err
		}
	}
	if len(roleIDs) > 0 {
		repository.AssignRoles(adminID, roleIDs)
	}
	return nil
}

// ========== 角色管理 ==========

// CreateRole 创建角色
func CreateRole(name, code, desc string, permissionIDs []uint64) (*model.Role, error) {
	role := &model.Role{
		Name:        name,
		Code:        code,
		Description: desc,
		Status:      1,
	}
	if err := repository.CreateRole(role, permissionIDs); err != nil {
		return nil, err
	}
	return role, nil
}

// UpdateRole 更新角色
func UpdateRole(roleID uint64, name, desc string, permissionIDs []uint64) error {
	if name != "" || desc != "" {
		updates := map[string]interface{}{}
		if name != "" {
			updates["name"] = name
		}
		if desc != "" {
			updates["description"] = desc
		}
		if err := repository.UpdateRole(roleID, updates); err != nil {
			return err
		}
	}
	if len(permissionIDs) > 0 {
		return repository.UpdateRolePermissions(roleID, permissionIDs)
	}
	return nil
}

// ChangePassword 修改密码
func ChangePassword(adminID uint64, oldPassword, newPassword string) error {
	admin, err := repository.FindAdminByID(adminID)
	if err != nil {
		return errors.New("管理员不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword)); err != nil {
		return errors.New("原密码错误")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return repository.UpdateAdminPassword(adminID, string(hashed))
}
