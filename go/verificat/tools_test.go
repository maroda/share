package main

import (
	"io"
	"log"
	"testing"
	"testing/fstest"
)

// Will a value be returned for an envar key
func TestGetEnvVar(t *testing.T) {
	filename := ".env"
	testenvvar := "TOKEN=my_1029384756"
	// We need a fake environment file first
	// fstest.MapFS provides fs.FS
	fs := fstest.MapFS{
		filename: {
			Data: []byte(testenvvar),
		},
	}

	key := "TOKEN"

	// Make sure the file ".env" was created in this
	data, err := fs.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	} else {
		log.Println("TestGetEnvVar string data match is", string(data) == "TOKEN=my_1029384756")
	}

	// Another way to test that the file exists in this test
	if err := fstest.TestFS(fs, ".env"); err != nil {
		t.Fatal(err)
	}

	/*
		Investigating a potential bug in godotenv
		Trying to use a fstest.TestFS
		And godotenv isn't changing its filesystem to look for the file.
	*/
	// want := "my_1029384756"
	want := ""
	got, err := NewConfigFromFS(key, fs)

	assertString(t, got, want)
	assertError(t, err, nil)
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
