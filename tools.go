//go:build tools

// Package main holds tool dependency pins so go mod tidy retains them
// before any production code imports them directly.
package main

import (
	_ "google.golang.org/adk/agent/llmagent"
	_ "google.golang.org/genai"
)
