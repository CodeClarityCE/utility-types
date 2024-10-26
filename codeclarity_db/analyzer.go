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
	OrganizationId string `bun:"organizationId"`
	CreatedById    string `bun:"createdById"`
}
