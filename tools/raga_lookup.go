// Package tools provides lookup and search tools for the Nāda Guru knowledge base.
package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/vpondala/nada-guru/knowledge"
)

var store *knowledge.KnowledgeStore

// Init sets the knowledge store used by all tool functions.
func Init(s *knowledge.KnowledgeStore) {
	store = s
}

// LookupRaga returns a Raga by exact ID or by alias (case-insensitive).
// Returns an error if not found, triggering Google Search fallback in the agent.
//
// ADK tool name: "lookup_raga"
// Parameters:    name (string) — raga name or alias
func LookupRaga(ctx context.Context, name string) (*knowledge.Raga, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}
	return store.LookupRaga(name)
}

// SearchRagasBySwara returns all ragas whose arohana or avarohana contains
// all the given swaras as a subsequence.
//
// ADK tool name: "search_ragas_by_swara"
// Parameters:    swaras ([]string) — e.g. ["S","R2","G3","P"]
func SearchRagasBySwara(ctx context.Context, swaras []string) ([]knowledge.Raga, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}
	if len(swaras) == 0 {
		return nil, fmt.Errorf("swaras must not be empty")
	}
	normalized := make([]string, len(swaras))
	for i, sw := range swaras {
		normalized[i] = strings.ToLower(strings.TrimSpace(sw))
	}
	results := store.SearchRagasBySwara(normalized)
	if len(results) == 0 {
		return nil, fmt.Errorf("no ragas found matching swaras %v", swaras)
	}
	return results, nil
}

// SearchRagasByMood returns ragas matching the given rasa and/or time of day.
//
// ADK tool name: "search_ragas_by_mood"
// Parameters:    rasa (string), timeOfDay (string) — either may be empty
func SearchRagasByMood(ctx context.Context, rasa, timeOfDay string) ([]knowledge.Raga, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}
	results := store.SearchRagasByMood(rasa, timeOfDay)
	if len(results) == 0 {
		return nil, fmt.Errorf("no ragas found matching rasa=%q timeOfDay=%q", rasa, timeOfDay)
	}
	return results, nil
}
