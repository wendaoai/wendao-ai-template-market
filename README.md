# Wendao AI Template Market

文到 AI 的官方声明式模板仓库。当前提供 51 套微信公众号正文模板。

客户端只读取：

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
