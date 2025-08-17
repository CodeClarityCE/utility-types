package knowledge

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// FriendsOfPHPAdvisory represents a security advisory from FriendsOfPHP Security Advisories Database
type FriendsOfPHPAdvisory struct {
	bun.BaseModel `bun:"table:friends_of_php,alias:fop"`
	Id            uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	AdvisoryId    string    `bun:"advisory_id,unique" json:"advisory_id"`
	Title         string    `bun:"title" json:"title"`
	CVE           string    `bun:"cve" json:"cve,omitempty"`
	Link          string    `bun:"link" json:"link"`
	Reference     string    `bun:"reference" json:"reference"`
	Composer      string    `bun:"composer" json:"composer,omitempty"`
	Description   string    `bun:"description" json:"description,omitempty"`
	Branches      map[string]AdvisoryBranch `bun:"branches,type:jsonb" json:"branches"`
	Published     string    `bun:"published" json:"published"`
	Modified      string    `bun:"modified" json:"modified"`
}

// AdvisoryBranch represents version information for a specific branch
type AdvisoryBranch struct {
	Versions []string `json:"versions"`
	Time     string   `json:"time"`
}

// FriendsOfPHPAffected represents affected packages for vulnerability matching
type FriendsOfPHPAffected struct {
	PackageName      string   `json:"package_name"`
	Ecosystem        string   `json:"ecosystem"`
	AffectedVersions []string `json:"affected_versions"`
}

// GetAffectedPackages converts branch information to affected packages format
func (f *FriendsOfPHPAdvisory) GetAffectedPackages() []FriendsOfPHPAffected {
	var affected []FriendsOfPHPAffected
	
	packageName := f.Composer
	if packageName == "" {
		// Extract from advisory ID if composer field is missing
		// Advisory IDs are typically in format "vendor/package/YYYY-MM-DD"
		parts := []rune(f.AdvisoryId)
		var vendor, pkg string
		slashCount := 0
		var current []rune
		
		for _, char := range parts {
			if char == '/' {
				slashCount++
				if slashCount == 1 {
					vendor = string(current)
					current = []rune{}
				} else if slashCount == 2 {
					pkg = string(current)
					break
				}
			} else {
				current = append(current, char)
			}
		}
		
		if vendor != "" && pkg != "" {
			packageName = vendor + "/" + pkg
		}
	}

	if packageName != "" {
		var allVersions []string
		for _, branch := range f.Branches {
			allVersions = append(allVersions, branch.Versions...)
		}
		
		affected = append(affected, FriendsOfPHPAffected{
			PackageName:      packageName,
			Ecosystem:        "Packagist",
			AffectedVersions: allVersions,
		})
	}

	return affected
}