package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/moli-xia/netupdown/internal/config"
	"github.com/moli-xia/netupdown/internal/model"
	"github.com/moli-xia/netupdown/internal/pkg/apperr"
	"github.com/moli-xia/netupdown/internal/pkg/cryptoutil"
	"github.com/moli-xia/netupdown/internal/repo"
	storagepkg "github.com/moli-xia/netupdown/internal/storage"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type StorageManager struct {
	store   *repo.Store
	dataDir string
	mu      sync.RWMutex
	drivers map[uint64]storagepkg.Driver
	sealer  *cryptoutil.Sealer
}

func NewStorageManager(store *repo.Store, dataDir string, sealer *cryptoutil.Sealer) *StorageManager {
	return &StorageManager{store: store, dataDir: dataDir, drivers: map[uint64]storagepkg.Driver{}, sealer: sealer}
}
func (s *StorageManager) List(ctx context.Context) ([]model.Storage, error) {
	var rows []model.Storage
	err := s.store.DB.WithContext(ctx).Order("is_default DESC, id ASC").Find(&rows).Error
	if err == nil {
		for i := range rows {
			plain, decryptErr := s.sealer.Decrypt(rows[i].Config)
			if decryptErr != nil {
				return nil, decryptErr
			}
			rows[i].Config = maskStorageConfig(rows[i].Driver, plain)
		}
	}
	return rows, err
}
func (s *StorageManager) Save(ctx context.Context, row *model.Storage) (bool, error) {
	if row.Name == "" || (row.Driver != "local" && row.Driver != "s3") {
		return false, apperr.Wrap(10008, 422, "存储驱动必须为 local 或 s3", nil)
	}
	plainConfig := row.Config
	created := row.ID == 0
	if !created && row.Driver == "s3" {
		var incoming storagepkg.S3Config
		if json.Unmarshal([]byte(plainConfig), &incoming) == nil && (incoming.AccessKeyID == "******" || incoming.SecretAccessKey == "******") {
			var existing model.Storage
			if err := s.store.DB.WithContext(ctx).First(&existing, row.ID).Error; err != nil {
				return false, apperr.NotFound
			}
			oldPlain, err := s.sealer.Decrypt(existing.Config)
			if err != nil {
				return false, err
			}
			var old storagepkg.S3Config
			if json.Unmarshal([]byte(oldPlain), &old) != nil {
				return false, apperr.Unprocessable
			}
			if incoming.AccessKeyID == "******" {
				incoming.AccessKeyID = old.AccessKeyID
			}
			if incoming.SecretAccessKey == "******" {
				incoming.SecretAccessKey = old.SecretAccessKey
			}
			raw, _ := json.Marshal(incoming)
			plainConfig = string(raw)
		}
	}
	if row.Driver == "local" {
		var cfg struct {
			Root string `json:"root"`
		}
		if json.Unmarshal([]byte(plainConfig), &cfg) != nil || cfg.Root == "" {
			return false, apperr.BadRequest
		}
	} else {
		var cfg storagepkg.S3Config
		if json.Unmarshal([]byte(plainConfig), &cfg) != nil || cfg.Bucket == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
			return false, apperr.BadRequest
		}
	}
	encrypted, err := s.sealer.Encrypt(plainConfig)
	if err != nil {
		return false, err
	}
	row.Config = encrypted
	err = s.store.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if row.IsDefault {
			if err := tx.Model(&model.Storage{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
				return err
			}
		}
		if created {
			return tx.Create(row).Error
		}
		return tx.Model(row).Select("name", "driver", "config", "is_default", "is_enabled", "remark").Updates(row).Error
	})
	if err == nil {
		s.mu.Lock()
		delete(s.drivers, row.ID)
		s.mu.Unlock()
	}
	row.Config = maskStorageConfig(row.Driver, plainConfig)
	return created, err
}
func (s *StorageManager) Delete(ctx context.Context, id uint64) error {
	var n int64
	if err := s.store.DB.WithContext(ctx).Model(&model.DownloadSource{}).Where("storage_id = ?", id).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return apperr.Conflict
	}
	s.mu.Lock()
	delete(s.drivers, id)
	s.mu.Unlock()
	return s.store.DB.WithContext(ctx).Delete(&model.Storage{}, id).Error
}
func (s *StorageManager) Default(ctx context.Context) (model.Storage, error) {
	var row model.Storage
	err := s.store.DB.WithContext(ctx).Where("is_default = ? AND is_enabled = ?", true, true).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return row, apperr.NotFound
	}
	return row, err
}
func (s *StorageManager) Driver(ctx context.Context, id uint64) (storagepkg.Driver, model.Storage, error) {
	s.mu.RLock()
	d := s.drivers[id]
	s.mu.RUnlock()
	var row model.Storage
	if err := s.store.DB.WithContext(ctx).First(&row, id).Error; err != nil {
		return nil, row, apperr.NotFound
	}
	if d != nil {
		return d, row, nil
	}
	plain, err := s.sealer.Decrypt(row.Config)
	if err != nil {
		return nil, row, err
	}
	var built storagepkg.Driver
	if row.Driver == "local" {
		var cfg struct {
			Root string `json:"root"`
		}
		if json.Unmarshal([]byte(plain), &cfg) != nil {
			return nil, row, apperr.Unprocessable
		}
		root := cfg.Root
		if !filepath.IsAbs(root) {
			root = filepath.Join(s.dataDir, root)
		}
		built, err = storagepkg.NewLocal(root)
	} else if row.Driver == "s3" {
		var cfg storagepkg.S3Config
		if json.Unmarshal([]byte(plain), &cfg) != nil {
			return nil, row, apperr.Unprocessable
		}
		built, err = storagepkg.NewS3(ctx, cfg)
	} else {
		return nil, row, apperr.Unprocessable
	}
	if err != nil {
		return nil, row, err
	}
	s.mu.Lock()
	s.drivers[id] = built
	s.mu.Unlock()
	return built, row, nil
}

func maskStorageConfig(driver, raw string) string {
	if driver != "s3" {
		return raw
	}
	var cfg map[string]any
	if json.Unmarshal([]byte(raw), &cfg) != nil {
		return "{}"
	}
	if _, ok := cfg["access_key_id"]; ok {
		cfg["access_key_id"] = "******"
	}
	if _, ok := cfg["secret_access_key"]; ok {
		cfg["secret_access_key"] = "******"
	}
	b, _ := json.Marshal(cfg)
	return string(b)
}
func (s *StorageManager) Test(ctx context.Context, id uint64) (time.Duration, error) {
	d, _, err := s.Driver(ctx, id)
	if err != nil {
		return 0, err
	}
	start := time.Now()
	key := ".netupdown-probe-" + xid.New().String()
	payload := strings.NewReader("netupdown-storage-probe")
	if err := d.Put(ctx, key, payload, int64(payload.Len())); err != nil {
		return 0, err
	}
	defer d.Delete(context.Background(), key)
	if _, err := d.Stat(ctx, key); err != nil {
		return 0, err
	}
	return time.Since(start), nil
}

type UploadService struct {
	cfg      config.Config
	store    *repo.Store
	storages *StorageManager
}
type UploadMeta struct {
	FileName    string `json:"file_name"`
	Size        int64  `json:"size"`
	SHA256      string `json:"sha256"`
	ChunkSize   int64  `json:"chunk_size"`
	ChunksTotal int    `json:"chunks_total"`
}

func NewUploadService(cfg config.Config, store *repo.Store, sm *StorageManager) *UploadService {
	return &UploadService{cfg: cfg, store: store, storages: sm}
}
func (s *UploadService) Init(ctx context.Context, name string, size int64, hash string, chunkSize int64) (map[string]any, error) {
	if name == "" || size <= 0 || size > s.cfg.Upload.MaxSizeMB*1024*1024 || !regexp.MustCompile(`^[a-fA-F0-9]{64}$`).MatchString(hash) {
		return nil, apperr.BadRequest
	}
	var source model.DownloadSource
	err := s.store.DB.WithContext(ctx).Joins("JOIN assets ON assets.id = download_sources.asset_id").Where("assets.sha256 = ? AND download_sources.source_type = ?", strings.ToLower(hash), model.SourceManaged).First(&source).Error
	if err == nil && source.StorageID != nil {
		d, _, e := s.storages.Driver(ctx, *source.StorageID)
		if e == nil {
			if info, e := d.Stat(ctx, source.ObjectKey); e == nil && info.Size == size {
				return map[string]any{"exists": true, "object": map[string]any{"storage_id": *source.StorageID, "object_key": source.ObjectKey, "size": size, "sha256": strings.ToLower(hash)}}, nil
			}
		}
	}
	if chunkSize <= 0 {
		chunkSize = s.cfg.Upload.ChunkSizeMB * 1024 * 1024
	}
	root := filepath.Join(s.cfg.DataDir, "tmp", "uploads")
	if entries, readErr := os.ReadDir(root); readErr == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			_, existing, sessionErr := s.session(entry.Name())
			if sessionErr != nil || existing.Size != size || !strings.EqualFold(existing.SHA256, hash) {
				continue
			}
			uploaded := make([]int, 0)
			for index := 0; index < existing.ChunksTotal; index++ {
				if _, statErr := os.Stat(filepath.Join(root, entry.Name(), fmt.Sprintf("%d.part", index))); statErr == nil {
					uploaded = append(uploaded, index)
				}
			}
			return map[string]any{"exists": false, "upload_id": entry.Name(), "chunk_size": existing.ChunkSize, "uploaded_chunks": uploaded}, nil
		}
	}
	uploadID := xid.New().String()
	dir := filepath.Join(s.cfg.DataDir, "tmp", "uploads", uploadID)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	meta := UploadMeta{FileName: sanitizeFilename(name), Size: size, SHA256: strings.ToLower(hash), ChunkSize: chunkSize, ChunksTotal: int((size + chunkSize - 1) / chunkSize)}
	b, _ := json.Marshal(meta)
	if err := os.WriteFile(filepath.Join(dir, "meta.json"), b, 0o600); err != nil {
		return nil, err
	}
	return map[string]any{"exists": false, "upload_id": uploadID, "chunk_size": chunkSize, "uploaded_chunks": []int{}}, nil
}
func (s *UploadService) Chunk(uploadID string, index int, r io.Reader, chunkHash string) error {
	dir, meta, err := s.session(uploadID)
	if err != nil {
		return err
	}
	if index < 0 || index >= meta.ChunksTotal {
		return apperr.BadRequest
	}
	tmp := filepath.Join(dir, fmt.Sprintf("%d.part.tmp", index))
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	h := sha256.New()
	limit := meta.ChunkSize
	if index == meta.ChunksTotal-1 {
		limit = meta.Size - int64(index)*meta.ChunkSize
	}
	n, copyErr := io.Copy(io.MultiWriter(f, h), io.LimitReader(r, limit+1))
	closeErr := f.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if n != limit {
		_ = os.Remove(tmp)
		return apperr.Unprocessable
	}
	if chunkHash != "" && !strings.EqualFold(chunkHash, hex.EncodeToString(h.Sum(nil))) {
		_ = os.Remove(tmp)
		return apperr.Unprocessable
	}
	return os.Rename(tmp, filepath.Join(dir, fmt.Sprintf("%d.part", index)))
}
func (s *UploadService) Complete(ctx context.Context, uploadID string, storageID *uint64, keyHint string) (map[string]any, error) {
	dir, meta, err := s.session(uploadID)
	if err != nil {
		return nil, err
	}
	var st model.Storage
	if storageID == nil {
		st, err = s.storages.Default(ctx)
		if err != nil {
			return nil, err
		}
		storageID = &st.ID
	}
	d, _, err := s.storages.Driver(ctx, *storageID)
	if err != nil {
		return nil, err
	}
	readers := make([]io.Reader, 0, meta.ChunksTotal)
	files := make([]*os.File, 0, meta.ChunksTotal)
	defer func() {
		for _, f := range files {
			_ = f.Close()
		}
	}()
	for i := 0; i < meta.ChunksTotal; i++ {
		f, e := os.Open(filepath.Join(dir, fmt.Sprintf("%d.part", i)))
		if e != nil {
			return nil, apperr.Wrap(10008, 422, fmt.Sprintf("缺少分片 %d", i), e)
		}
		files = append(files, f)
		readers = append(readers, f)
	}
	h := sha256.New()
	key := strings.Trim(strings.ReplaceAll(keyHint, "\\", "/"), "/")
	if key != "" {
		key += "/"
	}
	key += meta.FileName
	key, err = uniqueObjectKey(ctx, d, key)
	if err != nil {
		return nil, err
	}
	tee := io.TeeReader(io.MultiReader(readers...), h)
	if err := d.Put(ctx, key, tee, meta.Size); err != nil {
		return nil, err
	}
	actual := hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(actual, meta.SHA256) {
		_ = d.Delete(ctx, key)
		return nil, apperr.Wrap(10008, 422, "SHA256 校验失败", nil)
	}
	_ = os.RemoveAll(dir)
	return map[string]any{"storage_id": *storageID, "object_key": key, "size": meta.Size, "sha256": actual}, nil
}
func (s *UploadService) Abort(uploadID string) error {
	dir, _, err := s.session(uploadID)
	if err != nil {
		return err
	}
	return os.RemoveAll(dir)
}
func (s *UploadService) PutFile(ctx context.Context, name string, r io.Reader, size int64, storageID *uint64, keyHint string) (map[string]any, error) {
	if size <= 0 || size > s.cfg.Upload.MaxSizeMB*1024*1024 {
		return nil, apperr.Wrap(10007, 413, "文件超出大小限制", nil)
	}
	var st model.Storage
	var err error
	if storageID == nil {
		st, err = s.storages.Default(ctx)
		if err != nil {
			return nil, err
		}
		storageID = &st.ID
	}
	d, _, err := s.storages.Driver(ctx, *storageID)
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	key := strings.Trim(strings.ReplaceAll(keyHint, "\\", "/"), "/")
	if key != "" {
		key += "/"
	}
	key += sanitizeFilename(name)
	key, err = uniqueObjectKey(ctx, d, key)
	if err != nil {
		return nil, err
	}
	if err := d.Put(ctx, key, io.TeeReader(r, h), size); err != nil {
		return nil, err
	}
	return map[string]any{"storage_id": *storageID, "object_key": key, "size": size, "sha256": hex.EncodeToString(h.Sum(nil))}, nil
}
func (s *UploadService) session(id string) (string, UploadMeta, error) {
	var meta UploadMeta
	if !regexp.MustCompile(`^[a-z0-9]{20}$`).MatchString(id) {
		return "", meta, apperr.BadRequest
	}
	dir := filepath.Join(s.cfg.DataDir, "tmp", "uploads", id)
	b, err := os.ReadFile(filepath.Join(dir, "meta.json"))
	if os.IsNotExist(err) {
		return "", meta, apperr.NotFound
	}
	if err != nil {
		return "", meta, err
	}
	if json.Unmarshal(b, &meta) != nil {
		return "", meta, apperr.Unprocessable
	}
	return dir, meta, nil
}
func sanitizeFilename(name string) string {
	name = filepath.Base(strings.ReplaceAll(name, "\\", "/"))
	name = strings.Map(func(r rune) rune {
		if r < 32 || strings.ContainsRune(`/\\:*?\"<>|`, r) {
			return '-'
		}
		return r
	}, name)
	name = strings.TrimSpace(name)
	if name == "" {
		name = "download.bin"
	}
	rr := []rune(name)
	if len(rr) > 200 {
		name = string(rr[:200])
	}
	return name
}
func uniqueObjectKey(ctx context.Context, d storagepkg.Driver, key string) (string, error) {
	if _, err := d.Stat(ctx, key); err != nil {
		return key, nil
	}
	ext := filepath.Ext(key)
	base := strings.TrimSuffix(key, ext)
	for i := 1; i <= 999; i++ {
		candidate := fmt.Sprintf("%s-%d%s", base, i, ext)
		if _, err := d.Stat(ctx, candidate); err != nil {
			return candidate, nil
		}
	}
	return "", errors.New("cannot allocate unique object key")
}

type DownloadService struct {
	store    *repo.Store
	storages *StorageManager
	mu       sync.Mutex
	seen     map[string]time.Time
}
type DownloadTarget struct {
	App     model.App
	Release model.Release
	Asset   model.Asset
	Source  model.DownloadSource
	Driver  storagepkg.Driver
	Storage model.Storage
}

func NewDownloadService(store *repo.Store, sm *StorageManager) *DownloadService {
	return &DownloadService{store: store, storages: sm, seen: map[string]time.Time{}}
}
func (s *DownloadService) Resolve(ctx context.Context, assetID uint64, sourceID *uint64, ip, ua, referer string) (DownloadTarget, error) {
	var out DownloadTarget
	if err := s.store.DB.WithContext(ctx).Preload("Sources", func(db *gorm.DB) *gorm.DB { return db.Where("is_enabled = ?", true).Order("priority ASC") }).First(&out.Asset, assetID).Error; err != nil {
		return out, apperr.NotFound
	}
	if err := s.store.DB.WithContext(ctx).First(&out.Release, out.Asset.ReleaseID).Error; err != nil || out.Release.Status != model.StatusPublished {
		return out, apperr.NotFound
	}
	if err := s.store.DB.WithContext(ctx).First(&out.App, out.Release.AppID).Error; err != nil || out.App.Status != model.StatusPublished || out.App.Type != model.AppTypeSelf {
		return out, apperr.NotFound
	}
	for _, src := range out.Asset.Sources {
		if sourceID == nil || src.ID == *sourceID {
			out.Source = src
			break
		}
	}
	if out.Source.ID == 0 {
		return out, apperr.NotFound
	}
	if out.Source.SourceType == model.SourceManaged {
		if out.Source.StorageID == nil {
			return out, apperr.Unprocessable
		}
		d, st, err := s.storages.Driver(ctx, *out.Source.StorageID)
		if err != nil {
			return out, err
		}
		out.Driver = d
		out.Storage = st
	}
	s.count(ctx, &out, ip, ua, referer)
	return out, nil
}
func (s *DownloadService) count(ctx context.Context, t *DownloadTarget, ip, ua, referer string) {
	key := ip + ":" + strconv.FormatUint(t.Asset.ID, 10)
	s.mu.Lock()
	last, ok := s.seen[key]
	if ok && time.Since(last) < 10*time.Minute {
		s.mu.Unlock()
		return
	}
	if len(s.seen) > 8192 {
		for k, v := range s.seen {
			if time.Since(v) >= 10*time.Minute {
				delete(s.seen, k)
			}
		}
	}
	s.seen[key] = time.Now()
	s.mu.Unlock()
	_ = s.store.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, x := range []struct {
			m  any
			id uint64
		}{{&model.App{}, t.App.ID}, {&model.Release{}, t.Release.ID}, {&model.Asset{}, t.Asset.ID}, {&model.DownloadSource{}, t.Source.ID}} {
			if err := tx.Model(x.m).Where("id = ?", x.id).UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error; err != nil {
				return err
			}
		}
		return tx.Create(&model.DownloadLog{AppID: t.App.ID, ReleaseID: t.Release.ID, AssetID: t.Asset.ID, SourceID: t.Source.ID, IP: truncateIP(ip), UA: truncate(ua, 300), Referer: truncate(referer, 300)}).Error
	})
}
func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}
func truncateIP(raw string) string {
	ip := net.ParseIP(raw)
	if ip == nil {
		return ""
	}
	if v4 := ip.To4(); v4 != nil {
		return fmt.Sprintf("%d.%d.%d.0", v4[0], v4[1], v4[2])
	}
	return ip.Mask(net.CIDRMask(48, 128)).String()
}

type Settings struct{ store *repo.Store }

func NewSettings(store *repo.Store) *Settings { return &Settings{store: store} }
func (s *Settings) All(ctx context.Context) (map[string]string, error) {
	var rows []model.Setting
	if err := s.store.DB.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, r := range rows {
		out[r.Key] = r.Value
	}
	return out, nil
}
func (s *Settings) Update(ctx context.Context, values map[string]string) error {
	return s.store.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for k, v := range values {
			if !allowedSetting(k) {
				return apperr.BadRequest
			}
			if err := tx.Save(&model.Setting{Key: k, Value: v}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func allowedSetting(k string) bool {
	for _, p := range []string{"site.", "theme.", "download.", "content.", "seo.", "custom.", "privacy."} {
		if strings.HasPrefix(k, p) {
			return true
		}
	}
	return false
}
