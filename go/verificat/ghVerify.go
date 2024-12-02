package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	webTimeout = 10 * time.Second
	prefixGH   = "* @GhostGroup/"
	suffixGH   = "\n"
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

// Currently CODEOWNERS is the only thing we check in GitHub
// so this path is very specific. TODO: /ghGetPATH/ can be dynamic
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

	// answer, err := Fetch(urlCat(ghDomain, ghPreURI, svc, ghGetPATH))

	// Get the actual value from CODEOWNERS in the matching GitHub repos
	// This can take a map of URLs, but for now we only have one to give it.
	target := urlCat(ghDomain, ghPreURI, svc, ghGetPATH)
	urls := map[int]string{0: target}
	answer, err := MultiFetch(urls)
	if err != nil {
		slog.Error("Cannot Fetch", slog.Any("Error", err))
	}

	// Strip off the CODEOWNERS formatting
	reality := strings.TrimPrefix(strings.TrimSuffix(answer[0], suffixGH), prefixGH)

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
	_, err = fmt.Fprintf(w, string(returnOut))
	if err != nil {
		slog.Error("Failed to print JSON to Writer", slog.Any("Error", err))
	}
	return err
}

// MultiFetch is ConfiguredFetch for multiple urls in a []string
// It will return the "Answer" value for each URL in a []string with matching indexes
func MultiFetch(urls map[int]string) ([]string, error) {
	// ErrorGroup for catching multiple URL fetches,
	egrp := new(errgroup.Group)
	// using a limit of 3 concurrent fetches.
	// egrp.SetLimit(3)

	// This is the same as wwwFetch.Answer but for multiple URLs
	results := make([]string, len(urls))

	// Step through the list and fire off a check
	for i, url := range urls {
		egrp.Go(func() error {
			answer, err := getGitHub(url)
			results[i] = answer
			return err
		})
	}

	if err := egrp.Wait(); err != nil {
		return results, err
	}

	return results, nil
}

// getGitHub should take the url and a pointer to the results
// then update the pointer and return only an error
func getGitHub(currURL string) (string, error) {
	// Grab GH_TOKEN from the environment
	// if there's no EnvVar, log an error and go no further
	envVar := "GH_TOKEN"
	token := fillEnvVar(envVar)
	if token == "ENOENT" {
		slog.Error("Environment Variable not set", slog.String("Key", envVar), slog.String("Value", token))
		return token, nil
	}

	// Build the authHeader with the new token
	authHeader := "Bearer " + token + ""

	// Create a new HTTP request object
	// This will be passed to a new HTTP client below.
	req, err := http.NewRequest(http.MethodGet, currURL, nil)
	if err != nil {
		slog.Error("Could not create http client request", slog.String("URL", currURL), slog.Any("Error", err))
		return "", err
	}

	// Add Auth headers to the object.
	req.Header.Add("Accept", "application/vnd.github.raw+json")
	req.Header.Add("Authorization", authHeader)

	// Create a new Client pointer with a configured timeout
	client := &http.Client{Timeout: webTimeout}

	// Perform the actual Get.
	r, err := client.Do(req)
	if err != nil {
		slog.Error("Could not reach service", slog.String("URL", currURL), slog.Any("Error", err))
		r.Body.Close()
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		slog.Error("Non-200 Status", slog.String("URL", currURL), slog.Any("Status", r.StatusCode))
		return "", errors.New("non 200 Status")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Could not read value at", slog.String("URL", currURL), slog.Any("Error", err))
	}

	bodyString := string(body)
	return bodyString, err
}
