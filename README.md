# Wendao AI Template Market

文到AI 的官方声明式模板仓库。首版提供 51 套微信公众号正文模板：从应用迁移的 41 套主题，以及新增的 10 套主题。

客户端只读取：

```text
registry/v1/stable.json
registry/v1/stable.json.sig
```

仓库不发布 ZIP、预览图、HTML、CSS、JavaScript、字体或其他运行资源。模板只包含 文到AI 已知的颜色 token 和样式枚举，实际 HTML 始终由应用内渲染器生成。

## 校验

```bash
go run ./cmd/marketctl validate
go run ./cmd/marketctl build --check
go run ./cmd/marketctl verify
go test ./...
```

## 发布

1. 修改 `templates/wechat-article/{id}.json` 并提升模板版本。
2. 提升 `catalog.json` 的 `catalogVersion` 和 `publishedAt`。
3. 运行 `validate` 和 `build`。
4. 使用未提交的 Ed25519 私钥运行 `sign`。
5. 运行 `verify` 后提交源模板、`stable.json` 和 `stable.json.sig`。

```bash
go run ./cmd/marketctl build
go run ./cmd/marketctl sign --key-id cw-template-2026-01 --private-key .local-keys/cw-template-2026-01.key
go run ./cmd/marketctl verify
```

私钥目录 `.local-keys/` 已被忽略，绝不能提交。

## Raw 地址

```text
https://api.gitcode.com/api/v5/repos/wendao-ai/wendao-ai-template-market/raw/registry/v1/stable.json
https://api.gitcode.com/api/v5/repos/wendao-ai/wendao-ai-template-market/raw/registry/v1/stable.json.sig
```

模板数据许可见 `LICENSE-TEMPLATES`。第三方参考和署名见 `THIRD_PARTY_NOTICES.md`。

