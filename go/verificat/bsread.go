package main

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

// SvcCatalog. Methods for operating with the Service Catalog, e.g. Backstage API.
type SvcCat interface {
	ReadSvc() (string, error)
}

// Client Configuration
type SvcConfig struct {
	URL        string // URL is the Backstage API endpoint
	Service    string // Each Service is known as the "Component" in Backstage
	Datetime   int64  // Unix Epoch in seconds
	Owner      string // Should equal CODEOWNERS for this repo in GitHub
	LabelKey   string
	LabelValue string
}

// TODO: Load a list of BSSystems locally to check against to sanitize input
// If the desired test subject is not on the list, verificat can not run.
// This will also help with other kinds of validation, like "is it a string".
// So this is a new function.

// ReadSvc. Query Backstage for a chunk of data about a Component, i.e. WM Service.
// This is the primary Client interface for the Backstage API.
// Each method called for filling in data adds the entry to the SvcConfig struct.
func (sc *SvcConfig) ReadSvc() (string, error) {
	sc.Datetime = time.Now().Unix()
	c, _ := backstage.NewClient(sc.URL, "default", nil)

	// We only want the owner here, not the entire system struct
	owner, _, err := ReadSystemBS(sc.Service, c)
	if err != nil {
		slog.Error("Failed to fetch Owner from Backstage", slog.Any("Error", err))
	}

	// TODO: if owner comes back empty... stop?
	sc.Owner = owner

	return owner, err
}

// Test this service for Production Readiness
func ReadinessRead(i SvcCat) (string, error) {
	// Calling ReadSvc() initiates the source data struct, SvcConfig
	found, err := i.ReadSvc()
	if err != nil {
		fmt.Println("error: ", err)
		log.Fatal(err)
	}

	// Currently only returning the Owner, which is what ReadSvc() returns
	return found, err
}
