package catalog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func LoadSource(root string) (SourceCatalog, []SourceTemplate, error) {
	var source SourceCatalog
	data, err := os.ReadFile(filepath.Join(root, "catalog.json"))
	if err != nil {
		return source, nil, fmt.Errorf("读取 catalog.json 失败：%w", err)
	}
	if err := DecodeStrict(data, &source); err != nil {
		return source, nil, fmt.Errorf("解析 catalog.json 失败：%w", err)
	}

	files, err := filepath.Glob(filepath.Join(root, "templates", "wechat-article", "*.json"))
	if err != nil {
		return source, nil, fmt.Errorf("扫描模板失败：%w", err)
	}
	sort.Strings(files)
	templates := make([]SourceTemplate, 0, len(files))
	for _, file := range files {
		var item SourceTemplate
		data, err := os.ReadFile(file)
		if err != nil {
			return source, nil, fmt.Errorf("读取模板 %s 失败：%w", filepath.Base(file), err)
		}
		if err := DecodeStrict(data, &item); err != nil {
			return source, nil, fmt.Errorf("解析模板 %s 失败：%w", filepath.Base(file), err)
		}
		if filepath.Base(file) != item.ID+".json" {
			return source, nil, fmt.Errorf("模板文件名必须与 ID 一致：%s", filepath.Base(file))
		}
		templates = append(templates, item)
	}
	if err := ValidateSource(source, templates); err != nil {
		return source, nil, err
	}
	return source, templates, nil
}

func Build(root string) ([]byte, Catalog, error) {
	source, templates, err := LoadSource(root)
	if err != nil {
		return nil, Catalog{}, err
	}
	records := make([]TemplateRecord, len(templates))
	for index := range templates {
		record := templates[index].TemplateRecord
		record.ContentHash, err = ContentHash(record)
		if err != nil {
			return nil, Catalog{}, err
		}
		records[index] = record
	}
	sort.Slice(records, func(i, j int) bool { return records[i].ID < records[j].ID })
	result := Catalog{
		SchemaVersion:    source.SchemaVersion,
		CatalogID:        source.CatalogID,
		CatalogVersion:   source.CatalogVersion,
		Channel:          source.Channel,
		RendererContract: source.RendererContract,
		MinAppVersion:    source.MinAppVersion,
		GeneratedAt:      source.PublishedAt,
		Templates:        records,
	}
	if err := ValidateCatalog(result); err != nil {
		return nil, Catalog{}, err
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, Catalog{}, fmt.Errorf("生成 stable.json 失败：%w", err)
	}
	return append(data, '\n'), result, nil
}

func WriteBuild(root string, check bool) (Catalog, error) {
	data, result, err := Build(root)
	if err != nil {
		return Catalog{}, err
	}
	path := filepath.Join(root, "registry", "v1", "stable.json")
	if check {
		current, err := os.ReadFile(path)
		if err != nil {
			return Catalog{}, fmt.Errorf("读取已生成 stable.json 失败：%w", err)
		}
		if !bytes.Equal(current, data) {
			return Catalog{}, errors.New("registry/v1/stable.json 未更新，请运行 marketctl build")
		}
		return result, nil
	}
	if err := atomicWrite(path, data, 0o644); err != nil {
		return Catalog{}, err
	}
	return result, nil
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("创建目录失败：%w", err)
	}
	temporary := path + ".tmp"
	if err := os.WriteFile(temporary, data, mode); err != nil {
		return fmt.Errorf("写入临时文件失败：%w", err)
	}
	if err := os.Rename(temporary, path); err != nil {
		if removeErr := os.Remove(path); removeErr != nil && !os.IsNotExist(removeErr) {
			_ = os.Remove(temporary)
			return fmt.Errorf("移除旧文件失败：%w", removeErr)
		}
		if retryErr := os.Rename(temporary, path); retryErr != nil {
			_ = os.Remove(temporary)
			return fmt.Errorf("替换文件失败：%w", retryErr)
		}
	}
	return nil
}
