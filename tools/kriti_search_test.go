package tools

import (
	"context"
	"testing"

	"github.com/vpondala/nada-guru/knowledge"
)

func TestSearchKritis(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	results, err := SearchKritis(context.Background(), knowledge.KritiFilter{Composer: "tyagaraja"})
	if err != nil {
		t.Fatalf("SearchKritis failed: %v", err)
	}
	if len(results) < 5 {
		t.Fatalf("expected at least 5 kritis, got %d", len(results))
	}
}

func TestSearchKritis_RequiresFilter(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	_, err = SearchKritis(context.Background(), knowledge.KritiFilter{})
	if err == nil {
		t.Fatal("expected error when all filter fields are empty")
	}
}

func TestLookupKriti(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	k, err := LookupKriti(context.Background(), "endaro_mahanubhavulu")
	if err != nil {
		t.Fatalf("LookupKriti failed: %v", err)
	}
	if k.Ragam != "Sri" {
		t.Fatalf("expected ragam Sri, got %s", k.Ragam)
	}
}

func TestLookupKriti_NotFound(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	_, err = LookupKriti(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent kriti")
	}
}
