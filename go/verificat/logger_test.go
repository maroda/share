package main

import "testing"

/*
func TestLogger(t *testing.T) {
	t.Run("produces a log", func(t *testing.T) {

		// logErrToPass := NewError()
		want := ""
		got := Logger("There was an error", logErrToPass)

		assertLogBody(t, got, want)

	})
}
*/

func assertLogBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("log body is wrong, got %q want %q", got, want)
	}
}
