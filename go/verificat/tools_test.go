package main

import (
	"testing"
	"testing/fstest"
)

// Will a value be returned for an envar key
func TestgetEnvVar(t *testing.T) {
	// We need a fake environment file first
	// fstest.MapFS provides fs.FS
	fs := fstest.MapFS{
		".env": {Data: []byte("TOKEN=my_1029384756")},
	}
	key := "TOKEN"

	got := "my_1029384756"
	want, err := NewConfigFromFS(key, fs)

	assertString(t, got, want)
	assertError(t, err, nil)
}
