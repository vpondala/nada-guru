// Package tools provides lookup and search tools for the Nāda Guru knowledge base.
package tools

import (
	"context"
	"fmt"
	"strings"
)

const supportedTargetLang = "te"

// TransliterationResult holds the output of a transliteration.
type TransliterationResult struct {
	Original       string   `json:"original"`
	Transliterated string   `json:"transliterated"`
	SourceLang     string   `json:"source_lang"`
	TargetLang     string   `json:"target_lang"`
	Notes          []string `json:"notes"`
}

// TransliterateText transliterates the given text from the source script into
// the target script using Gemini 3.1 Pro with Carnatic phoneme preservation.
//
// ADK tool name: "transliterate_text"
// Parameters:
//   text       (string)  — original text in source script
//   sourceLang (string)  — BCP-47 script tag: "sa", "ta", "kn"
//   targetLang (string)  — currently only "te" (Telugu) is supported
//
// Returns: TransliterationResult
func TransliterateText(ctx context.Context, text, sourceLang, targetLang string) (*TransliterationResult, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("text must not be empty")
	}
	if targetLang != supportedTargetLang {
		return nil, fmt.Errorf("unsupported target language %q; only %q is supported", targetLang, supportedTargetLang)
	}
	switch sourceLang {
	case "te", "sa", "ta", "kn":
	default:
		return nil, fmt.Errorf("unsupported source language %q", sourceLang)
	}
	if sourceLang == targetLang {
		return &TransliterationResult{
			Original:       text,
			Transliterated: text,
			SourceLang:     sourceLang,
			TargetLang:     targetLang,
			Notes:          []string{"source and target language are identical; text returned unchanged"},
		}, nil
	}
	return nil, fmt.Errorf("transliteration requires Gemini API integration")
}
