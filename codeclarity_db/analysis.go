package codeclarity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Analysis struct {
	bun.BaseModel  `bun:"table:analysis,alias:analysis"`
	Id             uuid.UUID  `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	AnalyzerId     uuid.UUID  `bun:"analyzerId"`
	OrganizationId uuid.UUID  `bun:"organizationId"`
	ProjectId      *uuid.UUID `bun:"projectId"` // Pointer allows null values
	Config         map[string]any
	Stage          int
	Steps          [][]Step
	Status         AnalysisStatus
	Commit         string `bun:"commit_hash"`
	Branch         string
	// Results       []*result.Result `bun:"rel:has-many,join:id=analysisId"`
}

type Step struct {
	Name       string
	Version    string
	Config     map[string]any
	Status     AnalysisStatus
	Result     map[string]any
	Started_on string
	Ended_on   string
}

type AnalysisStatus string

const (
	SUCCESS     AnalysisStatus = "success"
	UPDATING_DB AnalysisStatus = "updating_db"
	ONGOING     AnalysisStatus = "ongoing"
	FAILURE     AnalysisStatus = "failure"
	COMPLETED   AnalysisStatus = "completed"
	STARTED     AnalysisStatus = "started"
)
