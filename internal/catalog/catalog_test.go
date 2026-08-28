package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildAndVerify(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root)
	if _, err := WriteBuild(root, false); err != nil {
		t.Fatalf("build: %v", err)
	}
	privatePath := filepath.Join(root, ".local-keys", "test-key.key")
	publicPath := filepath.Join(root, "keys", "test-key.pub")
	if err := GenerateKeyPair(privatePath, publicPath); err != nil {
		t.Fatalf("keygen: %v", err)
	}
	if _, err := Sign(root, "test-key", privatePath); err != nil {
		t.Fatalf("sign: %v", err)
	}
	result, signature, err := Verify(root)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if result.CatalogVersion != "1.0.0" || len(result.Templates) != 1 {
		t.Fatalf("unexpected catalog: %#v", result)
	}
	if signature.KeyID != "test-key" {
		t.Fatalf("unexpected key: %s", signature.KeyID)
	}
	if _, err := WriteBuild(root, true); err != nil {
		t.Fatalf("build check: %v", err)
	}
}

func TestVerifyRejectsTamperedCatalog(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root)
	if _, err := WriteBuild(root, false); err != nil {
		t.Fatalf("build: %v", err)
	}
	privatePath := filepath.Join(root, ".local-keys", "test-key.key")
	if err := GenerateKeyPair(privatePath, filepath.Join(root, "keys", "test-key.pub")); err != nil {
		t.Fatalf("keygen: %v", err)
	}
	if _, err := Sign(root, "test-key", privatePath); err != nil {
		t.Fatalf("sign: %v", err)
	}
	path := filepath.Join(root, "registry", "v1", "stable.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), "默认清爽", "已被篡改", 1))
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Verify(root); err == nil {
		t.Fatal("expected tampered catalog to fail")
	}
}

func TestDecodeStrictRejectsDuplicateAndUnknownFields(t *testing.T) {
	var source SourceCatalog
	if err := DecodeStrict([]byte(`{"schemaVersion":1,"schemaVersion":1}`), &source); err == nil {
		t.Fatal("expected duplicate field to fail")
	}
	if err := DecodeStrict([]byte(`{"schemaVersion":1,"unknown":true}`), &source); err == nil {
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
		SchemaVersion:    1,
		CatalogID:        CatalogID,
		CatalogVersion:   "1.0.0",
		Channel:          ChannelStable,
		RendererContract: RendererContract,
		MinAppVersion:    "0.0.64",
		PublishedAt:      "2026-08-28T12:00:00Z",
	}
	template := SourceTemplate{
		SchemaVersion:  1,
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
		Category: "recommended", Tags: []string{"通用"}, Accent: "#96743d",
		Background: "#fffefa", Surface: "#fffaf0", Text: "#25221d",
		Muted: "#756b5c", Border: "#e6dcc9", HeadingStyle: "plain",
		BlockStyle: "quiet", Density: "standard", CodeTheme: "github",
		Selectable: true, Recommended: true, RiskLevel: "stable",
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
