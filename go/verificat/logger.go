package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// Create a custom Handler with attributes we want included with context
type ContextHandler struct {
	slog.Handler
}

// slog.Handler with custom attributes as context values.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if app, ok := ctx.Value("app").(string); ok {
		r.AddAttrs(slog.String("app", app))
	}
	return h.Handler.Handle(ctx, r)
}

/*

	With this configuration,
	any existing method calls using /log/ - like 'log.Printf()'
	will adopt the "implicit" (i.e. no context) slog configuration.
	That's enough for a structured logging MVP without getting into
	testing a Logging() function... yet ;)

*/
// Setup slog with init()
func init() {
	// Variable for maintaining log levels
	var logLevel = new(slog.LevelVar)

	// Set default logging attributes
	defaultAttrs := []slog.Attr{
		slog.String("app", "verificat"),
	}

	// New logger, log level options using a logLevel that can be set dynamically
	// baselogger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}).WithAttrs(defaultAttrs))

	// Build a new logger using custom handler for adding attributes to context.
	//
	//	baseHandler includes any slog.Attr key/value pairs defined in defaultAttrs
	baseHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel}).WithAttrs(defaultAttrs)
	//  customHandler inspects the log being written and adds context if it contains a matching defaultAttrs
	customHandler := &ContextHandler{Handler: baseHandler}
	//  creating a new logger with this configuration ensures the custom attributes appear in every log.
	logger := slog.New(customHandler)
	slog.SetDefault(logger)

	// Set log level. This can be set dynamically anywhere using logLevel method Set().
	logLevel.Set(slog.LevelDebug)

	/* implicit log levels look like this, use contexts to set explicit log levels
	slog.Debug("Arcane", "Code", "Included")
	slog.Info("JSON Handler", "Content", "Logging in JSON")
	slog.Warn("Article", "Read", "Till the end")
	slog.Error("Achtung!", "There be", "Dragons")
	*/
}

// Future Logger function, currently WIP
func Logger() error {
	ctx := context.Background()
	slog.InfoContext(ctx, "Doing something", slog.String("method", "manual"))
	fmt.Println("craquemattic")
	return nil
}
