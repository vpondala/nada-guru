package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/vpondala/nada-guru/knowledge"
)

func TestGetLyrics(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	_, err = GetLyrics(context.Background(), "endaro_mahanubhavulu")
	if err == nil {
		t.Fatal("expected error when kriti has no lyrics_file mapping")
	}
	if !strings.Contains(err.Error(), "no lyrics file mapped") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetLyrics_NotFound(t *testing.T) {
	store, err := knowledge.New()
	if err != nil {
		t.Fatalf("knowledge.New() failed: %v", err)
	}
	Init(store)

	_, err = GetLyrics(context.Background(), "nonexistent_kriti")
	if err == nil {
		t.Fatal("expected error for nonexistent kriti")
	}
}

func TestScrapeLyrics(t *testing.T) {
	_, err := ScrapeLyrics(context.Background(), "test", "test")
	if err == nil {
		t.Fatal("expected error from ScrapeLyrics stub")
	}
}
