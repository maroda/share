package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

// SvcCatalog. Methods for operating with the Service Catalog, e.g. Backstage API.
type SvcCat interface {
	ReadSvc() (string, error)
	wmOwner() string
}

// Client Configuration
type SvcConfig struct {
	Origin     *backstage.ComponentEntityV1alpha1
	URL        string // URL is the Backstage API endpoint
	Service    string // Each Service is known as the "Component" in Backstage
	Datetime   int64  // Unix Epoch in seconds
	Owner      string // Should equal CODEOWNERS for this repo in GitHub
	LabelKey   string
	LabelValue string
}

// ReadSvc. Query Backstage for a chunk of data about a Component, i.e. WM Service.
// This is the primary Client interface for the Backstage API.
// Each method called for filling in data adds the entry to the SvcConfig struct.
func (sc *SvcConfig) ReadSvc() (string, error) {
	sc.Datetime = time.Now().Unix()
	c, _ := backstage.NewClient(sc.URL, "default", nil)

	// component is a pointer, Type: *backstage.ComponentEntityV1alpha1
	component, _, err := c.Catalog.Components.Get(context.Background(), sc.Service, "")
	if err != nil {
		log.Fatal(err)
	} else {
		// Really don't need the description in the struct,
		// so use it here anonymously to display for now.
		log.Printf("Backstage Component: %s ::: %s", sc.Service, component.Metadata.Description)
	}

	// Add the component pointer to the struct
	sc.Origin = component

	// Call the methods for filling in data from Backstage
	o := sc.wmOwner()

	// Next method to add:
	//label := "wm_audit_soc2"
	//sc.wmLabels(component, label)

	// This is only for validation. This method primarily writes data to a struct.
	return o, err
}

// wmOwner. Checks Relations.Type for "ownedBy"
// This is to be matched with CODEOWNERS in GitHub
func (s *SvcConfig) wmOwner() string {
	var ob string
	for _, r := range s.Origin.Relations {
		if r.Type == "ownedBy" {
			ob = r.Target.Name
		}
	}

	// Update struct pointer
	s.Owner = ob

	// Return the value of Owner for now to allow a test on a string.
	return ob
}

// NOT IN USE YET, but the code works
// wmLabels. Should return a label we want.
func (s *SvcConfig) wmLabels(c *backstage.ComponentEntityV1alpha1, lbl string) string {
	var kv string
	// Labels: map
	for j, w := range c.Metadata.Labels {
		if j == lbl {
			// Update test data
			s.LabelKey = j
			s.LabelValue = w
			kv = j + ":" + w
			break
		}
	}

	// Return value for printing
	return kv
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
