# NetUpDown

> 独立开发者的软件发布与更新中心 · [netupdown.com](https://netupdown.com)

NetUpDown 是一个用 **Go** 编写的单二进制自托管程序，面向独立开发者集中展示和发布自己的产品，提供版本管理、多平台安装包、下载统计与客户端在线更新检查。

## 特性

- 📦 **产品与版本管理**：自研产品、多版本、多渠道（stable/beta/alpha）、Markdown 更新日志
- 💾 **多存储后端**：本地磁盘、S3 兼容对象存储（Cloudflare R2 / 阿里云 OSS / 腾讯云 COS / MinIO）、WebDAV、网盘外链（含提取码）
- ⬇️ **智能下载分发**：每个文件多下载源、S3 预签名直链、本地断点续传、下载统计与去重防刷
- 🔄 **客户端更新检查 API**：自研应用一行代码接入在线更新（含强制更新、灰度渠道）
- 🎨 **主题系统**：前台模板可插拔（zip 安装 / 一键切换 / 主题配置项），内置默认主题 Aurora（简洁卡片列表、内容优先）
- 🌗 **亮 / 暗主题**：跟随系统 + 手动切换，无闪烁
- 🔍 **SEO 友好**：服务端渲染、sitemap、RSS、Open Graph、JSON-LD 结构化数据
- 🚀 **部署极简**：单二进制 + SQLite 零依赖启动，也支持 MySQL/PostgreSQL 与 Docker

## 技术栈

Go 1.25+ · Gin · GORM · SQLite（默认）· html/template 主题引擎 · Tailwind CSS · Alpine.js · Vue 3 + Naive UI（管理后台，go:embed 内嵌）

## 文档

完整开发文档见 [docs/](docs/README.md)。

## 当前实现

项目已具备可运行的首个开发版本：

- Go 单二进制、SQLite WAL、自动迁移与初始数据；
- 管理员 CLI、Argon2id、短效 JWT、轮换刷新令牌、登录限流与锁定；
- 应用、分类、版本、文件、下载源、单页、设置和审计 API；
- 本地与 S3 兼容存储，敏感配置 AES-GCM 加密，整文件及分片/断点上传；
- 本地 Range 下载、S3 预签名/公共域名跳转、外链与提取码落地页；
- 更新检查 API（SemVer、渠道、强更和 OS/架构匹配）；
- Aurora SSR 前台（自研产品发布站设计）、亮暗模式、响应式布局、搜索、RSS、sitemap 与 robots；
- Vue 3 + Naive UI 管理端，构建产物内嵌；
- 主题 ZIP 安装、校验、激活回滚、模板回退与运行时静态资源。

后续增强项仍按 [开发计划](docs/11-roadmap.md) 的 P2 Backlog 推进。

## 快速开始

```bash
cp config.example.yaml config.yaml
go build -o netupdown ./cmd/netupdown
printf 'your-strong-password\n' | ./netupdown admin create --username admin --password-stdin
./netupdown serve -c config.yaml
```

访问前台 `http://localhost:8080/`，管理端为 `http://localhost:8080/admin/`。

管理端源码位于 `web/admin`；完整构建可在 Windows 执行 `scripts/build.ps1`。

## 许可

程序代码采用 MIT 许可。
