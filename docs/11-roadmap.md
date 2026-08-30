# 11 · 开发计划与规范

> 状态：设计定稿 · 2026-08-27

## 1. 里程碑总览

按"业余时间每周约 10 小时"估算；全职可压缩到 1/3。每个里程碑结束时可运行、可演示。

| 里程碑 | 内容 | 预估 | 出口标准 |
|---|---|---|---|
| M0 工程初始化 | 仓库/lint/CI/配置加载/日志/cobra 骨架/healthz | 0.5 周 | `make build` 产出可跑二进制，CI 绿 |
| M1 核心域与后台基础 | 模型+AutoMigrate+种子；认证（JWT/argon2id/CLI 建号）；应用与分类 CRUD API；图片上传；admin SPA 骨架（登录/应用管理） | 2 周 | 后台能建应用、传图标 |
| M2 版本·存储·下载 | releases/assets/sources；local+s3 驱动+连通测试；整文件上传→分片+秒传；`/d` 分发全类型+异步计数；check-update API | 2.5 周 | 走通"发布→下载→客户端查更新" |
| M3 前台与默认主题 | 主题引擎（编译缓存/回退/dev 模式）；Aurora 全页面；亮暗切换；搜索；响应式打磨 | 2 周 | 前台达到可公开水准 |
| M4 完善 | 站点设置全表；SEO（sitemap/robots/OG/JSON-LD）；RSS；仪表盘统计+聚合任务；操作审计；限流与登录保护；单页管理 | 1.5 周 | v1.0 功能齐备 |
| M5 主题开放与上线 | 主题 zip 安装/切换/配置；THEME-DEV.md；发布向导打磨；Docker 镜像；部署 netupdown.com；备份演练 | 1.5 周 | **v1.0 上线** 🎉 |

**合计约 10 周**。上线检查单 = 09 册 §10 + 10 册 §9 恢复演练。

## 2. 各里程碑任务清单（可勾选）

### M0
- [ ] `go mod init` · 目录骨架（03 册 §2）· cobra `serve/version`
- [ ] koanf 配置加载 + env 覆盖 + config.example.yaml
- [ ] slog 初始化（stdout+文件轮转）· Recovery/RequestID/访问日志中间件
- [ ] golangci-lint 配置 · GitHub Actions（lint+test+govulncheck+build）
- [ ] Makefile / scripts/build.ps1（Windows 开发）

### M1
- [ ] 全部 GORM 模型 + AutoMigrate + PRAGMA + 种子数据
- [ ] `netupdown admin create` · authsvc（登录/刷新轮换/登出/会话表）
- [ ] apps/categories/tags service+repo+API（含状态机、slug 校验与保留字）
- [ ] 图片上传（白名单+魔数+xid 落盘）· `/uploads` 静态路由
- [ ] admin SPA：工程搭建、登录流（内存 token+单飞刷新）、应用列表/编辑页
- [ ] 单测：slug 校验、状态机、auth 轮换

### M2
- [ ] storage Driver 接口 + local（secureJoin+原子写）+ s3（含 presign/public_base_url）
- [ ] storages CRUD + 加密落库（AES-GCM/HKDF）+ 连通测试
- [ ] releases/assets/sources CRUD + 发布校验 + latest_release 冗余维护
- [ ] 上传：整文件 → 分片协议（init/chunk/complete/abort）+ 秒传 + 断点续传 + cron 清理
- [ ] `/d/:id` 全类型分发 + 异步计数器 + download_logs + 去重 LRU
- [ ] check-update（semver/version_code/强更/OS-arch 匹配 + 缓存）
- [ ] admin：存储管理页、发布向导（含 ChunkUploader/WebCrypto SHA256）
- [ ] 单测：下载决策 8 用例、semver 比较、分片合并、secureJoin；MinIO 集成测试

### M3
- [ ] 主题引擎：扫描/theme.json/编译缓存/回退链/dev 热重载/模板函数集
- [ ] 视图模型 viewmodel.go + 各页面 handler + 内存缓存（首页/分类）
- [ ] Aurora：base/partials/index/list/detail/releases/search/page/download/error + tokens + 亮暗无闪烁
- [ ] 搜索（LIKE）· 分页组件 · 移动端与无障碍走查
- [ ] 浏览量计数（去重 LRU）

### M4
- [ ] settings 注册表 + 缓存 + admin 设置页（分组 Tab）
- [ ] SEO：meta/OG/JSON-LD/canonical/sitemap/robots/feed
- [ ] statsvc + stat_daily 聚合任务 + 仪表盘（ECharts）
- [ ] operation_logs 埋点 + 日志页 · 限流中间件 + 登录锁定
- [ ] pages 单页 + 前台渲染

### M5
- [ ] 主题上传（zip 校验/zip-slip/覆盖升级）/激活/删除/配置表单
- [ ] THEME-DEV.md + 把 Aurora 抽为规范范例
- [ ] Dockerfile/compose/systemd 文档校验 · GitHub Actions 发 Release（多平台二进制+SHA256SUMS+镜像）
- [ ] 上线 netupdown.com：Caddy/HTTPS/备份 cron/恢复演练
- [ ] 全站压测（wrk 首页/详情/下载 302）与慢查询确认

## 3. v1.0 后路线（P2 Backlog，按价值排序）

1. **GitHub Releases 同步**（F-209）：应用绑定 repo，定时/手动拉取 tag→建 release+assets（外链 GitHub 或转存镜像到存储）；
2. **外链有效性巡检**（F-407）+ 失效通知（Webhook：Telegram/Bark/钉钉）；
3. Web 安装向导 · 备份工具化（后台一键导出/恢复）；
4. WebDAV 驱动 · 存储迁移工具；
5. 评论（内置或 Artalk 挂载）与多用户注册；
6. TOTP 两步验证 · 管理端 CSP；
7. 全文搜索（SQLite FTS5 优先，跨库则 bleve）；
8. i18n（英文前台）· PWA · 每应用独立 RSS；
9. 下载测速选源（前端探测各源延迟自动推荐）。

## 4. 工程规范

### Git 与版本
- 分支：`main` 保护 + 短生命周期 `feat/xxx`、`fix/xxx`；
- 提交：Conventional Commits（`feat: ...` `fix: ...` `docs: ...` `refactor: ...` `chore: ...`）；
- 版本：SemVer；`CHANGELOG.md` 按 Keep a Changelog 维护；tag 触发 CI 发布。

### Go 代码规范
- `gofumpt` 格式化；golangci-lint 启用：`govet, errcheck, staticcheck, revive, gosec, misspell, unconvert, gocritic, noctx, bodyclose`；
- 错误：向上包装 `fmt.Errorf("load theme %s: %w", id, err)`；业务错误用 `apperr` 带码；不吞错、不裸 panic（仅 bootstrap 可 fatal）；
- Context 贯穿：所有 service/repo/storage 方法首参 `ctx`；
- 命名对齐 04 册术语表；导出符号必有 doc comment。

### 测试策略
| 层 | 方式 | 覆盖要求 |
|---|---|---|
| pkg/service 纯逻辑 | 表驱动单测 | semver 比较、下载决策、secureJoin、slug、状态机 **必测** |
| repo | SQLite 内存库（`:memory:`） | 关键查询 |
| API | httptest 全栈（内存库+local 存储） | 登录、发布流、check-update、/d 各分支 |
| s3 驱动 | MinIO（docker，CI service）集成测试，本地可 `-short` 跳过 | Put/Range/Presign |
| 前台 | 模板渲染冒烟（每页 200 + 关键字断言） | 全页面 |
| E2E（P2） | Playwright | 发布→下载主流程 |

### 文档维护
代码与文档同 PR 更新；接口变更先改 05 册再实现；视图模型（主题 API）变更必须记录在 CHANGELOG 的 "Theme API" 小节。

## 5. 风险清单

| 风险 | 影响 | 缓解 |
|---|---|---|
| 大文件上传经反代踩坑（body 限制/超时/缓冲落盘） | 发布流程受阻 | 分片 5MiB 从根上规避；10 册给出各反代参数 |
| SQLite 并发写瓶颈 | 高峰下载计数写放大 | WAL + 异步批量计数；必要时切 MySQL（GORM 无痛） |
| 网盘外链失效 | 访客体验差 | P2 巡检 + 前台"报告失效"入口（P2） |
| 第三方主题质量差/恶意 | 白屏/XSS | 回退链 + zip 校验 + html/template 转义；主题只装可信来源 |
| 收录软件的版权/合规 | 站点风险 | 只链官方渠道或注明转存来源；免责声明；不碰破解 |
| 密钥丢失（master key） | 存储配置不可解 | 文档强提示离线备份；备份脚本包含 data/secret |
| 单人项目烂尾 | — | 里程碑各自可交付；M2 结束即已可自用（API+后台） |

## 6. 开发环境速查（Windows 本机）

```powershell
# 后端热重载
air
# 管理端
cd web/admin; pnpm dev
# 主题 CSS 监听编译
.\scripts\tailwind.exe -i web/themes/aurora/src/main.css -o web/themes/aurora/static/css/main.css --watch
# 测试与检查
go test ./... ; golangci-lint run
```
