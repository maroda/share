package main

import (
	"log/slog"
	"testing"
)

// completes with a validation that the level passed is Enabled()
func TestCreateLogger(t *testing.T) {
	level := slog.LevelDebug
	app := "verificat"
	got := createLogger(level, app)

	assertBool(t, got, true)
}

func assertBool(t testing.TB, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("log level is wrong, got %t want %t", got, want)
	}
}
