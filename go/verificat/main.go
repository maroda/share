package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
)

const (
	app        = "verificat"
	dbFileName = "almanac.db.json"
	runPort    = "4330"          // TODO: this should be configurable
	llvl       = slog.LevelDebug // TODO: this should be configurable
)

// Main connects a local JSON database to a running API service.
func main() {
	// Create a new logger at Debug
	createLogger(llvl, app)

	// Log our presence to the world
	slog.Info("Starting Verificat: Weedmaps Production Readiness Verification",
		slog.String("port", runPort),
	)

	// Open JSON Database file
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	// Connect File System Storage operations to JSON Database
	store, err := NewFSStore(db)
	if err != nil {
		log.Fatalf("problem creating file system service store, %v ", err)
	}

	// A NewVerificationServ is configured with the database on local disk
	server := NewVerificationServ(store)
	if err := http.ListenAndServe(":"+runPort, server); err != nil {
		slog.Error("Servercrash")
	}
}
