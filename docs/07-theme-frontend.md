# 07 · 主题与前台设计

> 状态：设计定稿 · 2026-08-27

本册定义两件事：**前台产品设计**（页面信息架构、默认主题 Aurora 的视觉规范、亮暗方案、SEO）与**主题系统机制**（主题包格式、渲染引擎、开发规范）。

## 1. 前台页面信息架构

| 页面 | 区块（自上而下） |
|---|---|
| 首页 `/` | 顶栏（站名/导航/搜索/主题切换）→ 重点产品（特色卡片）→ 最近更新（卡片列表 ×8）→ 产品系列 → 更新保障说明 → 页脚（自定义 HTML + 备案号 + RSS） |
| 产品中心 `/apps` | 面包屑 → 产品系列 Tab → 筛选栏（搜索 / 平台 / 排序）→ **卡片列表**（单列，每行一款自研产品）→ 分页 |
| 应用详情 `/apps/:slug` | 头部（图标 + 名称 + 自研/收录徽标 + tagline + 元信息一行：开发者/平台/更新时间/下载量/浏览量/文件数）→ 操作区（**主下载按钮**=默认源 + 平台 + 大小；次级按钮：全部下载源、官网、源码、历史版本）→ 截图横向滚动条 → 左栏 Markdown 介绍 + 最新更新日志 / 右栏 文件与下载源 + 应用信息 |
| 历史版本 `/apps/:slug/releases` | 版本卡片列表（版本号/渠道徽标/日期/日志折叠/文件行） |
| 搜索 `/search?q=` | 结果卡片列表 + 空态引导 |
| 分类 `/categories/:slug`、标签 `/tags/:slug` | 同列表页（预置筛选） |
| 单页 `/pages/:slug` | 文章排版（prose，窄栏） |
| 下载落地页 `/d/:id`（外链带提取码时） | 应用摘要卡 + 醒目提取码（一键复制）+「前往网盘」按钮 + 免责提示 |
| 404 / 错误页 | 状态码 + 说明 + 搜索框 + 返回首页 |

交互增强（约 4KB 原生 JS，全部渐进式，无 JS 也可用）：主题切换、移动端抽屉菜单、截图灯箱、SHA256/提取码复制、筛选下拉自动提交、`/` 聚焦搜索。

## 2. 默认主题 Aurora 设计规范

**气质关键词**：简洁、克制、内容优先。整站只用留白、字重与 1px 分隔线建立层次，不使用渐变、光斑、玻璃拟态与浮起阴影。核心组件是**卡片列表**：一行一条内容，横向铺满，纵向靠分隔线切分。

### 2.1 设计令牌（CSS Variables 契约）

所有主题**必须**定义下表变量（亮暗各一套），程序与主题共同遵守；Aurora 的取值如下：

| 变量 | Light | Dark | 用途 |
|---|---|---|---|
| `--bg` | `#ffffff` | `#0f1115` | 页面底色 |
| `--bg-soft` | `#f7f8fa` | `#14171d` | 分区底色（页脚等） |
| `--surface` | `#ffffff` | `#14171d` | 卡片/面板 |
| `--surface-soft` | `#f7f8fa` | `#191d24` | 卡片内浅底（hover / 图标位） |
| `--text-1` | `#16181d` | `#e6e8ec` | 主文本 |
| `--text-2` | `#596170` | `#9aa2ae` | 次文本 |
| `--text-3` | `#8a919d` | `#6f7783` | 弱文本/占位 |
| `--line` | `#e3e6ea` | `#272c35` | 主分隔线（卡片外框） |
| `--line-soft` | `#edeff2` | `#1f242c` | 次分隔线（列表行之间） |
| `--brand` | `#2f6feb` | `#6c9bff` | 主色（可被主题配置 `accent` 覆盖） |
| `--brand-hover` | `color-mix(brand 84%, #000)` | `color-mix(brand 84%, #fff)` | 主色悬停（由 `--brand` 推导，主色改动自动跟随） |
| `--brand-tint` | `color-mix(brand 10%, transparent)` | 同 | 主色浅底（徽标/选区） |
| `--brand-contrast` | `#ffffff` | `#0f1115` | 主色上的文字 |
| `--ok` / `--warn` / `--err` | `#15803d / #b45309 / #b91c1c` | `#4ade80 / #fbbf24 / #f87171` | 语义色 |
| `--code-bg` | `#f4f5f7` | `#1a1f27` | 代码块/哈希底色 |
| `--radius-s/m` | `6px / 10px` | 同 | 圆角（小控件 / 卡片） |

字体：`--font-sans: system-ui, -apple-system, "Segoe UI", "Microsoft YaHei", "PingFang SC", "Noto Sans CJK SC", sans-serif`；`--font-mono: ui-monospace, "Cascadia Code", Consolas, monospace`（SHA256/版本号用）。

> 3.0 起不再提供 `--g1/--g2/--g3` 渐变色与 `--shadow-*` 阴影令牌：层次一律用 `--line` / `--line-soft` 表达。

### 2.2 版式与组件

- 内容最大宽 `1040px`（正文类页面再收窄到 `760px`）；正文 15px/1.65；标题只用 22 / 16 / 15px 三级，靠字重而非字号拉开层次；
- **卡片列表 `.card-list`**：外层 1px 边框 + 10px 圆角，内部条目用 `--line-soft` 分隔，条目 hover 仅换浅底色，不位移不投影；
- **应用条目 `.app-item`**：44px 图标 + 名称 + 徽标（自研 / 分类）+ 两行截断简介 + 元信息行（平台图标 · 下载量 · 发布日期），右端一个箭头指示可进入；
- 徽标与 chip 统一为 1px 描边 + 浅底，只有「自研」用主色描边区分；
- 主下载按钮：纯 `--brand` 实心，副文案显示平台与大小；其余按钮一律描边样式；
- 首页直接进入产品内容，不使用 Hero 大图或独立简介块；导航下方不重复放置搜索框；
- 更新日志/介绍排版：prose 样式（引用、代码块、表格横向滚动容器）；
- 数字统计用等宽数字 `font-variant-numeric: tabular-nums`。

### 2.3 性能预算

首页与详情页（不含截图）：CSS ≤ 40KB gzip、JS ≤ 20KB gzip（当前 CSS≈26KB、JS≈4KB，均为未压缩体积）；图片全部 `loading="lazy"` + 显式宽高防 CLS；无滚动动画与装饰性动效，LCP < 2s（4G）。

## 3. 亮/暗主题实现（F-110）

三态：`light` / `dark` / `auto`（默认，跟随系统）。约定对所有主题生效：

1. `<html>` 上以 `data-theme="light|dark"` 标记显式选择；`auto` 时不落属性，靠 `prefers-color-scheme`；
2. 主题 CSS 写法（tokens 双定义）：

```css
:root { /* light 值 */ }
@media (prefers-color-scheme: dark) { :root:not([data-theme="light"]) { /* dark 值 */ } }
:root[data-theme="dark"] { /* dark 值（显式优先） */ }
```

3. `base.html` 的 `<head>` **最前**内联无闪烁脚本（先于任何 CSS 渲染决定主题）：

```html
<script>
(function(){try{
  var t = localStorage.getItem('nud-theme');           // light | dark | null(auto)
  if (t === 'light' || t === 'dark') document.documentElement.dataset.theme = t;
}catch(e){}})();
</script>
```

4. 切换按钮三态循环 auto→light→dark，写 localStorage（auto 时移除键），同步更新 `<meta name="theme-color">`（亮 `#ffffff` / 暗 `#0f1115`，用两条带 media 的 meta + JS 兜底）；
5. 图片适配：主题内装饰图用 CSS 变量/`<picture media="(prefers-color-scheme)">`；截图统一加细边框保证暗底可辨。

## 4. 主题系统机制

### 4.1 主题包结构

```text
themes/aurora/
├── theme.json               # 元数据 + 配置 schema（必须）
├── preview.png              # 后台展示的预览图（建议 1200×760）
├── templates/
│   ├── layouts/base.html    # 骨架：<head>、顶栏、页脚、block 定义（必须）
│   ├── partials/            # 片段：header/footer/app-card/pagination/…
│   ├── index.html           # 必须
│   ├── list.html            # 必须（列表/分类/标签/搜索共用或分拆）
│   ├── detail.html          # 必须
│   ├── releases.html        # 可选（缺失回退默认主题）
│   ├── search.html          # 可选（缺失用 list.html）
│   ├── page.html            # 可选
│   ├── download.html        # 可选（外链落地页）
│   └── error.html           # 可选（404/500，ctx.Status 区分）
└── static/                  # css/js/img，映射到 /themes/{id}/static/*
```

加载优先级：`data/themes/{id}`（用户安装）> 内嵌默认主题。用户可放置同 id=`aurora` 的目录来覆写默认主题。

### 4.2 theme.json

```json
{
  "id": "aurora",
  "name": "Aurora 简约",
  "version": "3.2.1",
  "author": "NetUpDown",
  "homepage": "https://netupdown.com",
  "description": "内置默认主题：面向自研产品发布的精致展示体验，亮暗双模式，信息密度友好",
  "preview": "preview.png",
  "min_app_version": "0.1.0",
  "settings": [
    { "key": "accent",        "label": "主色",         "type": "color",  "default": "#2f6feb" },
    { "key": "show_dev_band", "label": "显示更新 API 引导条", "type": "switch", "default": true }
  ]
}
```

`settings[].type` 支持：`text / textarea / color / switch / select / number / image`（image 值为 /uploads URL）。后台按 schema 自动渲染表单（08 册），保存到 settings 键 `theme.cfg.{id}`，模板经 `.Theme.Config.<key>` 读取。

### 4.3 渲染上下文（模板可用数据）

所有页面公共：

| 变量 | 类型 | 说明 |
|---|---|---|
| `.Site` | SiteCtx | `Title Subtitle Description Keywords LogoURL FaviconURL FooterHTML ICP PoliceICP BaseURL` 及注入 `HeadInject FootInject`（template.HTML，仅站长可控） |
| `.Theme` | ThemeCtx | `ID Version StaticBase Config(map)` |
| `.Page` | PageMeta | `Title Description Canonical OGImage OGType NoIndex`（`<head>` 用） |
| `.Nav` | []NavItem | 导航（分类 + 单页 + 自定义链接） |
| `.Year` | int | 页脚年份 |

页面专属（节选）：

| 模板 | 数据 |
|---|---|
| index.html | `.Featured []AppItem`、`.Latest []AppItem`、`.Categories []CategoryItem`、`.Stats {AppCount, TotalDownloads}` |
| list.html | `.Apps []AppItem`、`.Pagination {Page PageSize Total TotalPages HasPrev HasNext PrevURL NextURL Pages}`、`.Filter {Category Platform Type Sort Q}`、`.Categories` |
| detail.html | `.App AppDetail`（含 `.LatestRelease{Version ChangelogHTML PublishedAt AssetGroups}`；AssetGroup 按 OS 分组，内含 Asset 与 Sources）、`.Related []AppItem` |
| releases.html | `.App AppBrief`、`.Releases []ReleaseView`、`.Pagination` |
| download.html | `.App AppBrief`、`.Asset AssetView`、`.Source {Name ExternalURL ExtractCode}` |
| error.html | `.Status int`、`.Message` |

> 视图模型在 `internal/web/viewmodel.go` 集中定义并保持向后兼容 —— **这是主题的 API**，字段只加不减，删改需升 min_app_version。

### 4.4 模板函数

| 函数 | 示例 | 说明 |
|---|---|---|
| `asset` | `{{ asset "css/main.css" }}` | → `/themes/aurora/static/css/main.css?v=1.0.0`（版本号缓存穿透） |
| `absURL` | `{{ absURL .App.URL }}` | 拼 BaseURL 绝对地址（OG/JSON-LD 用） |
| `markdown` | `{{ markdown .App.Description }}` | 服务端渲染 + bluemonday 消毒，返回 template.HTML |
| `formatBytes` | `52428800` → `50.0 MB` | |
| `timeAgo` / `dateFmt` | `{{ dateFmt .PublishedAt "2006-01-02" }}` | 按站点时区 |
| `setting` | `{{ setting "site.title" }}` | 白名单键（site.* / seo.*），防止读敏感设置 |
| `osName` / `osIcon` | `windows` → `Windows` / SVG 图标 | 内置平台图标（Lucide） |
| `default` `add` `seq` `contains` `json` `dict` `truncate` | — | 常用工具；`dict` 用于给 partial 传参 |
| `vtLink` | sha256 → VirusTotal 查询 URL | 详情页"病毒扫描"外链 |

安全：不提供任意文件读取/命令执行类函数；`html/template` 自动上下文转义，唯二的 raw HTML 出口是 `markdown`（已消毒）与站长注入（自担）。

### 4.5 渲染引擎实现

- 启动/切换主题时：解析 theme.json → `template.New(page).Funcs(funcs)` 逐页面编译（每页面 = base + partials + 页面文件的集合），缓存 `map[page]*template.Template`；
- **回退链**：活动主题缺失某页面模板或执行 panic → 用内嵌 aurora 的同名模板渲染并记 error 日志（保证换劣质主题不白屏）；
- 开发模式 `theme.dev=true`：每请求重编译 + 响应禁缓存；
- 渲染写入 `bytes.Buffer` 成功后才写响应（避免半页输出），> 1MB 页面记 warn。

### 4.6 安装 / 切换 / 卸载

- 上传 zip（≤ 20MB）→ 校验：条目数 ≤ 500、解压后 ≤ 50MB、**zip-slip 防护**（逐条目 secureJoin）、仅允许扩展名白名单（html/css/js/json/png/jpg/webp/svg/woff2/txt/md）、theme.json 必须合法且 id `[a-z0-9-]{2,40}`；
- 解压到 `data/themes/{id}`（已存在则先备份为 `{id}.bak` 再覆盖，视为升级）；
- 切换：activate → 编译新主题成功后才原子切换指针，失败回滚并报错；
- 卸载：删除目录；当前激活或内置主题不可卸载；
- 主题静态资源响应头同 03 册 §3.2 缓存策略。

### 4.7 主题开发者指引（随仓库发布 THEME-DEV.md，要点）

最小主题 = theme.json + layouts/base.html + index/list/detail 三页；从复制 aurora 起步；本地 `netupdown serve --theme-dev` 热重载；遵守 §2.1 CSS 变量契约与 §3 亮暗约定即可自动获得亮暗支持。

## 5. SEO 规范（F-113）

- `<title>`：详情页 `{App.Name} 下载 - {最新版本} | {Site.Title}`；列表页 `{分类} - {Site.Title}`；可被 seo_title 覆盖；
- meta description：seo_description 或 tagline 自动截断 160 字；canonical 全站输出；
- Open Graph / Twitter Card：`og:type=website`、`og:image`= 应用图标或首图（缺省站点默认图 `seo.og_image_default`）；
- JSON-LD（详情页）：

```json
{
  "@context": "https://schema.org", "@type": "SoftwareApplication",
  "name": "MyApp", "operatingSystem": "Windows, macOS",
  "applicationCategory": "UtilitiesApplication",
  "softwareVersion": "1.4.1", "fileSize": "50MB",
  "datePublished": "2026-08-20",
  "offers": {"@type": "Offer", "price": "0", "priceCurrency": "CNY"},
  "downloadUrl": "https://netupdown.com/apps/myapp"
}
```

- `sitemap.xml`：首页/列表首屏/全部已发布应用（lastmod=最新版本时间）/分类/单页，>1000 条分索引；
- `robots.txt`：默认放行并指向 sitemap，禁抓 `/admin` `/api` `/d/`；内容可在设置覆写；
- RSS `/feed.xml`：最近 20 条版本发布，条目 = 「AppName v1.4.1 发布」+ 日志 HTML；
- 草稿/下架应用返回 404；未发布内容绝不出现在 sitemap/feed。

## 6. 无障碍

交互元素可聚焦 + `:focus-visible` 样式；图片 alt；颜色对比度 ≥ 4.5:1（tokens 已按此校准）；下载按钮用真实 `<a href>`（可右键另存/复制链接）。
