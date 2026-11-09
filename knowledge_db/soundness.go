package knowledge

import "time"

// Why a source stopped asserting a record.
const (
	RetiredRejected  = "rejected"
	RetiredWithdrawn = "withdrawn"
	RetiredRemoved   = "removed"
)

// What an OSV advisory states about a package.
const (
	AdvisoryKindVulnerability = "vulnerability"
	AdvisoryKindMalicious     = "malicious"
)

// Retirement records why and when a source stopped asserting a record.
type Retirement struct {
	RetiredReason string    `bun:"retired_reason,nullzero" json:"-"`
	RetiredAt     time.Time `bun:"retired_at,nullzero,type:timestamptz" json:"-"`
}
