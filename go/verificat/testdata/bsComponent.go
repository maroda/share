package main

import (
	"context"
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

// ReadComponentBS takes a weedmaps service and returns the catalog from Backstage
// Currently only used for Backstage integration testing in component_test.go
func ReadComponentBS(wms string, c *backstage.Client) (BSCE, error) {
	component, _, err := c.Catalog.Components.Get(context.Background(), wms, "")
	if err != nil {
		slog.Error("Failed to fetch Component name from Backstage", slog.Any("Error", err))
	}

	/*
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
	*/

	return component, err
}
