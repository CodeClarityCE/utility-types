package plugin

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Plugin struct {
	bun.BaseModel `bun:"table:plugin,alias:plugin"`
	Id            uuid.UUID      `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	Name          string         `bun:"name" json:"name"`
	Version       string         `bun:"version" json:"version"`
	DependsOn     []string       `bun:"depends_on" json:"depends_on"`
	Description   string         `bun:"description" json:"description"`
	Config        map[string]any `bun:"config" json:"config"`
}
