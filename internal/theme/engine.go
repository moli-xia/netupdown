package theme

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/moli-xia/netupdown/internal/assets"
	"github.com/moli-xia/netupdown/internal/model"
	"github.com/rs/xid"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

type Meta struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Version       string        `json:"version"`
	Author        string        `json:"author"`
	Homepage      string        `json:"homepage"`
	Description   string        `json:"description"`
	Preview       string        `json:"preview"`
	MinAppVersion string        `json:"min_app_version"`
	Settings      []SettingSpec `json:"settings"`
	Active        bool          `json:"active"`
	Builtin       bool          `json:"builtin"`
}
type SettingSpec struct {
	Key     string              `json:"key"`
	Label   string              `json:"label"`
	Type    string              `json:"type"`
	Default any                 `json:"default"`
	Options []map[string]string `json:"options,omitempty"`
}

type Engine struct {
	dataDir string
	dev     bool
	mu      sync.RWMutex
	active  string
	pages   map[string]*template.Template
}

func New(dataDir string, dev bool, active string) (*Engine, error) {
	if active == "" {
		active = "aurora"
	}
	e := &Engine{dataDir: dataDir, dev: dev, active: active}
	if err := e.reload(); err != nil {
		e.active = "aurora"
		if fallbackErr := e.reload(); fallbackErr != nil {
			return nil, errors.Join(err, fallbackErr)
		}
	}
	return e, nil
}
func (e *Engine) Active() string { e.mu.RLock(); defer e.mu.RUnlock(); return e.active }
func (e *Engine) Render(page string, data any) ([]byte, error) {
	if e.dev {
		if err := e.reload(); err != nil {
			return nil, err
		}
	}
	e.mu.RLock()
	t := e.pages[page]
	e.mu.RUnlock()
	if t == nil {
		return nil, fmt.Errorf("unknown theme page %q", page)
	}
	var b bytes.Buffer
	if err := t.ExecuteTemplate(&b, "base", data); err != nil {
		return nil, fmt.Errorf("render %s: %w", page, err)
	}
	return b.Bytes(), nil
}
func (e *Engine) Activate(id string) error {
	if !themeIDPattern.MatchString(id) {
		return errors.New("invalid theme id")
	}
	old := e.Active()
	e.mu.Lock()
	e.active = id
	e.mu.Unlock()
	if err := e.reload(); err != nil {
		e.mu.Lock()
		e.active = old
		e.mu.Unlock()
		_ = e.reload()
		return err
	}
	return nil
}

func (e *Engine) reload() error {
	funcs := template.FuncMap{"asset": func(path string) string { return "/themes/" + e.Active() + "/static/" + strings.TrimLeft(path, "/") }, "initial": func(s string) string {
		r := []rune(strings.TrimSpace(s))
		if len(r) == 0 {
			return "N"
		}
		return strings.ToUpper(string(r[0]))
	}, "join": strings.Join, "formatBytes": formatBytes, "markdown": markdown, "osName": osName, "osIcon": osIcon, "archName": archName, "kindName": kindName, "numFmt": numFmt, "channelName": channelName, "date": formatDate}
	e.mu.RLock()
	active := e.active
	e.mu.RUnlock()
	pageNames := []string{"index", "list", "detail", "releases", "download", "page", "error"}
	compiled := map[string]*template.Template{}
	for _, page := range pageNames {
		base, body, err := e.templatePair(active, page)
		if err != nil {
			return err
		}
		t, err := template.New(page).Funcs(funcs).Parse(string(base))
		if err == nil {
			t, err = t.Parse(string(body))
		}
		if err != nil {
			return fmt.Errorf("compile theme %s page %s: %w", active, page, err)
		}
		compiled[page] = t
	}
	e.mu.Lock()
	e.pages = compiled
	e.mu.Unlock()
	return nil
}
func (e *Engine) templatePair(id, page string) ([]byte, []byte, error) {
	builtinBase, err := fs.ReadFile(assets.Embedded, "aurora/templates/base.html")
	if err != nil {
		return nil, nil, err
	}
	builtinBody, err := fs.ReadFile(assets.Embedded, "aurora/templates/"+page+".html")
	if err != nil {
		return nil, nil, err
	}
	root := filepath.Join(e.dataDir, "themes", id, "templates")
	customBase, baseErr := os.ReadFile(filepath.Join(root, "layouts", "base.html"))
	if baseErr != nil {
		customBase, baseErr = os.ReadFile(filepath.Join(root, "base.html"))
	}
	customBody, bodyErr := os.ReadFile(filepath.Join(root, page+".html"))
	if baseErr == nil && bodyErr == nil {
		return customBase, customBody, nil
	}
	if id == "aurora" {
		if baseErr == nil {
			builtinBase = customBase
		}
		if bodyErr == nil {
			builtinBody = customBody
		}
	}
	return builtinBase, builtinBody, nil
}

func (e *Engine) List() ([]Meta, error) {
	builtinRaw, err := fs.ReadFile(assets.Embedded, "aurora/theme.json")
	if err != nil {
		return nil, err
	}
	var builtin Meta
	if err := json.Unmarshal(builtinRaw, &builtin); err != nil {
		return nil, err
	}
	builtin.Builtin = true
	items := map[string]Meta{"aurora": builtin}
	root := filepath.Join(e.dataDir, "themes")
	entries, _ := os.ReadDir(root)
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") || strings.HasSuffix(entry.Name(), ".bak") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(root, entry.Name(), "theme.json"))
		if err != nil {
			continue
		}
		var meta Meta
		if json.Unmarshal(raw, &meta) == nil && themeIDPattern.MatchString(meta.ID) {
			items[meta.ID] = meta
		}
	}
	out := make([]Meta, 0, len(items))
	active := e.Active()
	for _, item := range items {
		item.Active = item.ID == active
		item.Builtin = item.ID == "aurora"
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Builtin {
			return true
		}
		return out[i].Name < out[j].Name
	})
	return out, nil
}

var themeIDPattern = regexp.MustCompile(`^[a-z0-9-]{2,40}$`)
var allowedExt = map[string]bool{".html": true, ".css": true, ".js": true, ".json": true, ".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".svg": true, ".woff2": true, ".txt": true, ".md": true}

func (e *Engine) Install(zipPath string) (Meta, error) {
	var meta Meta
	info, err := os.Stat(zipPath)
	if err != nil {
		return meta, err
	}
	if info.Size() > 20*1024*1024 {
		return meta, errors.New("theme archive exceeds 20MB")
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		return meta, err
	}
	defer zr.Close()
	if len(zr.File) > 500 {
		return meta, errors.New("theme archive has too many entries")
	}
	var total uint64
	var metaRaw []byte
	for _, f := range zr.File {
		total += f.UncompressedSize64
		if total > 50*1024*1024 {
			return meta, errors.New("uncompressed theme exceeds 50MB")
		}
		clean := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(f.Name)), "./")
		if clean == "theme.json" {
			rc, err := f.Open()
			if err != nil {
				return meta, err
			}
			metaRaw, err = io.ReadAll(io.LimitReader(rc, 1024*1024))
			_ = rc.Close()
			if err != nil {
				return meta, err
			}
		}
		if !f.FileInfo().IsDir() && !allowedExt[strings.ToLower(filepath.Ext(clean))] {
			return meta, fmt.Errorf("theme file type is not allowed: %s", clean)
		}
	}
	if json.Unmarshal(metaRaw, &meta) != nil || !themeIDPattern.MatchString(meta.ID) || meta.Name == "" || meta.Version == "" {
		return meta, errors.New("theme.json is missing or invalid")
	}
	themeRoot := filepath.Join(e.dataDir, "themes")
	if err := os.MkdirAll(themeRoot, 0o700); err != nil {
		return meta, err
	}
	tmp := filepath.Join(themeRoot, ".install-"+meta.ID+"-"+xid.New().String())
	if err := os.MkdirAll(tmp, 0o700); err != nil {
		return meta, err
	}
	defer os.RemoveAll(tmp)
	for _, f := range zr.File {
		clean := filepath.Clean(filepath.FromSlash(f.Name))
		target := filepath.Join(tmp, clean)
		rel, err := filepath.Rel(tmp, target)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return meta, errors.New("theme archive contains unsafe path")
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o700); err != nil {
				return meta, err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return meta, err
		}
		rc, err := f.Open()
		if err != nil {
			return meta, err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
		if err == nil {
			_, err = io.Copy(out, io.LimitReader(rc, 51*1024*1024))
		}
		_ = rc.Close()
		if out != nil {
			_ = out.Close()
		}
		if err != nil {
			return meta, err
		}
	}
	final := filepath.Join(themeRoot, meta.ID)
	backup := final + ".bak"
	_ = os.RemoveAll(backup)
	if _, err := os.Stat(final); err == nil {
		if err := os.Rename(final, backup); err != nil {
			return meta, err
		}
	}
	if err := os.Rename(tmp, final); err != nil {
		_ = os.Rename(backup, final)
		return meta, err
	}
	_ = os.RemoveAll(backup)
	return meta, nil
}
func (e *Engine) Delete(id string) error {
	if id == "aurora" || id == e.Active() {
		return errors.New("built-in or active theme cannot be deleted")
	}
	if !themeIDPattern.MatchString(id) {
		return errors.New("invalid theme id")
	}
	return os.RemoveAll(filepath.Join(e.dataDir, "themes", id))
}
func (e *Engine) Static(id, path string) ([]byte, string, error) {
	if !themeIDPattern.MatchString(id) {
		return nil, "", fs.ErrNotExist
	}
	clean := filepath.Clean(filepath.FromSlash(path))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return nil, "", fs.ErrNotExist
	}
	if raw, err := os.ReadFile(filepath.Join(e.dataDir, "themes", id, "static", clean)); err == nil {
		return raw, contentType(clean), nil
	}
	if id == "aurora" {
		raw, err := fs.ReadFile(assets.Embedded, "aurora/static/"+filepath.ToSlash(clean))
		return raw, contentType(clean), err
	}
	return nil, "", fs.ErrNotExist
}
func contentType(path string) string {
	if value := mime.TypeByExtension(filepath.Ext(path)); value != "" {
		return value
	}
	return "application/octet-stream"
}
func markdown(raw string) template.HTML {
	var b bytes.Buffer
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	if err := md.Convert([]byte(raw), &b); err != nil {
		return template.HTML(template.HTMLEscapeString(raw))
	}
	return template.HTML(bluemonday.UGCPolicy().SanitizeBytes(b.Bytes()))
}
func formatBytes(n int64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	v := float64(n)
	i := 0
	for v >= 1024 && i < len(units)-1 {
		v /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%d %s", n, units[i])
	}
	return fmt.Sprintf("%.1f %s", v, units[i])
}
func osName(v string) string {
	m := map[string]string{"windows": "Windows", "macos": "macOS", "linux": "Linux", "android": "Android", "ios": "iOS", "web": "Web", "any": "全平台"}
	if x := m[v]; x != "" {
		return x
	}
	return v
}
func archName(v string) string {
	m := map[string]string{"amd64": "x64", "386": "x86", "arm64": "ARM64", "arm": "ARM", "universal": "通用", "any": "通用"}
	if x := m[v]; x != "" {
		return x
	}
	return v
}
func kindName(v int8) string {
	switch v {
	case 1:
		return "安装包"
	case 2:
		return "便携版"
	case 3:
		return "压缩包"
	case 4:
		return "补丁"
	default:
		return "文件"
	}
}
func numFmt(n int64) string {
	switch {
	case n >= 100_000_000:
		return strings.TrimSuffix(fmt.Sprintf("%.1f", float64(n)/1e8), ".0") + " 亿"
	case n >= 10_000:
		return strings.TrimSuffix(fmt.Sprintf("%.1f", float64(n)/1e4), ".0") + " 万"
	default:
		return fmt.Sprintf("%d", n)
	}
}

// formatDate accepts time.Time or *time.Time so templates can pass either.
func formatDate(v any) string {
	switch t := v.(type) {
	case *time.Time:
		if t == nil {
			return ""
		}
		return t.Format("2006-01-02")
	case time.Time:
		return t.Format("2006-01-02")
	default:
		return ""
	}
}

// osIcon returns a small inline SVG glyph for a platform value; decorative only.
func osIcon(v string) template.HTML {
	const head = `<svg class="os-icon" viewBox="0 0 16 16" width="14" height="14" fill="currentColor" aria-hidden="true">`
	body := map[string]string{
		"windows": `<path d="M2 3.6 7.4 2.8v4.8H2zM8.4 2.6 14 1.8v6H8.4zM2 8.6h5.4v4.8L2 12.6zM8.4 8.6H14v6l-5.6-.8z"/>`,
		"macos":   `<path d="M11.6 8.5c0-1.6 1.3-2.4 1.4-2.4-.8-1.1-2-1.3-2.4-1.3-1-.1-2 .6-2.5.6s-1.3-.6-2.1-.6C4.7 4.9 3.5 5.6 2.9 6.8c-1.3 2.3-.3 5.7 1 7.6.6.9 1.3 1.9 2.3 1.9.9 0 1.3-.6 2.4-.6s1.4.6 2.4.6 1.6-.9 2.2-1.8c.7-1 1-2 1-2.1-.1 0-1.9-.7-1.9-2.9zM9.9 3.7c.5-.6.8-1.4.7-2.2-.7 0-1.5.5-2 1.1-.4.5-.8 1.3-.7 2.1.8.1 1.6-.4 2-1z"/>`,
		"ios":     `<path d="M11.6 8.5c0-1.6 1.3-2.4 1.4-2.4-.8-1.1-2-1.3-2.4-1.3-1-.1-2 .6-2.5.6s-1.3-.6-2.1-.6C4.7 4.9 3.5 5.6 2.9 6.8c-1.3 2.3-.3 5.7 1 7.6.6.9 1.3 1.9 2.3 1.9.9 0 1.3-.6 2.4-.6s1.4.6 2.4.6 1.6-.9 2.2-1.8c.7-1 1-2 1-2.1-.1 0-1.9-.7-1.9-2.9zM9.9 3.7c.5-.6.8-1.4.7-2.2-.7 0-1.5.5-2 1.1-.4.5-.8 1.3-.7 2.1.8.1 1.6-.4 2-1z"/>`,
		"linux":   `<path d="M8 1.4c-1.9 0-3 1.5-3 3.4 0 1.2-.5 2-1 3-.6 1.1-1.2 2.3-1.2 3.6 0 .9.5 1.6 1.3 1.9.6.9 1.7 1.6 3.9 1.6s3.3-.7 3.9-1.6c.8-.3 1.3-1 1.3-1.9 0-1.3-.6-2.5-1.2-3.6-.5-1-1-1.8-1-3 0-1.9-1.1-3.4-3-3.4zM6.6 5.2a.7.7 0 1 1 0 1.4.7.7 0 0 1 0-1.4zm2.8 0a.7.7 0 1 1 0 1.4.7.7 0 0 1 0-1.4zM8 7.4l1.5.9-1.5.9-1.5-.9z"/>`,
		"android": `<path d="M4.6 6.2h6.8v4.4c0 .8-.6 1.4-1.4 1.4h-.4v2a.8.8 0 0 1-1.6 0v-2h-.4v2a.8.8 0 0 1-1.6 0v-2H6c-.8 0-1.4-.6-1.4-1.4zM3 6.5a.8.8 0 0 1 .8.8v2.9a.8.8 0 0 1-1.6 0V7.3a.8.8 0 0 1 .8-.8zm10 0a.8.8 0 0 1 .8.8v2.9a.8.8 0 0 1-1.6 0V7.3a.8.8 0 0 1 .8-.8zM5.4 3.1l-.7-1a.3.3 0 0 1 .5-.3l.8 1.1a4 4 0 0 1 4 0l.8-1.1a.3.3 0 0 1 .5.3l-.7 1c.8.6 1.4 1.5 1.4 2.5H4c0-1 .6-1.9 1.4-2.5zM6.4 4.4a.4.4 0 1 0 0 .8.4.4 0 0 0 0-.8zm3.2 0a.4.4 0 1 0 0 .8.4.4 0 0 0 0-.8z"/>`,
		"web":     `<path fill="none" stroke="currentColor" stroke-width="1.3" d="M8 1.8a6.2 6.2 0 1 0 0 12.4A6.2 6.2 0 0 0 8 1.8zm0 0c-1.7 0-3 2.8-3 6.2s1.3 6.2 3 6.2 3-2.8 3-6.2-1.3-6.2-3-6.2zM2 8h12"/>`,
	}
	b, ok := body[v]
	if !ok {
		b = `<path d="M4 2.5h8a1.5 1.5 0 0 1 1.5 1.5v8A1.5 1.5 0 0 1 12 13.5H4A1.5 1.5 0 0 1 2.5 12V4A1.5 1.5 0 0 1 4 2.5zm1 3.2h6M5 8h6M5 10.3h4" fill="none" stroke="currentColor" stroke-width="1.3" stroke-linecap="round"/>`
	}
	return template.HTML(head + b + `</svg>`)
}
func channelName(v model.Channel) string {
	switch v {
	case model.ChannelBeta:
		return "beta"
	case model.ChannelAlpha:
		return "alpha"
	default:
		return "stable"
	}
}
