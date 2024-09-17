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
/*

	For some reason, this function is still loading the file on-disk during tests.

*/
func getEnvVar(key, loc string) (string, error) {
	// Read the values from .env into a map
	var cvEnv map[string]string

	/*

		So even though I am sending a loc string ".env"
		this function somehow is defaulting to the regular filesystem every time.
		Compiled it works fine of course, but testing ... something is weird.

	*/
	cvEnv, err := godotenv.Read(loc)
	if err != nil {
		log.Fatal("Error loading value: ", err)
	}

	// These *values* are correct, but this function doesn't know there's a fstest filesystem.
	// fmt.Println(key)
	// fmt.Println(cvEnv)

	// without a real .env in-place, this will fail
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
			// fmt.Println(d)
			// returns different things here...
			// for the on-disk file it is: - .env
			// for the test fs it is: ---------- 19 0001-01-01 00:00:00 .env
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
		fmt.Println("NewConfigFromFS found:", fng)
		fmt.Println("OR:", k)
	*/

	// Retrieve the key value /k/ from the filename /foundEnv/
	/*

		While this works fine in reality,
		the fstest filesystem isn't being passed.
		so it isn't being used.

		This may be a limitation of godotenv.
		...which tells me I may be trying to test a tested thing
		and I need to pull my tests further out...
	*/
	return getEnvVar(k, foundEnv)
}
