package catalog

const (
	SchemaVersion      = 1
	CatalogID          = "wechat-article"
	ChannelStable      = "stable"
	RendererContract   = "wechat-theme-preset-v1"
	SignatureAlgorithm = "ed25519"
)

type SourceCatalog struct {
	SchemaVersion    int    `json:"schemaVersion"`
	CatalogID        string `json:"catalogId"`
	CatalogVersion   string `json:"catalogVersion"`
	Channel          string `json:"channel"`
	RendererContract string `json:"rendererContract"`
	MinAppVersion    string `json:"minAppVersion"`
	PublishedAt      string `json:"publishedAt"`
}

type SourceTemplate struct {
	SchemaVersion int `json:"schemaVersion"`
	TemplateRecord
}

type TemplateRecord struct {
	ID            string   `json:"id"`
	Version       string   `json:"version"`
	ContentHash   string   `json:"contentHash,omitempty"`
	ReleasedAt    string   `json:"releasedAt"`
	ChangeSummary string   `json:"changeSummary"`
	Status        string   `json:"status"`
	ReplacementID string   `json:"replacementId"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Source        string   `json:"source"`
	Category      string   `json:"category"`
	Tags          []string `json:"tags"`
	Accent        string   `json:"accent"`
	Background    string   `json:"background"`
	Surface       string   `json:"surface"`
	Text          string   `json:"text"`
	Muted         string   `json:"muted"`
	Border        string   `json:"border"`
	HeadingStyle  string   `json:"headingStyle"`
	BlockStyle    string   `json:"blockStyle"`
	Density       string   `json:"density"`
	CodeTheme     string   `json:"codeTheme"`
	Selectable    bool     `json:"selectable"`
	Recommended   bool     `json:"recommended"`
	RiskLevel     string   `json:"riskLevel"`
}

type Catalog struct {
	SchemaVersion    int              `json:"schemaVersion"`
	CatalogID        string           `json:"catalogId"`
	CatalogVersion   string           `json:"catalogVersion"`
	Channel          string           `json:"channel"`
	RendererContract string           `json:"rendererContract"`
	MinAppVersion    string           `json:"minAppVersion"`
	GeneratedAt      string           `json:"generatedAt"`
	Templates        []TemplateRecord `json:"templates"`
}

type DetachedSignature struct {
	SchemaVersion int    `json:"schemaVersion"`
	Algorithm     string `json:"algorithm"`
	KeyID         string `json:"keyId"`
	SignedFile    string `json:"signedFile"`
	SHA256        string `json:"sha256"`
	Signature     string `json:"signature"`
}
