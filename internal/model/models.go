package model

import (
	"time"

	"gorm.io/gorm"
)

type AppType int8

const (
	AppTypeSelf  AppType = 1
	AppTypeThird AppType = 2
)

type PublishStatus int8

const (
	StatusDraft     PublishStatus = 0
	StatusPublished PublishStatus = 1
	StatusOffline   PublishStatus = 2
)

type Channel int8

const (
	ChannelStable Channel = 1
	ChannelBeta   Channel = 2
	ChannelAlpha  Channel = 3
)

type SourceType int8

const (
	SourceManaged  SourceType = 1
	SourceExternal SourceType = 2
)

type User struct {
	ID           uint64     `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Email        string     `gorm:"size:120" json:"email"`
	Nickname     string     `gorm:"size:50" json:"nickname"`
	Avatar       string     `gorm:"size:255" json:"avatar"`
	Role         int8       `gorm:"not null;default:1" json:"role"`
	Status       int8       `gorm:"not null;default:1" json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	LastLoginIP  string     `gorm:"size:45" json:"last_login_ip"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
type UserToken struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"index;not null"`
	TokenHash string    `gorm:"size:64;uniqueIndex;not null"`
	UA        string    `gorm:"size:300"`
	IP        string    `gorm:"size:45"`
	ExpiresAt time.Time `gorm:"index"`
	RevokedAt *time.Time
	CreatedAt time.Time
}
type Category struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Slug        string    `gorm:"size:50;uniqueIndex;not null" json:"slug"`
	Icon        string    `gorm:"size:100" json:"icon"`
	Description string    `gorm:"size:255" json:"description"`
	Sort        int       `json:"sort"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type Tag struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:30;uniqueIndex;not null" json:"name"`
	Slug      string    `gorm:"size:30;uniqueIndex;not null" json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}
type App struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	Slug            string         `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Type            AppType        `gorm:"not null;default:1;index" json:"type"`
	CategoryID      *uint64        `gorm:"index" json:"category_id"`
	Category        *Category      `json:"category,omitempty"`
	Tagline         string         `gorm:"size:200" json:"tagline"`
	Description     string         `gorm:"type:text" json:"description"`
	Icon            string         `gorm:"size:255" json:"icon"`
	Screenshots     []string       `gorm:"serializer:json" json:"screenshots"`
	OfficialURL     string         `gorm:"size:255" json:"official_url"`
	RepoURL         string         `gorm:"size:255" json:"repo_url"`
	Developer       string         `gorm:"size:100" json:"developer"`
	License         string         `gorm:"size:50" json:"license"`
	Platforms       []string       `gorm:"serializer:json" json:"platforms"`
	Status          PublishStatus  `gorm:"not null;default:0;index:idx_apps_pub,priority:1" json:"status"`
	IsFeatured      bool           `json:"is_featured"`
	IsPinned        bool           `json:"is_pinned"`
	SortWeight      int            `json:"sort_weight"`
	LatestReleaseID *uint64        `json:"latest_release_id"`
	DownloadCount   int64          `json:"download_count"`
	ViewCount       int64          `json:"view_count"`
	SeoTitle        string         `gorm:"size:200" json:"seo_title"`
	SeoDescription  string         `gorm:"size:300" json:"seo_description"`
	SeoKeywords     string         `gorm:"size:200" json:"seo_keywords"`
	PublishedAt     *time.Time     `gorm:"index:idx_apps_pub,priority:2,sort:desc" json:"published_at"`
	Tags            []Tag          `gorm:"many2many:app_tags" json:"tags"`
	Releases        []Release      `json:"releases,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
type Release struct {
	ID                 uint64         `gorm:"primaryKey" json:"id"`
	AppID              uint64         `gorm:"uniqueIndex:idx_release_version;not null" json:"app_id"`
	Version            string         `gorm:"size:50;uniqueIndex:idx_release_version;not null" json:"version"`
	VersionCode        *int           `json:"version_code"`
	Channel            Channel        `gorm:"uniqueIndex:idx_release_version;not null;default:1" json:"channel"`
	Title              string         `gorm:"size:200" json:"title"`
	Changelog          string         `gorm:"type:text" json:"changelog"`
	MinRequiredVersion string         `gorm:"size:50" json:"min_required_version"`
	Status             PublishStatus  `gorm:"not null;default:0" json:"status"`
	DownloadCount      int64          `json:"download_count"`
	PublishedAt        *time.Time     `json:"published_at"`
	Assets             []Asset        `json:"assets,omitempty"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}
type Asset struct {
	ID            uint64           `gorm:"primaryKey" json:"id"`
	ReleaseID     uint64           `gorm:"index;not null" json:"release_id"`
	Name          string           `gorm:"size:200" json:"name"`
	FileName      string           `gorm:"size:255" json:"file_name"`
	OS            string           `gorm:"size:20" json:"os"`
	Arch          string           `gorm:"size:20" json:"arch"`
	Kind          int8             `json:"kind"`
	Size          int64            `json:"size"`
	SHA256        string           `gorm:"size:64;index" json:"sha256"`
	DownloadCount int64            `json:"download_count"`
	Sort          int              `json:"sort"`
	Sources       []DownloadSource `json:"sources,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}
type DownloadSource struct {
	ID            uint64     `gorm:"primaryKey" json:"id"`
	AssetID       uint64     `gorm:"index;not null" json:"asset_id"`
	Name          string     `gorm:"size:100" json:"name"`
	SourceType    SourceType `gorm:"not null" json:"source_type"`
	StorageID     *uint64    `gorm:"index" json:"storage_id,omitempty"`
	ObjectKey     string     `gorm:"size:500" json:"object_key,omitempty"`
	ExternalURL   string     `gorm:"size:1000" json:"external_url,omitempty"`
	ExtractCode   string     `gorm:"size:50" json:"extract_code,omitempty"`
	Priority      int        `json:"priority"`
	IsEnabled     bool       `gorm:"default:true" json:"is_enabled"`
	DownloadCount int64      `json:"download_count"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
type Storage struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Driver    string    `gorm:"size:20;not null" json:"driver"`
	Config    string    `gorm:"type:text;not null" json:"config"`
	IsDefault bool      `json:"is_default"`
	IsEnabled bool      `gorm:"default:true" json:"is_enabled"`
	Remark    string    `gorm:"size:255" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Setting struct {
	Key       string    `gorm:"size:100;primaryKey" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}
type Page struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	Title          string         `gorm:"size:200;not null" json:"title"`
	Slug           string         `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Content        string         `gorm:"type:text" json:"content"`
	Status         PublishStatus  `json:"status"`
	Sort           int            `json:"sort"`
	SeoDescription string         `gorm:"size:300" json:"seo_description"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
type DownloadLog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	AppID     uint64    `gorm:"index" json:"app_id"`
	ReleaseID uint64    `json:"release_id"`
	AssetID   uint64    `gorm:"index" json:"asset_id"`
	SourceID  uint64    `json:"source_id"`
	IP        string    `gorm:"size:64" json:"ip"`
	UA        string    `gorm:"size:300" json:"ua"`
	Referer   string    `gorm:"size:300" json:"referer"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
type StatDaily struct {
	ID        uint64 `gorm:"primaryKey"`
	Date      string `gorm:"size:10;uniqueIndex:idx_stat_date_app"`
	AppID     uint64 `gorm:"uniqueIndex:idx_stat_date_app"`
	Downloads int64
	Views     int64
}
type OperationLog struct {
	ID         uint64    `gorm:"primaryKey" json:"id"`
	UserID     uint64    `json:"user_id"`
	Action     string    `gorm:"size:50" json:"action"`
	TargetType string    `gorm:"size:30" json:"target_type"`
	TargetID   uint64    `json:"target_id"`
	Detail     string    `gorm:"type:text" json:"detail"`
	IP         string    `gorm:"size:45" json:"ip"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}
