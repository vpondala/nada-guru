package tools

import (
	"context"
	"testing"
)

func TestTransliterateText_SameLang(t *testing.T) {
	result, err := TransliterateText(context.Background(), "test", "te", "te")
	if err != nil {
		t.Fatalf("TransliterateText failed: %v", err)
	}
	if result.Transliterated != "test" {
		t.Fatalf("expected unchanged text, got %s", result.Transliterated)
	}
	if len(result.Notes) == 0 {
		t.Fatal("expected at least one note")
	}
}

func TestTransliterateText_UnsupportedTarget(t *testing.T) {
	_, err := TransliterateText(context.Background(), "text", "sa", "en")
	if err == nil {
		t.Fatal("expected error for unsupported target language")
	}
}

func TestTransliterateText_UnsupportedSource(t *testing.T) {
	_, err := TransliterateText(context.Background(), "text", "en", "te")
	if err == nil {
		t.Fatal("expected error for unsupported source language")
	}
}

func TestTransliterateText_EmptyText(t *testing.T) {
	_, err := TransliterateText(context.Background(), "", "sa", "te")
	if err == nil {
		t.Fatal("expected error for empty text")
	}
}
