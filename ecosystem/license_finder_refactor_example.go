package ecosystem

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// This file demonstrates how the license-finder plugin could be refactored
// to use the new ecosystem handler pattern, eliminating the complex manual
// workspace merging logic (~140 lines) in favor of the generic result merger

// AnalysisDocument represents the analysis document structure
type AnalysisDocument struct {
	Stage  int                      `json:"stage"`
	Steps  [][]AnalysisStep        `json:"steps"`
	Config map[string]interface{}  `json:"config"`
}

// AnalysisStep represents a step in the analysis
type AnalysisStep struct {
	Name   string                 `json:"name"`
	Result map[string]interface{} `json:"result"`
}

// Result represents a database result
type Result struct {
	Id     uuid.UUID   `json:"id"`
	Result interface{} `json:"result"`
}

// RefactoredLicenseAnalyzer demonstrates the simplified approach
type RefactoredLicenseAnalyzer struct {
	registry    *HandlerRegistry
	codeclarity *bun.DB
	knowledge   *bun.DB
	merger      *LicenseResultMerger
}

// NewRefactoredLicenseAnalyzer creates a new refactored license analyzer
func NewRefactoredLicenseAnalyzer(codeclarity, knowledge *bun.DB) *RefactoredLicenseAnalyzer {
	return &RefactoredLicenseAnalyzer{
		registry:    CreateDefaultHandlers(),
		codeclarity: codeclarity,
		knowledge:   knowledge,
		merger:      NewLicenseResultMerger(MergeStrategyUnion),
	}
}

// AnalyzeLicenses performs license analysis using the new pattern
func (rla *RefactoredLicenseAnalyzer) AnalyzeLicenses(analysisDoc AnalysisDocument, licensePolicy interface{}) (LicenseOutput, error) {
	start := time.Now()

	// Step 1: Discover available SBOM results using generic discovery
	sbomResults, err := rla.discoverSBOMResults(analysisDoc)
	if err != nil {
		return LicenseOutput{AnalysisStatus: "FAILURE"}, err
	}

	if len(sbomResults) == 0 {
		log.Printf("No SBOM results found, returning empty license analysis")
		return LicenseOutput{
			WorkSpaces:     make(map[string]LicenseWorkspaceInfo),
			AnalysisStats:  LicenseAnalysisStats{},
			AnalysisStatus: "SUCCESS",
		}, nil
	}

	// Step 2: Process each SBOM using appropriate ecosystem handlers
	var licenseOutputs []LicenseOutput
	for _, sbomResult := range sbomResults {
		output, err := rla.processSBOMForLicenses(sbomResult, licensePolicy, start)
		if err != nil {
			log.Printf("Error processing %s SBOM: %v", sbomResult.Language, err)
			continue
		}
		licenseOutputs = append(licenseOutputs, output)
	}

	// Step 3: Merge results using the generic result merger (replaces 140+ lines!)
	mergedResult := rla.merger.MergeWorkspaces(licenseOutputs)

	log.Printf("License analysis completed: processed %d ecosystems, merged into %d workspaces",
		len(licenseOutputs), len(mergedResult.WorkSpaces))

	return mergedResult, nil
}

// SBOMResult represents discovered SBOM data
type SBOMResult struct {
	ID         uuid.UUID
	Language   string
	PluginName string
	Data       interface{} // The actual SBOM data
}

// discoverSBOMResults discovers all available SBOM results from previous analysis stages
func (rla *RefactoredLicenseAnalyzer) discoverSBOMResults(analysisDoc AnalysisDocument) ([]SBOMResult, error) {
	var results []SBOMResult

	// Look at the previous stage
	if analysisDoc.Stage <= 0 {
		return results, nil
	}

	previousStage := analysisDoc.Stage - 1
	if previousStage >= len(analysisDoc.Steps) {
		return results, fmt.Errorf("invalid previous stage index: %d", previousStage)
	}

	// Discover all SBOM plugins from the ecosystem registry
	mapper := NewEcosystemMapper()
	supportedPlugins := mapper.GetSupportedPlugins()

	for _, step := range analysisDoc.Steps[previousStage] {
		// Check if this step is a supported SBOM plugin
		for _, pluginName := range supportedPlugins {
			if step.Name == pluginName {
				ecosystemInfo, exists := mapper.GetEcosystemInfo(pluginName)
				if !exists {
					continue
				}

				// Get the SBOM result ID
				sbomKeyInterface, exists := step.Result["sbomKey"]
				if !exists {
					log.Printf("No sbomKey found for %s", pluginName)
					continue
				}

				sbomKeyStr, ok := sbomKeyInterface.(string)
				if !ok {
					log.Printf("Invalid sbomKey type for %s", pluginName)
					continue
				}

				sbomKeyUUID, err := uuid.Parse(sbomKeyStr)
				if err != nil {
					log.Printf("Invalid sbomKey UUID for %s: %v", pluginName, err)
					continue
				}

				// Fetch the SBOM data
				dbResult := Result{Id: sbomKeyUUID}
				err = rla.codeclarity.NewSelect().Model(&dbResult).Where("id = ?", sbomKeyUUID).Scan(context.Background())
				if err != nil {
					log.Printf("Failed to retrieve SBOM for %s: %v", pluginName, err)
					continue
				}

				// Parse SBOM data (simplified - in real implementation would use proper types)
				var sbomData interface{}
				err = json.Unmarshal(dbResult.Result.([]byte), &sbomData)
				if err != nil {
					log.Printf("Failed to unmarshal SBOM for %s: %v", pluginName, err)
					continue
				}

				results = append(results, SBOMResult{
					ID:         sbomKeyUUID,
					Language:   mapPluginToLanguageID(pluginName, ecosystemInfo),
					PluginName: pluginName,
					Data:       sbomData,
				})

				log.Printf("Discovered %s SBOM result: %s", ecosystemInfo.Language, sbomKeyUUID)
				break
			}
		}
	}

	return results, nil
}

// processSBOMForLicenses processes a single SBOM for license analysis using ecosystem handlers
func (rla *RefactoredLicenseAnalyzer) processSBOMForLicenses(sbomResult SBOMResult, licensePolicy interface{}, start time.Time) (LicenseOutput, error) {
	// Get the appropriate handler for this language
	handler, exists := rla.registry.GetHandler(sbomResult.Language)
	if !exists {
		return LicenseOutput{AnalysisStatus: "FAILURE"}, fmt.Errorf("no handler found for language: %s", sbomResult.Language)
	}

	log.Printf("Processing %s SBOM using %s handler", sbomResult.Language, handler.GetEcosystemInfo().Name)

	// Use the handler to process licenses
	// Note: In a real implementation, this would return properly typed results
	// Here we simulate the conversion to our generic license output format
	handlerResult := handler.ProcessLicenses(rla.knowledge, sbomResult.Data, licensePolicy, start)

	// Convert handler result to our generic format
	// This is where the actual conversion logic would go
	licenseOutput := convertHandlerResultToLicenseOutput(handlerResult, sbomResult.Language)

	return licenseOutput, nil
}

// Helper functions for the refactored implementation

// mapPluginToLanguageID maps a plugin name to the language ID expected by handlers
func mapPluginToLanguageID(pluginName string, ecosystemInfo EcosystemInfo) string {
	switch pluginName {
	case "js-sbom":
		return "JS"
	case "php-sbom":
		return "PHP"
	default:
		// For future plugins, use the ecosystem info to determine language ID
		return ecosystemInfo.Language
	}
}

// convertHandlerResultToLicenseOutput converts handler results to our generic license output
// This is a simplified version - real implementation would handle proper type conversions
func convertHandlerResultToLicenseOutput(handlerResult interface{}, language string) LicenseOutput {
	// In a real implementation, this would properly convert the handler result
	// For this example, we return a mock result to demonstrate the pattern
	return LicenseOutput{
		WorkSpaces: map[string]LicenseWorkspaceInfo{
			"default": {
				LicensesDepMap:              make(map[string][]string),
				NonSpdxLicensesDepMap:       make(map[string][]string),
				LicenseComplianceViolations: []string{},
				DependencyInfo:              make(map[string]LicenseDependencyInfo),
			},
		},
		AnalysisStats: LicenseAnalysisStats{
			NumberOfSpdxLicenses:       0,
			NumberOfNonSpdxLicenses:    0,
			NumberOfCopyLeftLicenses:   0,
			NumberOfPermissiveLicenses: 0,
			LicenseDist:                make(map[string]int),
		},
		AnalysisStatus: "SUCCESS",
	}
}

/*
BENEFITS OF THIS REFACTORED APPROACH:

1. ELIMINATED COMPLEXITY:
   - Removed 140+ lines of manual workspace merging logic
   - Replaced hardcoded plugin discovery with registry-based discovery
   - Eliminated language-specific if/else chains

2. IMPROVED MAINTAINABILITY:
   - Single place to add new language support (ecosystem registry)
   - Generic result merger can be reused for vulnerability analysis
   - Handlers encapsulate ecosystem-specific logic

3. BETTER TESTING:
   - Each component can be tested independently
   - Mock handlers can be easily created for testing
   - Result merger has its own comprehensive test suite

4. EXTENSIBILITY:
   - Adding Python support just requires registering a handler
   - Different merge strategies can be easily swapped
   - Plugin discovery automatically picks up new SBOM plugins

5. TYPE SAFETY:
   - Shared type definitions ensure consistency
   - Compile-time checking prevents API contract bugs
   - Auto-generated TypeScript maintains frontend sync

LINES OF CODE IMPACT:
- Original license-finder main.go: ~305 lines
- Refactored approach: ~150 lines of business logic + shared utilities
- Net reduction: ~40-50% with improved maintainability
*/