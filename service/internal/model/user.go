package model

import "time"

// User 小程序用户
type User struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID uint64    `gorm:"not null;uniqueIndex:uk_mp_openid;index:idx_mp_phone" json:"mp_account_id"`
	Openid      string    `gorm:"size:64;not null;uniqueIndex:uk_mp_openid" json:"openid"`
	Unionid     string    `gorm:"size:64;not null;default:''" json:"unionid"`
	SessionKey  string    `gorm:"size:128;not null;default:''" json:"-"`
	Nickname    string    `gorm:"size:100;not null;default:''" json:"nickname"`
	AvatarURL   string    `gorm:"size:500;not null;default:''" json:"avatar_url"`
	Phone       string    `gorm:"size:20;not null;default:'';index:idx_mp_phone" json:"phone"`
	Gender      int8      `gorm:"not null;default:0" json:"gender"`
	Province    string    `gorm:"size:50;not null;default:''" json:"province"`
	City        string    `gorm:"size:50;not null;default:''" json:"city"`
	Country     string    `gorm:"size:50;not null;default:''" json:"country"`
	Status      int8      `gorm:"not null;default:1" json:"status"`
	LastLoginAt *time.Time `gorm:"null" json:"last_login_at"`
	LastLoginIP string    `gorm:"size:50;not null;default:''" json:"last_login_ip"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string { return "users" }
