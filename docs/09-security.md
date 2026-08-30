# 09 · 安全设计

> 状态：设计定稿 · 2026-08-27

## 1. 资产与威胁模型

| 资产 | 主要威胁 | 对策（下文详述） |
|---|---|---|
| 管理员账号 | 口令爆破、令牌窃取、CSRF | argon2id、登录限流锁定、短效 JWT + HttpOnly 轮换刷新令牌、Origin 校验 |
| 存储密钥（S3 AK/SK） | 数据库泄露、API 回显 | AES-256-GCM 加密落库、API 脱敏、根密钥走环境变量 |
| 托管文件 | 路径穿越读写、越权直读 | secureJoin、files 目录不静态暴露、上传仅管理员 |
| 站点完整性 | 恶意主题包、XSS、SQL 注入 | zip 校验白名单、html/template 自动转义 + bluemonday、ORM 参数化 |
| 可用性 | 下载/接口刷量、大文件拖垮 | 限流、计数去重、302 分发优先、代理下载并发上限 |
| 访客隐私 | 下载日志滥用 | IP 匿名化选项、保留期清理 |

## 2. 认证与会话

- **口令哈希**：argon2id，参数 `memory=64MiB, iterations=3, parallelism=2, salt=16B, key=32B`；升级参数时登录成功后透明重哈希；
- **访问令牌**：JWT HS256，2h 有效；claims 仅 `sub uid, role, iat, exp, jti`；secret 32B 随机，首启生成存 `data/secret/jwt.key`（0600）；
- **刷新令牌**：32B 随机不透明串，SHA256 后存 `user_tokens`；HttpOnly+Secure+SameSite=Lax Cookie，Path 限定 `/api/admin/auth`；**每次使用即轮换**（旧行 revoke、新行签发）；检测到已 revoke 令牌被重放 → 撤销该用户全部会话并记审计；
- **改密码**：作废除当前外全部会话；
- **登录保护**：IP 限流 5 次/分钟；账号连续失败 10 次锁定 15 分钟（内存计数即可）；登录成功/失败均记审计；
- **CSRF**：管理 API 认证走 Bearer 头（非 Cookie），天然免疫；唯一 Cookie 端点 `auth/refresh|logout` 校验 `Origin`/`Sec-Fetch-Site: same-origin`；
- **TOTP（P2）**：pquerna/otp，secret 加密落库，恢复码一次性。

## 3. 输入与输出防护

- **SQL 注入**：一律 GORM 参数化；禁止拼接 SQL（golangci-lint 自查 + code review）；
- **XSS**：
  - 模板输出：html/template 上下文自动转义；
  - Markdown（介绍/日志）：goldmark 渲染后经 bluemonday `UGCPolicy()` 扩展（允许 img/表格/代码类名，剥 script/事件属性/javascript: URL）；
  - `custom.head/foot` 与 `site.footer` 为站长自有注入，信任但在文档与 UI 中明示风险；
  - SVG 不在图片上传白名单（防内嵌脚本）；
- **参数校验**：binding tag + 显式业务校验；slug/版本号/文件名均白名单正则；
- **SSRF 面**：v1 无服务端拉取外部 URL 的功能（外链仅 302 给浏览器）；P2 的外链巡检/GH 同步实现时校验目标为 http(s) 公网地址并设超时。

## 4. 文件上传与下载安全

- 上传入口全部在管理 API 之后（仅管理员）；
- 图片：扩展名白名单 + 魔数嗅探（`http.DetectContentType`）双校验；重命名为 xid，杜绝用户可控路径；
- 软件包：扩展名黑名单不设（分发站本就传 exe），但**只进存储驱动、绝不落 web 可执行路径**；大小上限 `upload.max_size_mb`（默认 4096）；
- 下载响应：`Content-Disposition: attachment`、`X-Content-Type-Options: nosniff`、`Content-Type: application/octet-stream`（不信任存储侧 MIME）；
- 路径安全：`secureJoin` 统一实现并全覆盖测试（06 册 §8）；zip 解压逐条目防 zip-slip；
- 完整性：托管上传强制服务端计算 SHA256 并展示；外链源建议站长手填 SHA256。

## 5. HTTP 安全头与传输

中间件统一输出：

| 头 | 值 |
|---|---|
| `X-Content-Type-Options` | `nosniff` |
| `X-Frame-Options` | `SAMEORIGIN`（管理端 `DENY`） |
| `Referrer-Policy` | `strict-origin-when-cross-origin` |
| `Permissions-Policy` | 关闭 camera/mic/geolocation |
| `Strict-Transport-Security` | 由反代在 HTTPS 层加（10 册） |
| CSP | v1 不上全站 CSP（自定义注入与主题会碎）；管理端路径可上 `default-src 'self'` 级别策略，列入 P1 验证项 |

TLS 终止在 Caddy/Nginx；程序侧 `behind_proxy=true` 时才信任 `X-Forwarded-For`（取最右非私网地址），否则用 RemoteAddr —— 防伪造 IP 绕过限流。

## 6. 限流汇总

| 面 | 策略 | 默认 |
|---|---|---|
| `/api/v1/*` | token bucket / IP | 60 req/min，burst 20 |
| `/api/admin/auth/login` | 独立桶 / IP | 5 req/min |
| `/d/*` | 计数去重（非拒绝） | 10 分钟窗口 |
| WebDAV/代理下载 | 全局并发信号量 | 32 路 |
| 管理 API | 已认证不限流（单管理员） | — |

实现：`x/time/rate` + IP→limiter 的 LRU（上限 1 万条，防内存膨胀）。

## 7. 秘密管理

- 根密钥 `NETUPDOWN_MASTER_KEY`（32B base64）：加密 storages.config、totp_secret（HKDF 派生子密钥，见 06 册 §7）；
- JWT secret、master key 落盘文件均 0600 且在 `data/secret/`（备份需包含，泄露需轮换）；
- 日志与审计 detail 字段脱敏：绝不打印口令、AK/SK、令牌；
- 配置文件不含明文云密钥（密钥只进后台表单 → 加密落库）。

## 8. 审计与监控

- operation_logs 覆盖：登录成败、改密、应用/版本/文件/源/存储/主题/设置的增删改、发布/下架、会话撤销；
- 记录字段：who(uid)、when、action、target、关键 diff（脱敏）、ip；
- 异常信号（连续登录失败、刷新令牌重放）记 warn 并入审计，P2 接 Webhook 通知（TG/Bark）。

## 9. 供应链与发布

- CI 每次跑 `govulncheck ./...` 与 `go mod verify`；依赖升级走 Dependabot/Renovate PR；
- 二进制发布附 SHA256SUMS；Docker 镜像固定基础镜像 digest；
- pnpm 依赖锁定 + `pnpm audit` 于 CI。

## 10. 上线安全清单

- [ ] 管理员强口令（≥12 位），已改默认端口段访问策略（防火墙仅放 80/443）
- [ ] `NETUPDOWN_MASTER_KEY` 已设为环境变量并另行离线备份
- [ ] `behind_proxy` 与反代实际拓扑一致
- [ ] HTTPS 强制（HTTP 301）、HSTS 开启
- [ ] `data/` 目录权限 0700，属主为运行用户；数据库文件不在 web 目录
- [ ] 备份任务已配置并完成一次恢复演练（10 册）
- [ ] `/api/admin` 通过公网扫描确认未泄露调试信息（gin release 模式）
- [ ] 免责声明/隐私页已发布；`privacy.ip_mode` 已按需设置
