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

- HTML：173 个，93 native + 50 edition + 30 imported。
- MD：24 个。
- 图片：56 个，含 AIWriteX 50 个及工作台原有 6 个封面模板。

`data/catalogs/` 记录分类、版本、大小和 SHA-256；`previews/templates/` 提供效果图；`assets/templates/` 保存可安装配置与来源内容。浏览目录不会自动安装完整模板。模板记录的 attribution 标明来源，原始配置保留在下载资源中。

## HTML 模板验收与去重 1.3.2

当前保留173套。对1.3.1的174套逐套读取实际计算样式、几何和装饰，完成15051对比较，确认1组仅明暗配色不同的同版式：保留素粉叠卡（aiwritex-html-51277cdd70b7），移除夜粉叠卡（aiwritex-html-18137d753aff）。重复JSON、预览及独立HTML已删除；累计52套退休，此前51项删除映射完整保留。

- 中文名称全部使用4—8字视觉命名，英文slug全量唯一。
- `assets/templates/<slug>.json`：173份可安装资源，示例正文不含文章主标题。
- `previews/templates/<slug>.jpg`：173张真实前端完整正文预览。
- `documents/<slug>.html`：173份独立HTML，离线图片位于`sample-images/`。
- [完整画廊](gallery.html)与[最终命名表](naming-map.json)包含全部173套。
- [原174套去留及命名记录](audit/html-template-dedup-2026-10-24.json)，[退休映射](data/retired-html-templates.json)。

173套保留资源的JSON、HTML与JPEG均与上一轮验收版本字节一致。上一轮对原174套分别使用普通稿、长导语稿在1440px、1024px桌面窗口完成真实排版、前端、保存刷新与导出检查，共696组原尺寸全文证据；本轮逐份复核截图SHA-256及资源字节，沿用这些证据。26套样式修复、15套内嵌GIF及15种动画资产保留。近似构图但组件或装饰不同的设计变体保留，详见去重记录。

目录及资源版本更新为1.3.2，大小与SHA-256校验一致。稳定ID保留，命名、路径及引用统一。源生成入口排除退休项；MD、图片和旧版registry/v2协议保持独立。原始来源与署名保存在original/attribution字段。

本次为市场资源更新。新退休项对旧下载列表的屏蔽已在应用源码与临时Go/SQLite环境验证，随桌面应用后续版本生效；本次未构建或发布桌面安装包。
