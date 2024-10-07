package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	webTimeout time.Duration = 10 * time.Second
)

type SvcTest interface {
	TestItem(svc string) *TestReturn
}

// This is the results database.
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

// This is returning a test result to ReadinessDisplay
// Currently this represents the "Owner" test between Backstage and GitHub
func (s *SvcTestDB) TestItem(svc string) *TestReturn {
	// Check the owner field. If it's populated, return true. If not, return false.
	var present, works bool

	reality, err := Fetch(urlCat(ghDomain, ghPreURI, svc, ghGetPATH))
	if err != nil {
		log.Fatal(err)
	}

	// Check the Owner for any WMService in Backstage
	if s.Owner == "" {
		// Validation has failed, the field is empty
		present = false
		s.Score--
		slog.Error("Invalid Field", slog.String("Owner", s.Owner))
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

	// I have ScvTestDB here so I need to put a value in it

	// This will be included in the API return value
	return &TestReturn{Present: present, Owner: s.Owner, Reality: reality, Works: works, Score: s.Score}
}

// ReadinessDisplay. Take the data and run queries for processing and presentation.
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

// Fetch. Initiate data retrieval from test sources.
// `url` is an endpoint, typically an API request.
func Fetch(url string) (string, error) {
	// this should take the data interface?
	return ConfiguredFetch(url, 10*time.Second)
}

// ConfiguredFetch. Conduct parallel fetches for verification data.
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
	// Get GH_TOKEN from the environment
	// Making this available as a system-wide map might be more efficient
	// ...but it is also less secure, more of the code can see those values.
	// So for now, we check for just the one we need: GH_TOKEN
	fsys := os.DirFS(".")
	token, err := NewConfigFromFS("GH_TOKEN", fsys)
	if err != nil {
		log.Fatal("Token could not be loaded: ", err)
	}
	authHeader := "Bearer " + token + ""

	// make a channel for fetching
	www := make(chan wwwFetch)

	go func() {
		// Create a new HTTP request object
		// This will be passed to a new HTTP client below.
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			fmt.Errorf("Could not create request")
		}

		// Add Auth headers to the object.
		req.Header.Add("Accept", "application/vnd.github.raw+json")
		req.Header.Add("Authorization", authHeader)

		// Create a new Client pointer with a configured timeout
		client := &http.Client{Timeout: webTimeout}

		// Perform the actual Get.
		r, err := client.Do(req)
		if err != nil {
			fmt.Errorf("Could not reach service at %q\n", url)
		}
		defer r.Body.Close()

		// HTTP status code has to be 200 to work
		if r.StatusCode == http.StatusOK {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				fmt.Errorf("Could not read value at %q\n", url)
			}
			bodyString := string(body)
			rval := &wwwFetch{Answer: bodyString}
			www <- *rval
		} else {
			fmt.Errorf("Error! Non-200 Status: %q\n", r.StatusCode)
		}
		close(www)
	}()
	return www
}
