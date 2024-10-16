package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

// BSComponent is a local struct for carrying the contents of .Catalog.Components.Get
type BSComponent struct {
	Type           string   `json:"type" yaml:"type"`                                         // Type of component.
	Lifecycle      string   `json:"lifecycle" yaml:"lifecycle"`                               // Lifecycle state of the component.
	Owner          string   `json:"owner" yaml:"owner"`                                       // Owner is an entity reference to the owner of the component.
	SubcomponentOf string   `json:"subcomponentOf,omitempty" yaml:"subcomponentOf,omitempty"` // SubcomponentOf is an entity reference to another component of which the component is a part.
	ProvidesApis   []string `json:"providesApis,omitempty" yaml:"providesApis,omitempty"`     // ProvidesApis is an array of entity references to the APIs that are provided by the component.
	ConsumesApis   []string `json:"consumesApis,omitempty" yaml:"onsumesApis,omitempty"`      // ConsumesApis is an array of entity references to the APIs that are consumed by the component.
	DependsOn      []string `json:"dependsOn,omitempty" yaml:"dependsOn,omitempty"`           // DependsOn is an array of entity references to the components and resources that the component depends on.
	System         string   `json:"system,omitempty" yaml:"system,omitempty"`                 // System is an array of references to other entities that the component depends on to function.
}

const (
	backstageURL = "https://backstage.internal-weedmaps.com"
)

// BSCE is shorthand for the backstage.ComponentEntityV1alpha1 type (which is a struct)
type BSCE *backstage.ComponentEntityV1alpha1

// ReadBS takes a weedmaps service and returns the catalog from Backstage
// Currently only used for Backstage integration testing in component_test.go
func ReadBS(wms string) (string, BSCE, error) {
	c, err := backstage.NewClient(backstageURL, "default", nil)
	if err != nil {
		slog.Error("Failed to create Backstage client", slog.Any("Error", err))
	}

	c.Catalog.Components.Get(context.Background(), wms, "")
	if err != nil {
		slog.Error("Failed to fetch Backstage Components", slog.Any("Error", err))
	}

	component, _, err := c.Catalog.Components.Get(context.Background(), wms, "")
	if err != nil {
		slog.Error("Failed to fetch Component name from Backstage", slog.Any("Error", err))
	}

	var owner string
	for _, r := range component.Relations {
		switch r.Type {
		case "ownedBy":
			owner = r.Target.Name
			slog.Debug("Processing Component", slog.String("Type", r.Type))
			return owner, component, nil
		case "partOf":
			slog.Debug("Found Component", slog.String("Type", r.Type))
		default:
			slog.Error("Is Backstage data empty?")
			return "", component, errors.New("EmptyData")
		}
	}

	return "", component, nil
}

// These are the functions that used to call the Component.
// They may be useful in the future?

/*
	// component is a pointer, Type: *backstage.ComponentEntityV1alpha1
	component, _, err := c.Catalog.Components.Get(context.Background(), sc.Service, "")
	if err != nil {
		slog.Error("Failed to fetch Component from Backstage", slog.Any("Error", err))
		return "", err
	}

	// originC, _, err := ReadComponentBS(sc.Service, c)

	// Add the component pointer to the struct
	// This won't be needed long...
	sc.Origin = component
	slog.Info("Component Found", slog.String("Service", sc.Service), slog.String("Description", component.Metadata.Description))

	// Call the methods for filling in data from Backstage
	// This function only works with Components!

		o := sc.wmOwner()
		if owner != o {
			slog.Error("Owner mismatch", slog.String("System", owner), slog.String("Component", o))
		}
*/

// wmOwner. Checks Relations.Type for "ownedBy"
// This is to be matched with CODEOWNERS in GitHub
//
// This requires the presence of 'Origin     *backstage.ComponentEntityV1alpha1'
// in type SvcConfig
func (s *SvcConfig) wmOwner() string {
	var owner string
	for _, r := range s.Origin.Relations {
		if r.Type == "ownedBy" {
			owner = r.Target.Name
		}
	}

	// Update struct pointer
	s.Owner = owner

	// Return the value of Owner for now to allow a test on a string.
	return owner
}

// NOT IN USE YET, but the code works
// wmLabels. Should return a label we want.
/*
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
*/
