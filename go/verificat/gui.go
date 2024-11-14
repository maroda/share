package main

import (
	"embed"
	"io"
	"log/slog"
	"text/template"
)

// AlmanacWeb stores only data required for rendering the webpage
type AlmanacWeb struct {
	Title     string  // HTML Doc Title
	Content   string  // SVG XML
	FullScore Almanac // All I'm doing right now is printing the data, no fancy display yet.
}

// Load template directory
var (
	//go:embed templates/*
	htmlTmpl embed.FS
)

// RenderWeb takes an io.Writer, the content builder struct,
// the location of the HTML Templates to be used, and the target template for Execution.
func RenderWeb(w io.Writer, aw *AlmanacWeb, tmpldir, tdoc string) error {
	// Load the provided local directory for HTML templates
	tmpl, err := template.ParseFS(htmlTmpl, tmpldir)
	if err != nil {
		slog.Error("Templates could not be loaded", slog.Any("Error", err))
		return err
	}

	// Render and write the HTML using the Target Template doc and the AlmanacWeb struct for template fill.
	if err := tmpl.ExecuteTemplate(w, tdoc, aw); err != nil {
		slog.Error("Target Template Doc not found", slog.Any("Error", err))
		return err
	}

	return nil
}
