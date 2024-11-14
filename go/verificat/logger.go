package main

import (
	"context"
	"log/slog"
	"os"
)

// ContextHandler creates a custom Handler with attributes we want included with context
type ContextHandler struct {
	slog.Handler
}

// Handle is an slog.Handler with custom attributes as context values.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if app, ok := ctx.Value("app").(string); ok {
		r.AddAttrs(slog.String("app", app))
	}
	return h.Handler.Handle(ctx, r)
}

// Create a new logger with a configurable Level.
func createLogger(l slog.Level, a string) bool {
	// Variable for maintaining log levels
	var logLevel = new(slog.LevelVar)

	// Set log level.
	// This can be set dynamically anywhere using logLevel method Set().
	// Here we configure it with the passed-in slog.Level
	logLevel.Set(l)

	// Set default logging attributes
	defaultAttrs := []slog.Attr{
		// Name for the app as it appears in logs
		slog.String("app", a),
	}

	// Now we can configure these options in the Handler
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}

	// Build a new logger using custom handler for adding attributes to context.
	//
	//	baseHandler includes any slog.Attr key/value pairs defined in defaultAttrs
	baseHandler := slog.NewJSONHandler(os.Stdout, opts).WithAttrs(defaultAttrs)
	//
	//  customHandler inspects the log being written and adds context if it contains a matching defaultAttrs
	customHandler := &ContextHandler{Handler: baseHandler}
	ctx := context.Background()
	//
	//  creating a new logger with this configuration ensures the custom attributes appear in every log.
	logger := slog.New(customHandler)
	slog.SetDefault(logger)

	// Return a boolean of whether slog knows it is set to the passed-in level
	return slog.Default().Enabled(ctx, l)
}
