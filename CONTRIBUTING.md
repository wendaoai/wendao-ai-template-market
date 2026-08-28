# 模板贡献规则

## 新增模板

- 文件名和 `id` 一致，只使用小写字母、数字和中划线。
- `id` 一经发布永久稳定，不能复用或重命名。
- 只使用 Schema 允许的声明式字段，不添加 HTML、CSS、脚本、URL 或文件路径。
- 使用 `samples/wechat-article.md` 在 文到AI 中检查长文、列表、引用、表格和代码块。
- 新模板从 `1.0.0` 开始。

## 更新模板

- 任意字段变化都必须提升模板 SemVer。
- 同一 `id + version` 的内容不可变。
- 不直接删除已发布模板；设置 `status` 为 `deprecated`，必要时填写 `replacementId`。
- 每次目录变化都提升 `catalogVersion`，并更新 `publishedAt`。

## 提交前

```bash
go fmt ./...
go test ./...
go run ./cmd/marketctl validate
go run ./cmd/marketctl build --check
```

`stable.json` 必须由 `marketctl build` 生成，不要直接编辑。源模板无需填写 `contentHash`，该字段由构建工具自动生成。
