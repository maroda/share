package main

import (
	"context"
	"errors"
	"log/slog"

	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

// BSComponent is a local struct for carrying the contents of .Catalog.Systems.Get
type BSSystem struct {
	Owner  string `json:"owner" yaml:"owner"`                       // Owner is an entity reference to the owner of the system.
	Domain string `json:"domain,omitempty" yaml:"domain,omitempty"` // Domain is an entity reference to the domain that the system belongs to.
}

// BSCE is shorthand for the backstage.SystemEntityV1alpha1 type (which is a struct)
type BSSE *backstage.SystemEntityV1alpha1

var SystemNotRecognized = errors.New("System not recognized")

// ReadBS takes a weedmaps service and returns the service owner
func ReadSystemBS(wms string, c *backstage.Client) (string, BSSE, error) {
	// first get a list of systems
	services, err := bsSystemList(wms, c)
	if err != nil {
		return "", nil, err
	}

	// make sure the requested test belongs in the list
	// if it doesn't belong, the test will appear zeroed out.
	for _, service := range services {
		if wms == service {
			// When there is a match with the System List,
			// grab the System Entity itself and get the Owner.
			se, _, err := c.Catalog.Systems.Get(context.Background(), wms, "")
			if err != nil {
				slog.Error("Failed to fetch System", slog.Any("Error", err))
				return "", nil, err
			}
			owner := se.Spec.Owner
			slog.Info("Owner Found", slog.String("Owner", owner))
			return owner, se, err
		}
	}

	// If we've gotten this far, there wasn't a match.
	slog.Error("Failed to fetch System", slog.Any("Error", SystemNotRecognized))
	return "", nil, SystemNotRecognized
}

// bsSystemList queries Backstage for a list of all System definitions ("kind=system") and returns an array populated with the list.
// TODO: Make this a map, key = backstage system / value = github repo link
// This way, we can check the real repo for CODEOWNERS instead of defaulting to the 'service name' (which works for some, but not all)
// This is: .Metadata.Annotations.github.com/project-slug (e.g.: for `core` this is `GhostGroup/weedmaps`)
func bsSystemList(wms string, c *backstage.Client) ([]string, error) {
	var s []string

	if systems, _, err := c.Catalog.Entities.List(context.Background(), &backstage.ListEntityOptions{Filters: []string{"kind=system"}}); err != nil {
		slog.Error("Failed to get System List from Backstage", slog.Any("Error", err))
		s = append(s, "")
		return s, err
	} else {
		for _, e := range systems {
			s = append(s, e.Metadata.Name)
		}
		slog.Debug("System List Found", slog.Any("Systems", s))
		return s, err
	}
}
