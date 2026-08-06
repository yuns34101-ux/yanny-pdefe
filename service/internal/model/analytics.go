package model

import "time"

// ViewLog 视频播放记录（埋点原始数据）
type ViewLog struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID   uint64    `gorm:"not null;index:idx_mp_entity_date" json:"mp_account_id"`
	EntityID      uint64    `gorm:"not null;default:0;index:idx_mp_entity_date" json:"entity_id"`
	UserID        *uint64   `gorm:"null;index:idx_user_date" json:"user_id"`
	VideoID       uint64    `gorm:"not null;index:idx_video_date" json:"video_id"`
	CategoryID    uint64    `gorm:"not null;default:0;index:idx_category_date" json:"category_id"`
	WatchDuration uint      `gorm:"not null;default:0" json:"watch_duration"`
	IsComplete    int8      `gorm:"not null;default:0" json:"is_complete"`
	Source        string    `gorm:"size:30;not null;default:''" json:"source"`
	IP            string    `gorm:"size:50;not null;default:''" json:"ip"`
	Province      string    `gorm:"size:50;not null;default:'';index:idx_province_date" json:"province"`
	City          string    `gorm:"size:50;not null;default:''" json:"city"`
	Device        string    `gorm:"size:100;not null;default:''" json:"device"`
	OS            string    `gorm:"size:30;not null;default:''" json:"os"`
	CreatedAt     time.Time `gorm:"not null;autoCreateTime;index:idx_mp_entity_date;index:idx_video_date;index:idx_user_date;index:idx_category_date;index:idx_province_date" json:"created_at"`
}

func (ViewLog) TableName() string { return "view_logs" }

// ActionLog 用户行为事件（埋点原始数据）
type ActionLog struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID uint64    `gorm:"not null;index:idx_mp_type_date" json:"mp_account_id"`
	EntityID    uint64    `gorm:"not null;default:0" json:"entity_id"`
	UserID      *uint64   `gorm:"null;index:idx_user_date" json:"user_id"`
	EventType   string    `gorm:"size:30;not null;index:idx_mp_type_date" json:"event_type"`
	TargetType  string    `gorm:"size:30;not null;default:''" json:"target_type"`
	TargetID    uint64    `gorm:"not null;default:0" json:"target_id"`
	PagePath    string    `gorm:"size:200;not null;default:'';index:idx_page_date" json:"page_path"`
	ExtraData   string    `gorm:"type:json;null" json:"extra_data"`
	IP          string    `gorm:"size:50;not null;default:''" json:"ip"`
	Province    string    `gorm:"size:50;not null;default:''" json:"province"`
	City        string    `gorm:"size:50;not null;default:''" json:"city"`
	Device      string    `gorm:"size:100;not null;default:''" json:"device"`
	OS          string    `gorm:"size:30;not null;default:''" json:"os"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime;index:idx_mp_type_date;index:idx_user_date;index:idx_page_date" json:"created_at"`
}

func (ActionLog) TableName() string { return "action_logs" }

// StatsVideoDaily 视频每日统计
type StatsVideoDaily struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID      uint64    `gorm:"not null;index:idx_mp_entity_date" json:"mp_account_id"`
	EntityID         uint64    `gorm:"not null;index:idx_mp_entity_date" json:"entity_id"`
	VideoID          uint64    `gorm:"not null;uniqueIndex:uk_video_date" json:"video_id"`
	CategoryID       uint64    `gorm:"not null;default:0;index:idx_category_date" json:"category_id"`
	StatDate         string    `gorm:"type:date;not null;uniqueIndex:uk_video_date;index:idx_mp_entity_date;index:idx_category_date" json:"stat_date"`
	ViewCount        uint64    `gorm:"not null;default:0" json:"view_count"`
	ViewUsers        uint64    `gorm:"not null;default:0" json:"view_users"`
	AvgWatchDuration float64   `gorm:"type:decimal(8,2);not null;default:0" json:"avg_watch_duration"`
	CompleteCount    uint64    `gorm:"not null;default:0" json:"complete_count"`
	LikeCount        uint64    `gorm:"not null;default:0" json:"like_count"`
	CollectCount     uint64    `gorm:"not null;default:0" json:"collect_count"`
	ShareCount       uint64    `gorm:"not null;default:0" json:"share_count"`
	CommentCount     uint64    `gorm:"not null;default:0" json:"comment_count"`
	CreatedAt        time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (StatsVideoDaily) TableName() string { return "stats_video_daily" }

// StatsPlatformDaily 平台每日统计
type StatsPlatformDaily struct {
	ID               uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID      uint64    `gorm:"not null;uniqueIndex:uk_mp_entity_date" json:"mp_account_id"`
	EntityID         uint64    `gorm:"not null;uniqueIndex:uk_mp_entity_date" json:"entity_id"`
	StatDate         string    `gorm:"type:date;not null;uniqueIndex:uk_mp_entity_date" json:"stat_date"`
	TotalViews       uint64    `gorm:"not null;default:0" json:"total_views"`
	TotalViewUsers   uint64    `gorm:"not null;default:0" json:"total_view_users"`
	ActiveUsers      uint64    `gorm:"not null;default:0" json:"active_users"`
	NewUsers         uint64    `gorm:"not null;default:0" json:"new_users"`
	TotalUsers       uint64    `gorm:"not null;default:0" json:"total_users"`
	TotalLikes       uint64    `gorm:"not null;default:0" json:"total_likes"`
	TotalCollects    uint64    `gorm:"not null;default:0" json:"total_collects"`
	TotalShares      uint64    `gorm:"not null;default:0" json:"total_shares"`
	TotalComments    uint64    `gorm:"not null;default:0" json:"total_comments"`
	AvgOnlineMinutes float64   `gorm:"type:decimal(8,2);not null;default:0" json:"avg_online_minutes"`
	CreatedAt        time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (StatsPlatformDaily) TableName() string { return "stats_platform_daily" }

// StatsRegionDaily 地域每日统计
type StatsRegionDaily struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	MpAccountID uint64    `gorm:"not null;uniqueIndex:uk_mp_province_date;index:idx_entity_date" json:"mp_account_id"`
	EntityID    uint64    `gorm:"not null;index:idx_entity_date" json:"entity_id"`
	StatDate    string    `gorm:"type:date;not null;uniqueIndex:uk_mp_province_date;index:idx_entity_date" json:"stat_date"`
	Province    string    `gorm:"size:50;not null;uniqueIndex:uk_mp_province_date" json:"province"`
	ViewCount   uint64    `gorm:"not null;default:0" json:"view_count"`
	ViewUsers   uint64    `gorm:"not null;default:0" json:"view_users"`
	ActiveUsers uint64    `gorm:"not null;default:0" json:"active_users"`
	NewUsers    uint64    `gorm:"not null;default:0" json:"new_users"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;autoUpdateTime" json:"updated_at"`
}

func (StatsRegionDaily) TableName() string { return "stats_region_daily" }
