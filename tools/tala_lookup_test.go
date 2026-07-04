package tools

import (
	"context"
	"testing"

	"github.com/vpondala/nada-guru/knowledge"
)

func TestLookupTala(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	tala, err := LookupTala(context.Background(), "adi")
	if err != nil {
		t.Fatalf("LookupTala failed: %v", err)
	}
	if tala.Name != "Adi Talam" {
		t.Fatalf("expected Adi Talam, got %s", tala.Name)
	}
}

func TestLookupTala_NotFound(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	_, err = LookupTala(context.Background(), "unknown")
	if err == nil {
		t.Fatal("expected error for unknown tala")
	}
}

func TestSearchTalasByBeats(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	results, err := SearchTalasByBeats(context.Background(), 0)
	if err != nil {
		t.Fatalf("SearchTalasByBeats failed: %v", err)
	}
	if len(results) != 38 {
		t.Fatalf("expected 38 talas with 0 beats, got %d", len(results))
	}
}
