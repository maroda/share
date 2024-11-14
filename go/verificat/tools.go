package main

import (
	"io"
	"os"
)

// Return the value of a runtime Environment Variable
func fillEnvVar(ev string) string {
	// If the EnvVar doesn't exist return a default string
	value := os.Getenv(ev)
	if value == "" {
		value = "ENOENT"
	}
	return value
}

// When we write, we write from the beginning.
// This type takes a file and makes sure we're at the start.
type tape struct {
	file *os.File
}

// Return the file start location of 0 for writing.
func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, io.SeekStart)
	return t.file.Write(p)
}
