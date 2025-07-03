package knowledge

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type EPSS struct {
	bun.BaseModel `bun:"table:epss,alias:n"`
	Id            uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	CVE           string    `bun:"cve" json:"cve"`
	Score         float32   `bun:"score" json:"score"`
	Percentile    float32   `bun:"percentile" json:"percentile"`
}
