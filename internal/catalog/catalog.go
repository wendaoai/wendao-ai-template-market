package catalog

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

func Sign(root string, keyID string, privateKeyPath string) (DetachedSignature, error) {
	if !idPattern.MatchString(keyID) {
		return DetachedSignature{}, errors.New("key-id 不合法")
	}
	privateKey, err := readPrivateKey(privateKeyPath)
	if err != nil {
		return DetachedSignature{}, err
	}
	catalogPath := filepath.Join(root, "registry", "v1", "stable.json")
	data, err := os.ReadFile(catalogPath)
	if err != nil {
		return DetachedSignature{}, fmt.Errorf("读取 stable.json 失败：%w", err)
	}
	sum := sha256.Sum256(data)
	signature := DetachedSignature{
		SchemaVersion: SchemaVersion,
		Algorithm:     SignatureAlgorithm,
		KeyID:         keyID,
		SignedFile:    "stable.json",
		SHA256:        hex.EncodeToString(sum[:]),
		Signature:     base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, data)),
	}
	payload, err := json.MarshalIndent(signature, "", "  ")
	if err != nil {
		return DetachedSignature{}, err
	}
	path := filepath.Join(root, "registry", "v1", "stable.json.sig")
	if err := atomicWrite(path, append(payload, '\n'), 0o644); err != nil {
		return DetachedSignature{}, err
	}
	return signature, nil
}

func Verify(root string) (Catalog, DetachedSignature, error) {
	var (
		result    Catalog
		signature DetachedSignature
	)
	catalogPath := filepath.Join(root, "registry", "v1", "stable.json")
	data, err := os.ReadFile(catalogPath)
	if err != nil {
		return result, signature, fmt.Errorf("读取 stable.json 失败：%w", err)
	}
	signatureData, err := os.ReadFile(catalogPath + ".sig")
	if err != nil {
		return result, signature, fmt.Errorf("读取 stable.json.sig 失败：%w", err)
	}
	if err := DecodeStrict(data, &result); err != nil {
		return result, signature, fmt.Errorf("解析 stable.json 失败：%w", err)
	}
	if err := DecodeStrict(signatureData, &signature); err != nil {
		return result, signature, fmt.Errorf("解析 stable.json.sig 失败：%w", err)
	}
	if err := ValidateCatalog(result); err != nil {
		return result, signature, err
	}
	if signature.SchemaVersion != SchemaVersion || signature.Algorithm != SignatureAlgorithm || signature.SignedFile != "stable.json" {
		return result, signature, errors.New("签名元数据不受支持")
	}
	sum := sha256.Sum256(data)
	if signature.SHA256 != hex.EncodeToString(sum[:]) {
		return result, signature, errors.New("stable.json SHA-256 不匹配")
	}
	publicKey, err := readPublicKey(filepath.Join(root, "keys", signature.KeyID+".pub"))
	if err != nil {
		return result, signature, err
	}
	rawSignature, err := base64.StdEncoding.DecodeString(signature.Signature)
	if err != nil || len(rawSignature) != ed25519.SignatureSize {
		return result, signature, errors.New("Ed25519 签名格式错误")
	}
	if !ed25519.Verify(publicKey, data, rawSignature) {
		return result, signature, errors.New("stable.json 签名无效")
	}
	return result, signature, nil
}

func GenerateKeyPair(privateKeyPath string, publicKeyPath string) error {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("生成 Ed25519 密钥失败：%w", err)
	}
	seed := privateKey.Seed()
	if err := atomicWrite(privateKeyPath, []byte(base64.StdEncoding.EncodeToString(seed)+"\n"), 0o600); err != nil {
		return err
	}
	if err := atomicWrite(publicKeyPath, []byte(base64.StdEncoding.EncodeToString(publicKey)+"\n"), 0o644); err != nil {
		return err
	}
	return nil
}

func readPrivateKey(path string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 Ed25519 私钥失败：%w", err)
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		return nil, fmt.Errorf("解析 Ed25519 私钥失败：%w", err)
	}
	if len(raw) == ed25519.SeedSize {
		return ed25519.NewKeyFromSeed(raw), nil
	}
	if len(raw) != ed25519.PrivateKeySize {
		return nil, errors.New("Ed25519 私钥必须是 seed 或 private key 的 Base64")
	}
	return ed25519.PrivateKey(raw), nil
}

func readPublicKey(path string) (ed25519.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取 Ed25519 公钥失败：%w", err)
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data)))
	if err != nil {
		return nil, fmt.Errorf("解析 Ed25519 公钥失败：%w", err)
	}
	if len(raw) != ed25519.PublicKeySize {
		return nil, errors.New("Ed25519 公钥长度错误")
	}
	return ed25519.PublicKey(raw), nil
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
