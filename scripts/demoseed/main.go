// demoseed 向本地数据库写入一组演示应用，便于主题开发与前台预览。
// 用法：go run ./scripts/demoseed [-dsn data/netupdown.db] [-datadir data]
// 已存在同 slug 的应用会被跳过，可重复执行。
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/moli-xia/netupdown/internal/model"
	"gorm.io/gorm"
)

func fakeSHA(name string) string { h := sha256.Sum256([]byte(name)); return hex.EncodeToString(h[:]) }

func daysAgo(n int) *time.Time { t := time.Now().UTC().AddDate(0, 0, -n); return &t }

type demoAsset struct {
	Name, FileName, OS, Arch string
	Kind                     int8
	Size                     int64
	Downloads                int64
	Managed                  bool   // true 时挂到本地存储的演示文件
	ExtURL, ExtName, Code    string // 外链源
}

type demoApp struct {
	model.App
	Version, Changelog string
	ReleasedDaysAgo    int
	Assets             []demoAsset
	TagNames           []string
}

func main() {
	dsn := flag.String("dsn", "data/netupdown.db", "sqlite database path")
	dataDir := flag.String("datadir", "data", "data directory")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open(*dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	for _, stmt := range []string{"PRAGMA journal_mode=WAL", "PRAGMA busy_timeout=5000", "PRAGMA foreign_keys=ON"} {
		if err := db.Exec(stmt).Error; err != nil {
			log.Fatal(err)
		}
	}

	var existing []model.App
	db.Unscoped().Find(&existing)
	fmt.Printf("数据库现有应用 %d 个：", len(existing))
	for _, a := range existing {
		fmt.Printf(" %s(%d)", a.Slug, a.Status)
	}
	fmt.Println()

	cats := map[string]uint64{}
	var catRows []model.Category
	db.Find(&catRows)
	catDesc := map[string][2]string{
		"productivity": {"桌面效率", "专注、记录与日常工作流"},
		"development":  {"开发者工具", "为开发、调试与发布提速"},
		"media":        {"创作工具", "截图、图像与内容创作"},
		"system":       {"系统工具", "让设备保持轻快与可靠"},
		"network":      {"网络与服务", "连接、诊断与自托管服务"},
		"works":        {"移动应用", "随身可用的轻量应用"},
	}
	for _, c := range catRows {
		cats[c.Slug] = c.ID
		if d, ok := catDesc[c.Slug]; ok {
			db.Model(&model.Category{}).Where("id = ?", c.ID).Updates(map[string]any{"name": d[0], "description": d[1]})
		}
	}

	branding := map[string]string{
		"site.title":       "造物工坊",
		"site.subtitle":    "独立开发者的软件发布与更新中心",
		"site.description": "探索 造物工坊自研产品，获取可靠的多平台安装包、版本记录与更新说明。",
		"site.keywords":    "自研软件,独立开发,软件下载,应用更新",
	}
	for key, value := range branding {
		if err := db.Save(&model.Setting{Key: key, Value: value}).Error; err != nil {
			log.Fatal(err)
		}
	}
	result := db.Model(&model.App{}).Where("type = ? AND status = ?", model.AppTypeThird, model.StatusPublished).Update("status", model.StatusOffline)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	fmt.Printf("已下线第三方演示条目 %d 个（记录仍保留在后台）\n", result.RowsAffected)
	if err := writeSelfDevIcons(*dataDir); err != nil {
		log.Fatal(err)
	}

	// 本地存储演示文件（供托管源真实可下载）
	demoObject := "demo/netupdown-demo-payload.bin"
	payloadPath := filepath.Join(*dataDir, "files", filepath.FromSlash(demoObject))
	if err := os.MkdirAll(filepath.Dir(payloadPath), 0o700); err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(payloadPath); err != nil {
		if err := os.WriteFile(payloadPath, []byte("NetUpDown demo payload — 演示文件，仅用于本地预览下载链路。\n"), 0o600); err != nil {
			log.Fatal(err)
		}
	}
	var localStorage model.Storage
	if err := db.Where("driver = ?", "local").First(&localStorage).Error; err != nil {
		log.Fatal("未找到本地存储实例，请先启动一次服务完成初始化")
	}

	apps := make([]demoApp, 0, 10)
	for _, item := range demoApps(cats) {
		if item.Type == model.AppTypeSelf {
			apps = append(apps, item)
		}
	}
	apps = append(apps, selfDevApps(cats)...)
	for _, d := range apps {
		var n int64
		db.Unscoped().Model(&model.App{}).Where("slug = ?", d.Slug).Count(&n)
		if n > 0 {
			updated := d.App
			updated.Type = model.AppTypeSelf
			updated.Status = model.StatusPublished
			updated.PublishedAt = daysAgo(d.ReleasedDaysAgo)
			if err := db.Unscoped().Model(&model.App{}).Where("slug = ?", d.Slug).Select("name", "type", "category_id", "tagline", "description", "icon", "screenshots", "official_url", "repo_url", "developer", "license", "platforms", "status", "is_featured", "sort_weight", "published_at").Updates(&updated).Error; err != nil {
				log.Fatal(err)
			}
			fmt.Println("已更新:", d.Slug)
			continue
		}
		app := d.App
		app.Status = model.StatusPublished
		app.PublishedAt = daysAgo(d.ReleasedDaysAgo)
		for _, tn := range d.TagNames {
			var tag model.Tag
			db.Where("name = ?", tn).Attrs(model.Tag{Name: tn, Slug: tn}).FirstOrCreate(&tag)
			app.Tags = append(app.Tags, tag)
		}
		if err := db.Create(&app).Error; err != nil {
			log.Fatal(err)
		}
		rel := model.Release{AppID: app.ID, Version: d.Version, Channel: model.ChannelStable, Changelog: d.Changelog, Status: model.StatusPublished, PublishedAt: daysAgo(d.ReleasedDaysAgo), DownloadCount: app.DownloadCount / 2}
		if err := db.Create(&rel).Error; err != nil {
			log.Fatal(err)
		}
		db.Model(&model.App{}).Where("id = ?", app.ID).Update("latest_release_id", rel.ID)
		for i, da := range d.Assets {
			asset := model.Asset{ReleaseID: rel.ID, Name: da.Name, FileName: da.FileName, OS: da.OS, Arch: da.Arch, Kind: da.Kind, Size: da.Size, SHA256: fakeSHA(da.FileName), DownloadCount: da.Downloads, Sort: i}
			if err := db.Create(&asset).Error; err != nil {
				log.Fatal(err)
			}
			pri := 0
			if da.Managed {
				db.Create(&model.DownloadSource{AssetID: asset.ID, Name: "直链下载", SourceType: model.SourceManaged, StorageID: &localStorage.ID, ObjectKey: demoObject, Priority: pri, IsEnabled: true, DownloadCount: da.Downloads})
				pri++
			}
			if da.ExtURL != "" {
				db.Create(&model.DownloadSource{AssetID: asset.ID, Name: da.ExtName, SourceType: model.SourceExternal, ExternalURL: da.ExtURL, ExtractCode: da.Code, Priority: pri, IsEnabled: true})
			}
		}
		fmt.Println("已创建:", d.Slug)
	}
	fmt.Println("完成。")
}

func demoApps(cats map[string]uint64) []demoApp {
	cid := func(slug string) *uint64 {
		if id, ok := cats[slug]; ok {
			return &id
		}
		return nil
	}
	return []demoApp{
		{
			App: model.App{Name: "NetPing 网络体检", Slug: "netping", Type: model.AppTypeSelf, CategoryID: cid("network"),
				Tagline:     "一键诊断网络质量：延迟、丢包、DNS、路由追踪，全都讲人话",
				Description: "## 为什么做 NetPing\n\n每次网络卡顿，家人只会说「网不好」。NetPing 把 ping、traceroute、DNS 解析、带宽测速打包成一次点击，用普通人能看懂的语言给出结论和建议。\n\n## 主要功能\n\n- **一键体检**：30 秒生成网络健康报告，红黄绿三色结论\n- **持续监控**：后台记录延迟与丢包曲线，波动一目了然\n- **路由可视化**：逐跳定位是谁在拖慢你的网络\n- **DNS 竞速**：自动找出当前环境下最快的 DNS\n\n## 隐私承诺\n\n所有诊断在本地完成，不上传任何数据。开源可审计。\n",
				Icon:        "/uploads/demo/netping.png", Screenshots: []string{"/uploads/demo/shot-netping-1.png", "/uploads/demo/shot-netping-2.png", "/uploads/demo/shot-netping-3.png"},
				OfficialURL: "https://netupdown.com", Developer: "NetUpDown Studio", License: "GPL-3.0",
				Platforms: []string{"windows", "macos", "linux"}, IsFeatured: true, DownloadCount: 128934, ViewCount: 402331},
			Version: "2.3.1", ReleasedDaysAgo: 3, TagNames: []string{"网络", "诊断", "开源"},
			Changelog: "## 新增\n\n- 路由追踪支持 IPv6\n- 新增「游戏模式」：针对常见游戏服务器的专项测试\n\n## 修复\n\n- 修复部分 Wi-Fi 6 网卡下测速偏低的问题\n- 修复深色模式下图表文字对比度不足\n\n## 优化\n\n- 体检报告生成速度提升约 40%\n",
			Assets: []demoAsset{
				{Name: "Windows 64 位 安装版", FileName: "NetPing-2.3.1-win-x64-setup.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 48234496, Downloads: 80123, Managed: true, ExtURL: "https://wwi.lanzoup.com/demo123", ExtName: "蓝奏云", Code: "np23"},
				{Name: "macOS 通用版", FileName: "NetPing-2.3.1-macos-universal.dmg", OS: "macos", Arch: "universal", Kind: 1, Size: 52690944, Downloads: 30231, Managed: true},
				{Name: "Linux 便携版", FileName: "NetPing-2.3.1-linux-amd64.tar.gz", OS: "linux", Arch: "amd64", Kind: 3, Size: 41943040, Downloads: 18580, Managed: true},
			},
		},
		{
			App: model.App{Name: "MarkFlow 笔记", Slug: "markflow", Type: model.AppTypeSelf, CategoryID: cid("productivity"),
				Tagline:     "本地优先的 Markdown 笔记，块编辑 + 双链 + 全文检索",
				Description: "## 简介\n\nMarkFlow 是一款**本地优先**的 Markdown 笔记应用：数据以纯文本存放在你自己的文件夹里，永远不会被锁进私有格式。\n\n## 特性\n\n| 特性 | 说明 |\n| --- | --- |\n| 块编辑 | 拖拽重排、折叠、多级引用 |\n| 双向链接 | `[[wiki]]` 语法与关系图谱 |\n| 全文检索 | 毫秒级中文分词检索 |\n| 主题 | 内置 12 套主题，支持 CSS 自定义 |\n\n## 同步方案\n\n支持 iCloud / OneDrive / 坚果云等任意网盘同步，也可搭配 Syncthing 实现全平台点对点同步。\n",
				Icon:        "/uploads/demo/markflow.png", Screenshots: []string{"/uploads/demo/shot-markflow-1.png", "/uploads/demo/shot-markflow-2.png"},
				OfficialURL: "https://netupdown.com", Developer: "NetUpDown Studio", License: "免费",
				Platforms: []string{"windows", "macos"}, IsFeatured: true, DownloadCount: 76112, ViewCount: 210998},
			Version: "1.8.0", ReleasedDaysAgo: 11, TagNames: []string{"笔记", "Markdown", "效率"},
			Changelog: "## 1.8.0\n\n- 新增关系图谱视图（Beta）\n- 表格编辑器支持列拖拽\n- 修复若干中文输入法候选框遮挡问题\n",
			Assets: []demoAsset{
				{Name: "Windows 安装版", FileName: "MarkFlow-1.8.0-win-x64.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 68157440, Downloads: 50234, Managed: true, ExtURL: "https://pan.quark.cn/s/demo456", ExtName: "夸克网盘", Code: "mkfl"},
				{Name: "macOS 版", FileName: "MarkFlow-1.8.0-macos-arm64.dmg", OS: "macos", Arch: "arm64", Kind: 1, Size: 71303168, Downloads: 25878, Managed: true},
			},
		},
		{
			App: model.App{Name: "7-Zip", Slug: "7-zip", Type: model.AppTypeThird, CategoryID: cid("system"),
				Tagline:     "老牌开源压缩软件，7z 格式的发明者，免费无广告",
				Description: "## 简介\n\n7-Zip 是历史悠久的开源压缩工具，凭借出色的 7z 压缩比和完全免费的授权，成为装机必备。\n\n## 亮点\n\n- 7z 格式压缩比业界领先\n- 支持 zip、rar、tar、gz 等几十种格式解压\n- 干净纯粹：无广告、无捆绑、无后台\n- 支持 AES-256 加密\n\n> 收录说明：本页提供官方渠道与本站网盘转存，转存文件校验值与官方一致。\n",
				Icon:        "/uploads/demo/7zip.png", OfficialURL: "https://www.7-zip.org", Developer: "Igor Pavlov", License: "LGPL",
				Platforms: []string{"windows", "linux"}, DownloadCount: 45211, ViewCount: 98122},
			Version: "25.01", ReleasedDaysAgo: 25, TagNames: []string{"压缩", "开源", "装机必备"},
			Changelog: "- 官方 25.01 版本\n- 改进 zstd 解压性能\n- 修复若干安全问题（建议升级）\n",
			Assets: []demoAsset{
				{Name: "Windows 64 位 安装版", FileName: "7z2501-x64.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 1572864, Downloads: 40211, ExtURL: "https://www.7-zip.org/a/7z2501-x64.exe", ExtName: "官方直链"},
				{Name: "Windows ARM64 安装版", FileName: "7z2501-arm64.exe", OS: "windows", Arch: "arm64", Kind: 1, Size: 1468006, Downloads: 5000, ExtURL: "https://www.7-zip.org/a/7z2501-arm64.exe", ExtName: "官方直链"},
			},
		},
		{
			App: model.App{Name: "Everything", Slug: "everything", Type: model.AppTypeThird, CategoryID: cid("system"),
				Tagline:     "秒搜全盘文件，Windows 上最快的文件名搜索工具",
				Description: "## 简介\n\nEverything 通过读取 NTFS 索引实现毫秒级文件名搜索，索引百万文件也只需几秒钟，内存占用极低。\n\n## 使用建议\n\n- 搭配快捷键呼出，替代系统自带搜索\n- 支持正则、路径过滤等高级语法\n- 可开启 HTTP/ETP 服务远程搜索\n",
				Icon:        "/uploads/demo/everything.png", OfficialURL: "https://www.voidtools.com", Developer: "voidtools", License: "免费",
				Platforms: []string{"windows"}, DownloadCount: 38455, ViewCount: 76001},
			Version: "1.4.1.1028", ReleasedDaysAgo: 60, TagNames: []string{"搜索", "装机必备"},
			Changelog: "- 官方稳定版\n- 改进对长路径的支持\n",
			Assets: []demoAsset{
				{Name: "Windows 64 位 安装版", FileName: "Everything-1.4.1.1028.x64-Setup.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 1782579, Downloads: 30455, ExtURL: "https://www.voidtools.com/Everything-1.4.1.1028.x64-Setup.exe", ExtName: "官方直链"},
				{Name: "Windows 64 位 便携版", FileName: "Everything-1.4.1.1028.x64.zip", OS: "windows", Arch: "amd64", Kind: 2, Size: 1887436, Downloads: 8000, ExtURL: "https://www.voidtools.com/Everything-1.4.1.1028.x64.zip", ExtName: "官方直链"},
			},
		},
		{
			App: model.App{Name: "Visual Studio Code", Slug: "vscode", Type: model.AppTypeThird, CategoryID: cid("development"),
				Tagline:     "微软出品的免费代码编辑器，插件生态无出其右",
				Description: "## 简介\n\nVS Code 是目前最流行的代码编辑器：轻量启动、海量插件、内置终端与调试器，从前端到数据科学通吃。\n\n## 推荐插件\n\n- Chinese Language Pack — 中文界面\n- GitLens — Git 增强\n- Remote - SSH — 远程开发\n",
				Icon:        "/uploads/demo/vscode.png", OfficialURL: "https://code.visualstudio.com", RepoURL: "https://github.com/microsoft/vscode", Developer: "Microsoft", License: "MIT",
				Platforms: []string{"windows", "macos", "linux", "web"}, IsFeatured: true, DownloadCount: 52310, ViewCount: 120414},
			Version: "1.103.0", ReleasedDaysAgo: 8, TagNames: []string{"编辑器", "开发"},
			Changelog: "- 官方月度更新\n- 内置 AI 补全体验改进\n- 终端启动速度优化\n",
			Assets: []demoAsset{
				{Name: "Windows 64 位 安装版", FileName: "VSCodeUserSetup-x64-1.103.0.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 99614720, Downloads: 40310, ExtURL: "https://update.code.visualstudio.com/latest/win32-x64-user/stable", ExtName: "官方直链"},
				{Name: "macOS 通用版", FileName: "VSCode-darwin-universal.zip", OS: "macos", Arch: "universal", Kind: 3, Size: 157286400, Downloads: 12000, ExtURL: "https://update.code.visualstudio.com/latest/darwin-universal/stable", ExtName: "官方直链"},
			},
		},
		{
			App: model.App{Name: "PotPlayer", Slug: "potplayer", Type: model.AppTypeThird, CategoryID: cid("media"),
				Tagline:     "功能强大的本地视频播放器，解码能力与自定义程度俱佳",
				Description: "## 简介\n\nPotPlayer 以强大的内置解码器和几乎无限的自定义能力著称，是 Windows 平台最受欢迎的本地播放器之一。\n\n## 提示\n\n安装时请留意勾选项；建议在设置中关闭在线内容推荐。\n",
				Icon:        "/uploads/demo/potplayer.png", OfficialURL: "https://potplayer.daum.net", Developer: "Kakao", License: "免费",
				Platforms: []string{"windows"}, DownloadCount: 29877, ViewCount: 61022},
			Version: "250815", ReleasedDaysAgo: 15, TagNames: []string{"播放器", "影音"},
			Changelog: "- 更新内置解码器\n- 改进 HDR 到 SDR 的映射\n",
			Assets: []demoAsset{
				{Name: "Windows 64 位 安装版", FileName: "PotPlayerSetup64.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 47185920, Downloads: 29877, ExtURL: "https://t1.daumcdn.net/potplayer/PotPlayer/Version/Latest/PotPlayerSetup64.exe", ExtName: "官方直链", Code: ""},
			},
		},
		{
			App: model.App{Name: "uTools", Slug: "utools", Type: model.AppTypeThird, CategoryID: cid("productivity"),
				Tagline:     "新一代效率启动器：一个快捷键，唤起所有插件",
				Description: "## 简介\n\nuTools 用一个输入框整合了启动器、剪贴板历史、翻译、取色、OCR 等上百种插件能力，Alt+Space 即刻唤起。\n\n## 常用插件\n\n- 剪贴板历史\n- 屏幕取色 / 截图 OCR\n- 密码管理\n",
				Icon:        "/uploads/demo/utools.png", OfficialURL: "https://www.u-tools.cn", Developer: "uTools 团队", License: "免费增值",
				Platforms: []string{"windows", "macos", "linux"}, DownloadCount: 21003, ViewCount: 45023},
			Version: "7.0.0", ReleasedDaysAgo: 20, TagNames: []string{"效率", "启动器"},
			Changelog: "- 全新 7.0 界面\n- 插件市场改版\n",
			Assets: []demoAsset{
				{Name: "Windows 安装版", FileName: "uTools-7.0.0.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 89128960, Downloads: 15003, ExtURL: "https://www.u-tools.cn/download/", ExtName: "官方页面"},
				{Name: "macOS 版", FileName: "uTools-7.0.0.dmg", OS: "macos", Arch: "universal", Kind: 1, Size: 94371840, Downloads: 6000, ExtURL: "https://www.u-tools.cn/download/", ExtName: "官方页面"},
			},
		},
		{
			App: model.App{Name: "Snipaste", Slug: "snipaste", Type: model.AppTypeThird, CategoryID: cid("productivity"),
				Tagline:     "截图 + 贴图，重新定义你的截图工作流",
				Description: "## 简介\n\nSnipaste 的核心创意是「贴图」：把截图钉在屏幕上随时参考，配合像素级取色与标注，是设计师和开发者的心头好。\n\n- F1 截图，F3 贴图\n- 自动检测界面元素边界\n- 取色器支持多种颜色格式\n",
				Icon:        "/uploads/demo/snipaste.png", OfficialURL: "https://www.snipaste.com", Developer: "levie", License: "免费",
				Platforms: []string{"windows", "macos"}, DownloadCount: 18211, ViewCount: 39877},
			Version: "2.10.7", ReleasedDaysAgo: 40, TagNames: []string{"截图", "效率"},
			Changelog: "- 修复多显示器缩放下的取色偏移\n- 新增 WebP 导出\n",
			Assets: []demoAsset{
				{Name: "Windows 64 位 便携版", FileName: "Snipaste-2.10.7-x64.zip", OS: "windows", Arch: "amd64", Kind: 2, Size: 26214400, Downloads: 18211, ExtURL: "https://www.snipaste.com/download.html", ExtName: "官方页面", Code: ""},
			},
		},
	}
}

func writeSelfDevIcons(dataDir string) error {
	type iconSpec struct{ File, Mark, Color string }
	icons := []iconSpec{
		{"clipforge.svg", "CF", "#7c3aed"},
		{"devdock.svg", "DD", "#2563eb"},
		{"pixelmint.svg", "PM", "#db2777"},
		{"focusdeck.svg", "FD", "#ea580c"},
		{"cleansweep.svg", "CS", "#059669"},
		{"relaybox.svg", "RB", "#0891b2"},
		{"pocketledger.svg", "PL", "#16a34a"},
		{"launchkit.svg", "LK", "#4f46e5"},
	}
	root := filepath.Join(dataDir, "uploads", "demo")
	if err := os.MkdirAll(root, 0o700); err != nil {
		return err
	}
	for _, icon := range icons {
		svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 128 128"><rect width="128" height="128" rx="28" fill="%s"/><circle cx="96" cy="32" r="18" fill="#fff" opacity=".14"/><path d="M18 96 52 62l19 19 39-39" fill="none" stroke="#fff" stroke-width="8" stroke-linecap="round" stroke-linejoin="round" opacity=".2"/><text x="64" y="76" text-anchor="middle" fill="#fff" font-family="Segoe UI,Arial,sans-serif" font-size="38" font-weight="700" letter-spacing="-2">%s</text></svg>`, icon.Color, icon.Mark)
		if err := os.WriteFile(filepath.Join(root, icon.File), []byte(svg), 0o600); err != nil {
			return err
		}
	}
	return nil
}

func selfDevApps(cats map[string]uint64) []demoApp {
	cid := func(slug string) *uint64 {
		if id, ok := cats[slug]; ok {
			return &id
		}
		return nil
	}
	return []demoApp{
		{
			App: model.App{Name: "ClipForge 截图工坊", Slug: "clipforge", Type: model.AppTypeSelf, CategoryID: cid("media"),
				Tagline:     "从截图、标注到长图拼接，一次完成清晰表达",
				Description: "## 把想法直接画在屏幕上\n\nClipForge 是为产品、设计和技术支持场景打造的截图工具。它把区域截图、滚动长截图、贴图与标注整合在一个轻量工作流中。\n\n## 主要能力\n\n- 自动识别窗口与界面元素边界\n- 箭头、序号、聚光灯与敏感信息打码\n- 滚动截图自动拼接，支持导出 PNG、WebP 与 PDF\n- 截图历史保存在本地，可按应用和日期检索\n",
				Icon:        "/uploads/demo/clipforge.svg", Developer: "NetUpDown Studio", License: "免费",
				Platforms: []string{"windows", "macos"}, IsFeatured: true, SortWeight: 220, DownloadCount: 68320, ViewCount: 151204},
			Version: "3.2.0", ReleasedDaysAgo: 1, TagNames: []string{"截图", "标注", "创作"},
			Changelog: "## 3.2.0\n\n- 新增长截图断点续接\n- 标注工具支持统一调整颜色与线宽\n- 优化多显示器缩放下的边缘识别\n",
			Assets: []demoAsset{
				{Name: "Windows 64 位安装版", FileName: "ClipForge-3.2.0-win-x64.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 34603008, Downloads: 46210, Managed: true},
				{Name: "macOS 通用版", FileName: "ClipForge-3.2.0-macos.dmg", OS: "macos", Arch: "universal", Kind: 1, Size: 39845888, Downloads: 22110, Managed: true},
			},
		},
		{
			App: model.App{Name: "DevDock 开发环境管家", Slug: "devdock", Type: model.AppTypeSelf, CategoryID: cid("development"),
				Tagline:     "在项目之间快速切换 Node、Python、Go 与常用服务",
				Description: "## 一套项目，一套环境\n\nDevDock 读取项目配置并准备匹配的运行时、环境变量和本地服务，避免不同项目之间的版本互相干扰。\n\n- 自动识别常见语言与包管理器\n- 一键启动数据库、缓存与消息队列\n- 环境差异可视化，配置变更前提供预览\n- 所有凭据只保存在系统安全存储中\n",
				Icon:        "/uploads/demo/devdock.svg", Developer: "NetUpDown Studio", License: "MIT",
				Platforms: []string{"windows", "macos", "linux"}, IsFeatured: true, SortWeight: 210, DownloadCount: 54218, ViewCount: 128492},
			Version: "1.6.2", ReleasedDaysAgo: 2, TagNames: []string{"开发环境", "运行时", "效率"},
			Changelog: "- 新增 Go 1.25 工具链配置\n- 服务面板支持健康状态检测\n- 修复 Windows 长路径项目无法启动的问题\n",
			Assets: []demoAsset{
				{Name: "Windows 安装版", FileName: "DevDock-1.6.2-win-x64.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 57671680, Downloads: 32180, Managed: true},
				{Name: "macOS 通用版", FileName: "DevDock-1.6.2-macos.dmg", OS: "macos", Arch: "universal", Kind: 1, Size: 61865984, Downloads: 13721, Managed: true},
				{Name: "Linux AppImage", FileName: "DevDock-1.6.2-linux.AppImage", OS: "linux", Arch: "amd64", Kind: 2, Size: 55574528, Downloads: 8317, Managed: true},
			},
		},
		{
			App: model.App{Name: "PixelMint 图片轻压", Slug: "pixelmint", Type: model.AppTypeSelf, CategoryID: cid("media"),
				Tagline:     "批量压缩与转换图片，在清晰度和体积之间自动取平衡",
				Description: "## 面向真实工作流的图片优化\n\n拖入一个文件夹，PixelMint 会根据图片内容自动选择压缩参数，并保留目录结构和必要的元数据。\n\n- PNG、JPEG、WebP、AVIF 批量互转\n- 视觉对比模式与目标体积模式\n- 可监控文件夹并自动处理新增图片\n- 全程离线，不上传原始素材\n",
				Icon:        "/uploads/demo/pixelmint.svg", Developer: "NetUpDown Studio", License: "免费",
				Platforms: []string{"windows", "macos", "linux"}, IsFeatured: true, SortWeight: 200, DownloadCount: 48762, ViewCount: 102341},
			Version: "2.4.1", ReleasedDaysAgo: 4, TagNames: []string{"图片压缩", "格式转换", "离线"},
			Changelog: "- AVIF 编码速度提升约 35%\n- 新增按最长边批量缩放\n- 修复透明 PNG 边缘出现黑边的问题\n",
			Assets: []demoAsset{
				{Name: "Windows 64 位安装版", FileName: "PixelMint-2.4.1-win-x64.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 30408704, Downloads: 30551, Managed: true},
				{Name: "macOS Apple 芯片版", FileName: "PixelMint-2.4.1-macos-arm64.dmg", OS: "macos", Arch: "arm64", Kind: 1, Size: 32505856, Downloads: 12120, Managed: true},
				{Name: "Linux 便携版", FileName: "PixelMint-2.4.1-linux.tar.gz", OS: "linux", Arch: "amd64", Kind: 3, Size: 28311552, Downloads: 6091, Managed: true},
			},
		},
		{
			App: model.App{Name: "FocusDeck 专注面板", Slug: "focusdeck", Type: model.AppTypeSelf, CategoryID: cid("productivity"),
				Tagline:     "把今天最重要的三件事放到桌面，减少来回切换",
				Description: "## 少一点清单，多一点完成\n\nFocusDeck 用一个常驻桌面的小面板管理今日重点、专注计时和临时记录。没有复杂项目结构，打开即可开始。\n\n- 每日三项重点与番茄钟\n- 自动归档完成记录，生成周回顾\n- 支持全局快捷键和迷你悬浮模式\n- 本地数据库，可导出为 Markdown\n",
				Icon:        "/uploads/demo/focusdeck.svg", Developer: "NetUpDown Studio", License: "免费",
				Platforms: []string{"windows", "macos"}, IsFeatured: true, SortWeight: 190, DownloadCount: 41690, ViewCount: 93605},
			Version: "1.9.0", ReleasedDaysAgo: 5, TagNames: []string{"专注", "待办", "时间管理"},
			Changelog: "- 新增每周专注趋势\n- 迷你模式支持自动贴边\n- 提醒音量可独立于系统音量设置\n",
			Assets: []demoAsset{
				{Name: "Windows 安装版", FileName: "FocusDeck-1.9.0-win-x64.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 24117248, Downloads: 29430, Managed: true},
				{Name: "macOS 通用版", FileName: "FocusDeck-1.9.0-macos.dmg", OS: "macos", Arch: "universal", Kind: 1, Size: 27262976, Downloads: 12260, Managed: true},
			},
		},
		{
			App: model.App{Name: "CleanSweep 系统清理", Slug: "cleansweep", Type: model.AppTypeSelf, CategoryID: cid("system"),
				Tagline:     "先解释、再清理，安全找回被缓存与临时文件占用的空间",
				Description: "## 可解释的磁盘清理\n\nCleanSweep 会标明每一类文件来自哪里、删除后有什么影响，并在操作前生成可恢复清单。\n\n- 扫描系统缓存、开发工具缓存与大型临时文件\n- 文件类型和最近使用时间双重筛选\n- 默认进入回收站，支持创建恢复点\n- 不提供所谓的一键注册表优化\n",
				Icon:        "/uploads/demo/cleansweep.svg", Developer: "NetUpDown Studio", License: "免费",
				Platforms: []string{"windows"}, IsFeatured: true, SortWeight: 180, DownloadCount: 38914, ViewCount: 110482},
			Version: "2.1.3", ReleasedDaysAgo: 7, TagNames: []string{"磁盘清理", "存储", "安全"},
			Changelog: "- 新增 Docker 与 WSL 缓存分析\n- 大文件视图支持按最近访问时间排序\n- 修复部分目录权限不足时扫描中断的问题\n",
			Assets:    []demoAsset{{Name: "Windows 64 位安装版", FileName: "CleanSweep-2.1.3-win-x64.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 22020096, Downloads: 38914, Managed: true}},
		},
		{
			App: model.App{Name: "RelayBox 文件中转站", Slug: "relaybox", Type: model.AppTypeSelf, CategoryID: cid("network"),
				Tagline:     "把临时文件安全送到自己的设备，支持自托管中转节点",
				Description: "## 跨设备传文件，不经过陌生网盘\n\nRelayBox 优先在局域网内直连，外网环境可使用自己部署的中转节点。文件端到端加密，中转端无法读取内容。\n\n- Windows、macOS、Linux 与 Web 客户端\n- 链接、二维码和六位取件码\n- 传输过期后自动清理\n- 支持 Docker 一键部署中转节点\n",
				Icon:        "/uploads/demo/relaybox.svg", Developer: "NetUpDown Studio", License: "AGPL-3.0",
				Platforms: []string{"windows", "macos", "linux", "web"}, IsFeatured: true, SortWeight: 170, DownloadCount: 33402, ViewCount: 85620},
			Version: "1.3.0", ReleasedDaysAgo: 9, TagNames: []string{"文件传输", "自托管", "加密"},
			Changelog: "- 新增断点续传与失败分片重试\n- Web 客户端支持直接接收文件\n- 中转节点新增存储用量告警\n",
			Assets: []demoAsset{
				{Name: "桌面端 Windows", FileName: "RelayBox-1.3.0-win-x64.exe", OS: "windows", Arch: "amd64", Kind: 1, Size: 29360128, Downloads: 18040, Managed: true},
				{Name: "桌面端 macOS", FileName: "RelayBox-1.3.0-macos.dmg", OS: "macos", Arch: "universal", Kind: 1, Size: 31457280, Downloads: 9211, Managed: true},
				{Name: "Linux 服务端", FileName: "relaybox-server-1.3.0-linux.tar.gz", OS: "linux", Arch: "amd64", Kind: 3, Size: 15728640, Downloads: 6151, Managed: true},
			},
		},
		{
			App: model.App{Name: "PocketLedger 随手账", Slug: "pocketledger", Type: model.AppTypeSelf, CategoryID: cid("works"),
				Tagline:     "无需账户、没有广告的本地记账与预算应用",
				Description: "## 账本属于你自己\n\nPocketLedger 不要求注册账号，记录默认只保存在设备上。你可以选择通过自己的 WebDAV 服务同步。\n\n- 快速记账、周期账单与分类预算\n- 月度趋势和消费日历\n- CSV 与加密备份导入导出\n- 生物识别解锁，隐藏敏感金额\n",
				Icon:        "/uploads/demo/pocketledger.svg", Developer: "NetUpDown Studio", License: "免费",
				Platforms: []string{"android", "ios"}, SortWeight: 160, DownloadCount: 27185, ViewCount: 64210},
			Version: "1.5.1", ReleasedDaysAgo: 12, TagNames: []string{"记账", "预算", "隐私"},
			Changelog: "- 新增年度预算视图\n- 支持识别常见账单 CSV 格式\n- 优化大数据量账本的启动速度\n",
			Assets: []demoAsset{
				{Name: "Android 安装包", FileName: "PocketLedger-1.5.1-android.apk", OS: "android", Arch: "universal", Kind: 1, Size: 18874368, Downloads: 19420, Managed: true},
				{Name: "iOS 下载说明", FileName: "PocketLedger-iOS.txt", OS: "ios", Arch: "universal", Kind: 5, Size: 4096, Downloads: 7765, Managed: true},
			},
		},
		{
			App: model.App{Name: "LaunchKit 发布助手", Slug: "launchkit", Type: model.AppTypeSelf, CategoryID: cid("development"),
				Tagline:     "打包、签名、生成更新清单，把桌面应用可靠地交付给用户",
				Description: "## 让发布流程可以重复\n\nLaunchKit 把构建产物校验、代码签名、版本说明和多渠道上传组织成可审阅的发布流水线。\n\n- 支持 Windows、macOS 与 Linux 产物\n- 自动生成 SHA256 和更新清单\n- 发布前展示文件、版本与渠道差异\n- 可连接自托管对象存储和 WebDAV\n",
				Icon:        "/uploads/demo/launchkit.svg", Developer: "NetUpDown Studio", License: "MIT",
				Platforms: []string{"windows", "macos", "linux"}, SortWeight: 150, DownloadCount: 19684, ViewCount: 51782},
			Version: "0.9.4", ReleasedDaysAgo: 16, TagNames: []string{"发布", "自动化", "签名"},
			Changelog: "- 新增发布前差异摘要\n- 支持复用已有上传对象\n- Linux 包校验新增 AppImage 元数据检查\n",
			Assets: []demoAsset{
				{Name: "Windows 便携版", FileName: "LaunchKit-0.9.4-win-x64.zip", OS: "windows", Arch: "amd64", Kind: 2, Size: 26214400, Downloads: 10200, Managed: true},
				{Name: "macOS 通用版", FileName: "LaunchKit-0.9.4-macos.dmg", OS: "macos", Arch: "universal", Kind: 1, Size: 29360128, Downloads: 5800, Managed: true},
				{Name: "Linux 便携版", FileName: "LaunchKit-0.9.4-linux.tar.gz", OS: "linux", Arch: "amd64", Kind: 3, Size: 24117248, Downloads: 3684, Managed: true},
			},
		},
	}
}
