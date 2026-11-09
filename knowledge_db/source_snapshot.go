package knowledge

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Health of a knowledge source, as reported by the API status endpoint.
const (
	SourceHealthy = "healthy"
	SourceStale   = "stale"
	SourceFailed  = "failed"
	// SourceUnknown means the sync state could not be read.
	SourceUnknown = "unknown"
)

// AllSources lists every knowledge source in a fixed order.
var AllSources = []string{SourceCWE, SourceEPSS, SourceFriendsOfPHP, SourceGCVE, SourceLicenses, SourceNVD, SourceOSV}

// DefaultStaleAfterHours matches the API and the Grafana stale alerts.
var DefaultStaleAfterHours = map[string]float64{
	SourceCWE:          168,
	SourceEPSS:         24,
	SourceFriendsOfPHP: 168,
	SourceGCVE:         48,
	SourceLicenses:     168,
	SourceNVD:          48,
	SourceOSV:          24,
}

// StaleAfterHours reads KNOWLEDGE_STALE_AFTER_HOURS_<SOURCE>, falling back to the default when unset or invalid.
func StaleAfterHours(source string) float64 {
	raw := os.Getenv("KNOWLEDGE_STALE_AFTER_HOURS_" + strings.ToUpper(source))
	if value, err := strconv.ParseFloat(raw, 64); err == nil && value > 0 {
		return value
	}
	return DefaultStaleAfterHours[source]
}

// KnowledgeSnapshot is the state of every knowledge source when an analysis ran.
type KnowledgeSnapshot struct {
	TakenAt time.Time `json:"taken_at"`
	// KnowledgeAsOf is the analysis's as-of day (YYYY-MM-DD), absent for a normal analysis.
	KnowledgeAsOf *string          `json:"knowledge_asof,omitempty"`
	Sources       []SourceSnapshot `json:"sources"`
}

// SourceSnapshot is one knowledge source's state when an analysis ran.
// A source is comparable only when it was healthy.
type SourceSnapshot struct {
	Source          string     `json:"source"`
	Status          string     `json:"status"`
	Comparable      bool       `json:"comparable"`
	Watermark       *time.Time `json:"watermark"`
	RowCount        int64      `json:"row_count"`
	LastOutcome     *string    `json:"last_outcome"`
	StaleAfterHours float64    `json:"stale_after_hours"`
}

// TakeSnapshot records every source from its sync row. With an as-of day, freshness is judged at the end of that day.
func TakeSnapshot(rows []SourceSync, now time.Time, knowledgeAsOf *time.Time) KnowledgeSnapshot {
	reference := now
	if knowledgeAsOf != nil {
		if endOfDay := knowledgeAsOf.Truncate(24 * time.Hour).Add(24 * time.Hour); endOfDay.Before(now) {
			reference = endOfDay
		}
	}
	byName := make(map[string]SourceSync, len(rows))
	for _, row := range rows {
		byName[row.Source] = row
	}
	snapshot := newSnapshot(now, knowledgeAsOf)
	for _, source := range AllSources {
		row := byName[source]
		staleAfter := StaleAfterHours(source)
		status := health(row, staleAfter, reference)
		snapshot.Sources = append(snapshot.Sources, SourceSnapshot{
			Source:          source,
			Status:          status,
			Comparable:      status == SourceHealthy,
			Watermark:       row.Watermark,
			RowCount:        row.RowCount,
			LastOutcome:     row.LastOutcome,
			StaleAfterHours: staleAfter,
		})
	}
	return snapshot
}

// UnknownSnapshot records every source as unknown and not comparable, for when sync state cannot be read.
func UnknownSnapshot(now time.Time, knowledgeAsOf *time.Time) KnowledgeSnapshot {
	snapshot := newSnapshot(now, knowledgeAsOf)
	for _, source := range AllSources {
		snapshot.Sources = append(snapshot.Sources, SourceSnapshot{
			Source:          source,
			Status:          SourceUnknown,
			StaleAfterHours: StaleAfterHours(source),
		})
	}
	return snapshot
}

func newSnapshot(now time.Time, knowledgeAsOf *time.Time) KnowledgeSnapshot {
	snapshot := KnowledgeSnapshot{TakenAt: now, Sources: []SourceSnapshot{}}
	if knowledgeAsOf != nil {
		day := knowledgeAsOf.Format("2006-01-02")
		snapshot.KnowledgeAsOf = &day
	}
	return snapshot
}

// health follows the API: a failed run or validation is failed, and no recorded run or watermark is stale.
func health(row SourceSync, staleAfterHours float64, reference time.Time) string {
	if (row.LastOutcome != nil && *row.LastOutcome == SourceSyncFailure) || row.ValidationError != nil {
		return SourceFailed
	}
	if row.LastOutcome == nil || row.Watermark == nil {
		return SourceStale
	}
	if reference.Sub(*row.Watermark) > time.Duration(staleAfterHours*float64(time.Hour)) {
		return SourceStale
	}
	return SourceHealthy
}
