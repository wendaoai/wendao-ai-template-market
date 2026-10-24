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

## HTML 市场整合 1.2.2

本次以174套为完整目录，直接移除旧目录的51套配置及预览文件，新增50套独立设计，并升级31套结构导入资源；不保留可下载安装的下架占位项。MD、图片及旧版 registry/v2/stable.json 保持各自原有目录。

新增设计需要 inline-editorial-v1 渲染器，结构导入需要 imported-structure-v2（含 headerLead、textStyle、group/image）支持。31套导入明确标注为部分适配。模板市场JPEG与正文主标题分离目前尚未统一，个别文字配色与图片/表格宽度问题仍待修正，本次目录整合不代表完成这些样式修复。原始来源与署名保存在各资源 original/attribution 字段中。
