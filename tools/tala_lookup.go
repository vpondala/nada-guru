// Package tools provides lookup and search tools for the Nāda Guru knowledge base.
package tools

import (
	"context"
	"fmt"

	"github.com/vpondala/nada-guru/knowledge"
)

// LookupTala returns a Tala by ID or name (case-insensitive).
//
// ADK tool name: "lookup_tala"
// Parameters:    name (string)
func LookupTala(ctx context.Context, name string) (*knowledge.Tala, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}
	return store.LookupTala(name)
}

// SearchTalasByBeats returns all talas with the given total beat count.
//
// ADK tool name: "search_talas_by_beats"
// Parameters:    beats (int)
func SearchTalasByBeats(ctx context.Context, beats int) ([]knowledge.Tala, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}
	results := store.SearchTalasByBeats(beats)
	if len(results) == 0 {
		return nil, fmt.Errorf("no talas found with %d beats", beats)
	}
	return results, nil
}
