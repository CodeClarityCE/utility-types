package knowledge

import (
	"encoding/json"
	"testing"
)

const twoConfigurationCVE = `{"vulnerabilities":[{"cve":{"id":"CVE-2022-1011","vulnStatus":"Analyzed","configurations":[
{"nodes":[{"operator":"OR","negate":false,"cpeMatch":[{"vulnerable":true,"criteria":"cpe:2.3:o:linux:linux_kernel:5.0:*:*:*:*:*:*:*"}]}]},
{"operator":"AND","nodes":[
  {"operator":"OR","negate":false,"cpeMatch":[{"vulnerable":true,"criteria":"cpe:2.3:a:vendor:first_app:1.0:*:*:*:*:*:*:*","versionEndExcluding":"2.0"}]},
  {"operator":"OR","negate":false,"cpeMatch":[{"vulnerable":true,"criteria":"cpe:2.3:a:vendor:second_app:*:*:*:*:*:*:*:*"}]}
]}]}}]}`

func parseVulns(t *testing.T, body string) []NVDItem {
	t.Helper()
	var nvd NVD
	if err := json.Unmarshal([]byte(body), &nvd); err != nil {
		t.Fatal(err)
	}
	return GetVulns(nvd)
}

func products(sources []Sources) []string {
	names := make([]string, 0, len(sources))
	for _, source := range sources {
		names = append(names, source.CriteriaDict.Product)
	}
	return names
}

func TestGetVulnsReadsApplicationsFromEveryConfiguration(t *testing.T) {
	item := parseVulns(t, twoConfigurationCVE)[0]

	got := products(item.AffectedFlattened)
	if len(got) != 2 || got[0] != "first_app" || got[1] != "second_app" {
		t.Errorf("flattened products = %v, want first_app and second_app", got)
	}
	if len(item.Affected) != 2 {
		t.Errorf("affected entries = %d, want one per application node", len(item.Affected))
	}
}

func TestGetVulnsKeepsTheConfigurationsAsPublished(t *testing.T) {
	item := parseVulns(t, twoConfigurationCVE)[0]

	if len(item.Configurations) != 2 {
		t.Fatalf("configurations = %d, want 2", len(item.Configurations))
	}
	stored, err := json.Marshal(item.Configurations[1])
	if err != nil {
		t.Fatal(err)
	}
	var roundTrip Configuration
	if err := json.Unmarshal(stored, &roundTrip); err != nil {
		t.Fatal(err)
	}
	if roundTrip.Operator != "AND" || len(roundTrip.Nodes) != 2 {
		t.Errorf("configuration = %+v, want the AND operator and both nodes", roundTrip)
	}
	if got := roundTrip.Nodes[0].CpeMatch[0].VersionEndExcluding; got != "2.0" {
		t.Errorf("versionEndExcluding = %q, want it kept", got)
	}
}

func TestGetVulnsSkipsAMalformedCriteria(t *testing.T) {
	const body = `{"vulnerabilities":[{"cve":{"id":"CVE-1","configurations":[{"nodes":[{"operator":"OR","cpeMatch":[
{"vulnerable":true,"criteria":"cpe:2.3:a:truncated"},
{"vulnerable":true,"criteria":"cpe:2.3:a:vendor:good:1.0:*:*:*:*:*:*:*"}]}]}]}}]}`

	item := parseVulns(t, body)[0]

	if got := products(item.AffectedFlattened); len(got) != 1 || got[0] != "good" {
		t.Errorf("products = %v, want only the well-formed one", got)
	}
}

func TestNVDMetricsReadCVSS4(t *testing.T) {
	var metrics Metrics
	body := `{"cvssMetricV40":[{"source":"nvd@nist.gov","type":"Primary","cvssData":{"version":"4.0","baseScore":9.3,"baseSeverity":"CRITICAL","vectorString":"CVSS:4.0/AV:N/AC:L/AT:N/PR:N/UI:N/VC:H/VI:H/VA:H/SC:N/SI:N/SA:N"}}]}`
	if err := json.Unmarshal([]byte(body), &metrics); err != nil {
		t.Fatal(err)
	}
	if len(metrics.CvssMetricV40) != 1 {
		t.Fatalf("got %d CVSS 4.0 metrics, want 1", len(metrics.CvssMetricV40))
	}
	got := metrics.CvssMetricV40[0]
	if got.Type != "Primary" || got.CvssData.BaseScore != 9.3 || got.CvssData.Version != "4.0" {
		t.Errorf("metric = %+v", got)
	}
}
