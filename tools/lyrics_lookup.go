// Package tools provides lookup and search tools for the Nāda Guru knowledge base.
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/vpondala/nada-guru/knowledge"
)

// GetLyrics retrieves lyrics for a kriti, first from the embedded store,
// then from external sources if not found locally.
//
// ADK tool name: "get_lyrics"
// Parameters:    kritiID (string)
func GetLyrics(ctx context.Context, kritiID string) (*knowledge.Lyrics, error) {
	if store == nil {
		return nil, fmt.Errorf("knowledge store not initialized")
	}

	k, err := store.LookupKriti(kritiID)
	if err != nil {
		return nil, fmt.Errorf("kriti %q not found: %w", kritiID, err)
	}

	if k.LyricsFile == "" {
		return nil, fmt.Errorf("no lyrics file mapped for kriti %q", kritiID)
	}

	data, err := knowledge.ReadLyricsFile(k.LyricsFile)
	if err != nil {
		return nil, fmt.Errorf("read embedded lyrics %q: %w", k.LyricsFile, err)
	}

	var lyrics knowledge.Lyrics
	if err := json.Unmarshal(data, &lyrics); err != nil {
		return nil, fmt.Errorf("parse lyrics %q: %w", k.LyricsFile, err)
	}

	return &lyrics, nil
}

// ScrapeLyrics attempts to retrieve lyrics from external sources.
// Results are cached in session state for the duration of the session.
// Only called by GetLyrics when embedded lyrics are unavailable.
//
// ADK tool name: "scrape_lyrics"   (internal; not directly exposed to agents)
// Parameters:    kritiName (string), composerName (string)
func ScrapeLyrics(ctx context.Context, kritiName, composerName string) (*knowledge.Lyrics, error) {
	return nil, fmt.Errorf("lyrics scraping not yet implemented for %q by %q", kritiName, composerName)
}

// KritiIDFromFilename derives a kriti ID from a lyrics filename.
func KritiIDFromFilename(filename string) string {
	name := strings.TrimSuffix(filename, ".json")
	return strings.ReplaceAll(name, "-", "_")
}
