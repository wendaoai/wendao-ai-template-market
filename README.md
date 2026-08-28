# Wendao AI Template Market

文到AI 的官方声明式模板仓库。首版提供 51 套微信公众号正文模板：从应用迁移的 41 套主题，以及新增的 10 套主题。

客户端只读取：

```text
registry/v1/stable.json
```

仓库不发布 ZIP、预览图、HTML、CSS、JavaScript、字体或其他运行资源。模板只包含 文到AI 已知的颜色 token 和样式枚举，实际 HTML 始终由应用内渲染器生成。

这是由 文到AI 开发人员维护的公开只读数据仓库。当前版本不使用数字签名；客户端通过 HTTPS 从固定 GitCode 地址读取 JSON，并继续执行 JSON 格式、字段、枚举和版本校验。

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
4. 运行 `build --check` 和测试后，提交源模板与 `stable.json`。

```bash
go run ./cmd/marketctl build
go run ./cmd/marketctl build --check
go test ./...
```

## Raw 地址

```text
https://api.gitcode.com/api/v5/repos/wendao-ai/wendao-ai-template-market/raw/registry/v1/stable.json
```

模板数据许可见 `LICENSE-TEMPLATES`。第三方参考和署名见 `THIRD_PARTY_NOTICES.md`。
