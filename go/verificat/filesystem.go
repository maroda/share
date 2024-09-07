package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// It turns out *json.Encoder is all we need here
// because we're performing a lot of file ops in the constructor.
type FSStore struct {
	database *json.Encoder
	almanac  Almanac
}

// Constructor for FSStore
func NewFSStore(file *os.File) (*FSStore, error) {
	// initialize the database file for use with JSON
	err := initDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initialising service db file, %v", err)
	}

	// create a new Almanac object
	almanac, err := NewAlmanac(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading service store from file %s, %v", file.Name(), err)
	}

	// using the /tape/ type to encapsulate the database
	// allows us to always start at the beginning of the file
	return &FSStore{
		database: json.NewEncoder(&tape{file}),
		almanac:  almanac,
	}, nil
}

// initDBFile. If the database file is empty, initialize it for use with JSON.
func initDBFile(file *os.File) error {
	file.Seek(0, io.SeekStart)

	// get info and stat a file
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	// file.Stat returns stats on our file, which lets us check the size of the file.
	// If it's empty, we Write an empty JSON array and Seek back to the start.
	//
	// catch and format a zero-byte file
	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart)
	}
	return nil
}

// GetAlmanac. Provide a sorted list of all services and their LastID.
func (f *FSStore) GetAlmanac() Almanac {
	sort.Slice(f.almanac, func(i, j int) bool {
		return f.almanac[i].LastID > f.almanac[j].LastID
	})
	return f.almanac
}

// GetTriggerID. Lookup the LastID for a given name.
// The var /service/ is a WMService
func (f *FSStore) GetTriggerID(name string) int {
	service := f.almanac.Find(name)

	if service != nil {
		return service.LastID
	}

	return 0
}

// TriggerID. Increase LastID by one, providing a run count.
// If it's a new service, create them and start their tally at 1.
// The var /service/ is a WMService
func (f *FSStore) TriggerID(name string, score int) {
	service := f.almanac.Find(name)

	// TriggerID needs to set the Score, it is the "trigger" for things happening.

	// If the service is found, increment LastID.
	// If not, add them to this almanac with an initial LastID.
	if service != nil {
		service.LastID++
		service.Score = score
	} else {
		// Initialize the service in the almanac
		f.almanac = append(f.almanac, WMService{name, 1, 100})
	}

	// This seek call isn't needed here because we're seeking to the beginning
	// for writing thanks to the /tape/ type
	f.database.Encode(f.almanac)
}

// GetScore. Lookup the Score for a given name.
// The var /service/ is a WMService
func (f *FSStore) GetScore(name string) int {
	service := f.almanac.Find(name)

	if service != nil {
		return service.Score
	}

	return 0
}
