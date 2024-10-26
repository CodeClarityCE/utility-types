package codeclarity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:user,alias:user"`
	Id            uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
}
