package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/moli-xia/netupdown/internal/config"
	"github.com/moli-xia/netupdown/internal/model"
	"gorm.io/gorm"
)

func Open(cfg config.Config) (*gorm.DB, error) {
	if cfg.Database.Driver != "sqlite" {
		return nil, fmt.Errorf("database driver %q is not enabled in this build", cfg.Database.Driver)
	}
	if err := os.MkdirAll(filepath.Dir(cfg.Database.DSN), 0o700); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}
	db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	for _, stmt := range []string{"PRAGMA journal_mode=WAL", "PRAGMA busy_timeout=5000", "PRAGMA foreign_keys=ON", "PRAGMA synchronous=NORMAL"} {
		if err := db.Exec(stmt).Error; err != nil {
			return nil, fmt.Errorf("apply %s: %w", stmt, err)
		}
	}
	if err := db.AutoMigrate(&model.User{}, &model.UserToken{}, &model.Category{}, &model.Tag{}, &model.App{}, &model.Release{}, &model.Asset{}, &model.DownloadSource{}, &model.Storage{}, &model.Setting{}, &model.Page{}, &model.DownloadLog{}, &model.StatDaily{}, &model.OperationLog{}); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}
	if err := seed(db); err != nil {
		return nil, err
	}
	return db, nil
}

func seed(db *gorm.DB) error {
	categories := []model.Category{{Name: "桌面效率", Slug: "productivity", Sort: 10}, {Name: "开发者工具", Slug: "development", Sort: 20}, {Name: "创作工具", Slug: "media", Sort: 30}, {Name: "系统工具", Slug: "system", Sort: 40}, {Name: "网络与服务", Slug: "network", Sort: 50}, {Name: "移动应用", Slug: "works", Sort: 60}}
	for _, item := range categories {
		var n int64
		if err := db.Model(&model.Category{}).Where("slug = ?", item.Slug).Count(&n).Error; err != nil {
			return err
		}
		if n == 0 {
			if err := db.Create(&item).Error; err != nil {
				return err
			}
		}
	}
	settings := map[string]string{"site.title": "造物工坊", "site.subtitle": "独立开发者的软件发布与更新中心", "site.description": "探索 造物工坊自研产品，获取可靠的多平台安装包、版本记录与更新说明。", "site.keywords": "自研软件,独立开发,软件下载,应用更新", "site.footer": "", "site.icp": "", "theme.active": "aurora", "download.dedup_window_min": "10", "download.log_retention_days": "90", "download.show_hash": "true", "download.landing_for_external": "true", "content.apps_page_size": "24", "content.home_latest_limit": "8", "seo.sitemap_enabled": "true", "seo.feed_enabled": "true", "privacy.ip_mode": "truncate"}
	for key, value := range settings {
		var row model.Setting
		if err := db.Where("key = ?", key).Attrs(model.Setting{Key: key, Value: value}).FirstOrCreate(&row).Error; err != nil {
			return err
		}
	}
	var n int64
	if err := db.Model(&model.Storage{}).Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		if err := db.Create(&model.Storage{Name: "本地存储", Driver: "local", Config: `{"root":"files"}`, IsDefault: true, IsEnabled: true}).Error; err != nil {
			return err
		}
	}
	return nil
}
