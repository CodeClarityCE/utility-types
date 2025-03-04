package knowledge

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Package struct {
	bun.BaseModel `bun:"table:package,alias:p"`
	Id            uuid.UUID      `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	Name          string         `bun:"name"`
	Description   string         `bun:"description"`
	Homepage      string         `bun:"homepage"`
	LatestVersion string         `bun:"latest_version"`
	Versions      []Version      `bun:"versions,rel:has-many,join:id=packageId"`
	Time          time.Time      `bun:"time"`
	Keywords      []string       `bun:"keywords"`
	Source        Source         `bun:"source"`
	License       string         `json:"license"`
	Licenses      []LicenseNpm   `json:"licenses"`
	Extra         map[string]any `bun:"extra"`
}

type LicenseNpm struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Source struct {
	Url  string `json:"Url"`
	Type string `json:"Type"`
}

type Version struct {
	bun.BaseModel   `bun:"table:js_version,alias:pv"`
	Id              uuid.UUID         `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	PackageID       uuid.UUID         `bun:"packageId,type:uuid"`
	Version         string            `bun:"version"`
	Dependencies    map[string]string `bun:"dependencies"`
	DevDependencies map[string]string `bun:"dev_dependencies"`
	Extra           map[string]any    `bun:"extra"`
}

type VersionEdge struct{}
