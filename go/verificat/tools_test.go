package main

import (
	"io"
	"os"
	"testing"
)

// Use a regular Environment Variable to fill the value requested
func TestFillEnvVar(t *testing.T) {

	t.Run("returns a default value", func(t *testing.T) {
		// This 'want' value is a default
		ev := "ANYTHING"
		want := "ENOENT"
		got := fillEnvVar(ev)

		assertString(t, got, want)
	})

	t.Run("returns a set value", func(t *testing.T) {
		// This 'want' value is a default
		ev := "TOKEN"
		want := "ghp_0987654321"

		// Set an environment variable to check
		err := os.Setenv("TOKEN", "ghp_0987654321")
		if err != nil {
			t.Errorf("could not set environment variable")
		}

		got := fillEnvVar(ev)

		assertString(t, got, want)
	})
}

// Test that we can set a file to 0 for writing
func TestTape_Write(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()

	tape := &tape{file}

	tape.Write([]byte("abc"))

	file.Seek(0, io.SeekStart)
	newFileContents, _ := io.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
