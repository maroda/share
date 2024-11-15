package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	webTimeout time.Duration = 10 * time.Second
	prefixGH                 = "* @GhostGroup/"
	suffixGH                 = "\n"
)

type SvcTest interface {
	TestItem(svc string) *TestReturn
}

// SvcTestDB is the results database.
// It needs to work in conjunction with the query entry in SvcConfig.
// This is passed through ReadinessDisplay to operate TestReturn.
// Its values are then available in runVerification,
// which has access to this struct for adding scoring.
type SvcTestDB struct {
	Service  string // The service to test, e.g.: admin
	Datetime int64  // A start timestamp
	Owner    string // The retrieved Owner from Backstage
	Score    int    // Score out of 100 available test points
}

// TestReturn holds the answers for this test
type TestReturn struct {
	Present bool
	Owner   string
	Reality string
	Works   bool
	Score   int
}

const (
	ghDomain  = "https://api.github.com"
	ghPreURI  = "/repos/GhostGroup/"
	ghGetPATH = "/contents/.github/CODEOWNERS"
)

// urlCat is variadic, concatenating any set of strings into a URL.
// It can be used to embed a dynamic string alongside static parts of a URI.
// /u/ is a slice of strings used to build completeURL
func urlCat(u ...string) string {
	var completeURL string
	for _, p := range u {
		completeURL = completeURL + p
	}
	slog.Info("New Source Created", slog.String("URL", completeURL))
	return completeURL
}

// TestItem is returning a test result to ReadinessDisplay
// Currently this represents the "Owner" test between Backstage and GitHub
func (s *SvcTestDB) TestItem(svc string) *TestReturn {
	// Check the owner field.
	// Validation: If it's populated, return true.
	// Verification: If it's populated with the correct string, return true.

	var present, works bool

	// Get the actual value from CODEOWNERS in the matching GitHub repo
	fetchreal, err := Fetch(urlCat(ghDomain, ghPreURI, svc, ghGetPATH))
	if err != nil {
		slog.Error("Cannot Fetch", slog.Any("Error", err))
	}

	// Strip off the CODEOWNERS formatting
	reality := strings.TrimPrefix(strings.TrimSuffix(fetchreal, suffixGH), prefixGH)

	// Check the Owner for any WMService in Backstage
	if s.Owner == "" {
		// Validation has failed, the field is empty
		present = false
		s.Score--
		slog.Warn("Empty Field", slog.String("Owner", s.Owner))
		// Verification automatically fails
		s.Score--
		slog.Info("New Adjustment", slog.Int("Score", s.Score))
	} else {
		// Validation succeeds!
		present = true
		// Now check if it is equal to the retrieved source of truth
		if s.Owner != reality {
			// Verification has failed
			works = false
			slog.Warn("Unequal Field", slog.String("Owner", s.Owner), slog.String("Reality", reality))
			s.Score--
			slog.Info("New Adjustment", slog.Int("Score", s.Score))
		} else {
			// Verification succeeds!
			works = true
			slog.Info("Matching Field", slog.String("Owner", s.Owner), slog.String("Reality", reality), slog.Int("Score", s.Score))
		}
	}

	// This will be included in the API return value
	return &TestReturn{Present: present, Owner: s.Owner, Reality: reality, Works: works, Score: s.Score}
}

// ReadinessDisplay takes the data and runs queries for processing and presentation.
// The first arg /i/ is the catalog with its data.
// The second is which service is being tested.
// The third is where this output goes.
func ReadinessDisplay(i SvcTest, service string, w io.Writer) error {
	// Our first test is just to verify that the Codeowners field is populated in Backstage.
	returnedTest := i.TestItem(service)
	returnOut, err := json.Marshal(returnedTest)
	if err != nil {
		slog.Error("Failed to marshal struct to JSON", slog.Any("Error", err))
	}

	// Print the JSON bytestring and return
	fmt.Fprintf(w, string(returnOut))
	return err
}

// Fetch will initiate data retrieval from test sources.
// `url` is an endpoint, typically an API request.
func Fetch(url string) (string, error) {
	// this should take the data interface?
	return ConfiguredFetch(url, 10*time.Second)
}

// ConfiguredFetch conducts parallel fetches for verification data.
func ConfiguredFetch(url string, timeout time.Duration) (string, error) {
	// In the future, this can take more URLs,
	// so we can conduct the gets from here in parallel.
	select {
	case a := <-extractGitHub(url):
		return a.Answer, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s", url)
	}
}

// wwwFetch is a struct for values returned to Fetch over HTTP.
type wwwFetch struct {
	Answer string
}

// Extract data from GitHub
// (i.e. data with the need to set headers, which should be generalized)
func extractGitHub(url string) chan wwwFetch {
	// Grab GH_TOKEN from the environment
	envVar := "GH_TOKEN"
	token := fillEnvVar(envVar)

	// if there's no EnvVar, log an error and go no further
	if token == "ENOENT" {
		slog.Error("Environment Variable not set", slog.String("Key", envVar), slog.String("Value", token))
		return nil
	}

	authHeader := "Bearer " + token + ""

	// make a channel for fetching
	www := make(chan wwwFetch)

	go func() {
		// Create a new HTTP request object
		// This will be passed to a new HTTP client below.
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			slog.Error("Could not create http client request", slog.String("URL", url), slog.Any("Error", err))
		}

		// Add Auth headers to the object.
		req.Header.Add("Accept", "application/vnd.github.raw+json")
		req.Header.Add("Authorization", authHeader)

		// Create a new Client pointer with a configured timeout
		client := &http.Client{Timeout: webTimeout}

		// Perform the actual Get.
		r, err := client.Do(req)
		if err != nil {
			slog.Error("Could not reach service", slog.String("URL", url), slog.Any("Error", err))
		}
		defer r.Body.Close()

		// HTTP status code has to be 200 to work
		if r.StatusCode == http.StatusOK {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				slog.Error("Could not read value at", slog.String("URL", url), slog.Any("Error", err))
			}
			bodyString := string(body)
			rval := &wwwFetch{Answer: bodyString}
			www <- *rval
		} else {
			slog.Error("Non-200 Status", slog.String("URL", url), slog.Any("Status", r.StatusCode))
		}
		close(www)
	}()
	return www
}
