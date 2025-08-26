package ecosystem

import (
	"testing"
	"time"
)

func TestEcosystemMapper(t *testing.T) {
	mapper := NewEcosystemMapper()

	// Test getting ecosystem info for JavaScript
	jsInfo, exists := mapper.GetEcosystemInfo("js-sbom")
	if !exists {
		t.Error("Expected js-sbom to exist in ecosystem map")
	}
	if jsInfo.Language != "JavaScript" {
		t.Errorf("Expected JavaScript language, got %s", jsInfo.Language)
	}
	if jsInfo.Ecosystem != "npm" {
		t.Errorf("Expected npm ecosystem, got %s", jsInfo.Ecosystem)
	}

	// Test getting ecosystem info for PHP
	phpInfo, exists := mapper.GetEcosystemInfo("php-sbom")
	if !exists {
		t.Error("Expected php-sbom to exist in ecosystem map")
	}
	if phpInfo.Language != "PHP" {
		t.Errorf("Expected PHP language, got %s", phpInfo.Language)
	}
	if phpInfo.Ecosystem != "packagist" {
		t.Errorf("Expected packagist ecosystem, got %s", phpInfo.Ecosystem)
	}

	// Test getting supported plugins
	plugins := mapper.GetSupportedPlugins()
	if len(plugins) < 2 {
		t.Errorf("Expected at least 2 plugins, got %d", len(plugins))
	}

	// Test package manager mapping
	ecosystem, found := mapper.MapPackageManagerToEcosystem("npm")
	if !found {
		t.Error("Expected to find ecosystem for npm")
	}
	if ecosystem != "npm" {
		t.Errorf("Expected npm ecosystem, got %s", ecosystem)
	}

	// Test ecosystem validation
	if !mapper.IsValidEcosystem("npm") {
		t.Error("Expected npm to be a valid ecosystem")
	}
	if !mapper.IsValidEcosystem("packagist") {
		t.Error("Expected packagist to be a valid ecosystem")
	}
	if mapper.IsValidEcosystem("invalid") {
		t.Error("Expected invalid to not be a valid ecosystem")
	}
}

func TestHandlerRegistry(t *testing.T) {
	registry := NewHandlerRegistry()

	// Create test handlers
	jsHandler := NewBasicEcosystemHandler("JS", EcosystemInfo{
		Name:     "JavaScript Test",
		Language: "JavaScript",
	})
	phpHandler := NewBasicEcosystemHandler("PHP", EcosystemInfo{
		Name:     "PHP Test",
		Language: "PHP",
	})

	// Register handlers
	registry.RegisterHandler("JS", jsHandler)
	registry.RegisterHandler("PHP", phpHandler)

	// Test getting handlers
	retrievedJS, exists := registry.GetHandler("JS")
	if !exists {
		t.Error("Expected to find JS handler")
	}
	if retrievedJS.GetLanguageID() != "JS" {
		t.Errorf("Expected JS language ID, got %s", retrievedJS.GetLanguageID())
	}

	retrievedPHP, exists := registry.GetHandler("PHP")
	if !exists {
		t.Error("Expected to find PHP handler")
	}
	if retrievedPHP.GetLanguageID() != "PHP" {
		t.Errorf("Expected PHP language ID, got %s", retrievedPHP.GetLanguageID())
	}

	// Test getting all handlers
	allHandlers := registry.GetAllHandlers()
	if len(allHandlers) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(allHandlers))
	}

	// Test getting supported languages
	languages := registry.GetSupportedLanguages()
	if len(languages) != 2 {
		t.Errorf("Expected 2 languages, got %d", len(languages))
	}
}

func TestBasicEcosystemHandler(t *testing.T) {
	ecosystemInfo := EcosystemInfo{
		Name:      "Test Ecosystem",
		Language:  "TestLang",
		Ecosystem: "test",
	}

	handler := NewBasicEcosystemHandler("TEST", ecosystemInfo)

	// Test basic functionality
	if handler.GetLanguageID() != "TEST" {
		t.Errorf("Expected TEST language ID, got %s", handler.GetLanguageID())
	}

	info := handler.GetEcosystemInfo()
	if info.Name != "Test Ecosystem" {
		t.Errorf("Expected Test Ecosystem name, got %s", info.Name)
	}

	// Test language support
	if !handler.SupportsLanguageID("TEST") {
		t.Error("Expected handler to support TEST language")
	}
	if handler.SupportsLanguageID("OTHER") {
		t.Error("Expected handler not to support OTHER language")
	}

	// Test basic implementations (should return not implemented responses)
	startTime := time.Now()
	licenseResult := handler.ProcessLicenses(nil, nil, nil, startTime)
	if licenseResult == nil {
		t.Error("Expected license result to not be nil")
	}

	vulnResult := handler.ProcessVulnerabilities("", nil, nil, startTime)
	if vulnResult == nil {
		t.Error("Expected vulnerability result to not be nil")
	}
}

func TestCreateDefaultHandlers(t *testing.T) {
	registry := CreateDefaultHandlers()

	// Test that default handlers are created
	languages := registry.GetSupportedLanguages()
	if len(languages) < 1 {
		t.Error("Expected at least one default handler to be created")
	}

	// Test that JS handler exists
	jsHandler, exists := registry.GetHandler("JS")
	if !exists {
		t.Error("Expected JS handler to be created by default")
	}
	if jsHandler.GetEcosystemInfo().Language != "JavaScript" {
		t.Error("Expected JS handler to have JavaScript language")
	}

	// Test that PHP handler exists
	phpHandler, exists := registry.GetHandler("PHP")
	if !exists {
		t.Error("Expected PHP handler to be created by default")
	}
	if phpHandler.GetEcosystemInfo().Language != "PHP" {
		t.Error("Expected PHP handler to have PHP language")
	}
}
