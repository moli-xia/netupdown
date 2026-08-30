package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/moli-xia/netupdown/internal/assets"
	"github.com/moli-xia/netupdown/internal/auth"
	"github.com/moli-xia/netupdown/internal/config"
	appmiddleware "github.com/moli-xia/netupdown/internal/middleware"
	"github.com/moli-xia/netupdown/internal/model"
	"github.com/moli-xia/netupdown/internal/pkg/apperr"
	"github.com/moli-xia/netupdown/internal/pkg/cryptoutil"
	"github.com/moli-xia/netupdown/internal/pkg/resp"
	"github.com/moli-xia/netupdown/internal/repo"
	"github.com/moli-xia/netupdown/internal/service"
	"github.com/moli-xia/netupdown/internal/storage"
	"github.com/moli-xia/netupdown/internal/theme"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Server struct {
	cfg           config.Config
	version       string
	store         *repo.Store
	auth          *auth.Service
	catalog       *service.Catalog
	storages      *service.StorageManager
	uploads       *service.UploadService
	downloads     *service.DownloadService
	settings      *service.Settings
	theme         *theme.Engine
	router        *gin.Engine
	publicLimiter *appmiddleware.IPLimiter
	loginLimiter  *appmiddleware.IPLimiter
}

func New(cfg config.Config, version string, store *repo.Store, authSvc *auth.Service, storageSealer *cryptoutil.Sealer) (*Server, error) {
	var activeTheme model.Setting
	_ = store.DB.Where("key = ?", "theme.active").First(&activeTheme).Error
	themes, err := theme.New(cfg.DataDir, cfg.Theme.Dev, activeTheme.Value)
	if err != nil {
		return nil, err
	}
	sm := service.NewStorageManager(store, cfg.DataDir, storageSealer)
	s := &Server{cfg: cfg, version: version, store: store, auth: authSvc, catalog: service.NewCatalog(store, cfg.Server.BaseURL), storages: sm, settings: service.NewSettings(store), theme: themes, publicLimiter: appmiddleware.NewIPLimiter(cfg.RateLimit.PublicPerMin, 20), loginLimiter: appmiddleware.NewIPLimiter(cfg.RateLimit.LoginPerMin, 3)}
	s.uploads = service.NewUploadService(cfg, store, sm)
	s.downloads = service.NewDownloadService(store, sm)
	s.router = s.routes()
	return s, nil
}
func (s *Server) Handler() http.Handler { return s.router }

func (s *Server) routes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), s.requestMiddleware())
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "version": s.version}) })
	r.GET("/", s.home)
	r.GET("/apps", s.appList)
	r.GET("/search", s.appList)
	r.GET("/categories/:slug", s.appList)
	r.GET("/apps/:slug/preview", s.appPreview)
	r.GET("/apps/:slug", s.appDetail)
	r.GET("/apps/:slug/releases", s.releaseHistory)
	r.GET("/pages/:slug", s.page)
	r.GET("/d/:id", s.download)
	r.GET("/feed.xml", s.feed)
	r.GET("/sitemap.xml", s.sitemap)
	r.GET("/robots.txt", s.robots)
	r.GET("/themes/:theme/static/*path", s.themeStatic)
	_ = os.MkdirAll(filepath.Join(s.cfg.DataDir, "uploads"), 0o700)
	r.StaticFS("/uploads", gin.Dir(filepath.Join(s.cfg.DataDir, "uploads"), false))
	s.adminStatic(r)
	v1 := r.Group("/api/v1")
	v1.Use(s.limit(s.publicLimiter))
	{
		v1.GET("/apps", s.publicApps)
		v1.GET("/apps/:slug", s.publicApp)
		v1.GET("/apps/:slug/releases", s.publicReleases)
		v1.GET("/apps/:slug/releases/latest", s.latestRelease)
		v1.GET("/apps/:slug/check-update", s.checkUpdate)
		v1.GET("/categories", s.publicCategories)
	}
	api := r.Group("/api/admin")
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", s.limit(s.loginLimiter), s.login)
		authGroup.POST("/refresh", s.refresh)
		authGroup.POST("/logout", s.logout)
	}
	secured := api.Group("")
	secured.Use(s.requireAuth())
	{
		secured.GET("/auth/profile", s.profile)
		secured.PUT("/auth/profile", s.updateProfile)
		secured.PUT("/auth/password", s.changePassword)
		secured.GET("/auth/sessions", s.sessions)
		secured.DELETE("/auth/sessions/:id", s.revokeSession)
		secured.DELETE("/auth/sessions", s.revokeAllSessions)
		secured.GET("/apps", s.adminApps)
		secured.POST("/apps", s.saveApp)
		secured.GET("/apps/:id", s.getApp)
		secured.PUT("/apps/:id", s.saveApp)
		secured.DELETE("/apps/:id", s.deleteApp)
		secured.POST("/apps/:id/preview", s.previewApp)
		secured.POST("/apps/:id/publish", s.publishApp(true))
		secured.POST("/apps/:id/unpublish", s.publishApp(false))
		secured.GET("/apps/:id/releases", s.adminReleases)
		secured.POST("/apps/:id/releases", s.saveRelease)
		secured.PUT("/releases/:id", s.saveRelease)
		secured.DELETE("/releases/:id", s.deleteRelease)
		secured.POST("/releases/:id/publish", s.publishRelease)
		secured.POST("/releases/:id/assets", s.saveAsset)
		secured.PUT("/assets/:id", s.saveAsset)
		secured.DELETE("/assets/:id", s.deleteAsset)
		secured.POST("/assets/:id/sources", s.saveSource)
		secured.PUT("/sources/:id", s.saveSource)
		secured.DELETE("/sources/:id", s.deleteSource)
		secured.GET("/categories", s.adminCategories)
		secured.POST("/categories", s.saveCategory)
		secured.PUT("/categories/:id", s.saveCategory)
		secured.DELETE("/categories/:id", s.deleteCategory)
		secured.POST("/uploads/image", s.uploadImage)
		secured.POST("/uploads/file", s.uploadFile)
		secured.POST("/uploads/init", s.uploadInit)
		secured.PUT("/uploads/:id/chunks/:index", s.uploadChunk)
		secured.POST("/uploads/:id/complete", s.uploadComplete)
		secured.DELETE("/uploads/:id", s.uploadAbort)
		secured.GET("/storages", s.listStorages)
		secured.POST("/storages", s.saveStorage)
		secured.PUT("/storages/:id", s.saveStorage)
		secured.DELETE("/storages/:id", s.deleteStorage)
		secured.POST("/storages/:id/test", s.testStorage)
		secured.GET("/settings", s.getSettings)
		secured.PUT("/settings", s.updateSettings)
		secured.GET("/themes", s.listThemes)
		secured.POST("/themes/upload", s.uploadTheme)
		secured.POST("/themes/:id/activate", s.activateTheme)
		secured.PUT("/themes/:id/config", s.configureTheme)
		secured.DELETE("/themes/:id", s.deleteTheme)
		secured.GET("/stats/overview", s.statsOverview)
		secured.GET("/pages", s.adminPages)
		secured.POST("/pages", s.savePage)
		secured.PUT("/pages/:id", s.savePage)
		secured.DELETE("/pages/:id", s.deletePage)
		secured.GET("/logs/operations", s.operationLogs)
		secured.GET("/logs/downloads", s.downloadLogs)
	}
	return r
}

func (s *Server) requestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		id := c.GetHeader("X-Request-Id")
		if id == "" {
			id = xid.New().String()
		}
		c.Header("X-Request-Id", id)
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Next()
		slog.Info("request", "req_id", id, "method", c.Request.Method, "path", c.Request.URL.Path, "status", c.Writer.Status(), "latency_ms", time.Since(start).Milliseconds())
	}
}
func (s *Server) requireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if raw == "" {
			resp.Fail(c, apperr.Unauthorized)
			c.Abort()
			return
		}
		claims, err := s.auth.Parse(raw)
		if err != nil {
			resp.Fail(c, err)
			c.Abort()
			return
		}
		if claims.Role != 1 {
			resp.Fail(c, apperr.Forbidden)
			c.Abort()
			return
		}
		uid, _ := strconv.ParseUint(claims.Subject, 10, 64)
		c.Set("uid", uid)
		c.Next()
	}
}
func (s *Server) limit(limiter *appmiddleware.IPLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow(s.clientIP(c)) {
			resp.Fail(c, apperr.New(10009, http.StatusTooManyRequests, "请求过于频繁"))
			c.Abort()
			return
		}
		c.Next()
	}
}
func (s *Server) clientIP(c *gin.Context) string {
	if s.cfg.Server.BehindProxy {
		if x := c.GetHeader("X-Forwarded-For"); x != "" {
			parts := strings.Split(x, ",")
			return strings.TrimSpace(parts[0])
		}
	}
	return c.ClientIP()
}
func idParam(c *gin.Context) uint64 { id, _ := strconv.ParseUint(c.Param("id"), 10, 64); return id }

type siteCtx struct {
	Title, Subtitle, Description, Keywords, Logo, Favicon, ICP string
	// Footer 与注入代码由站长在后台维护，属可信内容，按 HTML 原样输出。
	Footer, HeadInject, FootInject template.HTML
}

func (s *Server) baseData(c *gin.Context, title, description string) gin.H {
	settings, _ := s.settings.All(c.Request.Context())
	site := siteCtx{Title: or(settings["site.title"], "造物工坊"), Subtitle: or(settings["site.subtitle"], "独立开发者的软件发布与更新中心"), Description: settings["site.description"], Keywords: settings["site.keywords"], Logo: settings["site.logo"], Favicon: settings["site.favicon"], ICP: settings["site.icp"], Footer: template.HTML(settings["site.footer"]), HeadInject: template.HTML(settings["custom.head"]), FootInject: template.HTML(settings["custom.foot"])}
	active := s.theme.Active()
	var themeConfig map[string]any
	_ = json.Unmarshal([]byte(settings["theme.cfg."+active]), &themeConfig)
	if themeConfig == nil {
		themeConfig = map[string]any{}
	}
	return gin.H{"Site": site, "Theme": gin.H{"ID": active, "StaticBase": "/themes/" + active + "/static", "Config": themeConfig}, "Year": time.Now().Year(), "PageTitle": title + " | " + site.Title, "Description": or(description, site.Description), "Canonical": s.cfg.Server.BaseURL + c.Request.URL.Path}
}
func (s *Server) render(c *gin.Context, status int, page string, data gin.H) {
	body, err := s.theme.Render(page, data)
	if err != nil {
		slog.Error("render failed", "page", page, "error", err)
		c.String(500, "render error")
		return
	}
	c.Data(status, "text/html; charset=utf-8", body)
}
func or(v, f string) string {
	if v == "" {
		return f
	}
	return v
}

func (s *Server) home(c *gin.Context) {
	ctx := c.Request.Context()
	all, _ := s.catalog.Apps(ctx, service.AppQuery{Page: 1, PageSize: 8, Type: "self"})
	featured, _ := s.catalog.Apps(ctx, service.AppQuery{Page: 1, PageSize: 6, Type: "self", Featured: true})
	apps := all.List.([]model.App)
	feat, _ := featured.List.([]model.App)
	cats, _ := s.catalog.Categories(ctx)
	var total int64
	_ = s.store.DB.WithContext(ctx).Model(&model.App{}).Where("status = ? AND type = ?", model.StatusPublished, model.AppTypeSelf).Select("COALESCE(SUM(download_count),0)").Scan(&total).Error
	d := s.baseData(c, "首页", "")
	d["Latest"] = apps
	d["Featured"] = feat
	d["Categories"] = cats
	d["AppCount"] = all.Total
	d["TotalDownloads"] = total
	s.render(c, 200, "index", d)
}

type pageLink struct {
	N       int
	URL     string
	Current bool
}

func (s *Server) appList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	category := c.Query("category")
	byPath := c.Param("slug") != ""
	if byPath {
		category = c.Param("slug")
	}
	q := service.AppQuery{Page: page, PageSize: 24, Category: category, Platform: c.Query("platform"), Type: "self", Sort: c.Query("sort"), Q: c.Query("q")}
	result, err := s.catalog.Apps(c.Request.Context(), q)
	if err != nil {
		s.renderError(c, 500, "加载失败")
		return
	}
	cats, _ := s.catalog.Categories(c.Request.Context())
	title := "全部产品"
	categoryName := ""
	for _, cat := range cats {
		if cat.Slug == category {
			categoryName = cat.Name
			title = cat.Name
		}
	}
	if q.Q != "" {
		title = "搜索「" + q.Q + "」"
	}
	d := s.baseData(c, title, "")
	d["Title"] = title
	d["Apps"] = result.List
	d["Categories"] = cats
	d["Query"] = q.Q
	d["Cat"] = category
	d["CategoryName"] = categoryName
	d["Platform"] = q.Platform
	d["Type"] = q.Type
	d["Sort"] = or(q.Sort, "latest")
	d["Total"] = result.Total
	d["Page"] = result.Page
	pages := int((result.Total + int64(result.PageSize) - 1) / int64(result.PageSize))
	d["TotalPages"] = pages
	d["HasPrev"] = page > 1
	d["HasNext"] = page < pages
	mkURL := func(p int) string {
		v := url.Values{}
		if q.Q != "" {
			v.Set("q", q.Q)
		}
		if !byPath && category != "" {
			v.Set("category", category)
		}
		if q.Platform != "" {
			v.Set("platform", q.Platform)
		}
		if q.Sort != "" {
			v.Set("sort", q.Sort)
		}
		if p > 1 {
			v.Set("page", strconv.Itoa(p))
		}
		if enc := v.Encode(); enc != "" {
			return c.Request.URL.Path + "?" + enc
		}
		return c.Request.URL.Path
	}
	d["PrevURL"] = mkURL(page - 1)
	d["NextURL"] = mkURL(page + 1)
	start := page - 2
	if start < 1 {
		start = 1
	}
	end := start + 4
	if end > pages {
		end = pages
		if start = end - 4; start < 1 {
			start = 1
		}
	}
	links := make([]pageLink, 0, 5)
	for i := start; i <= end; i++ {
		links = append(links, pageLink{N: i, URL: mkURL(i), Current: i == page})
	}
	d["Pages"] = links
	s.render(c, 200, "list", d)
}
func (s *Server) appDetail(c *gin.Context) {
	app, err := s.catalog.AppBySlug(c.Request.Context(), c.Param("slug"), true)
	if err != nil {
		s.renderError(c, 404, "产品不存在")
		return
	}
	d := s.baseData(c, app.Name+" 下载", app.Tagline)
	d["App"] = app
	if app.Icon != "" {
		d["OGImage"] = s.absURL(app.Icon)
	}
	if len(app.Releases) > 0 {
		d["Release"] = app.Releases[0]
		d["Assets"] = app.Releases[0].Assets
		if len(app.Releases[0].Assets) > 0 {
			d["PrimaryAsset"] = app.Releases[0].Assets[0]
		}
	}
	s.render(c, 200, "detail", d)
}
func (s *Server) absURL(path string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return s.cfg.Server.BaseURL + path
}
func (s *Server) appPreview(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Header("Referrer-Policy", "no-referrer")
	slug := c.Param("slug")
	raw := c.Query("token")
	if raw != "" {
		if _, err := s.auth.ParsePreview(raw, slug); err != nil {
			s.renderError(c, http.StatusNotFound, "预览链接无效或已过期")
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("nud_preview", raw, 10*60, "/apps/"+slug+"/preview", "", strings.HasPrefix(s.cfg.Server.BaseURL, "https://"), true)
		c.Redirect(http.StatusFound, "/apps/"+url.PathEscape(slug)+"/preview")
		return
	}
	cookie, err := c.Cookie("nud_preview")
	if err != nil {
		s.renderError(c, http.StatusNotFound, "预览链接无效或已过期")
		return
	}
	raw = cookie
	if _, err := s.auth.ParsePreview(raw, slug); err != nil {
		s.renderError(c, http.StatusNotFound, "预览链接无效或已过期")
		return
	}
	app, err := s.catalog.AppBySlug(c.Request.Context(), slug, false)
	if err != nil {
		s.renderError(c, http.StatusNotFound, "产品不存在")
		return
	}
	d := s.baseData(c, app.Name+" 预览", app.Tagline)
	d["App"] = app
	d["Preview"] = true
	if len(app.Releases) > 0 {
		d["Release"] = app.Releases[0]
		d["Assets"] = app.Releases[0].Assets
		if len(app.Releases[0].Assets) > 0 {
			d["PrimaryAsset"] = app.Releases[0].Assets[0]
		}
	}
	s.render(c, http.StatusOK, "detail", d)
}
func (s *Server) releaseHistory(c *gin.Context) {
	app, err := s.catalog.AppBySlug(c.Request.Context(), c.Param("slug"), true)
	if err != nil {
		s.renderError(c, 404, "产品不存在")
		return
	}
	d := s.baseData(c, app.Name+" 历史版本", app.Tagline)
	d["App"] = app
	d["Releases"] = app.Releases
	s.render(c, 200, "releases", d)
}
func (s *Server) page(c *gin.Context) {
	var p model.Page
	err := s.store.DB.WithContext(c.Request.Context()).Where("slug = ? AND status = ?", c.Param("slug"), model.StatusPublished).First(&p).Error
	if err != nil {
		s.renderError(c, 404, "页面不存在")
		return
	}
	d := s.baseData(c, p.Title, p.SeoDescription)
	d["Doc"] = p
	s.render(c, 200, "page", d)
}
func (s *Server) renderError(c *gin.Context, status int, message string) {
	d := s.baseData(c, message, message)
	d["Status"] = status
	d["Message"] = message
	s.render(c, status, "error", d)
}

func (s *Server) themeStatic(c *gin.Context) {
	raw, contentType, err := s.theme.Static(c.Param("theme"), strings.TrimPrefix(c.Param("path"), "/"))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Header("Cache-Control", "public, max-age=604800")
	c.Data(http.StatusOK, contentType, raw)
}

func (s *Server) download(c *gin.Context) {
	assetID := idParam(c)
	var sourceID *uint64
	if raw := c.Query("source"); raw != "" {
		id, _ := strconv.ParseUint(raw, 10, 64)
		sourceID = &id
	}
	target, err := s.downloads.Resolve(c.Request.Context(), assetID, sourceID, s.clientIP(c), c.Request.UserAgent(), c.Request.Referer())
	if err != nil {
		s.renderError(c, 404, "下载文件不存在")
		return
	}
	if target.Source.SourceType == model.SourceExternal {
		if target.Source.ExtractCode != "" {
			d := s.baseData(c, "下载 "+target.App.Name, target.Asset.Name)
			d["App"] = target.App
			d["Asset"] = target.Asset
			d["Source"] = target.Source
			s.render(c, 200, "download", d)
			return
		}
		c.Redirect(http.StatusFound, target.Source.ExternalURL)
		return
	}
	if local, ok := target.Driver.(*storage.Local); ok {
		path, err := local.AbsPath(target.Source.ObjectKey)
		if err != nil {
			s.renderError(c, 404, "文件不存在")
			return
		}
		f, err := os.Open(path)
		if err != nil {
			s.renderError(c, 404, "文件不存在")
			return
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil {
			s.renderError(c, 404, "文件不存在")
			return
		}
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+url.PathEscape(target.Asset.FileName))
		http.ServeContent(c.Writer, c.Request, target.Asset.FileName, info.ModTime(), f)
		return
	}
	signed, err := target.Driver.PresignURL(c.Request.Context(), target.Source.ObjectKey, target.Asset.FileName, 30*time.Minute)
	if err != nil {
		s.renderError(c, 503, "下载源暂不可用")
		return
	}
	c.Redirect(302, signed)
}

func (s *Server) publicApps(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	out, err := s.catalog.Apps(c.Request.Context(), service.AppQuery{Page: page, PageSize: size, Category: c.Query("category"), Platform: c.Query("platform"), Type: c.Query("type"), Sort: c.Query("sort"), Q: c.Query("q")})
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, out)
}
func (s *Server) publicApp(c *gin.Context) {
	row, err := s.catalog.AppBySlug(c.Request.Context(), c.Param("slug"), true)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	sanitizeSources(&row)
	resp.OK(c, row)
}
func sanitizeSources(app *model.App) {
	for ri := range app.Releases {
		for ai := range app.Releases[ri].Assets {
			for si := range app.Releases[ri].Assets[ai].Sources {
				src := &app.Releases[ri].Assets[ai].Sources[si]
				src.ObjectKey = ""
				src.ExternalURL = ""
				src.StorageID = nil
			}
		}
	}
}
func (s *Server) publicReleases(c *gin.Context) {
	app, err := s.catalog.AppBySlug(c.Request.Context(), c.Param("slug"), true)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	for i := range app.Releases {
		for j := range app.Releases[i].Assets {
			for k := range app.Releases[i].Assets[j].Sources {
				app.Releases[i].Assets[j].Sources[k].ObjectKey = ""
				app.Releases[i].Assets[j].Sources[k].ExternalURL = ""
			}
		}
	}
	resp.OK(c, gin.H{"list": app.Releases, "page": 1, "page_size": len(app.Releases), "total": len(app.Releases)})
}
func (s *Server) latestRelease(c *gin.Context) {
	app, err := s.catalog.AppBySlug(c.Request.Context(), c.Param("slug"), true)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	channel := parseChannel(c.DefaultQuery("channel", "stable"))
	for _, rel := range app.Releases {
		if rel.Channel == channel {
			resp.OK(c, rel)
			return
		}
	}
	resp.Fail(c, apperr.NotFound)
}
func (s *Server) checkUpdate(c *gin.Context) {
	current := c.Query("version")
	if current == "" {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	var code *int
	if raw := c.Query("version_code"); raw != "" {
		v, e := strconv.Atoi(raw)
		if e != nil {
			resp.Fail(c, apperr.BadRequest)
			return
		}
		code = &v
	}
	out, err := s.catalog.CheckUpdate(c.Request.Context(), c.Param("slug"), current, c.Query("os"), c.Query("arch"), parseChannel(c.DefaultQuery("channel", "stable")), code)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, out)
}
func parseChannel(v string) model.Channel {
	switch v {
	case "beta":
		return model.ChannelBeta
	case "alpha":
		return model.ChannelAlpha
	default:
		return model.ChannelStable
	}
}
func (s *Server) publicCategories(c *gin.Context) {
	rows, err := s.catalog.Categories(c.Request.Context())
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, rows)
}

func (s *Server) login(c *gin.Context) {
	var in struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if c.ShouldBindJSON(&in) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	user, access, refreshToken, err := s.auth.Authenticate(c.Request.Context(), strings.ToLower(in.Username), in.Password, s.clientIP(c), c.Request.UserAgent())
	if err != nil {
		resp.Fail(c, err)
		return
	}
	s.setRefreshCookie(c, refreshToken)
	resp.OK(c, gin.H{"access_token": access, "expires_in": 7200, "user": user})
}
func (s *Server) refresh(c *gin.Context) {
	raw, err := c.Cookie("nud_rt")
	if err != nil {
		resp.Fail(c, apperr.Unauthorized)
		return
	}
	user, access, next, err := s.auth.Rotate(c.Request.Context(), raw, s.clientIP(c), c.Request.UserAgent())
	if err != nil {
		resp.Fail(c, err)
		return
	}
	s.setRefreshCookie(c, next)
	resp.OK(c, gin.H{"access_token": access, "expires_in": 7200, "user": user})
}
func (s *Server) logout(c *gin.Context) {
	if raw, err := c.Cookie("nud_rt"); err == nil {
		_ = s.auth.Revoke(c.Request.Context(), raw)
	}
	c.SetCookie("nud_rt", "", -1, "/api/admin/auth", "", strings.HasPrefix(s.cfg.Server.BaseURL, "https://"), true)
	resp.OK(c, nil)
}
func (s *Server) setRefreshCookie(c *gin.Context, v string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("nud_rt", v, 30*24*3600, "/api/admin/auth", "", strings.HasPrefix(s.cfg.Server.BaseURL, "https://"), true)
}
func (s *Server) profile(c *gin.Context) {
	u, err := s.auth.User(c.Request.Context(), c.MustGet("uid").(uint64))
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, u)
}
func (s *Server) updateProfile(c *gin.Context) {
	var in struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
	}
	if c.ShouldBindJSON(&in) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	u, err := s.auth.UpdateProfile(c.Request.Context(), c.MustGet("uid").(uint64), in.Nickname, in.Email, in.Avatar)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, u)
}
func (s *Server) changePassword(c *gin.Context) {
	var in struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if c.ShouldBindJSON(&in) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if err := s.auth.ChangePassword(c.Request.Context(), c.MustGet("uid").(uint64), in.OldPassword, in.NewPassword); err != nil {
		resp.Fail(c, err)
		return
	}
	s.audit(c, "auth.password", "user", c.MustGet("uid").(uint64))
	resp.OK(c, nil)
}
func (s *Server) sessions(c *gin.Context) {
	rows, err := s.auth.Sessions(c.Request.Context(), c.MustGet("uid").(uint64))
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, rows)
}
func (s *Server) revokeSession(c *gin.Context) {
	if err := s.auth.RevokeSession(c.Request.Context(), c.MustGet("uid").(uint64), idParam(c)); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}
func (s *Server) revokeAllSessions(c *gin.Context) {
	if err := s.auth.RevokeAll(c.Request.Context(), c.MustGet("uid").(uint64)); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}

func (s *Server) adminApps(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	out, err := s.catalog.Apps(c.Request.Context(), service.AppQuery{Page: page, PageSize: size, Q: c.Query("q"), Admin: true})
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, out)
}
func (s *Server) getApp(c *gin.Context) {
	row, err := s.catalog.AppByID(c.Request.Context(), idParam(c))
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, row)
}
func (s *Server) saveApp(c *gin.Context) {
	var row model.App
	if c.ShouldBindJSON(&row) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if c.Param("id") != "" {
		row.ID = idParam(c)
	}
	created, err := s.catalog.SaveApp(c.Request.Context(), &row)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	s.audit(c, "app.save", "app", row.ID)
	if created {
		resp.Created(c, row)
	} else {
		resp.OK(c, row)
	}
}
func (s *Server) previewApp(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	app, err := s.catalog.AppByID(c.Request.Context(), idParam(c))
	if err != nil {
		resp.Fail(c, err)
		return
	}
	token, err := s.auth.PreviewToken(app.ID, app.Slug)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	previewURL := "/apps/" + url.PathEscape(app.Slug) + "/preview?token=" + url.QueryEscape(token)
	resp.OK(c, gin.H{"url": previewURL})
}
func (s *Server) deleteApp(c *gin.Context) {
	err := s.catalog.DeleteApp(c.Request.Context(), idParam(c))
	if err != nil {
		resp.Fail(c, err)
		return
	}
	s.audit(c, "app.delete", "app", idParam(c))
	resp.OK(c, nil)
}
func (s *Server) publishApp(publish bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := s.catalog.PublishApp(c.Request.Context(), idParam(c), publish); err != nil {
			resp.Fail(c, err)
			return
		}
		s.audit(c, "app.publish", "app", idParam(c))
		resp.OK(c, nil)
	}
}
func (s *Server) adminCategories(c *gin.Context) {
	rows, err := s.catalog.Categories(c.Request.Context())
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, gin.H{"list": rows, "page": 1, "page_size": len(rows), "total": len(rows)})
}
func (s *Server) saveCategory(c *gin.Context) {
	var row model.Category
	if c.ShouldBindJSON(&row) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if c.Param("id") != "" {
		row.ID = idParam(c)
	}
	created, err := s.catalog.SaveCategory(c.Request.Context(), &row)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	if created {
		resp.Created(c, row)
	} else {
		resp.OK(c, row)
	}
}
func (s *Server) deleteCategory(c *gin.Context) {
	if err := s.catalog.DeleteCategory(c.Request.Context(), idParam(c)); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}
func (s *Server) adminReleases(c *gin.Context) {
	rows, err := s.catalog.Releases(c.Request.Context(), idParam(c), true)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, gin.H{"list": rows, "page": 1, "page_size": len(rows), "total": len(rows)})
}
func (s *Server) saveRelease(c *gin.Context) {
	var row model.Release
	if c.ShouldBindJSON(&row) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if strings.Contains(c.FullPath(), "/apps/:id/releases") {
		row.AppID = idParam(c)
	} else {
		row.ID = idParam(c)
		var old model.Release
		if s.store.DB.First(&old, row.ID).Error == nil {
			row.AppID = old.AppID
		}
	}
	created, err := s.catalog.SaveRelease(c.Request.Context(), &row)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	if created {
		resp.Created(c, row)
	} else {
		resp.OK(c, row)
	}
}
func (s *Server) deleteRelease(c *gin.Context) {
	if err := s.catalog.DeleteRelease(c.Request.Context(), idParam(c)); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}
func (s *Server) publishRelease(c *gin.Context) {
	if err := s.catalog.PublishRelease(c.Request.Context(), idParam(c)); err != nil {
		resp.Fail(c, err)
		return
	}
	s.audit(c, "release.publish", "release", idParam(c))
	resp.OK(c, nil)
}
func (s *Server) saveAsset(c *gin.Context) {
	var row model.Asset
	if c.ShouldBindJSON(&row) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if strings.Contains(c.FullPath(), "releases/:id") {
		row.ReleaseID = idParam(c)
	} else {
		row.ID = idParam(c)
		var old model.Asset
		if s.store.DB.First(&old, row.ID).Error == nil {
			row.ReleaseID = old.ReleaseID
		}
	}
	created, err := s.catalog.SaveAsset(c.Request.Context(), &row)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	if created {
		resp.Created(c, row)
	} else {
		resp.OK(c, row)
	}
}
func (s *Server) deleteAsset(c *gin.Context) {
	if err := s.catalog.DeleteAsset(c.Request.Context(), idParam(c)); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}
func (s *Server) saveSource(c *gin.Context) {
	var row model.DownloadSource
	if c.ShouldBindJSON(&row) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if strings.Contains(c.FullPath(), "assets/:id") {
		row.AssetID = idParam(c)
	} else {
		row.ID = idParam(c)
		var old model.DownloadSource
		if s.store.DB.First(&old, row.ID).Error == nil {
			row.AssetID = old.AssetID
		}
	}
	created, err := s.catalog.SaveSource(c.Request.Context(), &row)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	if created {
		resp.Created(c, row)
	} else {
		resp.OK(c, row)
	}
}
func (s *Server) deleteSource(c *gin.Context) {
	if err := s.catalog.DeleteSource(c.Request.Context(), idParam(c)); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}

func (s *Server) uploadInit(c *gin.Context) {
	var in struct {
		FileName  string `json:"file_name"`
		Size      int64  `json:"size"`
		SHA256    string `json:"sha256"`
		ChunkSize int64  `json:"chunk_size"`
	}
	if c.ShouldBindJSON(&in) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	out, err := s.uploads.Init(c.Request.Context(), in.FileName, in.Size, in.SHA256, in.ChunkSize)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, out)
}
func (s *Server) uploadChunk(c *gin.Context) {
	idx, err := strconv.Atoi(c.Param("index"))
	if err != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if err := s.uploads.Chunk(c.Param("id"), idx, c.Request.Body, c.GetHeader("X-Chunk-Sha256")); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}
func (s *Server) uploadComplete(c *gin.Context) {
	var in struct {
		StorageID *uint64 `json:"storage_id"`
		KeyHint   string  `json:"key_hint"`
	}
	if c.ShouldBindJSON(&in) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	out, err := s.uploads.Complete(c.Request.Context(), c.Param("id"), in.StorageID, in.KeyHint)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, out)
}
func (s *Server) uploadAbort(c *gin.Context) {
	if err := s.uploads.Abort(c.Param("id")); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}
func (s *Server) uploadFile(c *gin.Context) {
	f, h, err := c.Request.FormFile("file")
	if err != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	defer f.Close()
	var storageID *uint64
	if raw := c.PostForm("storage_id"); raw != "" {
		v, _ := strconv.ParseUint(raw, 10, 64)
		storageID = &v
	}
	out, err := s.uploads.PutFile(c.Request.Context(), h.Filename, f, h.Size, storageID, c.PostForm("key_hint"))
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, out)
}
func (s *Server) uploadImage(c *gin.Context) {
	f, h, err := c.Request.FormFile("file")
	if err != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	defer f.Close()
	url, err := s.saveImage(f, h)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, gin.H{"url": url})
}
func (s *Server) saveImage(f multipart.File, h *multipart.FileHeader) (string, error) {
	if h.Size > s.cfg.Upload.ImageMaxSizeMB*1024*1024 {
		return "", apperr.Wrap(10007, 413, "图片超出大小限制", nil)
	}
	head := make([]byte, 512)
	n, _ := io.ReadFull(f, head)
	kind := http.DetectContentType(head[:n])
	exts := map[string]string{"image/png": "png", "image/jpeg": "jpg", "image/webp": "webp", "image/gif": "gif", "image/x-icon": "ico"}
	ext, ok := exts[kind]
	if !ok {
		return "", apperr.Unprocessable
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	now := time.Now().UTC()
	rel := filepath.Join(strconv.Itoa(now.Year()), fmt.Sprintf("%02d", now.Month()), xid.New().String()+"."+ext)
	path := filepath.Join(s.cfg.DataDir, "uploads", rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return "", err
	}
	out, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, f); err != nil {
		return "", err
	}
	return "/uploads/" + strings.ReplaceAll(rel, "\\", "/"), nil
}

func (s *Server) listStorages(c *gin.Context) {
	rows, err := s.storages.List(c.Request.Context())
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, rows)
}
func (s *Server) saveStorage(c *gin.Context) {
	var row model.Storage
	if c.ShouldBindJSON(&row) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if c.Param("id") != "" {
		row.ID = idParam(c)
	}
	created, err := s.storages.Save(c.Request.Context(), &row)
	if err != nil {
		resp.Fail(c, err)
		return
	}
	if created {
		resp.Created(c, row)
	} else {
		resp.OK(c, row)
	}
}
func (s *Server) deleteStorage(c *gin.Context) {
	if err := s.storages.Delete(c.Request.Context(), idParam(c)); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}
func (s *Server) testStorage(c *gin.Context) {
	d, err := s.storages.Test(c.Request.Context(), idParam(c))
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, gin.H{"ok": true, "latency_ms": d.Milliseconds()})
}
func (s *Server) getSettings(c *gin.Context) {
	v, err := s.settings.All(c.Request.Context())
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, v)
}
func (s *Server) updateSettings(c *gin.Context) {
	var v map[string]string
	if c.ShouldBindJSON(&v) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if err := s.settings.Update(c.Request.Context(), v); err != nil {
		resp.Fail(c, err)
		return
	}
	s.audit(c, "setting.update", "setting", 0)
	resp.OK(c, v)
}

func (s *Server) listThemes(c *gin.Context) {
	rows, err := s.theme.List()
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, rows)
}
func (s *Server) uploadTheme(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	defer file.Close()
	if header.Size > 20*1024*1024 {
		resp.Fail(c, apperr.Wrap(10007, 413, "主题包超出 20MB", nil))
		return
	}
	tmpDir := filepath.Join(s.cfg.DataDir, "tmp")
	if err := os.MkdirAll(tmpDir, 0o700); err != nil {
		resp.Fail(c, err)
		return
	}
	tmp, err := os.CreateTemp(tmpDir, "theme-*.zip")
	if err != nil {
		resp.Fail(c, err)
		return
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err = io.Copy(tmp, io.LimitReader(file, 20*1024*1024+1)); err == nil {
		err = tmp.Close()
	} else {
		_ = tmp.Close()
	}
	if err != nil {
		resp.Fail(c, err)
		return
	}
	meta, err := s.theme.Install(tmpPath)
	if err != nil {
		resp.Fail(c, apperr.Wrap(10008, 422, "主题包校验失败", err))
		return
	}
	s.audit(c, "theme.install", "theme", 0)
	resp.OK(c, meta)
}
func (s *Server) activateTheme(c *gin.Context) {
	id := c.Param("id")
	if err := s.theme.Activate(id); err != nil {
		resp.Fail(c, apperr.Wrap(10008, 422, "主题编译失败", err))
		return
	}
	if err := s.settings.Update(c.Request.Context(), map[string]string{"theme.active": id}); err != nil {
		_ = s.theme.Activate("aurora")
		resp.Fail(c, err)
		return
	}
	s.audit(c, "theme.activate", "theme", 0)
	resp.OK(c, gin.H{"active": id})
}
func (s *Server) configureTheme(c *gin.Context) {
	var values map[string]any
	if c.ShouldBindJSON(&values) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	raw, err := json.Marshal(values)
	if err != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if err := s.settings.Update(c.Request.Context(), map[string]string{"theme.cfg." + c.Param("id"): string(raw)}); err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, values)
}
func (s *Server) deleteTheme(c *gin.Context) {
	if err := s.theme.Delete(c.Param("id")); err != nil {
		resp.Fail(c, apperr.Wrap(10008, 422, err.Error(), err))
		return
	}
	s.audit(c, "theme.delete", "theme", 0)
	resp.OK(c, nil)
}
func (s *Server) statsOverview(c *gin.Context) {
	ctx := c.Request.Context()
	var appCount, releaseCount, total int64
	_ = s.store.DB.WithContext(ctx).Model(&model.App{}).Count(&appCount).Error
	_ = s.store.DB.WithContext(ctx).Model(&model.Release{}).Count(&releaseCount).Error
	_ = s.store.DB.WithContext(ctx).Model(&model.App{}).Select("COALESCE(SUM(download_count),0)").Scan(&total).Error
	resp.OK(c, gin.H{"app_count": appCount, "release_count": releaseCount, "total_downloads": total, "today_downloads": 0, "today_views": 0, "storage_used_bytes": 0})
}
func (s *Server) adminPages(c *gin.Context) {
	var rows []model.Page
	_ = s.store.DB.WithContext(c.Request.Context()).Order("sort,id").Find(&rows).Error
	resp.OK(c, gin.H{"list": rows, "page": 1, "page_size": len(rows), "total": len(rows)})
}
func (s *Server) savePage(c *gin.Context) {
	var row model.Page
	if c.ShouldBindJSON(&row) != nil {
		resp.Fail(c, apperr.BadRequest)
		return
	}
	if c.Param("id") != "" {
		row.ID = idParam(c)
	}
	var err error
	if row.ID == 0 {
		err = s.store.DB.WithContext(c.Request.Context()).Create(&row).Error
	} else {
		err = s.store.DB.WithContext(c.Request.Context()).Model(&row).Updates(row).Error
	}
	if err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, row)
}
func (s *Server) deletePage(c *gin.Context) {
	if err := s.store.DB.WithContext(c.Request.Context()).Delete(&model.Page{}, idParam(c)).Error; err != nil {
		resp.Fail(c, err)
		return
	}
	resp.OK(c, nil)
}
func (s *Server) operationLogs(c *gin.Context) {
	var rows []model.OperationLog
	s.store.DB.WithContext(c.Request.Context()).Order("created_at DESC").Limit(100).Find(&rows)
	resp.OK(c, gin.H{"list": rows, "page": 1, "page_size": 100, "total": len(rows)})
}
func (s *Server) downloadLogs(c *gin.Context) {
	var rows []model.DownloadLog
	s.store.DB.WithContext(c.Request.Context()).Order("created_at DESC").Limit(100).Find(&rows)
	resp.OK(c, gin.H{"list": rows, "page": 1, "page_size": 100, "total": len(rows)})
}
func (s *Server) audit(c *gin.Context, action, target string, id uint64) {
	uid, _ := c.Get("uid")
	detail, _ := json.Marshal(gin.H{"path": c.FullPath()})
	_ = s.store.DB.WithContext(c.Request.Context()).Create(&model.OperationLog{UserID: uid.(uint64), Action: action, TargetType: target, TargetID: id, Detail: string(detail), IP: s.clientIP(c)}).Error
}

func (s *Server) feed(c *gin.Context) {
	var releases []model.Release
	s.store.DB.WithContext(c.Request.Context()).Joins("JOIN apps ON apps.id = releases.app_id").Where("releases.status = ? AND apps.status = ? AND apps.type = ?", model.StatusPublished, model.StatusPublished, model.AppTypeSelf).Order("releases.published_at DESC").Limit(20).Find(&releases)
	c.Header("Content-Type", "application/rss+xml; charset=utf-8")
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?><rss version="2.0"><channel><title>造物工坊</title><link>` + s.cfg.Server.BaseURL + `</link>`)
	for _, r := range releases {
		fmt.Fprintf(&b, "<item><title>版本 %s 发布</title><guid>%s/releases/%d</guid></item>", r.Version, s.cfg.Server.BaseURL, r.ID)
	}
	b.WriteString("</channel></rss>")
	c.String(200, b.String())
}
func (s *Server) sitemap(c *gin.Context) {
	var apps []model.App
	s.store.DB.WithContext(c.Request.Context()).Where("status = ? AND type = ?", model.StatusPublished, model.AppTypeSelf).Find(&apps)
	c.Header("Content-Type", "application/xml; charset=utf-8")
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>` + s.cfg.Server.BaseURL + `/</loc></url>`)
	for _, a := range apps {
		fmt.Fprintf(&b, "<url><loc>%s/apps/%s</loc></url>", s.cfg.Server.BaseURL, a.Slug)
	}
	b.WriteString("</urlset>")
	c.String(200, b.String())
}
func (s *Server) robots(c *gin.Context) {
	c.String(200, "User-agent: *\nDisallow: /admin\nDisallow: /api\nDisallow: /d/\nSitemap: %s/sitemap.xml\n", s.cfg.Server.BaseURL)
}

func (s *Server) adminStatic(r *gin.Engine) {
	adminFS, _ := fs.Sub(assets.Embedded, "admin")
	handler := http.FileServer(http.FS(adminFS))
	r.GET("/admin", func(c *gin.Context) { c.Redirect(302, "/admin/") })
	r.GET("/admin/*path", func(c *gin.Context) {
		path := strings.TrimPrefix(c.Param("path"), "/")
		if path == "" || path == "index.html" {
			index, err := fs.ReadFile(adminFS, "index.html")
			if err != nil {
				c.String(500, "admin assets unavailable")
				return
			}
			c.Header("Cache-Control", "no-cache")
			c.Data(200, "text/html; charset=utf-8", index)
			return
		}
		if f, err := adminFS.Open(path); err == nil {
			_ = f.Close()
			if path == "index.html" {
				c.Header("Cache-Control", "no-cache")
			}
			c.Request.URL.Path = "/" + path
			handler.ServeHTTP(c.Writer, c.Request)
			return
		}
		index, err := fs.ReadFile(adminFS, "index.html")
		if err != nil {
			c.String(500, "admin assets unavailable")
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(200, "text/html; charset=utf-8", index)
	})
}

func RunHTTP(ctx context.Context, addr string, h http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 10 * time.Second}
	errCh := make(chan error, 1)
	go func() { slog.Info("server listening", "addr", addr); errCh <- srv.ListenAndServe() }()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}

var _ = gorm.ErrRecordNotFound
