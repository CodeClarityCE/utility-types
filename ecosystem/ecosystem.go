package ecosystem

import (
	"regexp"
	"time"

	"github.com/uptrace/bun"
)

// EcosystemInfo contains information about a package ecosystem
type EcosystemInfo struct {
	Name                  string   `json:"name"`
	Ecosystem             string   `json:"ecosystem"`
	Language              string   `json:"language"`
	PackageManagerPattern string   `json:"packageManagerPattern"` // Will be converted to RegExp in frontend
	DefaultPackageManager string   `json:"defaultPackageManager"`
	Icon                  string   `json:"icon"`
	Color                 string   `json:"color"`
	Website               string   `json:"website"`
	PurlType              string   `json:"purlType"`
	RegistryUrl           string   `json:"registryUrl"`
	Tools                 []string `json:"tools"`
}

// PluginEcosystemMap maps plugin names to their ecosystem information
type PluginEcosystemMap map[string]EcosystemInfo

// GetDefaultEcosystemMap returns the default mapping of plugins to ecosystems
func GetDefaultEcosystemMap() PluginEcosystemMap {
	return PluginEcosystemMap{
		"js-sbom": {
			Name:                  "JavaScript",
			Ecosystem:             "npm",
			Language:              "JavaScript",
			PackageManagerPattern: `(npm|yarn|pnpm|bun)`,
			DefaultPackageManager: "npm",
			Icon:                  "devicon:javascript",
			Color:                 "#F7DF1E",
			Website:               "https://www.npmjs.com",
			PurlType:              "npm",
			RegistryUrl:           "https://registry.npmjs.org",
			Tools:                 []string{"npm", "yarn", "pnpm", "bun"},
		},
		"php-sbom": {
			Name:                  "PHP",
			Ecosystem:             "packagist",
			Language:              "PHP",
			PackageManagerPattern: `composer`,
			DefaultPackageManager: "composer",
			Icon:                  "devicon:php",
			Color:                 "#777BB4",
			Website:               "https://packagist.org",
			PurlType:              "composer",
			RegistryUrl:           "https://packagist.org",
			Tools:                 []string{"composer"},
		},
		// Future language support can be added here
		"python-sbom": {
			Name:                  "Python",
			Ecosystem:             "pypi",
			Language:              "Python",
			PackageManagerPattern: `(pip|poetry|pipenv|conda)`,
			DefaultPackageManager: "pip",
			Icon:                  "devicon:python",
			Color:                 "#3776AB",
			Website:               "https://pypi.org",
			PurlType:              "pypi",
			RegistryUrl:           "https://pypi.org/simple",
			Tools:                 []string{"pip", "poetry", "pipenv", "conda"},
		},
	}
}

// DetectedLanguage represents a detected programming language in a project
type DetectedLanguage struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Color string `json:"color"`
}

// EcosystemHandler interface defines the contract for handling different package ecosystems
type EcosystemHandler interface {
	// GetLanguageID returns the language identifier (e.g., "JS", "PHP")
	GetLanguageID() string

	// GetEcosystemInfo returns the ecosystem information
	GetEcosystemInfo() EcosystemInfo

	// ProcessLicenses processes license analysis for this ecosystem
	// Using interface{} for now to avoid dependency issues, will be properly typed later
	ProcessLicenses(knowledgeDB *bun.DB, sbom interface{}, licensePolicy interface{}, start time.Time) interface{}

	// ProcessVulnerabilities processes vulnerability analysis for this ecosystem
	// Using interface{} for now to avoid dependency issues, will be properly typed later
	ProcessVulnerabilities(projectURL string, knowledgeDB *bun.DB, sbom interface{}, start time.Time) interface{}

	// SupportsLanguageID checks if this handler supports the given language ID
	SupportsLanguageID(languageID string) bool
}

// ResultMerger provides generic functionality for merging analysis results from multiple ecosystems
type ResultMerger[T any] interface {
	// MergeWorkspaces merges results from multiple language ecosystems into a unified result
	MergeWorkspaces(results []T) T

	// GetMergeStrategy returns the merge strategy used by this merger
	GetMergeStrategy() string
}

// MergeStrategy defines different strategies for merging multi-language results
type MergeStrategy string

const (
	MergeStrategyUnion        MergeStrategy = "union"        // Combine all results
	MergeStrategyIntersection MergeStrategy = "intersection" // Only results present in all languages
	MergeStrategyPriority     MergeStrategy = "priority"     // Prioritize results from specific languages
)

// EcosystemMapper provides utilities for mapping between plugins, ecosystems, and languages
type EcosystemMapper struct {
	ecosystemMap PluginEcosystemMap
}

// NewEcosystemMapper creates a new EcosystemMapper with the default ecosystem mapping
func NewEcosystemMapper() *EcosystemMapper {
	return &EcosystemMapper{
		ecosystemMap: GetDefaultEcosystemMap(),
	}
}

// NewEcosystemMapperWithCustomMap creates a new EcosystemMapper with a custom ecosystem mapping
func NewEcosystemMapperWithCustomMap(ecosystemMap PluginEcosystemMap) *EcosystemMapper {
	return &EcosystemMapper{
		ecosystemMap: ecosystemMap,
	}
}

// GetEcosystemInfo returns ecosystem information for a given plugin name
func (em *EcosystemMapper) GetEcosystemInfo(pluginName string) (EcosystemInfo, bool) {
	info, exists := em.ecosystemMap[pluginName]
	return info, exists
}

// GetSupportedPlugins returns all supported SBOM plugin names
func (em *EcosystemMapper) GetSupportedPlugins() []string {
	plugins := make([]string, 0, len(em.ecosystemMap))
	for plugin := range em.ecosystemMap {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// GetSupportedEcosystems returns all supported ecosystem names
func (em *EcosystemMapper) GetSupportedEcosystems() []string {
	ecosystems := make([]string, 0, len(em.ecosystemMap))
	for _, info := range em.ecosystemMap {
		ecosystems = append(ecosystems, info.Ecosystem)
	}
	return ecosystems
}

// MapPackageManagerToEcosystem maps a package manager string to an ecosystem
func (em *EcosystemMapper) MapPackageManagerToEcosystem(packageManager string) (string, bool) {
	for _, info := range em.ecosystemMap {
		matched, _ := regexp.MatchString(info.PackageManagerPattern, packageManager)
		if matched {
			return info.Ecosystem, true
		}
	}
	return "", false
}

// IsValidEcosystem checks if an ecosystem filter is supported
func (em *EcosystemMapper) IsValidEcosystem(ecosystem string) bool {
	for _, info := range em.ecosystemMap {
		if info.Ecosystem == ecosystem {
			return true
		}
	}
	return false
}

// GetEcosystemFromPurl extracts ecosystem from Package URL (PURL)
func (em *EcosystemMapper) GetEcosystemFromPurl(purl string) (string, bool) {
	if purl == "" || len(purl) < 5 || purl[:4] != "pkg:" {
		return "", false
	}

	// Split by '/' to get type
	parts := regexp.MustCompile(`[/:]`).Split(purl, -1)
	if len(parts) < 2 {
		return "", false
	}

	purlType := parts[1] // parts[0] is "pkg", parts[1] is the type

	// Map PURL types to our ecosystem names
	purlToEcosystem := map[string]string{
		"npm":      "npm",
		"composer": "packagist",
		"pypi":     "pypi",
		"cargo":    "cargo",
		"maven":    "maven",
		"nuget":    "nuget",
		"golang":   "go",
		"gem":      "rubygems",
	}

	if ecosystem, exists := purlToEcosystem[purlType]; exists {
		return ecosystem, true
	}

	return "", false
}

// DetectEcosystemFromName detects ecosystem from dependency name patterns
func (em *EcosystemMapper) DetectEcosystemFromName(name string) (string, bool) {
	// PHP Composer packages typically have vendor/package format
	if len(name) > 0 && regexp.MustCompile(`^[^@][^/]*/[^/]+$`).MatchString(name) {
		return "packagist", true
	}

	// npm scoped packages start with @
	if len(name) > 0 && name[0] == '@' {
		return "npm", true
	}

	// Go modules typically have domain/path format
	if regexp.MustCompile(`^[^/]+\.[^/]+/[^/]`).MatchString(name) {
		return "go", true
	}

	return "", false
}
