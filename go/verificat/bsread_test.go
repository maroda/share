package main

import (
	"testing"
)

// Mock data for running Production ReadinessRead
type mockSvcConfig struct {
	URL      string
	Service  string
	Datetime int64
	Owner    string
}

// Mock method that satisfies SvcCat{}
func (sc *mockSvcConfig) ReadSvc() (string, error) {
	return "code-owners-admin", nil
}

func (s *mockSvcConfig) wmOwner() string {
	return "code-owners-admin"
}

// ReadinessRead runs the full data fill from Backstage.
// Right now all it returns is an owner, so we mock that here.
//
//	A future test will expect to have a data struct built
func TestReadinessRead(t *testing.T) {
	mockC := &mockSvcConfig{URL: "blank", Service: "admin"}

	got, err := ReadinessRead(mockC)
	want := "code-owners-admin"

	assertString(t, got, want)
	assertError(t, err, nil)
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
