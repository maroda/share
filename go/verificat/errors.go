package main

import (
	"fmt"
)

type VError struct {
	Kind string // Where in the code is the error coming from?
	Err  error  // What is the error?
}

// An error function
func (e *VError) Error() string {
	// slog.Error("Achtung!", slog.String("Kind", e.Kind), slog.Any("Error", e.Err))
	return fmt.Sprintf("%v: %v", e.Kind, e.Err)
}
