# 05 · API 接口设计

> 状态：设计定稿 · 2026-08-27

## 1. 通用规范

### 1.1 响应信封（envelope）

所有 `/api/*` 接口统一返回：

```json
{ "code": 0, "message": "ok", "data": { } }
```

- `code=0` 成功；非 0 为业务错误码（下表）；
- HTTP 状态码同时保持语义（4xx/5xx），客户端可任选一种判断方式；
- 列表数据统一结构：`data: { "list": [...], "page": 1, "page_size": 20, "total": 137 }`。

### 1.2 错误码表

| HTTP | code | 含义 |
|---|---|---|
| 500 | 10000 | 服务器内部错误（不暴露细节） |
| 400 | 10001 | 参数校验失败（`message` 给出首个字段错误） |
| 401 | 10002 | 未认证 / 令牌无效 |
| 401 | 10003 | 访问令牌过期（客户端应走刷新流程后重试） |
| 403 | 10004 | 无权限 |
| 404 | 10005 | 资源不存在 |
| 409 | 10006 | 资源冲突（slug/版本号重复等） |
| 413 | 10007 | 文件超出大小限制 |
| 422 | 10008 | 业务校验失败（如 SHA256 不匹配、状态机非法流转） |
| 429 | 10009 | 请求过于频繁 |
| 503 | 10010 | 服务暂不可用（未完成初始化等） |

### 1.3 分页与排序

- 请求参数：`page`（1 起，默认 1）、`page_size`（默认 20，上限 100）；
- 排序用语义化参数（如 `sort=latest|hot|name`），不开放任意字段排序。

### 1.4 认证

- 管理 API：`Authorization: Bearer <access_token>`（JWT，HS256，有效期 2h）；
- 刷新令牌：不透明随机串，存于 **HttpOnly Cookie**（`nud_rt`，`Path=/api/admin/auth`，`SameSite=Lax`，`Secure`），仅刷新与登出端点可见；
- 公开 API 无认证，全局限流 60 req/min/IP（可配）；登录端点独立限流 5 req/min/IP。

### 1.5 其他约定

- 时间字段一律 RFC3339 UTC 字符串（如 `2026-08-27T03:00:00Z`）；
- 富文本字段（description/changelog）API 返回**原始 Markdown**，渲染由消费方决定；SSR 页面则由服务端渲染并消毒；
- 写接口幂等性：创建接口用唯一约束兜底返回 10006；上传分片接口天然可重试。

## 2. 公开 API（/api/v1）

### 2.1 应用列表

`GET /api/v1/apps`

| 参数 | 类型 | 说明 |
|---|---|---|
| page / page_size | int | 分页 |
| category | string | 分类 slug |
| tag | string | 标签 slug |
| platform | string | `windows/macos/linux/android/ios/web` |
| type | string | `self` / `third` |
| sort | string | `latest`(默认，发布时间) / `hot`(下载量) / `name` |
| q | string | 关键词（名称/简介/标签） |

响应 `data.list[]` 元素：

```json
{
  "id": 12, "name": "MyApp", "slug": "myapp", "type": "self",
  "tagline": "一句话简介", "icon": "/uploads/2026/08/xxx.png",
  "category": {"name": "效率工具", "slug": "productivity"},
  "platforms": ["windows", "macos"],
  "latest_version": "1.4.1",
  "download_count": 3521, "is_featured": true,
  "published_at": "2026-08-20T09:00:00Z", "updated_at": "..."
}
```

### 2.2 应用详情

`GET /api/v1/apps/{slug}` → 应用完整信息 + `latest_release`（含 assets 与可用 sources）+ `tags`。

```json
{
  "id": 12, "name": "MyApp", "slug": "myapp", "type": "self",
  "tagline": "...", "description": "# Markdown 介绍…",
  "icon": "...", "screenshots": ["..."],
  "official_url": "https://myapp.example", "repo_url": "https://github.com/me/myapp",
  "developer": "Kondor", "license": "MIT",
  "platforms": ["windows","macos","linux"],
  "category": {"name":"我的作品","slug":"works"}, "tags": [{"name":"效率","slug":"efficiency"}],
  "download_count": 3521, "view_count": 12034,
  "latest_release": {
    "id": 88, "version": "1.4.1", "channel": "stable",
    "changelog": "## 新增\n- …", "published_at": "2026-08-20T09:00:00Z",
    "assets": [
      {
        "id": 301, "name": "Windows 64 位 安装版",
        "file_name": "MyApp-1.4.1-win-x64-setup.exe",
        "os": "windows", "arch": "amd64", "kind": "installer",
        "size": 52428800, "sha256": "ab12…",
        "sources": [
          {"id": 501, "name": "直链下载", "type": "managed", "default": true},
          {"id": 502, "name": "蓝奏云", "type": "external", "extract_code": "abcd"}
        ],
        "download_url": "/d/301"
      }
    ]
  }
}
```

> 注意：托管源**不**返回 object_key/storage 信息；外链源不直接返回 URL（统一走 `/d/{asset}?source={id}` 以便计数），仅提取码可见。

### 2.3 版本历史

`GET /api/v1/apps/{slug}/releases?page=&channel=` → 已发布版本分页列表（含各自 assets 概要）。

`GET /api/v1/apps/{slug}/releases/latest?channel=stable&os=&arch=` → 最新单版本（客户端简单场景直取）。

### 2.4 更新检查（核心，F-406）

`GET /api/v1/apps/{slug}/check-update`

| 参数 | 必填 | 说明 |
|---|---|---|
| version | 是 | 客户端当前版本，如 `1.2.0`（可带 `v` 前缀） |
| os / arch | 否 | 用于挑选匹配的 asset；缺省不返回 asset |
| channel | 否 | 默认 `stable` |
| version_code | 否 | 数值版本；当 version 非法或站点配置了 version_code 时优先比较 |

**判定逻辑**：

1. 取该应用指定渠道最新已发布 release；
2. 版本比较：双方可解析 semver → semver 比较；否则若都有 version_code → 数值比较；再否则字符串相等判断（不等即视为有更新并记 warn）；
3. `mandatory = release.min_required_version 存在且 当前版本 < min_required_version`；
4. asset 匹配：`os` 精确匹配（asset.os=any 亦可）后按 `arch` 精确 > `universal` > `any` 择优，取 sort 最小者。

**响应示例（有更新）**：

```json
{
  "code": 0, "message": "ok",
  "data": {
    "update_available": true,
    "mandatory": false,
    "current": "1.2.0",
    "latest": {
      "version": "1.4.1", "version_code": 10401, "channel": "stable",
      "title": "夏季更新", "changelog": "## 新增\n- …",
      "published_at": "2026-08-20T09:00:00Z",
      "notes_url": "https://netupdown.com/apps/myapp",
      "asset": {
        "os": "windows", "arch": "amd64",
        "file_name": "MyApp-1.4.1-win-x64-setup.exe",
        "size": 52428800, "sha256": "ab12…",
        "url": "https://netupdown.com/d/301"
      }
    }
  }
}
```

**无更新**：`{"update_available": false, "current": "1.4.1", "latest": {"version": "1.4.1", …}}`（仍带 latest 便于客户端展示）。

客户端接入示例（Go）：

```go
resp, _ := http.Get("https://netupdown.com/api/v1/apps/myapp/check-update?version=" +
    cur + "&os=" + runtime.GOOS + "&arch=" + runtime.GOARCH)
// 解析 data.update_available / data.mandatory / data.latest.asset.url
```

该接口按 `(slug, channel)` 内存缓存，发布事件失效；限流从公开池扣除。

### 2.5 其他公开接口

| 接口 | 说明 |
|---|---|
| `GET /api/v1/categories` | 分类列表（含各分类已发布应用数） |
| `GET /api/v1/tags/hot?limit=20` | 热门标签 |
| `GET /healthz` | `{"status":"ok","version":"1.0.0"}`，不套 envelope |

下载入口 `GET /d/{assetID}` 属页面路由（可能渲染落地页），行为定义见 06 册 §5。

应用预览入口 `GET /apps/{slug}/preview` 仅接受管理 API 签发的短时预览凭证；凭证验证后会换成 HttpOnly Cookie 并重定向到不带凭证的干净地址，因此草稿内容不会通过公开应用页暴露。

## 3. 管理 API（/api/admin）

除 auth 组外均需 Bearer 认证 + admin 角色；所有写操作记入 operation_logs。

### 3.1 认证 auth

| 方法 路径 | 说明 |
|---|---|
| `POST /api/admin/auth/login` | 入参 `{username, password, totp_code?}`；出参 `{access_token, expires_in, user}` 并 Set-Cookie 刷新令牌。失败计数触发锁定（09 册） |
| `POST /api/admin/auth/refresh` | 凭 Cookie 刷新；轮换刷新令牌，返回新 access_token。校验 `Origin`/`Sec-Fetch-Site` 防 CSRF |
| `POST /api/admin/auth/logout` | 作废当前刷新令牌，清 Cookie |
| `GET /api/admin/auth/profile` | 当前用户信息 |
| `PUT /api/admin/auth/profile` | 改昵称/头像/邮箱 |
| `PUT /api/admin/auth/password` | `{old_password, new_password}`；成功后作废除当前外全部会话 |
| `GET /api/admin/auth/sessions` / `DELETE .../sessions/{id}` / `DELETE .../sessions` | 会话列表 / 踢单个 / 注销全部 |

### 3.2 应用 apps

| 方法 路径 | 说明 |
|---|---|
| `GET /api/admin/apps` | 列表：支持 `status/type/category/q` 筛选（含草稿与已下架） |
| `POST /api/admin/apps` | 创建（默认草稿）。入参为 3.5 节 App 完整字段子集 |
| `GET /api/admin/apps/{id}` | 详情（编辑回显） |
| `PUT /api/admin/apps/{id}` | 更新 |
| `DELETE /api/admin/apps/{id}` | 软删除（回收站，P2 提供恢复/彻底删除） |
| `POST /api/admin/apps/{id}/preview` | 签发 10 分钟预览地址（草稿和已发布应用均可） |
| `POST /api/admin/apps/{id}/publish` · `/unpublish` | 状态流转：草稿→发布、发布→下架 |
| `PUT /api/admin/apps/{id}/feature` | `{is_featured, is_pinned, sort_weight}` |

### 3.3 版本 releases 与文件 assets

| 方法 路径 | 说明 |
|---|---|
| `GET /api/admin/apps/{id}/releases` | 该应用版本列表（含草稿） |
| `POST /api/admin/apps/{id}/releases` | 创建版本 `{version, version_code?, channel, title?, changelog, min_required_version?}` |
| `PUT /api/admin/releases/{id}` | 更新版本信息 |
| `DELETE /api/admin/releases/{id}` | 删除（软删；若含托管文件询问是否连带删对象） |
| `POST /api/admin/releases/{id}/publish` | 发布：校验至少一个 asset 且各 asset 至少一个启用源 → 置 status=1、published_at、刷新 app.latest_release_id、失效缓存 |
| `POST /api/admin/releases/{id}/assets` | 添加文件 `{name, file_name, os, arch, kind, size?, sha256?, sort}` |
| `PUT /api/admin/assets/{id}` · `DELETE` | 编辑 / 删除文件 |
| `POST /api/admin/assets/{id}/sources` | 添加下载源。托管：`{source_type:"managed", storage_id, object_key, name}`（object_key 来自上传结果）；外链：`{source_type:"external", external_url, extract_code?, name}` |
| `PUT /api/admin/sources/{id}` · `DELETE` | 编辑 / 删除源（托管源删除时可选连删对象） |
| `PUT /api/admin/assets/{id}/sources/order` | `{ids:[...]}` 重排 priority |

### 3.4 上传 uploads

小文件（图片）：

| 方法 路径 | 说明 |
|---|---|
| `POST /api/admin/uploads/image` | multipart 字段 `file`；白名单 png/jpg/jpeg/webp/gif/ico，≤ `upload.image_max_size_mb`；存 `data/uploads/YYYY/MM/{xid}.{ext}`；返回 `{url}` |

软件包（走分片协议，F-305；MVP 期可先提供 `POST /api/admin/uploads/file` 整文件 multipart 直传，≤500MB）：

| 方法 路径 | 说明 |
|---|---|
| `POST /api/admin/uploads/init` | `{file_name, size, sha256, chunk_size?}`；若 sha256 已有托管对象 → `{exists:true, object:{storage_id, object_key, size}}`（秒传）；否则 `{upload_id, chunk_size, uploaded_chunks:[…]}`（断点续传时返回已有分片号） |
| `PUT /api/admin/uploads/{upload_id}/chunks/{index}` | 原始字节流 body；可并发、可重传 |
| `POST /api/admin/uploads/{upload_id}/complete` | `{storage_id?, key_hint?}`（缺省默认存储；key_hint 如 `myapp/1.4.1`）；服务端顺序合并 + 计算 SHA256 校验 → 不匹配 422/10008；成功 `{storage_id, object_key, size, sha256}` |
| `DELETE /api/admin/uploads/{upload_id}` | 放弃上传，清理临时分片 |

协议细节（分片大小、临时目录、对象键规则）见 06 册 §4。

### 3.5 分类 / 标签 / 页面 / 友链

标准 CRUD，从略：`/api/admin/categories`、`/api/admin/tags`、`/api/admin/pages`、`/api/admin/links`（均支持列表+增删改；分类/友链支持 `PUT .../order` 排序）。

### 3.6 存储 storages

| 方法 路径 | 说明 |
|---|---|
| `GET /api/admin/storages` | 列表；config 返回**脱敏**版本（secret 字段以 `******` 占位） |
| `POST /api/admin/storages` | `{name, driver, config, remark}`；config 结构因 driver 而异（见 06 册 §3 配置 schema） |
| `PUT /api/admin/storages/{id}` | 更新；secret 字段传 `******` 表示保持不变 |
| `DELETE /api/admin/storages/{id}` | 有引用（download_sources）时拒绝，返回 10006 |
| `POST /api/admin/storages/{id}/test` | 连通性测试：写→读→删探针对象，返回 `{ok, latency_ms, error?}` |
| `PUT /api/admin/storages/{id}/default` | 设为默认 |

### 3.7 主题 themes

| 方法 路径 | 说明 |
|---|---|
| `GET /api/admin/themes` | 已装主题列表：`{id, name, version, author, preview_url, active, builtin, settings_schema}` |
| `POST /api/admin/themes/upload` | multipart zip；校验见 07 册 §6（zip-slip、大小、theme.json） |
| `POST /api/admin/themes/{id}/activate` | 切换主题（重建模板缓存） |
| `PUT /api/admin/themes/{id}/config` | 保存主题配置 `{key: value}`（按 schema 校验），写 `settings: theme.cfg.{id}` |
| `DELETE /api/admin/themes/{id}` | 删除（内置主题与当前激活主题不可删） |

### 3.8 设置 settings

| 方法 路径 | 说明 |
|---|---|
| `GET /api/admin/settings?group=site` | 按组读取（`site/seo/download/comment/user/custom/privacy`），返回 `{key: value}` |
| `PUT /api/admin/settings` | 批量写 `{key: value}`；服务端按注册表校验键合法性与值类型，写后失效缓存 |

### 3.9 统计 stats 与日志 logs

| 方法 路径 | 说明 |
|---|---|
| `GET /api/admin/stats/overview` | `{app_count, release_count, total_downloads, today_downloads, today_views, storage_used_bytes}` |
| `GET /api/admin/stats/trend?days=30` | `[{date, downloads, views}]`（stat_daily） |
| `GET /api/admin/stats/top-apps?days=7&limit=10` | Top 下载应用 |
| `GET /api/admin/logs/operations?page=&action=&q=` | 审计日志 |
| `GET /api/admin/logs/downloads?page=&app_id=` | 下载明细 |

### 3.10 用户 users（P2，多用户开启后）与评论 comments（P2）

预留：`/api/admin/users` CRUD 与禁用；`/api/admin/comments` 列表/审核/删除/回复。v1 不实现，仅保留路由命名空间。

## 4. 内部路由防冲突规则

前台 slug 路由（`/apps/:slug` 等）与保留前缀（`/api` `/admin` `/d` `/themes` `/uploads` `/pages` `/feed.xml` 等）冲突由注册顺序保证；应用/页面 slug 校验时拒绝保留字列表：`api, admin, d, themes, uploads, static, feed, sitemap, robots, healthz, install, login`。
