package repository

import (
	"time"
	"yanny-service/internal/database"
	"yanny-service/internal/model"
)

// FindUserByOpenid 按 openid 查找用户
func FindUserByOpenid(mpAccountID uint64, openid string) (*model.User, error) {
	var user model.User
	err := database.DB.Where("mp_account_id = ? AND openid = ?", mpAccountID, openid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

// UpdateUserLoginInfo 更新用户登录信息
func UpdateUserLoginInfo(userID uint64, ip string, loginAt *time.Time) error {
	return database.DB.Model(&model.User{}).Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": loginAt,
			"last_login_ip": ip,
		}).Error
}

// ListUsers 用户列表
func ListUsers(mpAccountID uint64, phone string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	db := database.DB.Model(&model.User{})
	if mpAccountID > 0 {
		db = db.Where("mp_account_id = ?", mpAccountID)
	}
	if phone != "" {
		db = db.Where("phone LIKE ?", "%"+phone+"%")
	}
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&users).Error
	return users, total, err
}

// UpdateUserStatus 更新用户状态
func UpdateUserStatus(userID uint64, status int8) error {
	return database.DB.Model(&model.User{}).Where("id = ?", userID).Update("status", status).Error
}

// UpdateUserSessionKey 更新用户微信 session_key（每次登录轮换）
func UpdateUserSessionKey(userID uint64, sessionKey string) error {
	return database.DB.Model(&model.User{}).Where("id = ?", userID).Update("session_key", sessionKey).Error
}

// UpdateUserPhone 更新用户手机号
func UpdateUserPhone(userID uint64, phone string) error {
	return database.DB.Model(&model.User{}).Where("id = ?", userID).Update("phone", phone).Error
}

// UpdateUserInfo 更新用户昵称/头像
func UpdateUserInfo(userID uint64, nickname, avatarURL string) error {
	return database.DB.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"nickname":   nickname,
		"avatar_url": avatarURL,
	}).Error
}

// FindUserByID 按 ID 查找用户
func FindUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
