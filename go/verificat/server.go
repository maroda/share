package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const jsonContentType = "application/json"

// This defines the service and its final checklist score
// The final score is taken from 100.
// If the checklist comes back as true,
// then the score is 0, and this value remains 100.
// If the checklist comes back as false,
// then the score is 1, and this value becomes 99.
type WMService struct {
	Name   string // Weedmaps Service Name
	LastID int    // The last test ID
	Score  int    // The current score (100 - score)
}

type ServiceStore interface {
	GetTriggerID(name string) int     // Retrieve the count of tests done
	TriggerID(name string, score int) // The current run ID, its score
	GetAlmanac() Almanac              // A collection of all services and their scores
}

// The struct needs to reference the interface to use it
type VerificationServ struct {
	store ServiceStore
	http.Handler
}

// Control the launch of an HTTP service using routed endpoints.
func NewVerificationServ(store ServiceStore) *VerificationServ {
	v := new(VerificationServ)
	v.store = store

	// This will be assigned to the http.Handler in PlayerServer
	// so that the routing is done once at the start, not on every request.
	router := http.NewServeMux()

	// Set up each server endpoint and its associated handler function
	router.Handle("/v0/almanac", http.HandlerFunc(v.almanacHandler))
	router.Handle("/v0/", http.HandlerFunc(v.servicesHandler))

	v.Handler = router

	return v
}

// Return the full JSON almanac of WMServices and their verification scores.
func (p *VerificationServ) almanacHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetAlmanac())
}

// Handler for /v0/
func (p *VerificationServ) servicesHandler(w http.ResponseWriter, r *http.Request) {
	// extract this once here, then it's not necessary to pass http.Request
	service := strings.TrimPrefix(r.URL.Path, "/v0/")

	// Based on the method of the HTTP action, do different things.
	// These methods on VerificationServ can pass the handler interfaces around
	switch r.Method {
	case http.MethodPost:
		// Record the session ID.
		//p.processVerify(w, service)
		// Kick off the test, now including processVerify
		p.runVerification(w, service)
	case http.MethodGet:
		// Get last session ID from the database.
		p.showLastID(w, service)
	}
}

// showLastID will display the ID of the most recent verification run for this service.
func (p *VerificationServ) showLastID(w http.ResponseWriter, service string) {
	// GetTriggerID is a method available through the interface
	lastID := p.store.GetTriggerID(service)

	if lastID == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, lastID)
}

// TriggerID is a monotonic "ID" that is incremented each time POST is run against that service.
// e.g. curl -X POST http://service:5000/v0/admin
// increments TriggerID
/*
func (p *VerificationServ) processVerify(w http.ResponseWriter, service string) {
	p.store.TriggerID(service)
	w.WriteHeader(http.StatusAccepted)
}
*/

// runVerification. Takes a passed configuration and launches testing.
func (p *VerificationServ) runVerification(w http.ResponseWriter, service string) {
	w.WriteHeader(http.StatusAccepted)

	// Create a data object for the configuration.
	// TODO: Move backstage URL into .env
	svcconf := &SvcConfig{URL: "https://backstage.internal-weedmaps.com", Service: service}

	// Read the SVC and get the "owner" string back
	// We don't need a return, it updates the struct
	_, err := ReadinessRead(svcconf)
	if err != nil {
		fmt.Println("Error reading Service: ", err)
		log.Fatal(err)
	}

	// ReadinessDisplay expects an interface with this struct
	// These values have been filled in by ReadinessRead() above
	// Score is initialized to 100 each time,
	//	then decremented on each failed test
	//	that is handled by ReadinessDisplay.
	stests := &SvcTestDB{Datetime: svcconf.Datetime, Owner: svcconf.Owner, Score: 100}

	// Send test metadata to ReadinessDisplay, which launches tests and displays the results.
	// w == http.ResponseWriter, which satisfies io.Writer
	err = ReadinessDisplay(stests, service, w)
	if err != nil {
		fmt.Println("UNEXPECTED ERROR: ", err)
	}

	// Initiate the TriggerID sequence that is used to set WMService.Score in the database.
	p.store.TriggerID(service, stests.Score)
}
