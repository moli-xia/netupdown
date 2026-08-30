// opensource_seed 写入一组经过官方来源核对的热门开源软件。
// 用法：go run ./scripts/opensource_seed [-dsn data/netupdown.db] [-datadir data]
// 已存在同 slug 的应用会被跳过，可重复执行。
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/moli-xia/netupdown/internal/model"
	"gorm.io/gorm"
)

type openSourceApp struct {
	Name, Slug, Category, Tagline, Description string
	OfficialURL, RepoURL, ReleaseURL           string
	Developer, License                         string
	Version, IconSource, IconFile              string
	Stars                                      int64
	Platforms                                  []string
	Tags                                       []string
	Featured                                   bool
}

func main() {
	dsn := "data/netupdown.db"
	dataDir := "data"
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-dsn":
			if i+1 < len(os.Args) {
				dsn = os.Args[i+1]
				i++
			}
		case "-datadir":
			if i+1 < len(os.Args) {
				dataDir = os.Args[i+1]
				i++
			}
		}
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		fatal(err)
	}
	for _, stmt := range []string{"PRAGMA journal_mode=WAL", "PRAGMA busy_timeout=5000", "PRAGMA foreign_keys=ON"} {
		if err := db.Exec(stmt).Error; err != nil {
			fatal(err)
		}
	}
	apps := catalog()
	for _, item := range apps {
		var count int64
		if err := db.Unscoped().Model(&model.App{}).Where("slug = ?", item.Slug).Count(&count).Error; err != nil {
			fatal(err)
		}
		if count > 0 {
			fmt.Println("跳过（已存在）:", item.Slug)
			continue
		}

		categoryID := category(db, item.Category)
		iconPath := filepath.Join(dataDir, "uploads", "opensource", item.IconFile)
		if err := downloadIcon(iconPath, item.IconSource); err != nil {
			fatal(fmt.Errorf("download %s icon: %w", item.Slug, err))
		}

		tags := make([]model.Tag, 0, len(item.Tags))
		for _, name := range item.Tags {
			var tag model.Tag
			if err := db.Where("name = ?", name).Attrs(model.Tag{Name: name, Slug: name}).FirstOrCreate(&tag).Error; err != nil {
				fatal(err)
			}
			tags = append(tags, tag)
		}

		now := time.Now().UTC()
		app := model.App{
			Name:        item.Name,
			Slug:        item.Slug,
			Type:        model.AppTypeThird,
			CategoryID:  &categoryID,
			Tagline:     item.Tagline,
			Description: item.Description,
			Icon:        "/uploads/opensource/" + item.IconFile,
			OfficialURL: item.OfficialURL,
			RepoURL:     item.RepoURL,
			Developer:   item.Developer,
			License:     item.License,
			Platforms:   item.Platforms,
			Status:      model.StatusPublished,
			IsFeatured:  item.Featured,
			SortWeight:  100,
			PublishedAt: &now,
			Tags:        tags,
		}
		if err := db.Create(&app).Error; err != nil {
			fatal(err)
		}

		release := model.Release{
			AppID:       app.ID,
			Version:     item.Version,
			Channel:     model.ChannelStable,
			Title:       "官方最新稳定发布",
			Changelog:   "- 收录官方发布页\n- 版本与下载内容以项目官方页面为准\n",
			Status:      model.StatusPublished,
			PublishedAt: &now,
		}
		if err := db.Create(&release).Error; err != nil {
			fatal(err)
		}
		if err := db.Model(&model.App{}).Where("id = ?", app.ID).Update("latest_release_id", release.ID).Error; err != nil {
			fatal(err)
		}

		asset := model.Asset{
			ReleaseID: release.ID,
			Name:      "官方多平台下载页",
			FileName:  item.Slug + "-official-download",
			OS:        "any",
			Arch:      "any",
			Kind:      9,
			Sort:      0,
		}
		if err := db.Create(&asset).Error; err != nil {
			fatal(err)
		}
		source := model.DownloadSource{
			AssetID:     asset.ID,
			Name:        "官方发布页",
			SourceType:  model.SourceExternal,
			ExternalURL: item.ReleaseURL,
			Priority:    0,
			IsEnabled:   true,
		}
		if err := db.Create(&source).Error; err != nil {
			fatal(err)
		}
		fmt.Printf("已创建: %s（GitHub ★ %s）\n", item.Slug, starLabel(item.Stars))
	}
	fmt.Println("完成。")
}

func category(db *gorm.DB, slug string) uint64 {
	var row model.Category
	if err := db.Where("slug = ?", slug).First(&row).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			fatal(err)
		}
		name := map[string]string{
			"productivity": "效率工具",
			"development":  "开发编程",
			"media":        "影音图像",
			"system":       "系统增强",
			"network":      "网络工具",
		}[slug]
		row = model.Category{Name: name, Slug: slug}
		if err := db.Create(&row).Error; err != nil {
			fatal(err)
		}
	}
	return row.ID
}

func downloadIcon(path, source string) error {
	if info, err := os.Stat(path); err == nil && info.Size() > 100 {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, source, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "netupdown-catalog")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %s", resp.Status)
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(f, resp.Body)
	closeErr := f.Close()
	if copyErr != nil {
		return copyErr
	}
	if closeErr != nil {
		return closeErr
	}
	if err := os.Rename(tmp, path); err != nil {
		return err
	}
	return nil
}

func starLabel(n int64) string {
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func catalog() []openSourceApp {
	return []openSourceApp{
		{
			Name: "RustDesk", Slug: "rustdesk", Category: "network",
			Tagline: "开源自托管远程桌面，跨平台连接自己的设备 · GitHub ★ 121.9k（2026-08-27）",
			Description: `## 简介

RustDesk 是一款开源远程桌面工具，适合远程协助、跨设备访问和团队运维。它支持使用官方服务，也支持部署自己的中继与会合服务，让连接和数据更可控。

## 主要特点

- 支持 Windows、macOS、Linux、Android 和 iOS
- 支持远程控制、文件传输、剪贴板等常用能力
- 可自托管服务端，适合重视隐私与可控性的用户

## 热度与来源

GitHub Star 快照：121,960（查询于 2026-08-27）。图标来自 RustDesk 官方 GitHub 组织头像；下载按钮跳转项目官方发布页，本站不重新打包安装文件。
`,
			OfficialURL: "https://rustdesk.com/download", RepoURL: "https://github.com/rustdesk/rustdesk", ReleaseURL: "https://github.com/rustdesk/rustdesk/releases/latest",
			Developer: "RustDesk", License: "AGPL-3.0", Version: "1.4.9", IconSource: "https://avatars.githubusercontent.com/u/71636191?v=4", IconFile: "rustdesk.png",
			Stars: 121960, Platforms: []string{"windows", "macos", "linux", "android", "ios"}, Tags: []string{"远程桌面", "自托管", "开源"}, Featured: true,
		},
		{
			Name: "OBS Studio", Slug: "obs-studio", Category: "media",
			Tagline: "专业直播与录屏工具，场景、滤镜、编码一站式完成 · GitHub ★ 75.5k（2026-08-27）",
			Description: `## 简介

OBS Studio 是面向直播和录屏的开源创作工具，可以把桌面、窗口、摄像头、麦克风等内容组合成场景，再进行录制或推流。

## 主要特点

- 场景与来源自由组合，支持转场和实时预览
- 内置音视频混音、滤镜和编码控制
- 适用于 Windows、macOS 与 Linux，并拥有丰富的插件生态

## 热度与来源

GitHub Star 快照：75,511（查询于 2026-08-27）。软件采用 GPL-2.0；图标来自 OBS Project 官方 GitHub 组织头像，下载入口为官方发布页。
`,
			OfficialURL: "https://obsproject.com/download", RepoURL: "https://github.com/obsproject/obs-studio", ReleaseURL: "https://github.com/obsproject/obs-studio/releases/latest",
			Developer: "OBS Project", License: "GPL-2.0", Version: "32.2.2", IconSource: "https://avatars.githubusercontent.com/u/7725691?v=4", IconFile: "obs-studio.png",
			Stars: 75511, Platforms: []string{"windows", "macos", "linux"}, Tags: []string{"录屏", "直播", "视频"}, Featured: true,
		},
		{
			Name: "VLC 媒体播放器", Slug: "vlc", Category: "media",
			Tagline: "几乎什么都能播放的开源媒体播放器，支持本地文件、光盘与网络流 · GitHub ★ 19.5k（2026-08-27）",
			Description: `## 简介

VLC 是 VideoLAN 社区维护的开源媒体播放器与多媒体引擎，覆盖常见音视频文件、光盘、网络流和采集设备，也可以进行转码与串流。

## 主要特点

- 支持 Windows、macOS、Linux、Android 和 iOS 等平台
- 兼容大量媒体格式、字幕、网络协议和输入源
- 播放器之外还提供可嵌入应用的 libVLC 引擎

## 热度与来源

GitHub Star 快照：19,450（查询于 2026-08-27）。VLC 主开发协作位于 VideoLAN 生态；本条目链接到 VideoLAN 官网和 GitHub 源码镜像，下载请以官网页面为准。
`,
			OfficialURL: "https://www.videolan.org/vlc/", RepoURL: "https://github.com/videolan/vlc", ReleaseURL: "https://www.videolan.org/vlc/",
			Developer: "VideoLAN", License: "GPL-2.0 / LGPL-2.1", Version: "3.0.23", IconSource: "https://avatars.githubusercontent.com/u/1389585?v=4", IconFile: "vlc.png",
			Stars: 19450, Platforms: []string{"windows", "macos", "linux", "android", "ios"}, Tags: []string{"播放器", "影音", "跨平台"}, Featured: false,
		},
		{
			Name: "Joplin", Slug: "joplin", Category: "productivity",
			Tagline: "隐私优先的 Markdown 笔记与待办，离线可用并支持端到端加密同步 · GitHub ★ 56.1k（2026-08-27）",
			Description: `## 简介

Joplin 是一款开源笔记与待办应用，使用 Markdown 保存内容，支持笔记本、标签、附件和全文搜索。它采用离线优先思路，断网时也能继续访问自己的资料。

## 主要特点

- 支持 Windows、macOS、Linux、Android 和 iOS
- 可与多种云服务同步，并支持端到端加密
- 支持从 Evernote 等工具导入，也提供浏览器 Web Clipper

## 热度与来源

GitHub Star 快照：56,128（查询于 2026-08-27）。图标来自 Joplin 官方仓库的品牌资源，下载入口跳转官方安装说明页。
`,
			OfficialURL: "https://joplinapp.org/help/install/", RepoURL: "https://github.com/laurent22/joplin", ReleaseURL: "https://github.com/laurent22/joplin/releases/latest",
			Developer: "Joplin", License: "AGPL-3.0", Version: "3.6.16", IconSource: "https://raw.githubusercontent.com/laurent22/joplin/dev/Assets/SquareIcon512.png", IconFile: "joplin.png",
			Stars: 56128, Platforms: []string{"windows", "macos", "linux", "android", "ios"}, Tags: []string{"笔记", "Markdown", "隐私"}, Featured: true,
		},
		{
			Name: "HandBrake", Slug: "handbrake", Category: "media",
			Tagline: "简单易用的开源视频转码器，快速生成 MP4、MKV 或 WebM · GitHub ★ 24.2k（2026-08-27）",
			Description: `## 简介

HandBrake 是开源视频转码器，可以把已有视频转换为更适合手机、电视、游戏机、电脑或浏览器播放的格式。

## 主要特点

- 支持 Windows、macOS 和 Linux
- 支持常见视频文件、摄像机素材、屏幕录制以及 DVD / Blu-ray 来源
- 可输出 MP4、MKV、WebM，并利用现代编码器控制画质与体积

## 热度与来源

GitHub Star 快照：24,188（查询于 2026-08-27）。项目官方说明其采用 GPL-2.0；下载按钮只跳转 handbrake.fr 官方下载页。
`,
			OfficialURL: "https://handbrake.fr/downloads.php", RepoURL: "https://github.com/HandBrake/HandBrake", ReleaseURL: "https://handbrake.fr/downloads.php",
			Developer: "HandBrake Team", License: "GPL-2.0", Version: "1.11.2", IconSource: "https://avatars.githubusercontent.com/u/627269?v=4", IconFile: "handbrake.png",
			Stars: 24188, Platforms: []string{"windows", "macos", "linux"}, Tags: []string{"视频", "转码", "多平台"}, Featured: false,
		},
		{
			Name: "Audacity", Slug: "audacity", Category: "media",
			Tagline: "免费开源的多轨音频编辑与录音工具，适合剪辑、降噪和混音 · GitHub ★ 17.6k（2026-08-27）",
			Description: `## 简介

Audacity 是面向 Windows、macOS 和 Linux 的多轨音频编辑与录音软件，适合剪切拼接、录音、降噪、音量处理和导出常见音频格式。

## 主要特点

- 多轨录音与编辑，适合播客、配音和音乐草稿
- 提供效果器、频谱查看和降噪等常用工具
- 支持插件扩展，并可按需安装 FFmpeg 等附加组件

## 热度与来源

GitHub Star 快照：17,638（查询于 2026-08-27）。软件采用 GPL-3.0；图标来自 Audacity 官方 GitHub 组织头像，下载入口为官方页面。
`,
			OfficialURL: "https://www.audacityteam.org/download/", RepoURL: "https://github.com/audacity/audacity", ReleaseURL: "https://www.audacityteam.org/download/",
			Developer: "Audacity Team", License: "GPL-3.0", Version: "3.7.8", IconSource: "https://avatars.githubusercontent.com/u/11648186?v=4", IconFile: "audacity.png",
			Stars: 17638, Platforms: []string{"windows", "macos", "linux"}, Tags: []string{"音频", "录音", "编辑"}, Featured: false,
		},
		{
			Name: "qBittorrent", Slug: "qbittorrent", Category: "network",
			Tagline: "简洁、稳定的开源 BitTorrent 客户端，支持 Windows、macOS 和 Linux · GitHub ★ 39.7k（2026-08-27）",
			Description: `## 简介

qBittorrent 是基于 Qt 和 libtorrent 的开源 BitTorrent 客户端，目标是提供轻量、稳定且功能完整的下载体验。

## 主要特点

- 支持 Windows、macOS、Linux 以及多种发行版
- 提供队列、限速、种子管理、Web UI 和 RSS 等能力
- 官方下载页提供安装包、AppImage、源码与校验信息

## 热度与来源

GitHub Star 快照：39,733（查询于 2026-08-27）。官方下载页同时提供 PGP 签名与 SHA-256 校验信息；本站不托管安装包，只提供官方入口。
`,
			OfficialURL: "https://www.qbittorrent.org/download", RepoURL: "https://github.com/qbittorrent/qBittorrent", ReleaseURL: "https://www.qbittorrent.org/download",
			Developer: "qBittorrent Project", License: "GPL-2.0-or-later", Version: "5.2.3", IconSource: "https://avatars.githubusercontent.com/u/2131270?v=4", IconFile: "qbittorrent.png",
			Stars: 39733, Platforms: []string{"windows", "macos", "linux"}, Tags: []string{"下载", "网络", "跨平台"}, Featured: false,
		},
		{
			Name: "Neovim", Slug: "neovim", Category: "development",
			Tagline: "面向可扩展性与易用性的新一代 Vim，内置 LSP、终端和 Lua 扩展能力 · GitHub ★ 102.0k（2026-08-27）",
			Description: `## 简介

Neovim 是专注于可扩展性和易用性的 Vim 分支，保留 Vim 的编辑模型，同时提供现代化 API、异步架构和更适合插件开发的 Lua 配置方式。

## 主要特点

- 内置 LSP 客户端、终端和语法树解析能力
- 支持 Lua 插件、远程插件以及多个 UI 连接同一会话
- 适合开发者、终端用户和希望深度定制编辑器的用户

## 热度与来源

GitHub Star 快照：101,991（查询于 2026-08-27）。项目贡献代码主要采用 Apache-2.0；下载入口跳转官方 GitHub Release 页面。
`,
			OfficialURL: "https://neovim.io/", RepoURL: "https://github.com/neovim/neovim", ReleaseURL: "https://github.com/neovim/neovim/releases/latest",
			Developer: "Neovim", License: "Apache-2.0", Version: "0.12.5", IconSource: "https://avatars.githubusercontent.com/u/6471485?v=4", IconFile: "neovim.png",
			Stars: 101991, Platforms: []string{"windows", "macos", "linux"}, Tags: []string{"编辑器", "终端", "开发"}, Featured: true,
		},
	}
}
