package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tdabasinskas/go-backstage/v2/backstage"
)

const (
	DefaultNamespaceName = "default"
	userAgent            = "go-backstage"
)

/*
	Integration Testing for Backstage
*/

// TODO: A unit test version of this should use a fake BSCE object.
//
// Can we see all elements of every component entity?
func TestReadComponentBS(t *testing.T) {
	c, err := backstage.NewClient(backstageURL, "default", nil)
	assertError(t, err, nil)

	readTests := []struct {
		Name    string
		Service string
		Expect  string
		Client  *backstage.Client
	}{
		{"Admin", "admin", "code-owners-admin", c},
		{"CoreBraze", "core-braze", "code-owners-core", c},
		{"AdServer", "ad-server-cache-builder", "code-owners-wasp", c},
	}

	for _, tt := range readTests {
		component, err := ReadComponentBS(tt.Service, c)
		got := component.Spec.Owner
		want := tt.Expect
		assertError(t, err, nil)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Error(diff)
			t.Errorf("For '%v' the component looks like\n: %v", tt.Service, component)
		}
	}

	t.Run("Receives error with bad client", func(t *testing.T) {
		var c BSBadClient
		service := "admin"
		_, err := ReadComponentBS(service, c)
		assertNoError(t, err)
	})
}
