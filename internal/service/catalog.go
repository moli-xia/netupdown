package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	semver "github.com/Masterminds/semver/v3"
	"github.com/moli-xia/netupdown/internal/model"
	"github.com/moli-xia/netupdown/internal/pkg/apperr"
	"github.com/moli-xia/netupdown/internal/pkg/validate"
	"github.com/moli-xia/netupdown/internal/repo"
	"gorm.io/gorm"
)

type Catalog struct {
	store   *repo.Store
	baseURL string
}

func NewCatalog(store *repo.Store, baseURL string) *Catalog {
	return &Catalog{store: store, baseURL: strings.TrimRight(baseURL, "/")}
}

type PageResult struct {
	List     any   `json:"list"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}
type AppQuery struct {
	Page, PageSize                    int
	Category, Platform, Type, Sort, Q string
	Featured                          bool
	Admin                             bool
}

func paging(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}
func (s *Catalog) Apps(ctx context.Context, q AppQuery) (PageResult, error) {
	q.Page, q.PageSize = paging(q.Page, q.PageSize)
	db := s.store.DB.WithContext(ctx).Model(&model.App{}).Preload("Category").Preload("Tags")
	if !q.Admin {
		db = db.Where("apps.status = ? AND apps.type = ?", model.StatusPublished, model.AppTypeSelf)
		if q.Type == "third" {
			db = db.Where("1 = 0")
		}
	}
	if q.Category != "" {
		db = db.Joins("JOIN categories ON categories.id = apps.category_id").Where("categories.slug = ?", q.Category)
	}
	if q.Platform != "" {
		db = db.Where("apps.platforms LIKE ?", "%\""+q.Platform+"\"%")
	}
	if q.Admin && q.Type == "self" {
		db = db.Where("apps.type = ?", model.AppTypeSelf)
	} else if q.Admin && q.Type == "third" {
		db = db.Where("apps.type = ?", model.AppTypeThird)
	}
	if q.Featured {
		db = db.Where("apps.is_featured = ?", true)
	}
	if q.Q != "" {
		like := "%" + q.Q + "%"
		db = db.Where("apps.name LIKE ? OR apps.tagline LIKE ?", like, like)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return PageResult{}, err
	}
	switch q.Sort {
	case "hot":
		db = db.Order("apps.download_count DESC")
	case "name":
		db = db.Order("apps.name ASC")
	default:
		db = db.Order("apps.is_pinned DESC, apps.sort_weight DESC, apps.published_at DESC, apps.created_at DESC")
	}
	var rows []model.App
	err := db.Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize).Find(&rows).Error
	return PageResult{List: rows, Page: q.Page, PageSize: q.PageSize, Total: total}, err
}
func (s *Catalog) AppBySlug(ctx context.Context, slug string, publishedOnly bool) (model.App, error) {
	var row model.App
	db := s.store.DB.WithContext(ctx).Preload("Category").Preload("Tags").Preload("Releases", func(db *gorm.DB) *gorm.DB {
		if publishedOnly {
			return db.Where("status = ?", model.StatusPublished).Order("published_at DESC")
		}
		return db.Order("created_at DESC")
	}).Preload("Releases.Assets", func(db *gorm.DB) *gorm.DB { return db.Order("sort ASC") }).Preload("Releases.Assets.Sources", func(db *gorm.DB) *gorm.DB { return db.Where("is_enabled = ?", true).Order("priority ASC") }).Where("slug = ?", slug)
	if publishedOnly {
		db = db.Where("status = ? AND type = ?", model.StatusPublished, model.AppTypeSelf)
	}
	err := db.First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, apperr.NotFound
	}
	return row, err
}
func (s *Catalog) AppByID(ctx context.Context, id uint64) (model.App, error) {
	var row model.App
	err := s.store.DB.WithContext(ctx).Preload("Category").Preload("Tags").First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, apperr.NotFound
	}
	return row, err
}
func (s *Catalog) SaveApp(ctx context.Context, row *model.App) (bool, error) {
	if !validate.Slug(row.Slug) || strings.TrimSpace(row.Name) == "" {
		return false, apperr.BadRequest
	}
	row.Name = strings.TrimSpace(row.Name)
	row.Slug = strings.ToLower(strings.TrimSpace(row.Slug))
	created := row.ID == 0
	err := s.store.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tags := row.Tags
		row.Tags = nil
		if created {
			if err := tx.Create(row).Error; err != nil {
				return mapConflict(err)
			}
		} else {
			if err := tx.Model(&model.App{ID: row.ID}).Select("name", "slug", "type", "category_id", "tagline", "description", "icon", "screenshots", "official_url", "repo_url", "developer", "license", "platforms", "seo_title", "seo_description", "seo_keywords").Updates(row).Error; err != nil {
				return mapConflict(err)
			}
		}
		if tags != nil {
			return tx.Model(row).Association("Tags").Replace(tags)
		}
		return nil
	})
	return created, err
}
func (s *Catalog) DeleteApp(ctx context.Context, id uint64) error {
	return s.store.DB.WithContext(ctx).Delete(&model.App{}, id).Error
}
func (s *Catalog) PublishApp(ctx context.Context, id uint64, publish bool) error {
	var row model.App
	if err := s.store.DB.WithContext(ctx).First(&row, id).Error; err != nil {
		return apperr.NotFound
	}
	now := time.Now().UTC()
	updates := map[string]any{}
	if publish {
		updates["status"] = model.StatusPublished
		if row.PublishedAt == nil {
			updates["published_at"] = now
		}
	} else {
		updates["status"] = model.StatusOffline
	}
	return s.store.DB.WithContext(ctx).Model(&row).Updates(updates).Error
}
func (s *Catalog) Categories(ctx context.Context) ([]model.Category, error) {
	var rows []model.Category
	err := s.store.DB.WithContext(ctx).Order("sort ASC, name ASC").Find(&rows).Error
	return rows, err
}
func (s *Catalog) SaveCategory(ctx context.Context, row *model.Category) (bool, error) {
	if !validate.Slug(row.Slug) || strings.TrimSpace(row.Name) == "" {
		return false, apperr.BadRequest
	}
	created := row.ID == 0
	var err error
	if created {
		err = s.store.DB.WithContext(ctx).Create(row).Error
	} else {
		err = s.store.DB.WithContext(ctx).Model(row).Updates(row).Error
	}
	return created, mapConflict(err)
}
func (s *Catalog) DeleteCategory(ctx context.Context, id uint64) error {
	var n int64
	if err := s.store.DB.WithContext(ctx).Model(&model.App{}).Where("category_id = ?", id).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return apperr.Conflict
	}
	return s.store.DB.WithContext(ctx).Delete(&model.Category{}, id).Error
}
func (s *Catalog) Releases(ctx context.Context, appID uint64, admin bool) ([]model.Release, error) {
	var rows []model.Release
	db := s.store.DB.WithContext(ctx).Preload("Assets", func(db *gorm.DB) *gorm.DB { return db.Order("sort ASC") }).Preload("Assets.Sources", func(db *gorm.DB) *gorm.DB { return db.Order("priority ASC") }).Where("app_id = ?", appID)
	if !admin {
		db = db.Where("status = ?", model.StatusPublished)
	}
	err := db.Order("published_at DESC, created_at DESC").Find(&rows).Error
	return rows, err
}
func normalizeVersion(v string) string { return strings.TrimPrefix(strings.TrimSpace(v), "v") }
func (s *Catalog) SaveRelease(ctx context.Context, row *model.Release) (bool, error) {
	row.Version = normalizeVersion(row.Version)
	if row.AppID == 0 || row.Version == "" || row.Channel < 1 || row.Channel > 3 {
		return false, apperr.BadRequest
	}
	created := row.ID == 0
	var err error
	if created {
		err = s.store.DB.WithContext(ctx).Create(row).Error
	} else {
		err = s.store.DB.WithContext(ctx).Model(row).Select("version", "version_code", "channel", "title", "changelog", "min_required_version").Updates(row).Error
	}
	return created, mapConflict(err)
}
func (s *Catalog) PublishRelease(ctx context.Context, id uint64) error {
	return s.store.Tx(func(tx *repo.Store) error {
		var rel model.Release
		if err := tx.DB.WithContext(ctx).Preload("Assets.Sources", func(db *gorm.DB) *gorm.DB { return db.Where("is_enabled = ?", true) }).First(&rel, id).Error; err != nil {
			return apperr.NotFound
		}
		if len(rel.Assets) == 0 {
			return apperr.Wrap(10008, 422, "版本至少需要一个文件", nil)
		}
		for _, a := range rel.Assets {
			if len(a.Sources) == 0 {
				return apperr.Wrap(10008, 422, "每个文件至少需要一个可用下载源", nil)
			}
		}
		now := time.Now().UTC()
		if err := tx.DB.WithContext(ctx).Model(&rel).Updates(map[string]any{"status": model.StatusPublished, "published_at": now}).Error; err != nil {
			return err
		}
		if rel.Channel == model.ChannelStable {
			if err := tx.DB.WithContext(ctx).Model(&model.App{ID: rel.AppID}).Updates(map[string]any{"latest_release_id": rel.ID, "updated_at": now}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func (s *Catalog) DeleteRelease(ctx context.Context, id uint64) error {
	return s.store.DB.WithContext(ctx).Delete(&model.Release{}, id).Error
}
func (s *Catalog) SaveAsset(ctx context.Context, row *model.Asset) (bool, error) {
	if row.ReleaseID == 0 || row.Name == "" || row.FileName == "" {
		return false, apperr.BadRequest
	}
	created := row.ID == 0
	var err error
	if created {
		err = s.store.DB.WithContext(ctx).Create(row).Error
	} else {
		err = s.store.DB.WithContext(ctx).Model(row).Updates(row).Error
	}
	return created, err
}
func (s *Catalog) DeleteAsset(ctx context.Context, id uint64) error {
	return s.store.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("asset_id = ?", id).Delete(&model.DownloadSource{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Asset{}, id).Error
	})
}
func (s *Catalog) SaveSource(ctx context.Context, row *model.DownloadSource) (bool, error) {
	if row.AssetID == 0 || row.Name == "" {
		return false, apperr.BadRequest
	}
	if row.SourceType == model.SourceManaged && (row.StorageID == nil || row.ObjectKey == "") {
		return false, apperr.Unprocessable
	}
	if row.SourceType == model.SourceExternal {
		u, err := url.Parse(row.ExternalURL)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
			return false, apperr.Unprocessable
		}
	}
	created := row.ID == 0
	var err error
	if created {
		err = s.store.DB.WithContext(ctx).Create(row).Error
	} else {
		err = s.store.DB.WithContext(ctx).Model(row).Updates(row).Error
	}
	return created, err
}
func (s *Catalog) DeleteSource(ctx context.Context, id uint64) error {
	return s.store.DB.WithContext(ctx).Delete(&model.DownloadSource{}, id).Error
}

type UpdateResult struct {
	UpdateAvailable bool   `json:"update_available"`
	Mandatory       bool   `json:"mandatory"`
	Current         string `json:"current"`
	Latest          any    `json:"latest"`
}

func (s *Catalog) CheckUpdate(ctx context.Context, slug, current, osName, arch string, channel model.Channel, currentCode *int) (UpdateResult, error) {
	app, err := s.AppBySlug(ctx, slug, true)
	if err != nil {
		return UpdateResult{}, err
	}
	var latest *model.Release
	for i := range app.Releases {
		if app.Releases[i].Channel == channel {
			latest = &app.Releases[i]
			break
		}
	}
	if latest == nil {
		return UpdateResult{}, apperr.NotFound
	}
	available := versionLess(current, currentCode, latest.Version, latest.VersionCode)
	mandatory := latest.MinRequiredVersion != "" && versionLess(current, currentCode, latest.MinRequiredVersion, nil)
	var chosen *model.Asset
	if osName != "" {
		candidates := []model.Asset{}
		for _, a := range latest.Assets {
			if a.OS == osName || a.OS == "any" {
				candidates = append(candidates, a)
			}
		}
		sort.SliceStable(candidates, func(i, j int) bool {
			score := func(a model.Asset) int {
				if a.Arch == arch {
					return 0
				}
				if a.Arch == "universal" {
					return 1
				}
				if a.Arch == "any" {
					return 2
				}
				return 9
			}
			si, sj := score(candidates[i]), score(candidates[j])
			if si == sj {
				return candidates[i].Sort < candidates[j].Sort
			}
			return si < sj
		})
		if len(candidates) > 0 && ((candidates[0].Arch == arch) || candidates[0].Arch == "universal" || candidates[0].Arch == "any") {
			chosen = &candidates[0]
		}
	}
	latestData := map[string]any{"version": latest.Version, "version_code": latest.VersionCode, "channel": channelName(latest.Channel), "title": latest.Title, "changelog": latest.Changelog, "published_at": latest.PublishedAt, "notes_url": s.baseURL + "/apps/" + app.Slug}
	if chosen != nil {
		latestData["asset"] = map[string]any{"os": chosen.OS, "arch": chosen.Arch, "file_name": chosen.FileName, "size": chosen.Size, "sha256": chosen.SHA256, "url": fmt.Sprintf("%s/d/%d", s.baseURL, chosen.ID)}
	}
	return UpdateResult{UpdateAvailable: available, Mandatory: mandatory, Current: current, Latest: latestData}, nil
}
func versionLess(current string, currentCode *int, target string, targetCode *int) bool {
	c, ce := semver.NewVersion(normalizeVersion(current))
	t, te := semver.NewVersion(normalizeVersion(target))
	if ce == nil && te == nil {
		return c.LessThan(t)
	}
	if currentCode != nil && targetCode != nil {
		return *currentCode < *targetCode
	}
	return normalizeVersion(current) != normalizeVersion(target)
}
func channelName(c model.Channel) string {
	switch c {
	case model.ChannelBeta:
		return "beta"
	case model.ChannelAlpha:
		return "alpha"
	default:
		return "stable"
	}
}
func mapConflict(err error) error {
	if err == nil {
		return nil
	}
	low := strings.ToLower(err.Error())
	if strings.Contains(low, "unique") || strings.Contains(low, "duplicate") {
		return apperr.Conflict
	}
	return err
}
