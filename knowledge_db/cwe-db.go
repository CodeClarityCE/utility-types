package knowledge

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type CWEEntry struct {
	bun.BaseModel        `bun:"table:cwe,alias:c"`
	Id                   uuid.UUID             `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	CWEId                string                `json:"id" bun:"cwe_id"`
	Name                 string                `bun:"name"`
	Abstraction          string                `bun:"abstraction"`
	Structure            string                `bun:"structure"`
	Status               string                `bun:"status"`
	Description          string                `bun:"description"`
	ExtendedDescription  string                `bun:"extended_description"`
	RelatedWeaknesses    []RelatedWeakness     `bun:"related_weaknesses"`
	ModesOfIntroduction  []ModesOfIntroduction `bun:"modes_of_introduction"`
	CommonConsequences   []CommonConsequence   `bun:"common_consequences"`
	DetectionMethods     []DetectionMethod     `bun:"detection_methods"`
	PotentialMitigations []PotentialMitigation `bun:"potential_mitigations"`
	TaxonomyMappings     []TaxonomyMapping     `bun:"taxonomy_mappings"`
	LikelihoodOfExploit  string                `bun:"likelihood_of_exploit"`
	ObservedExamples     []ObservedExamples    `bun:"observed_examples"`
	AlternateTerms       []AlternateTerm       `bun:"alternate_terms"`
	AffectedResources    []string              `bun:"affected_resources"`
	FunctionalAreas      []string              `bun:"functional_areas"`
	Categories           []CategorySimplified  `bun:"categories"`
	ApplicablePlatforms  ApplicablePlatform    `bun:"applicable_platforms"`
}

type CategorySimplified struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

type AlternateTerm struct {
	Term        string `json:"Term"`
	Description string `json:"Description"`
}

type ObservedExamples struct {
	Reference   string `json:"Reference"`
	Description string `json:"Description"`
	Link        string `json:"Link"`
}

type RelatedAttackPattern struct {
	CapecID string `json:"CapecID"`
}

type TaxonomyMapping struct {
	TaxonomyName string `json:"TaxonomyName"`
	EntryID      string `json:"EntryID"`
	EntryName    string `json:"EntryName"`
	MappingFit   string `json:"MappingFit"`
}

type PotentialMitigation struct {
	Phases      []string `json:"Phase"`
	Description string   `json:"Description"`
}

type CommonConsequence struct {
	Scope      []string `json:"Scope"`
	Impact     []string `json:"Impact"`
	Note       string   `json:"Note"`
	Likelihood string   `json:"Likelihood"`
}

type ModesOfIntroduction struct {
	Phase string `json:"Phase"`
	Note  string `json:"Note"`
}

type RelatedWeakness struct {
	Nature  string `json:"Nature"`
	CWEID   string `json:"CWE_ID"`
	ViewID  string `json:"View_ID"`
	Ordinal string `json:"Ordinal"`
	ChainID string `json:"Chain_ID"`
}

type ApplicablePlatform struct {
	Language        []ApplicablePlatformEntry `json:"Language,omitempty"`
	Technology      []ApplicablePlatformEntry `json:"Technology,omitempty"`
	OperatingSystem []ApplicablePlatformEntry `json:"OperatingSystem,omitempty"`
	Architecture    []ApplicablePlatformEntry `json:"Architecture,omitempty"`
}

type ApplicablePlatformEntry struct {
	Name       string `json:"Name"`
	Prevalence string `json:"Prevalence"`
	Class      string `json:"Class"`
}

type DetectionMethod struct {
	Method      string
	Description string
}
