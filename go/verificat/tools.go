package main

import (
	"io/fs"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Use an .env file to load configuration variables into a map
// This is currently the only way this works,
// this func does not read shell ENV VARs or touch the os ENV.
func getEnvVar(key, loc string) (string, error) {
	// Read the values from .env into a map
	var cvEnv map[string]string
	cvEnv, err := godotenv.Read(loc)
	if err != nil {
		log.Fatal("Error loading value: ", err)
	}
	return cvEnv[key], err
}

// This takes the new fs.FS interface as fileSystem,
// preceded by the key being requested.
func NewConfigFromFS(k string, fileSystem fs.FS) (string, error) {
	var foundEnv string

	// Walk the directory and find .env
	fs.WalkDir(fileSystem, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == ".env" {
			foundEnv = filepath.Ext(p)
		}
		return nil
	})

	// DEBUG ::: Useful to read the contents of the found file.
	/*
		f, err := fs.ReadFile(fileSystem, foundEnv)
		if err != nil {
			return "", err
		}
		fng := string(f[:])
	*/

	// Retrieve the key value /k/ from the filename /foundEnv/
	return getEnvVar(k, foundEnv)
}
