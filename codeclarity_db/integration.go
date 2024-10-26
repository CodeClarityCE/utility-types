package codeclarity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Integration struct {
	bun.BaseModel   `bun:"table:integration,alias:integration"`
	Id              uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	IntegrationType string    `bun:"integration_type"`
	AccessToken     string    `bun:"access_token"`
	TokenType       string    `bun:"token_type"`
	Invalid         bool      `bun:"invalid"`
	AddedOn         string    `bun:"added_on"`
	AddedBy         string    `bun:"ownerId"`
	ServiceDomain   string    `bun:"service_domain"`
}
