package knowledge

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// GCVEItem represents a vulnerability from CIRCL's vulnerability-lookup
// in CVE Record v5.x format.
type GCVEItem struct {
	bun.BaseModel     `bun:"table:gcve,alias:g"`
	Id                uuid.UUID         `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	GCVEId            string            `bun:"gcve_id" json:"gcve_id"`
	CVEId             string            `bun:"cve_id" json:"cve_id"`
	DataVersion       string            `bun:"data_version" json:"dataVersion"`
	State             string            `bun:"state" json:"state"`
	DatePublished     string            `bun:"date_published" json:"datePublished"`
	DateUpdated       string            `bun:"date_updated" json:"dateUpdated"`
	AssignerOrgId     string            `bun:"assigner_org_id" json:"assignerOrgId"`
	Descriptions      []GCVEDescription `bun:"descriptions,type:jsonb" json:"descriptions"`
	Affected          []GCVEAffected    `bun:"affected,type:jsonb" json:"affected"`
	AffectedFlattened []GCVEProduct     `bun:"affected_flattened,type:jsonb" json:"affected_flattened"`
	Metrics           []GCVEMetricEntry `bun:"metrics,type:jsonb" json:"metrics"`
	ProblemTypes      []GCVEProblemType `bun:"problem_types,type:jsonb" json:"problemTypes"`
	References        []GCVEReference   `bun:"\"references\",type:jsonb" json:"references"`
	ADPEnrichments    []GCVEAdp         `bun:"adp_enrichments,type:jsonb" json:"adp_enrichments"`
	Cwes              []string          `bun:"cwes,type:text[]" json:"cwes"`
	VlaiScore         string            `bun:"vlai_score"`
	VlaiConfidence    float64           `bun:"vlai_confidence"`
}

type GCVEDescription struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type GCVEAffected struct {
	Vendor        string        `json:"vendor"`
	Product       string        `json:"product"`
	DefaultStatus string        `json:"defaultStatus,omitempty"`
	Versions      []GCVEVersion `json:"versions"`
	Platforms     []string      `json:"platforms,omitempty"`
}

type GCVEVersion struct {
	Version         string `json:"version"`
	Status          string `json:"status"`
	LessThan        string `json:"lessThan,omitempty"`
	LessThanOrEqual string `json:"lessThanOrEqual,omitempty"`
	VersionType     string `json:"versionType,omitempty"`
}

// GCVEProduct is a flattened vendor+product pair for efficient JSONB containment queries.
type GCVEProduct struct {
	Vendor  string `json:"vendor"`
	Product string `json:"product"`
}

// GCVEMetricEntry represents a single metrics entry from the CVE Record format.
// A CVE record can have multiple metric entries with different CVSS versions.
type GCVEMetricEntry struct {
	CvssV31 *GCVECvssScore `json:"cvssV3_1,omitempty"`
	CvssV30 *GCVECvssScore `json:"cvssV3_0,omitempty"`
	CvssV40 *GCVECvssScore `json:"cvssV4_0,omitempty"`
	CvssV2  *GCVECvssScore `json:"cvssV2_0,omitempty"`
}

type GCVECvssScore struct {
	Version               string  `json:"version"`
	BaseScore             float64 `json:"baseScore"`
	BaseSeverity          string  `json:"baseSeverity"`
	VectorString          string  `json:"vectorString"`
	AttackComplexity      string  `json:"attackComplexity,omitempty"`
	AttackVector          string  `json:"attackVector,omitempty"`
	AvailabilityImpact    string  `json:"availabilityImpact,omitempty"`
	ConfidentialityImpact string  `json:"confidentialityImpact,omitempty"`
	IntegrityImpact       string  `json:"integrityImpact,omitempty"`
	PrivilegesRequired    string  `json:"privilegesRequired,omitempty"`
	Scope                 string  `json:"scope,omitempty"`
	UserInteraction       string  `json:"userInteraction,omitempty"`
}

type GCVEProblemType struct {
	Descriptions []GCVEProblemTypeDescription `json:"descriptions"`
}

type GCVEProblemTypeDescription struct {
	CweId       string `json:"cweId,omitempty"`
	Type        string `json:"type"`
	Lang        string `json:"lang"`
	Description string `json:"description"`
}

type GCVEReference struct {
	URL  string   `json:"url"`
	Name string   `json:"name,omitempty"`
	Tags []string `json:"tags,omitempty"`
}

type GCVEAdp struct {
	ProviderOrgId string            `json:"providerOrgId"`
	ShortName     string            `json:"shortName,omitempty"`
	Title         string            `json:"title,omitempty"`
	Affected      []GCVEAffected    `json:"affected,omitempty"`
	Metrics       []GCVEMetricEntry `json:"metrics,omitempty"`
}
