package web

import "embed"

// StaticFiles holds the embedded web UI static files.
// In production the dist directory is populated by the web app build step.
//
//go:embed all:dist
var StaticFiles embed.FS
