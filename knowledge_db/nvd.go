package knowledge

import (
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type NVD struct {
	Vulnerabilities []map[string]NVDItem `json:"vulnerabilities"`
}

type NVDItem struct {
	bun.BaseModel     `bun:"table:nvd,alias:n"`
	Id                uuid.UUID       `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	NVDId             string          `bun:"nvd_id" json:"id"`
	SourceIdentifier  string          `bun:"sourceIdentifier"`
	Published         string          `bun:"published"`
	LastModified      string          `bun:"lastModified"`
	VulnStatus        string          `bun:"vulnStatus"`
	Descriptions      []Descriptions  `bun:"descriptions"`
	Vlai_score        string          `bun:"vlai_score"`
	Vlai_confidence   float64         `bun:"vlai_confidence"`
	Metrics           Metrics         `bun:"metrics"`
	Weaknesses        []Weaknesses    `bun:"weaknesses"`
	Configurations    []Configuration `bun:"configurations"`
	AffectedFlattened []Sources       `bun:"affectedFlattened"`
	Affected          []NVDAffected   `bun:"affected"`
	References        []References    `bun:"references"`
}
type Descriptions struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type Configuration struct {
	Nodes []Node `json:"nodes"`
}
type Node struct {
	Operator string    `json:"operator"`
	Negate   bool      `json:"negate"`
	CpeMatch []Sources `json:"cpematch"`
	Children []Node    `json:"children"`
}
type CVSS3 struct {
	AttackComplexity              string  `json:"attackComplexity,omitempty"`
	AttackVector                  string  `json:"attackVector,omitempty"`
	AvailabilityImpact            string  `json:"availabilityImpact,omitempty"`
	AvailabilityRequirement       string  `json:"availabilityRequirement,omitempty"`
	BaseScore                     float64 `json:"baseScore"`
	BaseSeverity                  string  `json:"baseSeverity"`
	ConfidentialityImpact         string  `json:"confidentialityImpact,omitempty"`
	ConfidentialityRequirement    string  `json:"confidentialityRequirement,omitempty"`
	EnvironmentalScore            float64 `json:"environmentalScore,omitempty"`
	EnvironmentalSeverity         string  `json:"environmentalSeverity,omitempty"`
	ExploitCodeMaturity           string  `json:"exploitCodeMaturity,omitempty"`
	IntegrityImpact               string  `json:"integrityImpact,omitempty"`
	IntegrityRequirement          string  `json:"integrityRequirement,omitempty"`
	ModifiedAttackComplexity      string  `json:"modifiedAttackComplexity,omitempty"`
	ModifiedAttackVector          string  `json:"modifiedAttackVector,omitempty"`
	ModifiedAvailabilityImpact    string  `json:"modifiedAvailabilityImpact,omitempty"`
	ModifiedConfidentialityImpact string  `json:"modifiedConfidentialityImpact,omitempty"`
	ModifiedIntegrityImpact       string  `json:"modifiedIntegrityImpact,omitempty"`
	ModifiedPrivilegesRequired    string  `json:"modifiedPrivilegesRequired,omitempty"`
	ModifiedScope                 string  `json:"modifiedScope,omitempty"`
	ModifiedUserInteraction       string  `json:"modifiedUserInteraction,omitempty"`
	PrivilegesRequired            string  `json:"privilegesRequired,omitempty"`
	RemediationLevel              string  `json:"remediationLevel,omitempty"`
	ReportConfidence              string  `json:"reportConfidence,omitempty"`
	Scope                         string  `json:"scope,omitempty"`
	TemporalScore                 float64 `json:"temporalScore,omitempty"`
	TemporalSeverity              string  `json:"temporalSeverity,omitempty"`
	UserInteraction               string  `json:"userInteraction,omitempty"`
	VectorString                  string  `json:"vectorString"`
	Version                       string  `json:"version"`
}

type CVSSV2 struct {
	AccessComplexity           string  `json:"accessComplexity,omitempty"`
	AccessVector               string  `json:"accessVector,omitempty"`
	Authentication             string  `json:"authentication,omitempty"`
	AvailabilityImpact         string  `json:"availabilityImpact,omitempty"`
	AvailabilityRequirement    string  `json:"availabilityRequirement,omitempty"`
	BaseScore                  float64 `json:"baseScore"`
	CollateralDamagePotential  string  `json:"collateralDamagePotential,omitempty"`
	ConfidentialityImpact      string  `json:"confidentialityImpact,omitempty"`
	ConfidentialityRequirement string  `json:"confidentialityRequirement,omitempty"`
	EnvironmentalScore         float64 `json:"environmentalScore,omitempty"`
	Exploitability             string  `json:"exploitability,omitempty"`
	IntegrityImpact            string  `json:"integrityImpact,omitempty"`
	IntegrityRequirement       string  `json:"integrityRequirement,omitempty"`
	RemediationLevel           string  `json:"remediationLevel,omitempty"`
	ReportConfidence           string  `json:"reportConfidence,omitempty"`
	TargetDistribution         string  `json:"targetDistribution,omitempty"`
	TemporalScore              float64 `json:"temporalScore,omitempty"`
	VectorString               string  `json:"vectorString"`
	Version                    string  `json:"version"`
}

type ImpactMetricV2 struct {
	AcInsufInfo             bool    `json:"acInsufInfo"`
	BaseSeverity            string  `json:"baseSeverity"`
	CvssData                CVSSV2  `json:"cvssData"`
	ExploitabilityScore     float64 `json:"exploitabilityScore"`
	ImpactScore             float64 `json:"impactScore"`
	ObtainAllPrivilege      bool    `json:"obtainAllPrivilege"`
	ObtainOtherPrivilege    bool    `json:"obtainOtherPrivilege"`
	ObtainUserPrivilege     bool    `json:"obtainUserPrivilege"`
	Source                  string  `json:"source"`
	Type                    string  `json:"type"`
	UserInteractionRequired bool    `json:"userInteractionRequired"`
}

type ImpactMetricV31 struct {
	CvssData            CVSS3   `json:"cvssData"`
	ExploitabilityScore float64 `json:"exploitabilityScore"`
	ImpactScore         float64 `json:"impactScore"`
	Source              string  `json:"source"`
	Type                string  `json:"type"`
}

type ImpactMetricV30 struct {
	CvssData            CVSS3   `json:"cvssData"`
	ExploitabilityScore float64 `json:"exploitabilityScore"`
	ImpactScore         float64 `json:"impactScore"`
	Source              string  `json:"source"`
	Type                string  `json:"type"`
}

type Metrics struct {
	CvssMetricV2  []ImpactMetricV2  `json:"cvssMetricV2"`
	CvssMetricV30 []ImpactMetricV30 `json:"cvssMetricV30"`
	CvssMetricV31 []ImpactMetricV31 `json:"cvssMetricV31"`
}
type Weaknesses struct {
	Source       string         `json:"source"`
	Type         string         `json:"type"`
	Descriptions []Descriptions `json:"descriptions"`
}
type CriteriaDict struct {
	Part      string `json:"part"`
	Vendor    string `json:"vendor"`
	Product   string `json:"product"`
	Version   string `json:"version"`
	Update    string `json:"update"`
	Edition   string `json:"edition"`
	Language  string `json:"language"`
	SwEdition string `json:"sw_edition"`
	TargetSw  string `json:"target_sw"`
	TargetHw  string `json:"target_hw"`
	Other     string `json:"other"`
}
type Sources struct {
	Vulnerable            bool         `json:"vulnerable"`
	Criteria              string       `json:"criteria"`
	MatchCriteriaID       string       `json:"matchCriteriaId"`
	VersionEndIncluding   string       `json:"versionEndIncluding"`
	VersionEndExcluding   string       `json:"versionEndExcluding"`
	VersionStartIncluding string       `json:"versionStartIncluding"`
	VersionStartExcluding string       `json:"versionStartExcluding"`
	CriteriaDict          CriteriaDict `json:"criteriaDict"`
}
type NVDAffected struct {
	Sources                   []Sources `json:"sources"`
	RunningOn                 []Sources `json:"running-on"`
	RunningOnApplicationsOnly []Sources `json:"running-on-applications-only"`
}
type References struct {
	URL    string   `json:"url"`
	Source string   `json:"source"`
	Tags   []string `json:"tags"`
}

func GetVulns(nvd NVD) []NVDItem {
	var vulns []NVDItem

	// We iterate over the vulnerabilities and create a new CVE object
	for key := range nvd.Vulnerabilities {
		cve := nvd.Vulnerabilities[key]["cve"]
		// We set the key to the id so we can use it as a key in the database
		// cve.Key = cve.Id
		cve.Affected = createAffected(cve)

		// We flatten the affected array so we can easily query it
		var flattened []Sources
		if len(cve.Affected) != 0 {
			flattened = append(flattened, cve.Affected[0].Sources...)
			flattened = append(flattened, cve.Affected[0].RunningOn...)
			flattened = append(flattened, cve.Affected[0].RunningOnApplicationsOnly...)
		}
		cve.AffectedFlattened = flattened

		for i, reference := range cve.References {
			if reference.Tags == nil {
				cve.References[i].Tags = make([]string, 0)
			}
		}

		// We dont need the configurations anymore
		cve.Configurations = nil
		vulns = append(vulns, cve)
	}

	return vulns
}

func createAffected(cve NVDItem) []NVDAffected {
	var affected []NVDAffected

	// See why configurations is now an array
	if len(cve.Configurations) > 0 {
		for _, config := range cve.Configurations[0].Nodes {
			// Three entries to fill: the actual source that is vulnerable, secondly what its running on and lastly running on but only applications
			// Example:
			//   source: bootstrap
			//   running-on: django, windows
			//   running-on-applicaitons-only: django

			if config.Operator == "AND" {
				if len(config.Children) < 2 {
					if len(config.CpeMatch) > 0 {
						sources := config.CpeMatch

						for source := range sources {
							sources[source].CriteriaDict = parseConfig(sources[source])
						}

						if validateLibrary(sources) {
							affected = append(affected, NVDAffected{
								Sources:                   filterCpe(sources),
								RunningOn:                 []Sources{},
								RunningOnApplicationsOnly: []Sources{},
							})
						}
					}
				} else {
					running_on := config.Children[1].CpeMatch
					sources := config.Children[0].CpeMatch

					for run_on := range running_on {
						running_on[run_on].CriteriaDict = parseConfig(running_on[run_on])
					}

					for source := range sources {
						sources[source].CriteriaDict = parseConfig(sources[source])
					}

					// We only insert the affected object into the report if the report is about a library / application that is vulnerable
					// We dont care about vulnerabilities about hardware systems or operating systems
					if validateLibrary(sources) {
						affected = append(affected, NVDAffected{
							Sources:                   filterCpe(sources),
							RunningOn:                 running_on,
							RunningOnApplicationsOnly: filterCpe(running_on),
						})
					}
				}
			} else if config.Operator == "OR" {
				sources := config.CpeMatch
				for source := range sources {
					sources[source].CriteriaDict = parseConfig(sources[source])
				}

				// We only insert the affected object into the report if the report is about a library / application that is vulnerable
				// We dont care about vulnerabilities about hardware systems or operating systems
				if validateLibrary(sources) {
					affected = append(affected, NVDAffected{
						Sources:                   filterCpe(sources),
						RunningOn:                 []Sources{},
						RunningOnApplicationsOnly: []Sources{},
					})
				}
			}
		}
	}

	return affected
}

func parseConfig(config Sources) CriteriaDict {
	// parsed_cpe = cpe_parser(config["cpe23Uri"])
	// config["cpe23Wfn"] = parsed_cpe.as_wfn()
	criteria_string := strings.Split(config.Criteria, ":")
	criteria := CriteriaDict{
		Part:      criteria_string[2],
		Vendor:    criteria_string[3],
		Product:   criteria_string[4],
		Version:   criteria_string[5],
		Update:    criteria_string[6],
		Edition:   criteria_string[7],
		Language:  criteria_string[8],
		SwEdition: criteria_string[9],
		TargetSw:  criteria_string[10],
		TargetHw:  criteria_string[11],
		Other:     criteria_string[12],
	}

	return criteria
}

// checks if the vulnerability is for a library or not
func validateLibrary(sources []Sources) bool {

	for _, source := range sources {
		// a stands for application
		// o stands for operating system
		// h stands for hardware
		if source.CriteriaDict.Part == "a" {
			return true
		}
	}
	return false
}

func filterCpe(sources []Sources) []Sources {
	var application_cpes []Sources

	for _, source := range sources {
		// a stands for application
		// o stands for operating system
		// h stands for hardware
		if source.CriteriaDict.Part == "a" {
			application_cpes = append(application_cpes, source)
		}
	}
	return application_cpes
}
