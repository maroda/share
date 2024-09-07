package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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

// This is returning a test result to ReadinessDisplay
// Currently this represents the "Owner" test between Backstage and GitHub
func (s *SvcTestDB) TestItem(svc string) *TestReturn {
	// Check the owner field. If it's populated, return true. If not, return false.
	var present, works bool

	// REFACTOR ::: Should the Fetch() method be on SvcTest?
	// 		This way... it can be used here: i.Fetch() whatever...
	// Build the source address
	// TODO: This should actually belong in the github function to make TestItem more general
	source := "https://api.github.com"
	preURI := "/repos/GhostGroup/"
	getPATH := "/contents/.github/CODEOWNERS"
	fetchURL := source + preURI + svc + getPATH
	reality, err := Fetch(fetchURL)
	if err != nil {
		log.Fatal(err)
	}

	// Verify the Owner for any WMService in Backstage
	if s.Owner == "" {
		present = false
		s.Score--
		log.Println("ERROR: Owner field is empty in Backstage!")
		log.Println("INFO: Score set to: ", s.Score)
	} else {
		present = true
		if s.Owner != reality {
			works = false
			s.Score--
			log.Printf("WARNING: Backstage value %q does not equal source %q", s.Owner, reality)
			log.Println("INFO: Score set to: ", s.Score)
		} else {
			log.Printf("MATCH: Source %q in Backstage is: %q", reality, s.Owner)
			works = true
			log.Println("INFO: Score remains at: ", s.Score)
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

	// TODO: This should be a struct marshalled into json
	fmt.Fprintf(w, "%t, %q, %q, %t, %d", returnedTest.Present, returnedTest.Owner, returnedTest.Reality, returnedTest.Works, returnedTest.Score)

	// currently this is data manipulation, nothing here throws an error (yet)
	return nil
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

// This extracts data from a URL.
// It doesn't do any auth, it assumes the URL will be available.
func extract(url string) chan wwwFetch {
	// make a channel for fetching
	www := make(chan wwwFetch)

	// fetch the value
	go func() {
		var client http.Client

		r, err := client.Get(url)
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
		}
		close(www)
	}()
	return www
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
		// TODO: make this timeout a global setting/const?
		client := &http.Client{Timeout: 10 * time.Second}

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
