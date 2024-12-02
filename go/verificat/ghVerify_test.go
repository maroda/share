package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Mock data for sending through ReadinessDisplay
type mockSvcTestDB struct {
	Service  string
	Datetime int64
	Owner    string
	Score    int
}

// These 'false' values are what comes through the mock call
// func (db *mockSvcTestDB) TestItem(svc string) (bool, string, string, bool) {
func (db *mockSvcTestDB) TestItem(svc string) *TestReturn {
	return &TestReturn{
		Present: true,
		Owner:   "mock-admin-group",
		Reality: "mock-developer-group",
		Works:   false,
		Score:   100}
}

// Is Readiness Display correctly writing results?
func TestReadinessDisplay(t *testing.T) {
	service := "admin"
	buffer := bytes.Buffer{}
	mockRD := &mockSvcTestDB{Service: service, Datetime: 1724367242, Owner: "code-owners-admin", Score: 0}

	t.Run("Is Readiness Display correctly writing results?", func(t *testing.T) {
		// ReadinessDisplay calls TestItem, which needs to send us more data
		err := ReadinessDisplay(mockRD, service, &buffer)
		got := buffer.String()
		want := "{\"Present\":true,\"Owner\":\"mock-admin-group\",\"Reality\":\"mock-developer-group\",\"Works\":false,\"Score\":100}"

		// What we're comparing is the buffer string, not the structs.
		if diff := cmp.Diff(got, want); diff != "" {
			t.Error(diff)
		}

		assertError(t, err, nil)
	})
}

// This is an integration test to check that a GitHub URL is reachable.
// getGitHub checks for the token in the GH_TOKEN EnvVar
func TestGetGitHub(t *testing.T) {
	t.Run("Integration: GitHub auth is working", func(t *testing.T) {
		svc := "admin"
		var got string
		want := "* @GhostGroup/js-developers\n"
		url := ghDomain + ghPreURI + svc + ghGetPATH
		got, err := getGitHub(url)

		assertError(t, err, nil)
		assertString(t, got, want)
	})
}

func TestMultiFetch(t *testing.T) {
	t.Run("Returns an answer from multiple endpoints", func(t *testing.T) {
		// We need a set of test server URLs
		mockWWW := make(map[int]*httptest.Server)
		urlsWWW := make(map[int]string)

		// Fill mockWWW with httptest.Servers
		// Fill urlsWWW with those endpoints
		for i := 0; i <= 2; i++ {
			mockWWW[i] = makeMockWebServ(0 * time.Millisecond)
			urlsWWW[i] = mockWWW[i].URL
		}

		// Now we can send the list to MultiFetch
		want := []string{"ownership", "ownership", "ownership"}
		got, err := MultiFetch(urlsWWW)
		assertMultiString(t, got, want)
		assertError(t, err, nil)

		// Clean up resources
		for i := 0; i <= 2; i++ {
			mockWWW[i].Close()
		}
	})
}

func assertMultiString(t *testing.T, got, want []string) {
	t.Helper()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Error(diff)
	}
}

// Build a URL takes an arbitrary set of pieces and combines them into a browsable URL.
func TestUrlCat(t *testing.T) {
	WebDomain := "craque.bandcamp.com"
	URIPre := "/track/"
	t.Run("Returns a URL from static strings", func(t *testing.T) {
		URIDyna := "relaxant" // This should be tested as a var that changes, too
		URIPost := ""

		want := "craque.bandcamp.com/track/relaxant"
		got := urlCat(WebDomain, URIPre, URIDyna, URIPost)

		assertString(t, got, want)
	})

	t.Run("Returns a URL from dynamic strings inside static strings", func(t *testing.T) {
		URIPost := "/listen"
		three := []string{"relaxant", "manifold", "synapse"}

		for _, h := range three {
			want := "craque.bandcamp.com/track/" + h + "/listen"
			got := urlCat(WebDomain, URIPre, h, URIPost)

			assertString(t, got, want)
		}
	})
}

// Mock responder for external API calls
func makeMockWebServ(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testAnswer := []byte("ownership")
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		_, err := w.Write(testAnswer)
		if err != nil {
			log.Fatalf("ERROR: Could not write to output.")
		}
	}))
}
