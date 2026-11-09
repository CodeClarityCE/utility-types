package knowledge

import (
	"time"

	"github.com/uptrace/bun"
)

// Knowledge source names, used as the primary key of SourceSync.
const (
	SourceCWE          = "cwe"
	SourceEPSS         = "epss"
	SourceFriendsOfPHP = "friends_of_php"
	SourceGCVE         = "gcve"
	SourceLicenses     = "licenses"
	SourceNVD          = "nvd"
	SourceOSV          = "osv"
)

// Values allowed in SourceSync.LastOutcome.
const (
	SourceSyncSuccess = "success"
	SourceSyncFailure = "failure"
)

// SourceSync is the sync state of one knowledge source.
// Watermark is the point up to which the source is completely synchronized.
type SourceSync struct {
	bun.BaseModel     `bun:"table:knowledge_source_sync,alias:kss"`
	Source            string     `bun:"source,pk" json:"source"`
	Watermark         *time.Time `bun:"watermark,type:timestamptz" json:"watermark"`
	LastRunStartedAt  *time.Time `bun:"last_run_started_at,type:timestamptz" json:"last_run_started_at"`
	LastRunFinishedAt *time.Time `bun:"last_run_finished_at,type:timestamptz" json:"last_run_finished_at"`
	LastOutcome       *string    `bun:"last_outcome" json:"last_outcome"`
	LastError         *string    `bun:"last_error" json:"last_error"`
	RowCount          int64      `bun:"row_count,notnull" json:"row_count"`
	// BulkCompletedAt is set while a bulk import is done but its enrichment is
	// not. It holds the time the imported data is current to.
	BulkCompletedAt *time.Time `bun:"bulk_completed_at,type:timestamptz" json:"bulk_completed_at"`
}
