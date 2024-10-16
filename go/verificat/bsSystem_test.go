package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

/*
Integration Testing for Backstage
*/
const (
	backstageURL = "https://backstage.internal-weedmaps.com"
)

// TODO: A unit test version of this could use a fake BSSE object.
//
// Can we see match elements of each service entity?
func TestReadSystemBS(t *testing.T) {
	t.Run("reads Backstage catalog and matches system names", func(t *testing.T) {
		c, err := backstage.NewClient(backstageURL, "default", nil)
		assertError(t, err, nil)

		readTests := []struct {
			Name    string
			Service string
			Expect  string
			Client  *backstage.Client
		}{
			{"Admin", "admin", "code-owners-admin", c},
			{"Core", "core", "code-owners-core", c},
			{"AdServer", "ad-server", "code-owners-wasp", c},
			{"WeedmapsAPI", "weedmaps-api", "code-owners-api", c},
		}

		for _, tt := range readTests {
			got, service, err := ReadSystemBS(tt.Service, c)
			want := tt.Expect
			assertError(t, err, nil)
			if diff := cmp.Diff(got, want); diff != "" {
				t.Error(diff)
				t.Errorf("For '%v' the service looks like\n: %v", tt.Service, service)
			}
		}
	})

	t.Run("Handles only Systems it knows about", func(t *testing.T) {
		c, err := backstage.NewClient(backstageURL, "default", nil)
		assertError(t, err, nil)

		readTests := []struct {
			Name    string
			Service string
			Expect  string
			Client  *backstage.Client
		}{
			{"Core-App", "core-app", "", c},
		}

		for _, tt := range readTests {
			got, service, err := ReadSystemBS(tt.Service, c)
			want := tt.Expect
			assertError(t, err, nil)
			if diff := cmp.Diff(got, want); diff != "" {
				t.Error(diff)
				t.Errorf("For '%v' the service looks like\n: %v", tt.Service, service)
			}
		}
	})

	// Use a bad client value to trigger an error
	/*
		t.Run("Receives error with bad client", func(t *testing.T) {
			_, err := backstage.NewClient(backstageURL, "default", nil)
			assertNoError(t, err)
		})
	*/
}
