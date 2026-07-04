// Package server implements the HTTP REST server for Nāda Guru.
// It exposes POST /chat and GET /health endpoints.
package server

// Start begins listening on PORT (default 8080) and serves the REST API.
// It exits cleanly on SIGINT / SIGTERM.
//
// NOTE: Full implementation is in Task 6.2. This stub satisfies the Task 1.1
// acceptance criterion (go build ./... succeeds; --mode=server exits cleanly).
func Start() error {
	return nil
}
