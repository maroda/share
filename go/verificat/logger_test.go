package main

import (
	"context"
	"log/slog"
	"testing"
)

// WIP Logger function
func TestLogInfoMsg(t *testing.T) {
	ctx := context.Background()

	t.Run("produces a log", func(t *testing.T) {
		logMessage := "Craque Log Message"

		LogInfoMsg(ctx, logMessage)
	})
}

// completes with a validation that the level passed is Enabled()
func TestCreateLogger(t *testing.T) {
	level := slog.LevelDebug
	app := "verificat"
	want := true
	got := createLogger(level, app)

	assertBool(t, got, want)
}

func assertBool(t testing.TB, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("log level is wrong, got %t want %t", got, want)
	}
}
