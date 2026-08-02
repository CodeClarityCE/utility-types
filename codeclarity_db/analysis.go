package codeclarity

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Analysis struct {
	bun.BaseModel  `bun:"table:analysis,alias:analysis"`
	Id             uuid.UUID  `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	Created_on     time.Time  `bun:"created_on"`
	AnalyzerId     uuid.UUID  `bun:"analyzerId"`
	OrganizationId uuid.UUID  `bun:"organizationId"`
	ProjectId      *uuid.UUID `bun:"projectId"` // Pointer allows null values
	Config         map[string]any
	Stage          int
	Steps          [][]Step
	Status         AnalysisStatus
	Commit         string `bun:"commit_hash"`
	Branch         string
	// Schedule_type mirrors the API column: "once" for a single run (default),
	// or "daily"/"weekly" for a recurring template. Recurring templates sit in a
	// pre-dispatch status indefinitely and are driven by the scheduler (which
	// clones them into "once" executions), so the reaper must never re-drive them.
	Schedule_type string `bun:"schedule_type"`
	// FailureReason records why the analysis was marked FAILURE (e.g. the
	// download error), truncated by the writer to fit the API's varchar(500)
	// column; NULL for analyses that never failed.
	FailureReason string `bun:"failure_reason,nullzero"`
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
	CANCELLED   AnalysisStatus = "cancelled"
	// REQUESTED and TRIGGERED are the API's pre-dispatch statuses (see the API's
	// AnalysisStatus enum). An analysis is created REQUESTED before the dispatcher
	// consumes its api_request message and flips it to STARTED; if that message is
	// lost (e.g. RabbitMQ wiped on restart), the analysis is stranded REQUESTED.
	// The reaper treats these as non-terminal so such orphans are recovered.
	REQUESTED AnalysisStatus = "requested"
	TRIGGERED AnalysisStatus = "triggered"
)
