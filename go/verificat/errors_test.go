package main

import (
	"errors"
	"testing"
)

// Handle errors
func TestError(t *testing.T) {

	my := &VError{Kind: "Test", Err: errors.New("craquemattic")}
	want := "Test: craquemattic"
	got := my.Error()

	assertString(t, got, want)
}
