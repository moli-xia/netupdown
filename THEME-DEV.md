# NetUpDown 主题开发

主题是运行时加载的 ZIP 包。最小结构：

```text
my-theme/
├── theme.json
├── templates/
│   ├── layouts/base.html    # 或直接放 templates/base.html
│   ├── index.html
│   ├── list.html
│   └── detail.html
└── static/
    ├── css/main.css
    └── js/main.js
```

`theme.json` 的 `id` 只能使用小写字母、数字和连字符，长度 2–40；必须包含 `name` 与 `version`。可选 `settings` 支持 `text`、`textarea`、`color`、`switch`、`select`、`number`、`image`，保存后经 `.Theme.Config.<key>` 读取。

## 页面与数据

页面共 7 个：`index`、`list`、`detail`、`releases`、`download`、`page`、`error`。缺失的页面自动回退到内置主题 Aurora 的同名模板。渲染入口是 `base` 模板，页面模板需 `{{define "content"}}`。

公共上下文：

| 变量 | 说明 |
|---|---|
| `.Site` | `Title` `Subtitle` `Description` `Keywords` `Logo` `Favicon` `ICP`，以及按 HTML 原样输出的 `Footer` `HeadInject` `FootInject` |
| `.Theme` | `ID` `StaticBase` `Config`（主题配置 map） |
| `.PageTitle` `.Description` `.Canonical` `.Year` | 头部信息 |
| `.OGImage` | 详情页存在，绝对地址的分享图 |

页面专属（节选）：`index` 有 `.Latest` `.Featured` `.Categories` `.AppCount` `.TotalDownloads`；`list` 有 `.Apps` `.Categories` `.Query` `.Cat` `.CategoryName` `.Platform` `.Type` `.Sort` `.Total` `.Title` 与分页数据 `.Page` `.TotalPages` `.HasPrev` `.HasNext` `.PrevURL` `.NextURL` `.Pages`（元素含 `N` `URL` `Current`，链接已保留全部筛选参数）；`detail` 有 `.App` `.Release` `.Assets` `.PrimaryAsset`（预览模式另有 `.Preview`）；`releases` 有 `.App` `.Releases`；`download` 有 `.App` `.Asset` `.Source`；`page` 有 `.Doc`；`error` 有 `.Status` `.Message`。完整设计见 [主题与前台设计](docs/07-theme-frontend.md)。

## 模板函数

| 函数 | 示例 | 说明 |
|---|---|---|
| `asset` | `{{asset "css/main.css"}}` | 主题静态资源路径（建议手动追加 `?v=主题版本` 破缓存） |
| `markdown` | `{{markdown .App.Description}}` | 渲染并消毒 Markdown |
| `formatBytes` | `52428800` → `46.0 MB` | 文件大小 |
| `numFmt` | `128934` → `12.9 万` | 中文习惯的计数缩写 |
| `osName` / `archName` / `kindName` | `windows`→`Windows`、`amd64`→`x64`、`1`→`安装包` | 枚举中文名 |
| `osIcon` | `{{osIcon "windows"}}` | 平台内联 SVG 小图标 |
| `channelName` | 渠道枚举 → `stable/beta/alpha` | |
| `date` | `{{date .PublishedAt}}` | 接受 `time.Time` 或 `*time.Time`，输出 `2006-01-02` |
| `join` / `initial` | — | 字符串工具（`initial` 取首字符大写，用于无图标占位） |

## 亮暗契约

内置无闪烁脚本会在 `<html>` 上写 `data-theme="light|dark"`（缺省跟随系统）。主题 CSS 按三段式定义令牌：`:root` 放亮色全量，`@media (prefers-color-scheme: dark)` 内用 `:root:not([data-theme="light"])` 覆写，再用 `:root[data-theme="dark"]` 覆写一遍。参考 Aurora 的 `static/css/main.css` 第 1 节。

## 开发流程

```yaml
theme:
  dev: true
```

把主题目录放到 `data/themes/<id>/`，在后台启用。开发模式每次请求重新编译模板。需要演示数据时执行 `go run ./scripts/demoseed`（写入 8 个示例应用，含图标、截图、多下载源与提取码外链，可重复执行）。

打包时让 `theme.json` 位于 ZIP 根目录。安装器限制 20MB 压缩包、500 个条目、50MB 解压大小，并拒绝路径穿越及非白名单扩展名。
