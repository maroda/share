package main

import (
	"log"
	"os"
	"testing"
)

// By returning a function that closes the file,
// the test doesn't need to worry about file details.
// The caller should assign the returned function to a variable
// and then defer that function variable,
// which effectively closes the file at the end of the calling function.
func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()

	tmpfile, err := os.CreateTemp("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func TestFileSystemStore(t *testing.T) {

	// Test that an almanac is returned from a database call
	t.Run("almanac from a reader", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Mattic", "LastID": 10, "Score": 99},
			{"Name": "Craque", "LastID": 33, "Score": 98}]`)
		defer cleanDatabase()

		store, err := NewFSStore(database)
		if err != nil {
			log.Fatalf("Almanac from reader: problem creating file system service store, %v ", err)
		}

		assertNoError(t, err)

		got := store.GetAlmanac()

		// This comes back sorted
		want := []WMService{
			{"Craque", 33, 98},
			{"Mattic", 10, 99},
		}

		/*	 OK! TIME TO FIX THE ALMANAC HERE	*/
		// it needs to return the three values,
		// and that needs to be set with runVerification
		// when it handles SvcTestDB and can add to WMService
		assertAlmanac(t, got, want)

		// read again
		got = store.GetAlmanac()
		assertAlmanac(t, got, want)

	})

	// Get a service LastID from the almanac database
	t.Run("get service LastID", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Mattic", "LastID": 10, "Score": 99},
			{"Name": "Craque", "LastID": 33, "Score": 98}]`)
		defer cleanDatabase()

		store, err := NewFSStore(database)
		if err != nil {
			log.Fatalf("Get LastID: problem creating file system service store, %v ", err)
		}

		assertNoError(t, err)

		got := store.GetTriggerID("Craque")
		want := 33
		assertIDEquals(t, got, want)
	})

	// This now needs to take WMService.Score
	t.Run("store LastID for existing services", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Mattic", "LastID": 10, "Score": 99},
			{"Name": "Craque", "LastID": 33, "Score": 98}]`)
		defer cleanDatabase()

		//store := FSStore{database}
		store, err := NewFSStore(database)
		if err != nil {
			log.Fatalf("Store LastID: problem creating file system service store, %v ", err)
		}

		assertNoError(t, err)

		craqueScore := 99
		store.TriggerID("Craque", craqueScore)

		got := store.GetTriggerID("Craque")
		want := 34
		assertIDEquals(t, got, want)
	})

	t.Run("store Score for existing services", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Mattic", "LastID": 10, "Score": 99},
			{"Name": "Craque", "LastID": 33, "Score": 98}]`)
		defer cleanDatabase()

		store, err := NewFSStore(database)
		if err != nil {
			log.Fatalf("Store LastID: problem creating file system service store, %v ", err)
		}

		assertNoError(t, err)

		craqueScore := 99
		store.TriggerID("Craque", craqueScore)

		got := store.GetScore("Craque")
		want := 99
		assertIDEquals(t, got, want)
	})

	t.Run("store LastID for new services", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Mattic", "LastID": 10, "Score": 99},
			{"Name": "Craque", "LastID": 33, "Score": 98}]`)
		defer cleanDatabase()

		store, err := NewFSStore(database)
		if err != nil {
			log.Fatalf("New services: problem creating file system service store, %v ", err)
		}

		store.TriggerID("Pepper", 0)

		got := store.GetTriggerID("Pepper")
		want := 1
		assertIDEquals(t, got, want)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		// "" creates an existing but empty file
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFSStore(database)

		assertNoError(t, err)
	})

	t.Run("almanac sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
			{"Name": "Mattic", "LastID": 10, "Score": 5},
			{"Name": "Craque", "LastID": 33, "Score": 4}]`)
		defer cleanDatabase()

		store, err := NewFSStore(database)

		assertNoError(t, err)

		got := store.GetAlmanac()
		want := Almanac{
			{"Craque", 33, 4},
			{"Mattic", 10, 5},
		}

		assertAlmanac(t, got, want)

		// read again
		got = store.GetAlmanac()
		assertAlmanac(t, got, want)

	})
}

func assertIDEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Unable to match ID, %d does not equal %d", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}
