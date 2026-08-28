package catalog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	idPattern     = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	semverPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z.-]+)?$`)
	colorPattern  = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	hashPattern   = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

var allowedValues = map[string]map[string]bool{
	"channel":      {"stable": true},
	"status":       {"active": true, "deprecated": true},
	"source":       {"wemd": true, "md": true, "neurapress": true, "cw": true},
	"category":     {"recommended": true, "knowledge": true, "opinion": true, "tech": true, "lifestyle": true, "visual": true, "code": true},
	"headingStyle": {"plain": true, "bar": true, "card": true, "marker": true, "center": true},
	"blockStyle":   {"quiet": true, "line": true, "boxed": true, "note": true, "bold": true},
	"density":      {"compact": true, "standard": true, "loose": true},
	"codeTheme":    {"github": true, "github-dark": true, "dracula": true, "monokai": true, "solarized-light": true, "nord": true},
	"riskLevel":    {"stable": true, "experimental": true},
}

func ValidateSource(source SourceCatalog, templates []SourceTemplate) error {
	if source.SchemaVersion != SchemaVersion {
		return fmt.Errorf("catalog.json schemaVersion 必须是 %d", SchemaVersion)
	}
	if source.CatalogID != CatalogID {
		return fmt.Errorf("catalogId 必须是 %s", CatalogID)
	}
	if !allowedValues["channel"][source.Channel] {
		return fmt.Errorf("channel 不受支持：%s", source.Channel)
	}
	if source.RendererContract != RendererContract {
		return fmt.Errorf("rendererContract 必须是 %s", RendererContract)
	}
	if !semverPattern.MatchString(source.CatalogVersion) {
		return errors.New("catalogVersion 必须是 SemVer")
	}
	if !semverPattern.MatchString(source.MinAppVersion) {
		return errors.New("minAppVersion 必须是 SemVer")
	}
	if _, err := time.Parse(time.RFC3339, source.PublishedAt); err != nil {
		return fmt.Errorf("publishedAt 必须是 RFC3339：%w", err)
	}
	if len(templates) == 0 {
		return errors.New("至少需要一个模板")
	}

	ids := make(map[string]bool, len(templates))
	for index := range templates {
		if templates[index].SchemaVersion != SchemaVersion {
			return fmt.Errorf("模板 %s schemaVersion 必须是 %d", templates[index].ID, SchemaVersion)
		}
		if err := ValidateTemplate(templates[index].TemplateRecord, false); err != nil {
			return err
		}
		if ids[templates[index].ID] {
			return fmt.Errorf("模板 ID 重复：%s", templates[index].ID)
		}
		ids[templates[index].ID] = true
	}
	for _, item := range templates {
		if item.ReplacementID != "" && !ids[item.ReplacementID] {
			return fmt.Errorf("模板 %s 的 replacementId 不存在：%s", item.ID, item.ReplacementID)
		}
		if item.ReplacementID == item.ID {
			return fmt.Errorf("模板 %s 不能替代自己", item.ID)
		}
	}
	return nil
}

func ValidateCatalog(value Catalog) error {
	source := SourceCatalog{
		SchemaVersion:    value.SchemaVersion,
		CatalogID:        value.CatalogID,
		CatalogVersion:   value.CatalogVersion,
		Channel:          value.Channel,
		RendererContract: value.RendererContract,
		MinAppVersion:    value.MinAppVersion,
		PublishedAt:      value.GeneratedAt,
	}
	templates := make([]SourceTemplate, len(value.Templates))
	for index := range value.Templates {
		record := value.Templates[index]
		record.ContentHash = ""
		templates[index] = SourceTemplate{SchemaVersion: SchemaVersion, TemplateRecord: record}
	}
	if err := ValidateSource(source, templates); err != nil {
		return err
	}
	if !sort.SliceIsSorted(value.Templates, func(i, j int) bool { return value.Templates[i].ID < value.Templates[j].ID }) {
		return errors.New("stable.json 模板必须按 ID 排序")
	}
	for _, item := range value.Templates {
		if err := ValidateTemplate(item, true); err != nil {
			return err
		}
		expected, err := ContentHash(item)
		if err != nil {
			return err
		}
		if item.ContentHash != expected {
			return fmt.Errorf("模板 %s contentHash 不匹配", item.ID)
		}
	}
	return nil
}

func ValidateTemplate(item TemplateRecord, requireHash bool) error {
	if !idPattern.MatchString(item.ID) || len(item.ID) > 80 {
		return fmt.Errorf("模板 ID 不合法：%s", item.ID)
	}
	if !semverPattern.MatchString(item.Version) {
		return fmt.Errorf("模板 %s version 必须是 SemVer", item.ID)
	}
	if _, err := time.Parse(time.RFC3339, item.ReleasedAt); err != nil {
		return fmt.Errorf("模板 %s releasedAt 必须是 RFC3339", item.ID)
	}
	if requireHash && !hashPattern.MatchString(item.ContentHash) {
		return fmt.Errorf("模板 %s contentHash 不合法", item.ID)
	}
	if !requireHash && item.ContentHash != "" {
		return fmt.Errorf("模板源文件 %s 不能手工填写 contentHash", item.ID)
	}
	if !allowedValues["status"][item.Status] {
		return fmt.Errorf("模板 %s status 不受支持", item.ID)
	}
	if item.Status == "active" && item.ReplacementID != "" {
		return fmt.Errorf("启用模板 %s 不能填写 replacementId", item.ID)
	}
	if item.Status == "deprecated" && item.Selectable {
		return fmt.Errorf("停用模板 %s 不能继续设为可选择", item.ID)
	}
	if !idOrEmpty(item.ReplacementID) {
		return fmt.Errorf("模板 %s replacementId 不合法", item.ID)
	}
	if err := validateText(item.ID, "name", item.Name, 1, 40); err != nil {
		return err
	}
	if err := validateText(item.ID, "description", item.Description, 1, 160); err != nil {
		return err
	}
	if err := validateText(item.ID, "changeSummary", item.ChangeSummary, 1, 200); err != nil {
		return err
	}
	if !allowedValues["source"][item.Source] || !allowedValues["category"][item.Category] {
		return fmt.Errorf("模板 %s source 或 category 不受支持", item.ID)
	}
	if len(item.Tags) == 0 || len(item.Tags) > 8 {
		return fmt.Errorf("模板 %s tags 数量必须为 1 到 8", item.ID)
	}
	seenTags := map[string]bool{}
	for _, tag := range item.Tags {
		if strings.TrimSpace(tag) != tag || len([]rune(tag)) == 0 || len([]rune(tag)) > 16 {
			return fmt.Errorf("模板 %s 包含不合法标签", item.ID)
		}
		if seenTags[tag] {
			return fmt.Errorf("模板 %s 包含重复标签：%s", item.ID, tag)
		}
		seenTags[tag] = true
	}
	for name, color := range map[string]string{
		"accent": item.Accent, "background": item.Background, "surface": item.Surface,
		"text": item.Text, "muted": item.Muted, "border": item.Border,
	} {
		if !colorPattern.MatchString(color) {
			return fmt.Errorf("模板 %s 的 %s 不是完整十六进制颜色", item.ID, name)
		}
	}
	for name, value := range map[string]string{
		"headingStyle": item.HeadingStyle, "blockStyle": item.BlockStyle, "density": item.Density,
		"codeTheme": item.CodeTheme, "riskLevel": item.RiskLevel,
	} {
		if !allowedValues[name][value] {
			return fmt.Errorf("模板 %s 的 %s 不受支持：%s", item.ID, name, value)
		}
	}
	return nil
}

func ContentHash(item TemplateRecord) (string, error) {
	item.ContentHash = ""
	payload, err := json.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("编码模板 %s 失败：%w", item.ID, err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

func DecodeStrict(data []byte, target any) error {
	if err := RejectDuplicateKeys(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := ensureJSONEOF(decoder); err != nil {
		return err
	}
	return nil
}

func RejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := walkJSONValue(decoder); err != nil {
		return err
	}
	return ensureJSONEOF(decoder)
}

func walkJSONValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delim, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delim {
	case '{':
		keys := map[string]bool{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return errors.New("JSON 对象键不是字符串")
			}
			if keys[key] {
				return fmt.Errorf("JSON 包含重复字段：%s", key)
			}
			keys[key] = true
			if err := walkJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err = decoder.Token()
		return err
	case '[':
		for decoder.More() {
			if err := walkJSONValue(decoder); err != nil {
				return err
			}
		}
		_, err = decoder.Token()
		return err
	default:
		return errors.New("JSON 分隔符不合法")
	}
}

func ensureJSONEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}
	return errors.New("JSON 末尾包含多余内容")
}

func validateText(id string, field string, value string, min int, max int) error {
	length := len([]rune(strings.TrimSpace(value)))
	if strings.TrimSpace(value) != value || length < min || length > max {
		return fmt.Errorf("模板 %s 的 %s 长度必须为 %d 到 %d", id, field, min, max)
	}
	if strings.ContainsAny(value, "<>") {
		return fmt.Errorf("模板 %s 的 %s 不能包含 HTML 字符", id, field)
	}
	return nil
}

func idOrEmpty(value string) bool {
	return value == "" || (idPattern.MatchString(value) && len(value) <= 80)
}
