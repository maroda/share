package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// Anything that used to return a []WMService
// can instead return a Almanac
type Almanac []WMService

// Will take an io.Reader like /database/ and decode the result.
func NewAlmanac(rdr io.Reader) ([]WMService, error) {
	var almanac []WMService
	err := json.NewDecoder(rdr).Decode(&almanac)
	if err != nil {
		err = fmt.Errorf("problem parsing almanac, %v", err)
	}
	return almanac, err
}

// Takes a Almanac and searches for a service name
func (l Almanac) Find(name string) *WMService {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}
