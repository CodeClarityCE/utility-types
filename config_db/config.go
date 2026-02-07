package config

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Config struct {
	bun.BaseModel `bun:"table:config,alias:config"`
	Id            uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	NvdLast       time.Time `bun:"nvd_last,type:timestamptz"`
	NpmLast       string    `bun:"npm_last"`
	GcveLast      time.Time `bun:"gcve_last,type:timestamptz"`
}
