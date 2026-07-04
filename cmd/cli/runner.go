// Package cli implements the interactive CLI read-eval-print loop for Nāda Guru.
package cli

// Run starts the interactive CLI loop.
// It reads user input from stdin, sends each line to the root orchestrator
// agent, and prints the response with a "Guru > " prefix.
// It exits cleanly on "quit" / "exit" input or SIGINT.
//
// NOTE: Full implementation is in Task 6.1. This stub satisfies the Task 1.1
// acceptance criterion (go build ./... succeeds; --mode=cli exits cleanly).
func Run() error {
	return nil
}
