# 08 · 管理后台设计

> 状态：设计定稿 · 2026-08-27

技术栈见 02 册 §2.3：Vue 3 + TS + Vite + Pinia + Naive UI + ECharts，构建产物 go:embed 内嵌，挂载于 `/admin`。

## 1. 工程结构

```text
web/admin/
├── index.html
├── vite.config.ts            # base: '/admin/'; dev proxy: /api → http://localhost:8080
├── src/
│   ├── main.ts / App.vue
│   ├── api/                  # Axios 实例 + 各域 API 封装（envelope 解包、401 自动刷新重放）
│   ├── stores/               # auth（access token 仅存内存）、settings、theme
│   ├── router/               # 路由 + 登录守卫
│   ├── layouts/              # AdminLayout（侧栏 + 顶栏 + 面包屑 + 暗色切换）
│   ├── views/                # 页面（见 §2）
│   ├── components/           # UploadDropzone、ChunkUploader、StorageConfigForm、
│   │                         # ThemeSettingForm、MarkdownEditor、ImagePicker、ConfirmButton…
│   └── utils/                # 格式化、sha256(WebCrypto 流式)、常量映射
└── package.json / pnpm-lock.yaml
```

## 2. 页面与路由

| 路由 | 页面 | 要点 |
|---|---|---|
| `/login` | 登录 | 用户名密码（+TOTP 预留位）；错误提示剩余锁定时间 |
| `/` | 仪表盘 | 总览卡片（应用/版本/累计下载/今日下载/今日浏览/存储用量）；30 天趋势双线图；Top10 下载榜 |
| `/apps` | 应用列表 | 表格：图标/名称/类型/分类/状态/最新版/下载量/操作；筛选与搜索；发布/下架/置顶快捷开关 |
| `/apps/new` `/apps/:id/edit` | 应用编辑 | 分组表单：基本信息（名称/slug 自动生成可改/类型/分类/标签/tagline）· 介绍（Markdown 编辑器带预览与图片粘贴上传）· 媒体（图标裁剪 1:1、截图排序）· 链接（官网/仓库/开发者/许可）· SEO（三字段 + 生成预览）· 平台勾选；已保存应用可打开前台预览，草稿也支持 |
| `/apps/:id/releases` | 版本列表 | 时间线式；渠道徽标；发布/删除 |
| `/apps/:id/releases/new` | **发布向导** | 见 §3 |
| `/releases/:id/edit` | 版本编辑 | 信息 + 文件与源管理（同向导后两步的复用组件） |
| `/categories` `/tags` | 分类/标签 | 行内编辑 + 拖拽排序 |
| `/pages` | 单页管理 | 列表 + Markdown 编辑 |
| `/links` | 友链（P2） | |
| `/storages` | 存储管理 | 卡片列表（驱动徽标/默认标记/用量）；新增抽屉：选驱动 → 动态表单（§4）；「测试连接」按钮实时反馈延迟 |
| `/themes` | 主题管理 | 预览图卡片墙；当前主题高亮；上传 zip；「配置」抽屉按 settings schema 动态表单；启用前可打开新窗口预览前台 |
| `/settings` | 站点设置 | Tab 分组（§5），批量保存 |
| `/logs/operations` `/logs/downloads` | 日志 | 筛选 + 分页表格 |
| `/profile` | 个人中心 | 改资料/密码；会话设备列表与踢出 |

## 3. 发布向导（F-208）

四步 Steps，可保存草稿随时中断：

1. **版本信息**：版本号（校验 semver 并提示；查重）、渠道、更新日志（Markdown）、可选强更阈值 min_required_version；
2. **上传文件**：拖拽多文件 → 每个文件卡片选择 OS/架构/形态（按文件名智能预填：含 `win`/`.exe` → windows，`arm64` → arm64…）→ ChunkUploader 组件：WebCrypto 流式算 SHA256 → init（秒传命中直接完成）→ 并发 3 路分片 PUT（进度条、失败重试、断点续传）→ complete；
3. **下载源**：每个文件自动生成托管源「直链下载」；可添加外链源（名称下拉常用网盘 + URL + 提取码）、拖拽排序定默认；也可整文件跳过上传只配外链（收录类应用常态）；
4. **确认发布**：预览详情页摘要 → 发布 / 存为草稿。

## 4. 存储动态表单

`StorageConfigForm` 按 driver 渲染字段（与 06 册 §3 schema 一一对应），前端仅做必填与格式校验，真校验靠「测试连接」。S3 表单附「常用服务预设」下拉（R2/OSS/COS/MinIO/B2）自动填 endpoint 模板与 path-style。secret 类输入框回显 `******`，未修改则不上传该字段。

## 5. 站点设置项全表（settings 注册表）

前端 Tab 与后端注册表（settingsvc 内定义：键、类型、默认值、校验）保持一致：

| Tab | 键 | 默认 | 说明 |
|---|---|---|---|
| 基本 | `site.title` | 造物工坊 | 站点名 |
| | `site.subtitle` | 个人应用发布与软件分发 | 副标题 |
| | `site.logo` / `site.favicon` | "" | /uploads URL |
| | `site.description` / `site.keywords` | "" | 默认 SEO |
| | `site.footer` | "" | 页脚 HTML |
| | `site.icp` / `site.police_icp` | "" | 备案号 |
| 外观 | `theme.active` | aurora | 当前主题 |
| | `theme.cfg.{id}` | {} | 各主题配置 JSON |
| 下载 | `download.dedup_window_min` | 10 | 计数去重窗口（分钟） |
| | `download.log_retention_days` | 90 | 明细保留天数 |
| | `download.require_referer` | false | 防盗链开关 |
| | `download.allowed_referers` | [] | 白名单域名 |
| | `download.show_hash` | true | 前台展示 SHA256 |
| | `download.landing_for_external` | true | 有码外链走落地页 |
| 内容 | `content.apps_page_size` | 24 | 列表每页 |
| | `content.home_latest_limit` | 8 | 首页最新数 |
| SEO | `seo.og_image_default` | "" | 默认分享图 |
| | `seo.sitemap_enabled` / `seo.feed_enabled` | true | |
| | `seo.robots_txt` | ""（用内置模板） | 覆写 robots |
| 注入 | `custom.head` / `custom.foot` | "" | 全站注入（统计脚本等，站长自担风险） |
| 隐私 | `privacy.ip_mode` | truncate | 下载日志 IP：plain/truncate/hash |
| 评论(P2) | `comment.enabled` 等 | false | 预留 |
| 用户(P2) | `user.register_enabled` | false | 预留 |

## 6. 认证与请求约定

- access token 只存 Pinia 内存（刷新页面后靠 Cookie refresh 换新），**不落 localStorage**（XSS 面最小化）；
- Axios 响应拦截：`code=10003` → 排队进刷新单飞（single-flight）→ 重放原请求；刷新失败跳登录；
- 全局错误 toast 统一处理非 0 code；表单校验错误就地展示；
- 上传等长任务用独立进度 UI，页面可离开（向导有未完成上传时离开提示）。

## 7. 暗色与视觉

- Naive UI `darkTheme` + `useOsTheme`，三态与前台一致（auto/light/dark，localStorage 键 `nud-admin-theme`）；
- 主色对齐 Aurora `--brand`（`themeOverrides.common.primaryColor = #4f6df5`）；
- 仪表盘 ECharts 注册亮暗两套主题并随全局切换。

## 8. 构建与内嵌

```makefile
admin:            ## 构建管理端并同步到内嵌目录
	cd web/admin && pnpm install --frozen-lockfile && pnpm build
	rm -rf internal/assets/admin && cp -r web/admin/dist internal/assets/admin

build: admin      ## 完整构建
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o bin/netupdown ./cmd/netupdown
```

- `internal/assets/admin` 提交一份构建产物到仓库（保证纯 Go 环境也能 `go build` 出完整程序；CI 里重建覆盖）；
- Go 侧：`//go:embed all:admin` + `fs.Sub`，路由 `/admin/*` 未命中文件回退 `index.html`；`index.html` 响应 `Cache-Control: no-cache`，hash 资源 `max-age=31536000, immutable`；
- 开发：`pnpm dev`（5173，proxy `/api`）+ `air`（8080）双进程。
