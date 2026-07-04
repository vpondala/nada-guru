package tools

import (
	"context"
	"testing"

	"github.com/vpondala/nada-guru/knowledge"
)

func TestLookupRaga(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	r, err := LookupRaga(context.Background(), "mechakalyani")
	if err != nil {
		t.Fatalf("LookupRaga failed: %v", err)
	}
	if r.Name != "Mechakalyani" {
		t.Fatalf("expected Mechakalyani, got %s", r.Name)
	}
}

func TestLookupRaga_NotFound(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	_, err = LookupRaga(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent raga")
	}
}

func TestSearchRagasBySwara(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	results, err := SearchRagasBySwara(context.Background(), []string{"R2", "G3", "M2"})
	if err != nil {
		t.Fatalf("SearchRagasBySwara failed: %v", err)
	}
	found := false
	for _, r := range results {
		if r.Name == "Mechakalyani" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected Mechakalyani in swara search results")
	}
}

func TestSearchRagasByMood(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	results, err := SearchRagasByMood(context.Background(), "bhakti", "morning")
	if err != nil {
		t.Fatalf("SearchRagasByMood failed: %v", err)
	}
	if len(results) < 3 {
		t.Fatalf("expected at least 3 ragas for mood bhakti/morning, got %d", len(results))
	}
}
