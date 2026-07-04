package tools

import (
	"context"
	"testing"

	"github.com/vpondala/nada-guru/knowledge"
)

func TestLookupComposer(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	c, err := LookupComposer(context.Background(), "tyagaraja")
	if err != nil {
		t.Fatalf("LookupComposer failed: %v", err)
	}
	if c.Name != "Tyagaraja" {
		t.Fatalf("expected Tyagaraja, got %s", c.Name)
	}
}

func TestLookupComposer_NotFound(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	_, err = LookupComposer(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent composer")
	}
}

func TestSearchComposersByLanguage(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	_, err = SearchComposersByLanguage(context.Background(), "Klingon")
	if err == nil {
		t.Fatal("expected error for unknown language")
	}
}
