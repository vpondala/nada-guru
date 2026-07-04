// Package tools provides lookup and search tools for the Nāda Guru knowledge base.
package tools

import (
	"context"
	"fmt"

	"github.com/vpondala/nada-guru/knowledge"
)

// LookupComposer returns a Composer by ID or name (case-insensitive).
//
// ADK tool name: "lookup_composer"
// Parameters:    name (string)
func LookupComposer(ctx context.Context, name string) (*knowledge.Composer, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}
	return store.LookupComposer(name)
}

// SearchComposersByLanguage returns all composers who composed in the given language.
//
// ADK tool name: "search_composers_by_language"
// Parameters:    language (string) — e.g. "Telugu", "Sanskrit"
func SearchComposersByLanguage(ctx context.Context, language string) ([]knowledge.Composer, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}
	results := store.SearchComposersByLanguage(language)
	if len(results) == 0 {
		return nil, fmt.Errorf("no composers found for language %q", language)
	}
	return results, nil
}
