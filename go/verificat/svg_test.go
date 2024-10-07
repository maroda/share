package main

import (
	"strings"
	"testing"
)

func TestBuildSVG(t *testing.T) {
	// Create a test database for drawing SVGs
	database, cleanDatabase := createTempFile(t, `[
			{"Name": "Mattic", "LastID": 10, "Score": 99},
			{"Name": "Craque", "LastID": 4, "Score": 98}]`)
	defer cleanDatabase()
	store, err := NewFSStore(database)
	assertNoError(t, err)
	drawData := store.GetAlmanac()

	// This will configure draw output margins and offsets
	// Change these, and the /wants/ strings below will change.
	sc := &SVGCfg{Gutter: 3, TxtOff: 8, Spacer: 14}

	// Build the XML blob
	xml := BuildSVG(&drawData, sc)

	// DEBUG ::: fmt.Println(xml)

	// The test fails if the known value of a dynamically built string doesn't match.
	// These values are the known outcomes from the test data if BuildSVG is working.
	wants := []string{
		"M5 17h99v10H5z",
		"M5 31h98v10H5z",
		"y=\"25\"",
		"y=\"39\"",
	}

	for _, want := range wants {
		assertXML(t, xml, want)
	}
}

func assertXML(t *testing.T, xml, want string) {
	t.Helper()
	if !strings.Contains(xml, want) {
		t.Errorf("Expected to find Path definition '%+v' in XML:\n%+v", want, xml)
	}
}
