package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// This is a data stub for mocking activities on the server
// It is built within each testing function
type StubServiceStore struct {
	IDs         map[string]int
	verifyCalls []string
	almanac     []WMService
}

func (s *StubServiceStore) GetTriggerID(name string) int {
	lastID := s.IDs[name]
	return lastID
}

func (s *StubServiceStore) TriggerID(name string, score int) {
	s.verifyCalls = append(s.verifyCalls, name)
}

func (s *StubServiceStore) GetAlmanac() Almanac {
	return s.almanac
}

// Test /almanac endpoint with JSON output
func TestAlmanac(t *testing.T) {

	t.Run("it returns the almanac table as JSON", func(t *testing.T) {
		/*
			wantedAlmanac := []WMService{
				{"core", 50},
				{"admin", 60},
				{"moonshot", 70},
			}
		*/
		// Can the Almanac return three values?
		// name, lastid, score
		wantedAlmanac := []WMService{
			{"svcA", 3, 50},
			{"svcB", 5, 60},
			{"svcC", 8, 70},
		}

		store := StubServiceStore{nil, nil, wantedAlmanac}
		server := NewVerificationServ(&store)

		request := newAlmanacRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getAlmanacFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertAlmanac(t, got, wantedAlmanac)
		assertContentType(t, response, jsonContentType)
	})
}

func TestGETServices(t *testing.T) {
	// Make a new stub "store" to use in testing
	store := StubServiceStore{
		map[string]int{
			"admin":  20,
			"Craque": 10,
		},
		nil, nil,
	}
	// A new struct with an internal reference to the interface
	server := NewVerificationServ(&store)

	t.Run("returns admin's LastID", func(t *testing.T) {
		request := newGetTriggerIDReq("admin")
		// Use ResponseRecorder for a canned response
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Craque's LastID", func(t *testing.T) {
		request := newGetTriggerIDReq("Craque")
		// Use ResponseRecorder for a canned response
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing services", func(t *testing.T) {
		request := newGetTriggerIDReq("Mattic")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

}

func TestStoreIDs(t *testing.T) {
	store := StubServiceStore{
		map[string]int{},
		nil, nil,
	}
	server := NewVerificationServ(&store)

	t.Run("it returns accepted on POST", func(t *testing.T) {
		service := "admin"
		request := newPostIDReq(service)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.verifyCalls) != 1 {
			t.Errorf("got %d calls to TriggerID want %d", len(store.verifyCalls), 1)
		}

		if store.verifyCalls[0] != service {
			t.Errorf("did not store correct winner got %q want %q", store.verifyCalls[0], service)
		}
	})
}

// This saves us from not having to test the temporary InMemoryServiceStore
// and this code integration test can be reused with some other value for /store/
func TestRecordingIDsAndRetrievingThem(t *testing.T) {
	log.Printf("Begin Database Integration Test")

	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()

	store, err := NewFSStore(database)
	if err != nil {
		log.Fatalf("Integration test: problem creating file system service store, %v ", err)
	}

	server := NewVerificationServ(store)
	service := "admin"

	server.ServeHTTP(httptest.NewRecorder(), newPostIDReq(service))
	server.ServeHTTP(httptest.NewRecorder(), newPostIDReq(service))
	server.ServeHTTP(httptest.NewRecorder(), newPostIDReq(service))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetTriggerIDReq(service))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}

func newGetTriggerIDReq(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/v0/%s", name), nil)
	return req
}

func newPostIDReq(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/v0/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func getAlmanacFromResponse(t testing.TB, body io.Reader) (almanac []WMService) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&almanac)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of WMService, '%v'", body, err)
	}

	return
}

func assertAlmanac(t testing.TB, got, want []WMService) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Almanac returns %v want %v", got, want)
	}
}

func newAlmanacRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/v0/almanac", nil)
	return req
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
