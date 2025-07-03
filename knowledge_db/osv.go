package knowledge

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type OSVItem struct {
	bun.BaseModel    `bun:"table:osv,alias:o"`
	Id               uuid.UUID      `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	OSVId            string         `json:"id" bun:"osv_id"`
	Schema_version   string         `json:"schema_version" bun:"schema_version"`
	Vlai_score       string         `bun:"vlai_score"`
	Vlai_confidence  float64        `bun:"vlai_confidence"`
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
	Type    string   `json:"type"`
}

// Custom unmarshaler for Credit to handle both single string and slice of strings
func (c *Credit) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name    string   `json:"name"`
		Contact []string `json:"contact"`
		Type    string   `json:"type"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid format for Credit: %s", string(data))
	}

	c.Name = raw.Name
	c.Contact = raw.Contact
	c.Type = raw.Type
	return nil
}
