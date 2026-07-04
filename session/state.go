// Package session provides per-session conversational state for Nāda Guru.
package session

import (
	"context"
	"fmt"

	"github.com/vpondala/nada-guru/knowledge"
	"google.golang.org/adk/session"
)

// State holds all per-session conversational context (REQ-010).
type State struct {
	SessionID            string `json:"session_id"`
	LastRagaID           string `json:"last_raga_id"`
	LastTalaID           string `json:"last_tala_id"`
	LastComposerID       string `json:"last_composer_id"`
	LastKritiID          string `json:"last_kriti_id"`
	LastLyricsKritiID    string `json:"last_lyrics_kriti_id"`
	PreferredScriptTarget string `json:"preferred_script_target"`
	TurnCount            int    `json:"turn_count"`
}

const (
	DefaultScriptTarget = "te"
)

// Key constants for ADK session state map
const (
	KeyLastRaga     = "nada.last_raga_id"
	KeyLastTala     = "nada.last_tala_id"
	KeyLastComposer = "nada.last_composer_id"
	KeyLastKriti    = "nada.last_kriti_id"
	KeyLastLyrics   = "nada.last_lyrics_kriti_id"
	KeyScriptTarget = "nada.preferred_script_target"
)

// FromADKSession reads the current session state from an ADK session.
// It returns a zero-value State (not a panic) when the session is empty.
func FromADKSession(ctx context.Context, s session.Session) State {
	var st State
	st.SessionID = s.ID()
	raw, err := s.State().Get(KeyScriptTarget)
	if err == nil {
		if v, ok := raw.(string); ok && v != "" {
			st.PreferredScriptTarget = v
		}
	}
	if st.PreferredScriptTarget == "" {
		st.PreferredScriptTarget = DefaultScriptTarget
	}
	if v, err := s.State().Get(KeyLastRaga); err == nil {
		st.LastRagaID = fmt.Sprint(v)
	}
	if v, err := s.State().Get(KeyLastTala); err == nil {
		st.LastTalaID = fmt.Sprint(v)
	}
	if v, err := s.State().Get(KeyLastComposer); err == nil {
		st.LastComposerID = fmt.Sprint(v)
	}
	if v, err := s.State().Get(KeyLastKriti); err == nil {
		st.LastKritiID = fmt.Sprint(v)
	}
	if v, err := s.State().Get(KeyLastLyrics); err == nil {
		st.LastLyricsKritiID = fmt.Sprint(v)
	}
	return st
}

// ToADKSession writes the session state back to an ADK session.
func ToADKSession(ctx context.Context, s session.Session, st State) error {
	state := s.State()
	if err := state.Set(KeyScriptTarget, st.PreferredScriptTarget); err != nil {
		return fmt.Errorf("set script target: %w", err)
	}
	if err := state.Set(KeyLastRaga, st.LastRagaID); err != nil {
		return fmt.Errorf("set last raga: %w", err)
	}
	if err := state.Set(KeyLastTala, st.LastTalaID); err != nil {
		return fmt.Errorf("set last tala: %w", err)
	}
	if err := state.Set(KeyLastComposer, st.LastComposerID); err != nil {
		return fmt.Errorf("set last composer: %w", err)
	}
	if err := state.Set(KeyLastKriti, st.LastKritiID); err != nil {
		return fmt.Errorf("set last kriti: %w", err)
	}
	if err := state.Set(KeyLastLyrics, st.LastLyricsKritiID); err != nil {
		return fmt.Errorf("set last lyrics kriti: %w", err)
	}
	return nil
}

// ResolveRaga resolves a raga ID or name using the knowledge store.
func (st State) ResolveRaga(store *knowledge.KnowledgeStore) (*knowledge.Raga, error) {
	if st.LastRagaID != "" {
		if r, err := store.LookupRaga(st.LastRagaID); err == nil {
			return r, nil
		}
	}
	return nil, fmt.Errorf("no previous raga in session")
}
