package knowledge

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// KEV is one entry of the CISA Known Exploited Vulnerabilities catalog, keyed by its CVE.
type KEV struct {
	bun.BaseModel `bun:"table:kev,alias:kev"`
	Id            uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()" json:"-"`
	CVE           string    `bun:"cve,notnull" json:"cve"`
	VendorProject string    `bun:"vendor_project" json:"vendor_project"`
	Product       string    `bun:"product" json:"product"`
	Name          string    `bun:"vulnerability_name" json:"vulnerability_name"`
	DateAdded     time.Time `bun:"date_added,type:date,notnull" json:"date_added"`
	// DueDate is when US federal agencies must have remediated it; nil when the catalog gives none.
	DueDate            *time.Time `bun:"due_date,type:date" json:"due_date"`
	ShortDescription   string     `bun:"short_description" json:"short_description"`
	RequiredAction     string     `bun:"required_action" json:"required_action"`
	Notes              string     `bun:"notes" json:"notes"`
	KnownRansomwareUse bool       `bun:"known_ransomware_use,notnull" json:"known_ransomware_use"`
	CWEs               []string   `bun:"cwes,type:jsonb,notnull" json:"cwes"`
}
