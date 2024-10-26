package codeclarity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Result struct {
	bun.BaseModel `bun:"table:result,alias:r"`
	Id            uuid.UUID   `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	Result        interface{} `bun:"result"`
	AnalysisId    uuid.UUID   `bun:"analysisId"`
	Plugin        string      `bun:"plugin"`
}
