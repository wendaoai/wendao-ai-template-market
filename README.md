# Wendao AI Template Market

文到 AI 的官方声明式模板仓库。当前提供 51 套微信公众号正文模板。

旧版正文主题客户端读取：

```text
registry/v2/stable.json
```

当前协议是 WeMD 风格的 v2-only 模板协议：

- `rendererContract` 固定为 `wechat-theme-css-v1`。
- 模板运行时只使用 `css`。
- 可视化模板保存 `css + designerVariables`。
- 手写 CSS 模板只保存 `css`。
- 不保留 v1 registry，也不兼容旧 token 顶层字段。
- `templates/wechat-article/*.json` 是源模板，`registry/v2/stable.json` 由工具生成。

## 校验

```bash
go run ./cmd/marketctl validate
go run ./cmd/marketctl build --check
go test ./...
```

## 发布

1. 修改 `templates/wechat-article/{id}.json` 并提升模板版本。
2. 提升 `catalog.json` 的 `catalogVersion` 和 `publishedAt`。
3. 运行 `validate` 和 `build`。
4. 运行 `build --check` 和测试后，提交源模板与 `registry/v2/stable.json`。

```bash
go run ./cmd/marketctl build
go run ./cmd/marketctl build --check
go test ./...
```

## Raw 地址

```text
https://api.gitcode.com/api/v5/repos/wendao-ai/wendao-ai-template-market/raw/registry/v2/stable.json?ref=master
```

模板数据许可见 `LICENSE-TEMPLATES`。第三方参考和署名见 `THIRD_PARTY_NOTICES.md`。

## 按需安装模板中心

新版工作台使用 `data/manifest.json`，按 HTML、MD、图片分组浏览，下载后安装到本地使用。该目录与既有 `registry/v2/stable.json` 独立，旧客户端继续使用原协议。

- HTML：174 个，含保留原生设计 93 个、新增 edition 50 个、结构导入 31 个。
- MD：24 个。
- 图片：56 个，含 AIWriteX 50 个及工作台原有 6 个封面模板。

`data/catalogs/` 记录分类、版本、大小和 SHA-256；`previews/templates/` 提供效果图；`assets/templates/` 保存可安装配置与来源内容。浏览目录不会自动安装完整模板。模板记录的 attribution 标明来源，原始配置保留在下载资源中。

## HTML 模板全量验收 1.3.1

174 个稳定资源 ID 保持不变；每套完成普通图文稿、长导语稿在 1440px 和 1024px 桌面窗口的全文渲染、保存刷新和导出检查，合计 696 组全文证据。26 套具体样式问题已修复，保留各套原设计、章节、图片与 GIF 装饰。

- 全部中文名称使用 4—8 字视觉命名，英文 slug 全量唯一。
- `assets/templates/<slug>.json`：可安装资源，示例 HTML 不含文章主标题。
- `previews/templates/<slug>.jpg`：从本次真实前端正文截取的完整预览。
- `documents/<slug>.html`：174 份独立 HTML，图片位于 `sample-images/`。
- [完整画廊](gallery.html) 与 [旧新命名对照](naming-map.json) 包含全部 174 套。

同步更新版本、文件大小与 SHA-256；删除被替换的旧生成路径。此前 51 套删除结果保留，MD、图片和旧版 registry/v2 协议保持各自目录。原始来源与署名保存在 original/attribution 字段中。

渲染依赖 inline-editorial-v1 / imported-structure-v2，以及 Content Workbench 本轮导语容器、颜色继承和独立 HTML 导出修复。本次为市场资源更新，桌面应用发布包另行发布；验收结论限定于本次双稿、双宽用例。

1.3.1 补齐钴蓝圆角导语的必要内边距；174 套、696 组均增加圆角上下内距复核，原配色和圆角保留。
