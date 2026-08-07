package model

import "time"

// MediaAsset 媒体资源去重记录
type MediaAsset struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID uint64    `gorm:"not null;default:0;uniqueIndex:uk_mp_hash" json:"mp_account_id"`
	DirType     string    `gorm:"size:20;not null" json:"dir_type"`
	ContentHash string    `gorm:"size:64;not null;uniqueIndex:uk_mp_hash" json:"content_hash"`
	ClientHash  string    `gorm:"size:64;not null;default:'';index:idx_client_hash" json:"client_hash"`
	ObjectKey   string    `gorm:"size:300;not null" json:"object_key"`
	URL         string    `gorm:"size:500;not null" json:"url"`
	FileSize    uint64    `gorm:"not null;default:0" json:"file_size"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (MediaAsset) TableName() string { return "media_assets" }
