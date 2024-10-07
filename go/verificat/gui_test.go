package main

import (
	"bytes"
	"os"
	"testing"

	approvals "github.com/approvals/go-approval-tests"
)

func TestMain(m *testing.M) {
	// Configure local dir for approval golden copies
	approvals.UseFolder("testdata")

	// Run tests
	exitVal := m.Run()

	// Post
	os.Exit(exitVal)
}

// Almanac Rendered to a Web Page
func TestRender(t *testing.T) {
	var (
		buf       = bytes.Buffer{}
		targetDoc = "almanac.gohtml"
		tmpldir   = "templates/*.gohtml"
		aWeb      = &AlmanacWeb{
			Title:   "Most Recent Almanac",
			Content: "This is content, whether you like it or not.",
		}
	)

	// Right now just getting functionality working, so a string is all I'm printing out.
	t.Run("renders string as HTML", func(t *testing.T) {
		if err := RenderWeb(&buf, aWeb, tmpldir, targetDoc); err != nil {
			t.Fatal(err)
		}
		// Use a golden copy of the good working template to compare and test against.
		approvals.VerifyString(t, buf.String())
	})

	t.Run("bad template location returns an error", func(t *testing.T) {
		tmpldir = "something/*else.html"
		if err := RenderWeb(&buf, aWeb, tmpldir, targetDoc); err == nil {
			t.Errorf("Expected an error but did not get one: %+v", err)
		}
	})

	t.Run("bad target document returns an error", func(t *testing.T) {
		targetDoc = "somethingelse.html"
		if err := RenderWeb(&buf, aWeb, tmpldir, targetDoc); err == nil {
			t.Errorf("Expected an error but did not get one: %+v", err)
		}
	})
}
