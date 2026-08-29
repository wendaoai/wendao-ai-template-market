package catalog

import "encoding/json"

const (
	SchemaVersion    = 2
	CatalogID        = "wechat-article"
	ChannelStable    = "stable"
	RendererContract = "wechat-theme-css-v1"
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
	ID                string          `json:"id"`
	Version           string          `json:"version"`
	ContentHash       string          `json:"contentHash,omitempty"`
	ReleasedAt        string          `json:"releasedAt"`
	ChangeSummary     string          `json:"changeSummary"`
	Status            string          `json:"status"`
	ReplacementID     string          `json:"replacementId"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Source            string          `json:"source"`
	Category          string          `json:"category"`
	Tags              []string        `json:"tags"`
	Selectable        bool            `json:"selectable"`
	Recommended       bool            `json:"recommended"`
	RiskLevel         string          `json:"riskLevel"`
	EditorMode        string          `json:"editorMode"`
	CSS               string          `json:"css"`
	DesignerVariables json.RawMessage `json:"designerVariables,omitempty"`
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
