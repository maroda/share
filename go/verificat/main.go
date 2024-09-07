package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
)

const dbFileName = "almanac.db.json"

func main() {
	runPort := "4330"

	// Logging setup
	defaultAttrs := []slog.Attr{
		slog.String("app", "verificat"),
		slog.String("env", "production"),
	}

	// Create logger with default attributes
	baseHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}).WithAttrs(defaultAttrs)

	// Custom attributes are then added by passing through ContextHandler
	// This determines if there are matching fields and adds them as log attributes if present.
	customHandler := &ContextHandler{Handler: baseHandler}

	// Now create the new logger
	logger := slog.New(customHandler)

	// Logging is a WIP.
	// This line is the only thing being logged presently.
	logger.Info("Starting Readiness Verification Service",
		slog.String("port", runPort),
		//slog.Group("payload",
		//	slog.String("username", "craquemattic"),
		//	slog.String("auth_method", "token"),
		//),
	)

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := NewFSStore(db)
	if err != nil {
		log.Fatalf("problem creating file system service store, %v ", err)
	}

	// VerificationServ is configured with a data storage source to hold run counts.
	server := NewVerificationServ(store)
	if err := http.ListenAndServe(":"+runPort, server); err != nil {
		slog.Error("Servercrash")
	}
}
