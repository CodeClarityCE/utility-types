package knowledge

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"github.com/uptrace/bun"
)

type LicenseList struct {
	LicenseListVersion string    `json:"licenseListVersion"`
	Licenses           []License `json:"licenses"`
	Date               string    `json:"releaseDate"`
}

type License struct {
	bun.BaseModel         `bun:"table:licenses,alias:l"`
	Id                    uuid.UUID `bun:",pk,autoincrement,type:uuid,default:uuid_generate_v4()"`
	Reference             string    `bun:"reference"`
	IsDeprecatedLicenseID bool      `bun:"isDeprecatedLicenseId"`
	DetailsURL            string    `bun:"detailsUrl"`
	Details               Details   `bun:"details"`
	ReferenceNumber       int       `bun:"referenceNumber"`
	Name                  string    `bun:"name"`
	LicenseID             string    `bun:"licenseId"`
	SeeAlso               []string  `bun:"seeAlso"`
	IsOsiApproved         bool      `bun:"isOsiApproved"`
}

type CrossRef struct {
	IsLive        bool      `json:"isLive"`
	IsValid       bool      `json:"isValid"`
	IsWayBackLink bool      `json:"isWayBackLink"`
	Match         string    `json:"match"`
	Order         int       `json:"order"`
	Timestamp     time.Time `json:"timestamp"`
	URL           string    `json:"url"`
}

type Details struct {
	CrossRef                    []CrossRef        `json:"crossRef"`
	IsDeprecatedLicenseID       bool              `json:"isDeprecatedLicenseId"`
	IsOsiApproved               bool              `json:"isOsiApproved"`
	LicenseID                   string            `json:"licenseId"`
	LicenseText                 string            `json:"licenseText"`
	LicenseTextHTML             string            `json:"licenseTextHtml"`
	LicenseTextNormalized       string            `json:"licenseTextNormalized"`
	LicenseTextNormalizedDigest string            `json:"licenseTextNormalizedDigest"`
	Name                        string            `json:"name"`
	SeeAlso                     []string          `json:"seeAlso"`
	StandardLicenseTemplate     string            `json:"standardLicenseTemplate"`
	Description                 string            `json:"description"`
	Classification              string            `json:"classification"`
	LicenseProperties           LicenseProperties `json:"licenseProperties"`
}

type LicenseProperties struct {
	Permissions []string `json:"permissions"`
	Conditions  []string `json:"conditions"`
	Limitations []string `json:"limitations"`
}

type LinkLicensePackage struct {
	FromKey    string `json:"packageKey"`
	LicenseKey string `json:"licenseKey"`
}

type LicensePolicy struct {
	DisallowedLicense []string `json:"disallowed_licenses"`
}

func GetDetails(licenses []License) []License {
	var wg sync.WaitGroup
	maxGoroutines := 100
	guard := make(chan struct{}, maxGoroutines)
	// Configure progression bar
	var length int64 = int64(len(licenses))
	bar := progressbar.Default(length)

	for key := range licenses {
		wg.Add(1)
		guard <- struct{}{}
		go func(wg *sync.WaitGroup, key int) {
			defer wg.Done()
			defer bar.Add(1)
			url := licenses[key].DetailsURL
			details, err := getBasicDetails(url)
			if err != nil {
				return
			}
			licenses[key].Details = details
			<-guard
		}(&wg, key)
	}
	wg.Wait()
	return licenses
}

func getBasicDetails(url string) (Details, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("No response from request")
		return Details{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte
	if err != nil {
		log.Println("Error reading body")
		return Details{}, err
	}

	var result Details
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		log.Println("Can not unmarshal JSON", url)
		return Details{}, err
	}
	return result, nil
}
