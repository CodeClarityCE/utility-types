package ecosystem

import (
	"github.com/uptrace/bun"
	"time"
)

// BasicEcosystemHandler provides a basic implementation of EcosystemHandler
// that can be extended for specific ecosystems
type BasicEcosystemHandler struct {
	languageID    string
	ecosystemInfo EcosystemInfo
}

// NewBasicEcosystemHandler creates a new basic ecosystem handler
func NewBasicEcosystemHandler(languageID string, ecosystemInfo EcosystemInfo) *BasicEcosystemHandler {
	return &BasicEcosystemHandler{
		languageID:    languageID,
		ecosystemInfo: ecosystemInfo,
	}
}

// GetLanguageID returns the language identifier
func (h *BasicEcosystemHandler) GetLanguageID() string {
	return h.languageID
}

// GetEcosystemInfo returns the ecosystem information
func (h *BasicEcosystemHandler) GetEcosystemInfo() EcosystemInfo {
	return h.ecosystemInfo
}

// ProcessLicenses processes license analysis for this ecosystem
// This is a basic implementation that should be overridden by specific handlers
func (h *BasicEcosystemHandler) ProcessLicenses(knowledgeDB *bun.DB, sbom interface{}, licensePolicy interface{}, start time.Time) interface{} {
	// Basic implementation - should be overridden by specific ecosystem handlers
	return map[string]interface{}{
		"status":  "not_implemented",
		"message": "License processing not implemented for " + h.languageID,
	}
}

// ProcessVulnerabilities processes vulnerability analysis for this ecosystem
// This is a basic implementation that should be overridden by specific handlers
func (h *BasicEcosystemHandler) ProcessVulnerabilities(projectURL string, knowledgeDB *bun.DB, sbom interface{}, start time.Time) interface{} {
	// Basic implementation - should be overridden by specific ecosystem handlers
	return map[string]interface{}{
		"status":  "not_implemented",
		"message": "Vulnerability processing not implemented for " + h.languageID,
	}
}

// SupportsLanguageID checks if this handler supports the given language ID
func (h *BasicEcosystemHandler) SupportsLanguageID(languageID string) bool {
	return h.languageID == languageID
}

// HandlerRegistry manages ecosystem handlers
type HandlerRegistry struct {
	handlers map[string]EcosystemHandler
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]EcosystemHandler),
	}
}

// RegisterHandler registers an ecosystem handler for a language
func (r *HandlerRegistry) RegisterHandler(languageID string, handler EcosystemHandler) {
	r.handlers[languageID] = handler
}

// GetHandler returns the handler for a specific language
func (r *HandlerRegistry) GetHandler(languageID string) (EcosystemHandler, bool) {
	handler, exists := r.handlers[languageID]
	return handler, exists
}

// GetAllHandlers returns all registered handlers
func (r *HandlerRegistry) GetAllHandlers() map[string]EcosystemHandler {
	return r.handlers
}

// GetSupportedLanguages returns all supported language IDs
func (r *HandlerRegistry) GetSupportedLanguages() []string {
	languages := make([]string, 0, len(r.handlers))
	for lang := range r.handlers {
		languages = append(languages, lang)
	}
	return languages
}

// CreateDefaultHandlers creates default handlers for supported ecosystems
func CreateDefaultHandlers() *HandlerRegistry {
	registry := NewHandlerRegistry()
	ecosystemMap := GetDefaultEcosystemMap()

	// Create JavaScript handler
	if jsInfo, exists := ecosystemMap["js-sbom"]; exists {
		jsHandler := NewBasicEcosystemHandler("JS", jsInfo)
		registry.RegisterHandler("JS", jsHandler)
	}

	// Create PHP handler
	if phpInfo, exists := ecosystemMap["php-sbom"]; exists {
		phpHandler := NewBasicEcosystemHandler("PHP", phpInfo)
		registry.RegisterHandler("PHP", phpHandler)
	}

	return registry
}
