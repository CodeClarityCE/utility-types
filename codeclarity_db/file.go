package codeclarity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type File struct {
	bun.BaseModel `bun:"table:file,alias:file"`
	Id            uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	Added_on      time.Time
	Type          string
	Name          string
	ProjectId     uuid.UUID `bun:"projectId"`
	Project       Project   `bun:"rel:belongs-to,join:'projectId'=id"`
	AddedById     uuid.UUID `bun:"addedById"`
	AddedBy       User      `bun:"rel:belongs-to,join:'addedById'=id"`
}
