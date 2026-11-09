package knowledge

import (
	"testing"
	"time"
)

func at(value string) *time.Time {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}
	return &parsed
}

func text(value string) *string { return &value }

func bySource(snapshot KnowledgeSnapshot) map[string]SourceSnapshot {
	out := map[string]SourceSnapshot{}
	for _, source := range snapshot.Sources {
		out[source.Source] = source
	}
	return out
}

func TestTakeSnapshotRecordsEverySourceWithWatermarkAndRowCount(t *testing.T) {
	now := *at("2026-11-09T12:00:00Z")
	rows := []SourceSync{
		{Source: SourceOSV, Watermark: at("2026-11-09T06:00:00Z"), LastOutcome: text(SourceSyncSuccess), RowCount: 238664},
	}

	snapshot := TakeSnapshot(rows, now, nil)

	if !snapshot.TakenAt.Equal(now) {
		t.Fatalf("taken_at = %v, want %v", snapshot.TakenAt, now)
	}
	got := bySource(snapshot)
	if len(got) != 7 {
		t.Fatalf("snapshot has %d sources, want all 7: %v", len(got), snapshot.Sources)
	}
	osv := got[SourceOSV]
	if osv.Status != SourceHealthy || !osv.Comparable || osv.RowCount != 238664 || !osv.Watermark.Equal(*at("2026-11-09T06:00:00Z")) {
		t.Fatalf("osv = %+v, want healthy, comparable, 238664 rows, watermark kept", osv)
	}
	if osv.StaleAfterHours != 24 {
		t.Fatalf("osv stale_after_hours = %v, want 24", osv.StaleAfterHours)
	}
	nvd := got[SourceNVD]
	if nvd.Status != SourceStale || nvd.Comparable || nvd.Watermark != nil || nvd.RowCount != 0 {
		t.Fatalf("source with no sync row = %+v, want stale, not comparable, no watermark, 0 rows", nvd)
	}
}

func TestTakeSnapshotMarksStaleAndFailedSourcesNotComparable(t *testing.T) {
	now := *at("2026-11-09T12:00:00Z")
	rows := []SourceSync{
		// 49 hours old against a 48 hour threshold.
		{Source: SourceNVD, Watermark: at("2026-11-07T11:00:00Z"), LastOutcome: text(SourceSyncSuccess), RowCount: 10},
		// 47 hours old: still fresh.
		{Source: SourceGCVE, Watermark: at("2026-11-07T13:00:00Z"), LastOutcome: text(SourceSyncSuccess), RowCount: 10},
		{Source: SourceOSV, Watermark: at("2026-11-09T11:00:00Z"), LastOutcome: text(SourceSyncFailure), RowCount: 10},
		{Source: SourceEPSS, Watermark: at("2026-11-09T11:00:00Z"), LastOutcome: text(SourceSyncSuccess), ValidationError: text("row count dropped"), RowCount: 10},
		// A backfilled watermark with no recorded run proves nothing.
		{Source: SourceCWE, Watermark: at("2026-11-09T11:00:00Z"), RowCount: 10},
	}

	got := bySource(TakeSnapshot(rows, now, nil))

	want := map[string]string{
		SourceNVD:  SourceStale,
		SourceGCVE: SourceHealthy,
		SourceOSV:  SourceFailed,
		SourceEPSS: SourceFailed,
		SourceCWE:  SourceStale,
	}
	for source, status := range want {
		if got[source].Status != status {
			t.Errorf("%s status = %s, want %s", source, got[source].Status, status)
		}
		if got[source].Comparable != (status == SourceHealthy) {
			t.Errorf("%s comparable = %v with status %s", source, got[source].Comparable, status)
		}
	}
}

func TestTakeSnapshotJudgesFreshnessAtTheKnowledgeAsOfDay(t *testing.T) {
	now := *at("2026-11-09T12:00:00Z")
	asOf := *at("2026-06-01T00:00:00Z")
	rows := []SourceSync{
		// Months old today, but fresh at the end of the as-of day.
		{Source: SourceOSV, Watermark: at("2026-06-01T20:00:00Z"), LastOutcome: text(SourceSyncSuccess)},
		// 72 hours before the end of the as-of day, over the 48 hour threshold.
		{Source: SourceNVD, Watermark: at("2026-05-30T00:00:00Z"), LastOutcome: text(SourceSyncSuccess)},
	}

	snapshot := TakeSnapshot(rows, now, &asOf)

	if snapshot.KnowledgeAsOf == nil || *snapshot.KnowledgeAsOf != "2026-06-01" {
		t.Fatalf("knowledge_asof = %v, want 2026-06-01", snapshot.KnowledgeAsOf)
	}
	got := bySource(snapshot)
	if got[SourceOSV].Status != SourceHealthy {
		t.Errorf("osv = %s, want healthy when fresh at the as-of day", got[SourceOSV].Status)
	}
	if got[SourceNVD].Status != SourceStale {
		t.Errorf("nvd = %s, want stale", got[SourceNVD].Status)
	}
}

func TestStaleAfterHoursReadsTheSameOverrideAsTheAPI(t *testing.T) {
	t.Setenv("KNOWLEDGE_STALE_AFTER_HOURS_OSV", "6")
	t.Setenv("KNOWLEDGE_STALE_AFTER_HOURS_NVD", "not a number")

	if got := StaleAfterHours(SourceOSV); got != 6 {
		t.Errorf("osv = %v, want the override 6", got)
	}
	if got := StaleAfterHours(SourceNVD); got != 48 {
		t.Errorf("nvd = %v, want the default 48 for an invalid override", got)
	}
	if got := StaleAfterHours(SourceCWE); got != 168 {
		t.Errorf("cwe = %v, want 168", got)
	}
}

func TestUnreadableSnapshotMarksEverySourceUnknown(t *testing.T) {
	now := *at("2026-11-09T12:00:00Z")

	snapshot := UnknownSnapshot(now, nil)

	if len(snapshot.Sources) != 7 {
		t.Fatalf("snapshot has %d sources, want 7", len(snapshot.Sources))
	}
	for _, source := range snapshot.Sources {
		if source.Status != SourceUnknown || source.Comparable {
			t.Errorf("%s = %+v, want unknown and not comparable", source.Source, source)
		}
	}
}
