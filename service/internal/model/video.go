package model

import (
	"time"

	"gorm.io/gorm"
)

// VideoCategory 视频分类
type VideoCategory struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	EntityID    uint64         `gorm:"not null;index:idx_entity_mp_sort" json:"entity_id"`
	MpAccountID uint64         `gorm:"not null;index:idx_entity_mp_sort;index:idx_mp" json:"mp_account_id"`
	Name        string         `gorm:"size:50;not null" json:"name"`
	IconURL     string         `gorm:"size:500;not null;default:''" json:"icon_url"`
	SortOrder   int            `gorm:"not null;default:0;index:idx_entity_mp_sort" json:"sort_order"`
	Status      int8           `gorm:"not null;default:1;index:idx_mp" json:"status"`
	CreatedAt   time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (VideoCategory) TableName() string { return "video_categories" }

// Video 视频
type Video struct {
	ID             uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID    uint64         `gorm:"not null;index:idx_mp_status;index:idx_mp_category" json:"mp_account_id"`
	EntityID       uint64         `gorm:"not null;index:idx_entity" json:"entity_id"`
	CategoryID     uint64         `gorm:"not null;default:0;index:idx_mp_category" json:"category_id"`
	Title          string         `gorm:"size:200;not null;default:''" json:"title"`
	Description    string         `gorm:"size:1000;not null;default:''" json:"description"`
	CoverURL       string         `gorm:"size:500;not null" json:"cover_url"`
	VideoURL       string         `gorm:"size:500;not null" json:"video_url"`
	Duration       uint           `gorm:"not null;default:0" json:"duration"`
	Width          uint           `gorm:"not null;default:0" json:"width"`
	Height         uint           `gorm:"not null;default:0" json:"height"`
	FileSize       uint64         `gorm:"not null;default:0" json:"file_size"`
	Tags           string         `gorm:"size:500;not null;default:''" json:"tags"`
	Status         int8           `gorm:"not null;default:0;index:idx_mp_status;index:idx_published" json:"status"`
	IsRecommended  int8           `gorm:"not null;default:0" json:"is_recommended"`
	ViewCount      uint64         `gorm:"not null;default:0" json:"view_count"`
	LikeCount      uint64         `gorm:"not null;default:0" json:"like_count"`
	CollectCount   uint64         `gorm:"not null;default:0" json:"collect_count"`
	ShareCount     uint64         `gorm:"not null;default:0" json:"share_count"`
	CommentCount   uint64         `gorm:"not null;default:0" json:"comment_count"`
	PublishedAt    *time.Time     `gorm:"null;index:idx_mp_status;index:idx_published" json:"published_at"`
	CreatedAt      time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Video) TableName() string { return "videos" }
