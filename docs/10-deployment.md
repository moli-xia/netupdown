# 10 · 部署与运维

> 状态：设计定稿 · 2026-08-27

## 1. 运行要求

| 项 | 最低 | 推荐 |
|---|---|---|
| 服务器 | 1 核 1G（SQLite） | 2 核 2G |
| 系统 | Linux amd64/arm64（也支持 Windows/macOS 自托管） | Debian 12 / Ubuntu 24.04 |
| 磁盘 | 视托管文件量；系统盘 ≥10G | 文件走对象存储则本地占用极小 |
| 域名 | netupdown.com 已解析 | 建议同时解析 `www`（301 归一） |

> **备案提示**：主机在中国大陆则域名需 ICP 备案（页脚备案号在后台"站点设置"填写）；不想备案可选香港/海外主机 + 国内下载走 R2/CDN 缓解速度。

## 2. 数据目录布局（唯一状态目录）

```text
data/
├── netupdown.db          # SQLite（WAL: 同目录 -wal/-shm）
├── secret/               # master.key / jwt.key（0600）
├── uploads/              # 公开图片（图标/截图），静态路由暴露
├── files/                # 本地存储驱动根（绝不静态暴露）
├── themes/               # 用户安装主题
├── tmp/uploads/          # 分片上传临时区
└── logs/                 # app.log 轮转
```

迁移服务器 = 拷贝二进制 + `config.yaml` + 整个 `data/`。

## 3. 配置文件（config.example.yaml 全文）

```yaml
server:
  addr: ":8080"                    # 监听地址
  base_url: "https://netupdown.com" # 绝对链接/OG/RSS 用，末尾不带 /
  behind_proxy: true               # 有反代时开启（信任 X-Forwarded-For）
  admin_path: "/admin"             # 管理端挂载路径

database:
  driver: "sqlite"                 # sqlite | mysql | postgres
  dsn: "data/netupdown.db"         # mysql 例: user:pass@tcp(127.0.0.1:3306)/netupdown?charset=utf8mb4&parseTime=true
                                   # postgres 例: host=127.0.0.1 user=.. password=.. dbname=netupdown sslmode=disable

data_dir: "data"

log:
  level: "info"                    # debug | info | warn | error
  file: "data/logs/app.log"        # 置空则仅 stdout
  max_size_mb: 50
  max_backups: 5

upload:
  max_size_mb: 4096                # 单软件包上限
  chunk_size_mb: 5
  image_max_size_mb: 10

ratelimit:
  public_per_min: 60
  login_per_min: 5

theme:
  dev: false                       # 主题开发热重载

timezone: "Asia/Shanghai"          # 展示时区
```

环境变量覆盖：`NETUPDOWN_` 前缀、`__` 表层级，如 `NETUPDOWN_SERVER__ADDR=":9000"`；秘密类只走环境变量：`NETUPDOWN_MASTER_KEY`。

## 4. 首次初始化

```bash
./netupdown admin create --username kondor
```

交互输入口令（或 `--password-stdin`）。P1 版本提供 Web 安装向导：数据库为空时访问任意页面重定向 `/install`（设站点名 + 建管理员，一次性）。

## 5. systemd 部署（裸机）

`/etc/systemd/system/netupdown.service`：

```ini
[Unit]
Description=NetUpDown
After=network-online.target
Wants=network-online.target

[Service]
User=netupdown
Group=netupdown
WorkingDirectory=/opt/netupdown
ExecStart=/opt/netupdown/netupdown serve -c /opt/netupdown/config.yaml
Environment=NETUPDOWN_MASTER_KEY=<base64-32bytes>
Restart=on-failure
RestartSec=3
LimitNOFILE=65536
# 加固
NoNewPrivileges=true
ProtectSystem=strict
ReadWritePaths=/opt/netupdown/data
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

```bash
sudo useradd -r -s /usr/sbin/nologin netupdown
sudo systemctl daemon-reload && sudo systemctl enable --now netupdown
```

## 6. Docker 部署

### Dockerfile（多阶段）

```dockerfile
# 1) 管理端
FROM node:22-alpine AS admin
WORKDIR /src/web/admin
COPY web/admin/package.json web/admin/pnpm-lock.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY web/admin .
RUN pnpm build

# 2) Go
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=admin /src/web/admin/dist internal/assets/admin
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/netupdown ./cmd/netupdown

# 3) 运行
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 1000 netupdown
USER netupdown
WORKDIR /app
COPY --from=build /out/netupdown .
VOLUME ["/app/data"]
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1
ENTRYPOINT ["./netupdown", "serve"]
```

### docker-compose.yml

```yaml
services:
  netupdown:
    image: ghcr.io/yourname/netupdown:latest   # 或 build: .
    container_name: netupdown
    restart: unless-stopped
    ports:
      - "127.0.0.1:8080:8080"    # 仅本机，公网走反代
    environment:
      - NETUPDOWN_MASTER_KEY=${NETUPDOWN_MASTER_KEY}
    volumes:
      - ./data:/app/data
      - ./config.yaml:/app/config.yaml:ro
```

初始化：`docker compose exec netupdown ./netupdown admin create --username kondor`。

## 7. 反向代理与 HTTPS

### Caddy（推荐，自动 HTTPS）

```caddyfile
netupdown.com {
    encode zstd gzip
    reverse_proxy 127.0.0.1:8080
    request_body {
        max_size 5GB      # 覆盖分片上传（单片小，此值宽松即可）
    }
    @static path /themes/* /uploads/*
    header @static Cache-Control "public, max-age=604800"
    header Strict-Transport-Security "max-age=31536000"
}
www.netupdown.com {
    redir https://netupdown.com{uri} permanent
}
```

### Nginx 要点

```nginx
server {
    listen 443 ssl http2;
    server_name netupdown.com;
    client_max_body_size 64m;          # 分片 5MiB + 余量；整文件直传接口需调大
    proxy_read_timeout 300s;
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

> 大文件经程序直发下载时，确保反代 buffering 不落盘：`proxy_buffering off;` 作用于 `/d/`。

## 8. CDN / 加速建议

- 对象存储选 **Cloudflare R2 + 自定义域**（`public_base_url=https://dl.netupdown.com`）：下载 302 到 CDN 边缘，零出口费；
- 站点本体也可套 Cloudflare（注意 `/api/admin` 不缓存、`behind_proxy` 下正确解析 CF 头）；
- 国内主力访客且已备案：静态与下载走国内 CDN + OSS/COS。

## 9. 备份与恢复

`scripts/backup.sh`（cron 每日 04:00）：

```bash
#!/usr/bin/env bash
set -euo pipefail
SRC=/opt/netupdown; DST=/backup/netupdown; TS=$(date +%Y%m%d-%H%M)
mkdir -p "$DST"
# SQLite 在线一致备份
sqlite3 "$SRC/data/netupdown.db" ".backup '$DST/db-$TS.db'"
# 文件与配置（uploads/themes/secret/config；本地托管文件 files/ 视体量决定是否入每日备份）
tar -czf "$DST/data-$TS.tar.gz" -C "$SRC" config.yaml data/secret data/uploads data/themes
find "$DST" -mtime +7 -delete
# 建议再 rclone 同步到异地对象存储，实现 3-2-1
```

恢复：还原 db 文件与 tar 内容到位 → 启动 → 后台核对存储连通性。**上线前必须演练一次恢复**。

## 10. 升级

1. 备份（§9）；2. 替换二进制（或 `docker compose pull && up -d`）；3. 启动时 AutoMigrate 自动升 schema（只增不破坏）；4. 验证 `/healthz` 与关键页面。破坏性迁移会在 Release Notes 标注并要求先升到中间版本。

## 11. 监控与日志

- `/healthz` 接 UptimeRobot / 自建探针；
- 日志：`data/logs/app.log`（lumberjack 轮转）；`journalctl -u netupdown` / `docker logs`；
- 关注指标（日志可查）：慢请求 warn、计数器丢弃 warn、主题回退 error、存储探针失败；
- 访客分析：后台注入 Umami/Plausible 脚本（`custom.head`）。

## 12. 常见问题

| 症状 | 排查 |
|---|---|
| 上传大文件 413 | 反代 `client_max_body_size` / `request_body max_size` |
| 下载不能断点续传 | 经反代时确认未启用改写 `Range` 的缓存层；程序直发路径原生支持 |
| SQLite `database is locked` | 确认 WAL 生效、单实例运行；NFS 上勿放 SQLite |
| 换机后无法解密存储配置 | `data/secret/` 未随迁或 MASTER_KEY 未带走 |
| Cloudflare 后台看到全是 CF IP | `behind_proxy: true` 并确认取最右可信 XFF |
