// Package tools provides lookup and search tools for the Nāda Guru knowledge base.
package tools

import (
	"context"
	"fmt"

	"github.com/vpondala/nada-guru/knowledge"
)

// SearchKritis returns all kritis matching the given filter.
// At least one filter field must be non-empty.
//
// ADK tool name: "search_kritis"
// Parameters:    filter (KritiFilter)
func SearchKritis(ctx context.Context, filter knowledge.KritiFilter) ([]knowledge.Kriti, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}
	results, err := store.SearchKritis(filter)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no kritis found matching filter %+v", filter)
	}
	return results, nil
}

// LookupKriti returns a single Kriti by exact ID (case-insensitive).
//
// ADK tool name: "lookup_kriti"
// Parameters:    id (string)
func LookupKriti(ctx context.Context, id string) (*knowledge.Kriti, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}
	return store.LookupKriti(id)
}
