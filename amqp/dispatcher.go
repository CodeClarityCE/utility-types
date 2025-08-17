package amqp

import "github.com/google/uuid"

type Date struct {
	Date         string `json:"date"`
	TimezoneType string `json:"timezonetype"`
}

// Symfony -> Dispatcher
type SymfonyDispatcherMessage struct {
	Analysis           int      `json:"analysis"`
	Analyzers          []string `json:"analyzers"`
	Date               Date     `json:"date"`
	Uid                int      `json:"uid"`
	Project            string   `json:"project"`
	DisallowedLicenses []string `json:"disallowed_licenses"`
}

type SymfonyDispatcherContext struct {
	Context SymfonyDispatcherMessage `json:"context"`
}

// Sbom -> Dispatcher
type PluginDispatcherMessage struct {
	AnalysisId uuid.UUID `json:"analysis_id"`
	Plugin     string    `json:"sbom"`
}

// Sbom -> PackageFollower
type SbomPackageFollowerMessage struct {
	AnalysisId    uuid.UUID `json:"analysis_id"`
	PackagesNames []string  `json:"package_name"`
	Language      string    `json:"language"` // "javascript", "php", etc.
}

// API -> Dispatcher
type APIDispatcherMessage struct {
	AnalysisId     uuid.UUID         `json:"analysis_id"`
	ProjectId      uuid.UUID         `json:"project_id"`
	IntegrationId  uuid.UUID         `json:"integration_id"`
	OrganizationId uuid.UUID         `json:"organization_id"`
	Config         map[string]Config `json:"config"`
}

// Dispatcher -> Downloader
type DispatcherDownloaderMessage struct {
	AnalysisId     uuid.UUID `json:"analysis_id"`
	ProjectId      uuid.UUID `json:"project_id"`
	IntegrationId  uuid.UUID `json:"integration_id"`
	OrganizationId uuid.UUID `json:"organization_id"`
}

// Downloader -> Dispatcer
// TODO change type
type DownloaderDispatcherMessage struct {
	AnalysisId     uuid.UUID `json:"analysis_id"`
	ProjectId      uuid.UUID `json:"project_id"`
	IntegrationId  uuid.UUID `json:"integration_id"`
	OrganizationId uuid.UUID `json:"organization_id"`
}

type Config struct {
	Data map[string]any `json:"data"`
}

type DispatcherPluginMessage struct {
	Data           any       `json:"data"`
	AnalysisId     uuid.UUID `json:"analysis_id"`
	OrganizationId uuid.UUID `json:"organization_id"`
}
