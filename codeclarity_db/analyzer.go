package codeclarity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Analyzer struct {
	bun.BaseModel  `bun:"table:analyzer,alias:analyzer"`
	Id             uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	Name           string
	Description    string
	CreatedOn      string `bun:"created_on"`
	Steps          [][]Step
	LanguageConfig map[string]LanguagePlugins `bun:"language_config,type:jsonb"`
	OrganizationId string                     `bun:"organizationId"`
	CreatedById    string                     `bun:"createdById"`
}

// LanguagePlugins lists the plugins that run for one language.
type LanguagePlugins struct {
	Plugins []string `json:"plugins"`
}
