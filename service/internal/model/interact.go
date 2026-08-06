package model

import "time"

// Comment 评论
type Comment struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID    uint64     `gorm:"not null;index:idx_video_root" json:"mp_account_id"`
	VideoID        uint64     `gorm:"not null;index:idx_video_root;index:idx_video_parent" json:"video_id"`
	UserID         uint64     `gorm:"not null;index:idx_user" json:"user_id"`
	ParentID       *uint64    `gorm:"null;index:idx_video_parent" json:"parent_id"`
	RootID         *uint64    `gorm:"null;index:idx_video_root" json:"root_id"`
	ReplyToUserID  *uint64    `gorm:"null" json:"reply_to_user_id"`
	Content        string     `gorm:"size:1000;not null" json:"content"`
	LikeCount      uint       `gorm:"not null;default:0" json:"like_count"`
	ReplyCount     uint       `gorm:"not null;default:0" json:"reply_count"`
	Status         int8       `gorm:"not null;default:1" json:"status"`
	IsTop          int8       `gorm:"not null;default:0" json:"is_top"`
	CreatedAt      time.Time  `gorm:"not null;autoCreateTime;index:idx_video_root;index:idx_video_parent" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"not null;autoUpdateTime" json:"updated_at"`

	// 关联（非数据库字段）
	User       *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ReplyToUser *User    `gorm:"foreignKey:ReplyToUserID" json:"reply_to_user,omitempty"`
	Replies    []Comment `gorm:"-" json:"replies,omitempty"`
}

func (Comment) TableName() string { return "comments" }

// Like 点赞
type Like struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID uint64    `gorm:"not null;uniqueIndex:uk_user_target" json:"mp_account_id"`
	UserID      uint64    `gorm:"not null;uniqueIndex:uk_user_target" json:"user_id"`
	TargetType  string    `gorm:"size:20;not null;uniqueIndex:uk_user_target;index:idx_target" json:"target_type"`
	TargetID    uint64    `gorm:"not null;uniqueIndex:uk_user_target;index:idx_target" json:"target_id"`
	Status      int8      `gorm:"not null;default:1" json:"status"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (Like) TableName() string { return "likes" }

// Favorite 收藏
type Favorite struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID uint64    `gorm:"not null;uniqueIndex:uk_user_video;index:idx_user" json:"mp_account_id"`
	UserID      uint64    `gorm:"not null;uniqueIndex:uk_user_video;index:idx_user" json:"user_id"`
	VideoID     uint64    `gorm:"not null;uniqueIndex:uk_user_video" json:"video_id"`
	Status      int8      `gorm:"not null;default:1;index:idx_user" json:"status"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (Favorite) TableName() string { return "favorites" }

// Share 分享记录
type Share struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID uint64    `gorm:"not null;index:idx_mp_date" json:"mp_account_id"`
	UserID      uint64    `gorm:"not null;index:idx_user" json:"user_id"`
	VideoID     uint64    `gorm:"not null;index:idx_video" json:"video_id"`
	ShareType   string    `gorm:"size:20;not null;default:''" json:"share_type"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime;index:idx_video;index:idx_user;index:idx_mp_date" json:"created_at"`
}

func (Share) TableName() string { return "shares" }
