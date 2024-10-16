package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

/*
	Integration Testing for Backstage
*/

// Can we see all elements of every component entity?
// Currently this only checks Backstage, because I'm figuring out a data problem.
// TODO: A unit test version of this should use a fake BSCE object.
func TestReadBS(t *testing.T) {
	readTests := []struct {
		Name    string
		Service string
		Expect  string
	}{
		{Name: "Admin", Service: "admin", Expect: "code-owners-admin"},
		{Name: "CoreBraze", Service: "core-braze", Expect: "code-owners-core"},
		{Name: "AdServer", Service: "ad-server-cache-builder", Expect: "code-owners-wasp"},
	}

	for _, tt := range readTests {
		got, component, err := ReadBS(tt.Service)
		want := tt.Expect
		assertError(t, err, nil)
		if diff := cmp.Diff(got, want); diff != "" {
			t.Error(diff)
			t.Errorf("For '%v' the component looks like\n: %v", tt.Service, component)
		}
	}
}
