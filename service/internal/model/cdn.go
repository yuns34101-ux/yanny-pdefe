package model

import "time"

// CdnConfig CDN 配置（七牛云等）
type CdnConfig struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID  uint64    `gorm:"not null;index:idx_mp_account" json:"mp_account_id"`
	Provider     string    `gorm:"size:30;not null;default:qiniu" json:"provider"`
	AccessKey    string    `gorm:"size:200;not null" json:"-"`
	SecretKey    string    `gorm:"size:200;not null" json:"-"`
	Bucket       string    `gorm:"size:100;not null" json:"bucket"`
	Domain       string    `gorm:"size:200;not null" json:"domain"`
	Region       string    `gorm:"size:50;not null;default:''" json:"region"`
	CallbackURL  string    `gorm:"size:300;not null;default:''" json:"callback_url"`
	Status       int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (CdnConfig) TableName() string { return "cdn_configs" }
