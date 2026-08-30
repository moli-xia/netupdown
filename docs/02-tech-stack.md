# 02 · 技术选型

> 状态：设计定稿 · 2026-08-27

## 1. 选型原则

1. **单二进制交付**：所有资源（默认主题、管理端构建产物、初始数据）经 `go:embed` 打进二进制；
2. **纯 Go、无 CGO**：`CGO_ENABLED=0` 可交叉编译到任意平台，SQLite 使用纯 Go 驱动；
3. **少而精的依赖**：优先标准库，第三方库须是维护活跃、口碑成熟的项目；
4. **前台资源自托管**：不引用境外 CDN，保证国内访问；
5. **为"主题可换"服务**：前台必须运行时加载模板（决定了模板技术路线，见 DR-4）。

## 2. 选型总表

### 2.1 后端（Go）

| 领域 | 选型 | 包路径 | 版本基线 | 备选 | 选择理由 |
|---|---|---|---|---|---|
| 语言 | Go | — | 1.25+（用当时最新稳定版） | — | 单二进制、交叉编译、高并发 IO |
| Web 框架 | Gin | `github.com/gin-gonic/gin` | v1.10+ | Echo / chi | 生态最大、中间件全、团队熟悉成本低 |
| ORM | GORM | `gorm.io/gorm` | v1.30+ | ent / sqlx | 多数据库方言、AutoMigrate、软删除内建 |
| SQLite 驱动 | glebarez/sqlite | `github.com/glebarez/sqlite` | 最新 | mattn（CGO） | 纯 Go（基于 modernc），免 CGO，支持 FTS5 |
| MySQL/PG 驱动 | GORM 官方 | `gorm.io/driver/mysql`、`gorm.io/driver/postgres` | 最新 | — | 可选数据库 |
| 配置加载 | koanf | `github.com/knadh/koanf/v2` | v2 | viper | 轻量，YAML + 环境变量合并 |
| CLI | cobra | `github.com/spf13/cobra` | v1.8+ | urfave/cli | 子命令 serve/admin/version |
| 日志 | slog + lumberjack | 标准库、`gopkg.in/natefinch/lumberjack.v2` | — | zap | 标准库结构化日志足够；lumberjack 负责轮转 |
| JWT | golang-jwt | `github.com/golang-jwt/jwt/v5` | v5 | — | 事实标准 |
| 密码哈希 | argon2id | `github.com/alexedwards/argon2id` | 最新 | bcrypt | OWASP 推荐算法，封装友好 |
| Markdown 渲染 | goldmark | `github.com/yuin/goldmark` | v1.7+ | blackfriday | CommonMark 兼容、扩展多（GFM 表格/删除线） |
| HTML 消毒 | bluemonday | `github.com/microcosm-cc/bluemonday` | 最新 | — | Markdown 渲染结果防 XSS |
| 代码高亮 | chroma | `github.com/alecthomas/chroma/v2` | v2 | — | 更新日志中的代码块高亮（服务端） |
| S3 SDK | AWS SDK v2 | `github.com/aws/aws-sdk-go-v2`（s3、s3/presign、credentials） | v2 | minio-go | 一套 SDK 通吃所有 S3 兼容服务（R2/OSS/COS/MinIO/B2） |
| WebDAV 客户端 | gowebdav | `github.com/studio-b12/gowebdav` | 最新 | — | P2 驱动 |
| 定时任务 | cron | `github.com/robfig/cron/v3` | v3 | — | 统计聚合、清理任务 |
| 版本比较 | semver | `github.com/Masterminds/semver/v3` | v3 | — | 更新检查的核心 |
| 参数校验 | validator | `github.com/go-playground/validator/v10` | v10 | — | Gin binding 内置 |
| 图片处理 | imaging | `github.com/disintegration/imaging` | 最新 | — | 图标/截图缩放（纯 Go） |
| 限流 | x/time/rate | `golang.org/x/time/rate` | — | ulule/limiter | 标准扩展库 + 自建 IP LRU 即可 |
| 短 ID | xid | `github.com/rs/xid` | 最新 | uuid | 上传会话 ID 等，短且可排序 |
| TOTP（P2） | otp | `github.com/pquerna/otp` | 最新 | — | 两步验证 |
| 测试 | testify | `github.com/stretchr/testify` | 最新 | — | 断言与 mock |

### 2.2 前台（访客侧，SSR）

| 领域 | 选型 | 说明 |
|---|---|---|
| 模板引擎 | 标准库 `html/template` | 运行时解析主题模板（换主题的前提），自动上下文转义防 XSS |
| CSS | Tailwind CSS v4（**standalone CLI**） | 主题开发期用独立二进制编译，产物 CSS 提交仓库；运行时零 Node 依赖 |
| 交互 JS | Alpine.js 3（本地引入，~15KB） | 主题切换、移动端菜单、复制提取码等轻交互；无构建步骤 |
| 图标 | Lucide（静态 SVG 内联/雪碧图） | 不引字体图标，体积可控 |
| 字体 | 系统字体栈 | 中文走系统字体，零下载；站长可经 F-504 注入自定义字体 |

### 2.3 管理后台（SPA）

| 领域 | 选型 | 版本基线 |
|---|---|---|
| 框架 | Vue 3 + `<script setup>` + TypeScript | Vue 3.5+ / TS 5.x |
| 构建 | Vite | 7.x |
| UI 组件库 | Naive UI | 2.x（TS 原生、树摇、内置暗色主题，气质契合"时尚精美"） |
| 状态 | Pinia | 3.x |
| 路由 | Vue Router | 4.x |
| 图表 | Apache ECharts（按需引入） | 5.x |
| HTTP | Axios（封装拦截器：envelope 解包、401 自动刷新） | 1.x |
| 包管理 | pnpm | 10.x |
| Node | Node.js LTS | 22.x |

管理后台构建产物复制到 `internal/assets/admin/` 后 `go:embed` 内嵌，线上零 Node 依赖。

## 3. 关键决策记录（ADR）

### DR-1 前台采用 Go SSR，而非 SPA
下载站的命脉是 SEO 与首屏速度。html/template 服务端渲染让每个应用详情页都是完整 HTML，搜索引擎可直接收录；同时避免访客侧加载大 JS 包。交互增强用 Alpine.js 渐进补充。

### DR-2 管理后台采用 Vue SPA，而非模板
后台交互密度高（发布向导、分片上传进度、动态存储表单、图表），SPA 的开发效率与体验显著更好；SEO 无关紧要。代价是引入 Node 构建链，用 go:embed + Makefile 把复杂度锁在构建期。

### DR-3 默认数据库 SQLite，可选 MySQL/PostgreSQL
个人站点写少读多，SQLite（WAL 模式）绰绰有余且零运维；GORM 抽象保证想换就换。约束：SQLite 单写者，所有写路径经服务层即可，无需额外队列。

### DR-4 模板引擎必须运行时加载 —— 否定 templ/编译期方案
「更换模板」意味着主题是用户数据（zip 上传、放在 `data/themes/`），必须运行时解析渲染。`templ` 等编译期模板无法承载第三方主题，故选标准库 `html/template`（另有自动转义安全红利）。

### DR-5 纯 Go、禁 CGO
换用 `glebarez/sqlite` 纯 Go 驱动后全链路无 CGO，`CGO_ENABLED=0` 一条命令交叉编译 Linux/amd64+arm64 产物，Docker 镜像可用极小基础镜像。

### DR-6 单仓库（monorepo）
Go 服务、管理端源码（`web/admin`）、默认主题源码（`web/themes/aurora`）、文档同仓，Makefile 统一编排构建。个人项目降低协同成本。

### DR-7 对象存储统一走 S3 协议
R2/OSS/COS/MinIO/B2 全部提供 S3 兼容端点，用 AWS SDK v2 一个驱动覆盖（配置差异仅 endpoint/region/path-style）。不为各云厂商单独引 SDK。

## 4. go.mod 依赖基线（预估）

```text
require (
    github.com/gin-gonic/gin v1.10.x
    gorm.io/gorm v1.30.x
    github.com/glebarez/sqlite vX
    gorm.io/driver/mysql vX          // 可选构建
    gorm.io/driver/postgres vX       // 可选构建
    github.com/knadh/koanf/v2 v2.x
    github.com/spf13/cobra v1.8.x
    gopkg.in/natefinch/lumberjack.v2 v2.x
    github.com/golang-jwt/jwt/v5 v5.x
    github.com/alexedwards/argon2id vX
    github.com/yuin/goldmark v1.7.x
    github.com/microcosm-cc/bluemonday v1.0.x
    github.com/alecthomas/chroma/v2 v2.x
    github.com/aws/aws-sdk-go-v2/... vX
    github.com/robfig/cron/v3 v3.x
    github.com/Masterminds/semver/v3 v3.x
    github.com/disintegration/imaging v1.6.x
    golang.org/x/time vX
    github.com/rs/xid v1.x
    github.com/stretchr/testify v1.9.x
)
```

> 实际版本以开发时 `go get` 最新稳定版为准；每次升级跑 `govulncheck ./...`（见 09 册）。

## 5. 开发工具链

| 工具 | 用途 |
|---|---|
| `golangci-lint` | 静态检查（配置见 11 册） |
| `air` | Go 热重载开发 |
| `tailwindcss`（standalone） | 主题 CSS 编译（`--watch` 开发 / `--minify` 产出） |
| `make`（或 PowerShell 脚本） | 统一构建入口；Windows 下可用 Git Bash / `scripts/*.ps1` |
| GitHub Actions | CI：lint + test + govulncheck + 交叉编译 + Docker 镜像 |
