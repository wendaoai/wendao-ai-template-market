package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildAndCheck(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root)
	if _, err := WriteBuild(root, false); err != nil {
		t.Fatalf("build: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "registry", "v2", "stable.json"))
	if err != nil {
		t.Fatal(err)
	}
	var result Catalog
	if err := DecodeStrict(data, &result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if err := ValidateCatalog(result); err != nil {
		t.Fatalf("validate generated catalog: %v", err)
	}
	if _, err := WriteBuild(root, true); err != nil {
		t.Fatalf("build check: %v", err)
	}
	if result.CatalogVersion != "1.0.0" || len(result.Templates) != 1 {
		t.Fatalf("unexpected catalog: %#v", result)
	}
}

func TestValidateCatalogRejectsTamperedContent(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root)
	if _, err := WriteBuild(root, false); err != nil {
		t.Fatalf("build: %v", err)
	}
	path := filepath.Join(root, "registry", "v2", "stable.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), "默认清爽", "已被篡改", 1))
	var result Catalog
	if err := DecodeStrict(data, &result); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCatalog(result); err == nil {
		t.Fatal("expected changed template content to fail contentHash validation")
	}
}

func TestDecodeStrictRejectsDuplicateAndUnknownFields(t *testing.T) {
	var source SourceCatalog
	if err := DecodeStrict([]byte(`{"schemaVersion":2,"schemaVersion":2}`), &source); err == nil {
		t.Fatal("expected duplicate field to fail")
	}
	if err := DecodeStrict([]byte(`{"schemaVersion":2,"unknown":true}`), &source); err == nil {
		t.Fatal("expected unknown field to fail")
	}
}

func TestValidateTemplateRejectsSelectableDeprecatedTemplate(t *testing.T) {
	item := validTemplateRecord()
	item.Status = "deprecated"
	item.Selectable = true
	if err := ValidateTemplate(item, false); err == nil || !strings.Contains(err.Error(), "不能继续设为可选择") {
		t.Fatalf("expected selectable deprecated template to fail, got %v", err)
	}
	item.Selectable = false
	if err := ValidateTemplate(item, false); err != nil {
		t.Fatalf("non-selectable deprecated template was rejected: %v", err)
	}
}

func writeFixture(t *testing.T, root string) {
	t.Helper()
	source := SourceCatalog{
		SchemaVersion:    SchemaVersion,
		CatalogID:        CatalogID,
		CatalogVersion:   "1.0.0",
		Channel:          ChannelStable,
		RendererContract: RendererContract,
		MinAppVersion:    "0.0.64",
		PublishedAt:      "2026-08-28T12:00:00Z",
	}
	template := SourceTemplate{
		SchemaVersion:  SchemaVersion,
		TemplateRecord: validTemplateRecord(),
	}
	writeJSON(t, filepath.Join(root, "catalog.json"), source)
	writeJSON(t, filepath.Join(root, "templates", "wechat-article", template.ID+".json"), template)
}

func validTemplateRecord() TemplateRecord {
	return TemplateRecord{
		ID: "cw-clean-reading", Version: "1.0.0", ReleasedAt: "2026-08-28T12:00:00Z",
		ChangeSummary: "首次发布。", Status: "active", Name: "默认清爽",
		Description: "留白充足，适合大多数公众号正文。", Source: "cw",
		Category: "recommended", Tags: []string{"通用"},
		Selectable: true, Recommended: true, RiskLevel: "stable",
		EditorMode: "visual", CSS: "#wemd { color: #25221d; }",
		DesignerVariables: json.RawMessage(`{"primaryColor":"#96743d","paragraphColor":"#25221d"}`),
	}
}

func writeJSON(t *testing.T, path string, value any) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}
