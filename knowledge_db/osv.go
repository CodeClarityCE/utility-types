package knowledge

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OSVItem struct {
	bun.BaseModel    `bun:"table:osv,alias:o"`
	Id               uuid.UUID      `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	OSVId            string         `json:"id" bun:"osv_id"`
	Schema_version   string         `json:"schema_version" bun:"schema_version"`
	Modified         string         `json:"modified" bun:"modified"`
	Published        string         `json:"published" bun:"published"`
	Withdrawn        string         `json:"withdrawn" bun:"withdrawn"`
	Aliases          []string       `json:"aliases" bun:"aliases"`
	Related          []string       `json:"related" bun:"related"`
	Summary          string         `json:"summary" bun:"summary"`
	Details          string         `json:"details" bun:"details"`
	Severity         []Severity     `json:"severity" bun:"severity"`
	Affected         []Affected     `json:"affected" bun:"affected"`
	References       []Reference    `json:"references" bun:"references"`
	Credits          []Credit       `json:"credits" bun:"credits"`
	DatabaseSpecific map[string]any `json:"database_specific" bun:"database_specific"`
	Cwes             []string       `json:"cwes" bun:"cwes"`
	Cve              string         `json:"cve" bun:"cve"`
}

type Severity struct {
	Type  string `json:"type"`
	Score string `json:"score"`
}

type Affected struct {
	Package           OSVPackage     `json:"package"`
	Severity          []Severity     `json:"severity"`
	Ranges            []Range        `json:"ranges"`
	Versions          []string       `json:"versions"`
	EcosystemSpecific map[string]any `json:"ecosystem_specific"`
	DatabaseSpecific  map[string]any `json:"database_specific"`
}

type OSVPackage struct {
	Ecosystem string `json:"ecosystem"`
	Name      string `json:"name"`
	Purl      string `json:"purl"`
}

type Range struct {
	Type             string         `json:"type"`
	Repo             string         `json:"repo"`
	Events           []Event        `json:"events"`
	DatabaseSpecific map[string]any `json:"database_specific"`
}

type Event struct {
	Introduced    string `json:"introduced"`
	Fixed         string `json:"fixed"`
	Last_affected string `json:"last_affected"`
	Limit         string `json:"limit"`
}

type Reference struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Credit struct {
	Name    string   `json:"name"`
	Contact []string `json:"contact"`
	Type    []string `json:"type"`
}
