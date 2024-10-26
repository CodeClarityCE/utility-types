package codeclarity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Project struct {
	bun.BaseModel  `bun:"table:project,alias:project"`
	Id             uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	Name           string
	Description    string
	Integration_id string `bun:"integrationId"`
	Type           string
	Url            string
}
