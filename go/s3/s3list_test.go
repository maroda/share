package main

import (
	"flag"
	"testing"
)

// Go Test already runs `flag.Parse()` so
// explictly defining the flag means it will be set in the test to the default
var localProfile = flag.String("profile", "test", "local profile name (default: test)")

// Test if localProfile defaults to "test" if the flag is not used
func TestFlag(t *testing.T) {
	want := "test"
	got := *localProfile

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// NewCData is a struct constructor function.
// Test that it returns a proper struct with
// variables used for S3 actions.
func TestNewCData(t *testing.T) {
	want := struct {
		region string
		bucket string
		key    string
	}{
		region: "a",
		bucket: "b",
		key:    "c",
	}
	got := *NewCData("a", "b", "c")

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// mockCData is an identical struct to CData...
type mockCData struct {
	region string
	bucket string
	key    string
}

// ... only mockCData's method doesn't do a search
// it only returns the expected string
func (cd *mockCData) SearchO() (string, error) {
	return "userneeds.png", nil
}

// Test the Find function
// Anything deeper than this will probably need AWS SDK mocking
func TestFind(t *testing.T) {
	// Do not use the constructor here
	mockF := &mockCData{region: "a", bucket: "b", key: "c"}

	got, err := Find(mockF)
	want := "userneeds.png"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	assertError(t, err, nil)
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got error %q want %q", got, want)
	}
}
