package main

import (
	"log/slog"
	"time"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

// SvcCat contains methods for operating with the Service Catalog, e.g. Backstage API.
type SvcCat interface {
	ReadSvc() (string, error)
}

// SvcConfig is the Client Configuration
type SvcConfig struct {
	URL      string // URL is the Backstage API endpoint
	Service  string // Each Service is known as the "Component" in Backstage
	Datetime int64  // Unix Epoch in seconds
	Owner    string // Should equal CODEOWNERS for this repo in GitHub
}

// ReadSvc can query Backstage for a chunk of data about a System,
// i.e. the "top-level" Weedmaps Service.
// Each method called for filling in data adds the entry to the SvcConfig struct.
func (sc *SvcConfig) ReadSvc() (string, error) {
	sc.Datetime = time.Now().Unix()
	c, _ := backstage.NewClient(sc.URL, "default", nil)

	// We only want the owner here, not the entire system struct
	owner, _, err := ReadSystemBS(sc.Service, c)
	sc.Owner = owner
	slog.Debug("Owner Set", slog.String("Owner", sc.Owner))
	return sc.Owner, err
}

// ReadinessRead is the function that tests this service for Production Readiness
func ReadinessRead(i SvcCat) (string, error) {
	// Calling ReadSvc() initiates the source data struct, SvcConfig
	// Currently only returning the Owner, which is what ReadSvc() returns
	return i.ReadSvc()
}
