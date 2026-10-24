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

## HTML 模板多段开场更正 1.3.3

[完整画廊](gallery.html) · [逐套更正与截图摘要](audit/html-template-intro-correction-2026-10-24.json) · [原174套去留记录](audit/html-template-dedup-2026-10-24.json)。

原174套中，夜粉叠卡与素粉叠卡确认重复，保留素粉叠卡；最终173套。此前51项退休映射与新增1项去重结果均保留，累计52项。中文名称4—8字，名称及英文语义slug全量唯一。

- [173套命名表](naming-map.json)。
- `documents/<slug>.html`：173份独立HTML，保留章节、图片、装饰和原设计。
- `assets/templates/<slug>.json`：173份可安装资源及来源署名。
- `previews/templates/<slug>.jpg`：173张真实前端完整正文预览。
- `data/retired-html-templates.json`：52项退休替换映射。
- `sample-images/`：当前为1张离线共享示例照片（1.3.4）。

1.3.3修复连续多段开场只有第一段进入导语容器的问题，并定点修复5套模板的内边距或文字对比度。新增三段开场普通稿、长导语压力稿，逐套在1440px与1024px真实前端渲染，保存692份原尺寸全文截图；校验预览、套用、保存刷新及独立HTML导出。人工复核用户问题与异常标记，不能据此声称人工逐张通读全部长图，也不能推断任意文章均无问题。

目录、资源、文件摘要、预览摘要与画廊缓存标识同步更新为1.3.3。内嵌GIF与原设计保留；MD24套、图片56套维持原有独立目录。

本次市场资源更新不包含桌面安装包。多段导语排版器与旧下载名称显示修复需要含本次代码的客户端；历史已保存HTML快照需在更新后的客户端重新套用模板。未执行Wails原生窗口验收，无章节文章仅有Go回归测试，未另建前端截图矩阵。

## HTML模板共享示例图 1.3.4

173套HTML模板现在共用1张1536×1024通用阅读静物图，用于展示排版效果。JPEG保留原分辨率，180112字节；相较生成PNG减少91.71%。保留每套原有图框、圆角、图文位置、章节和GIF，真实文章继续使用文章自己的图片。旧5张示例照片已移除。

普通稿、长导语稿×1440px、1024px共692项真实前端复验及原尺寸全文截图已完成，重新生成173份独立HTML和173张完整预览，并同步摘要与画廊缓存标识。自动检查全文内容、图片数量/加载/比例、保存刷新与导出，人工查看成图与异常标记截图；不能据此声称人工逐张通读全部长图。173个唯一ID和52项退休映射、15份GIF保持不变；MD与图片提示词资源未修改。

验收摘要见 `audit/html-template-shared-image-2026-10-24.json`。本次为市场资源更新，不包含桌面安装包。
