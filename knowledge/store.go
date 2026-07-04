// Package knowledge provides the in-memory knowledge store built from embedded JSON.
package knowledge

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"
)

// KnowledgeStore is the in-memory index built at startup from embedded JSON.
type KnowledgeStore struct {
	Ragas     []Raga
	Talas     []Tala
	Kritis    []Kriti
	Composers []Composer

	RagaByID      map[string]*Raga
	RagaByAlias   map[string]*Raga
	TalaByID      map[string]*Tala
	KritiByID     map[string]*Kriti
	ComposerByID  map[string]*Composer

	KritisByRaga     map[string][]string
	KritisByTala     map[string][]string
	KritisByComposer map[string][]string
	KritisByLanguage map[string][]string
}

// New loads and validates all embedded JSON, builds indexes, and returns a Store.
func New() (*KnowledgeStore, error) {
	store := &KnowledgeStore{}

	if err := store.loadRagas(); err != nil {
		return nil, fmt.Errorf("load ragas: %w", err)
	}
	if err := store.loadTalas(); err != nil {
		return nil, fmt.Errorf("load talas: %w", err)
	}
	if err := store.loadKritis(); err != nil {
		return nil, fmt.Errorf("load kritis: %w", err)
	}
	if err := store.loadComposers(); err != nil {
		return nil, fmt.Errorf("load composers: %w", err)
	}

	store.buildIndexes()
	if err := store.Validate(); err != nil {
		return nil, err
	}

	slog.Info("knowledge store loaded",
		"ragas", len(store.Ragas),
		"talas", len(store.Talas),
		"kritis", len(store.Kritis),
		"composers", len(store.Composers),
	)

	return store, nil
}

func (s *KnowledgeStore) loadRagas() error {
	var ragas []Raga
	if err := json.Unmarshal(RagasJSON, &ragas); err != nil {
		return err
	}
	s.Ragas = ragas
	return nil
}

func (s *KnowledgeStore) loadTalas() error {
	var talas []Tala
	if err := json.Unmarshal(TalasJSON, &talas); err != nil {
		return err
	}
	s.Talas = talas
	return nil
}

func (s *KnowledgeStore) loadKritis() error {
	var kritis []Kriti
	if err := json.Unmarshal(KritisJSON, &kritis); err != nil {
		return err
	}
	s.Kritis = kritis
	return nil
}

func (s *KnowledgeStore) loadComposers() error {
	var composers []Composer
	if err := json.Unmarshal(ComposersJSON, &composers); err != nil {
		return err
	}
	s.Composers = composers
	return nil
}

func (s *KnowledgeStore) buildIndexes() {
	s.RagaByID = make(map[string]*Raga)
	s.RagaByAlias = make(map[string]*Raga)
	for i := range s.Ragas {
		r := &s.Ragas[i]
		s.RagaByID[strings.ToLower(r.ID)] = r
		s.RagaByAlias[strings.ToLower(r.Name)] = r
		for _, a := range r.Aliases {
			s.RagaByAlias[strings.ToLower(a)] = r
		}
	}

	s.TalaByID = make(map[string]*Tala)
	for i := range s.Talas {
		t := &s.Talas[i]
		s.TalaByID[strings.ToLower(t.ID)] = t
		s.TalaByID[strings.ToLower(t.Name)] = t
	}

	s.KritiByID = make(map[string]*Kriti)
	s.KritisByRaga = make(map[string][]string)
	s.KritisByTala = make(map[string][]string)
	s.KritisByComposer = make(map[string][]string)
	s.KritisByLanguage = make(map[string][]string)
	for i := range s.Kritis {
		k := &s.Kritis[i]
		key := strings.ToLower(k.ID)
		s.KritiByID[key] = k
		s.KritisByRaga[strings.ToLower(k.Ragam)] = append(s.KritisByRaga[strings.ToLower(k.Ragam)], k.ID)
		s.KritisByTala[strings.ToLower(k.Talam)] = append(s.KritisByTala[strings.ToLower(k.Talam)], k.ID)
		s.KritisByComposer[strings.ToLower(k.Composer)] = append(s.KritisByComposer[strings.ToLower(k.Composer)], k.ID)
		s.KritisByLanguage[strings.ToLower(k.Language)] = append(s.KritisByLanguage[strings.ToLower(k.Language)], k.ID)
	}
	sort.Strings(s.KritisByRaga[strings.ToLower("Sri")])

	s.ComposerByID = make(map[string]*Composer)
	for i := range s.Composers {
		c := &s.Composers[i]
		s.ComposerByID[strings.ToLower(c.ID)] = c
	}
}

// Validate checks counts and required fields per REQ-012.
func (s *KnowledgeStore) Validate() error {
	if len(s.Ragas) < 72 {
		return fmt.Errorf("expected at least 72 ragas, got %d", len(s.Ragas))
	}
	for _, r := range s.Ragas {
		if r.Type == "melakarta" && r.Arohana == nil {
			return fmt.Errorf("melakarta raga %q has empty arohana", r.ID)
		}
	}
	for _, k := range s.Kritis {
		if k.Ragam == "" {
			return fmt.Errorf("kriti %q has empty ragam", k.ID)
		}
	}
	return nil
}

// LookupRaga returns a Raga by exact ID or alias (case-insensitive).
func (s *KnowledgeStore) LookupRaga(name string) (*Raga, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if r, ok := s.RagaByID[key]; ok {
		return r, nil
	}
	if r, ok := s.RagaByAlias[key]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("raga %q not found", name)
}

// LookupTala returns a Tala by ID or name (case-insensitive).
func (s *KnowledgeStore) LookupTala(name string) (*Tala, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if t, ok := s.TalaByID[key]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("tala %q not found", name)
}

// LookupKriti returns a Kriti by exact ID (case-insensitive).
func (s *KnowledgeStore) LookupKriti(id string) (*Kriti, error) {
	key := strings.ToLower(strings.TrimSpace(id))
	if k, ok := s.KritiByID[key]; ok {
		return k, nil
	}
	return nil, fmt.Errorf("kriti %q not found", id)
}

// LookupComposer returns a Composer by ID or name (case-insensitive).
func (s *KnowledgeStore) LookupComposer(name string) (*Composer, error) {
	key := strings.ToLower(strings.TrimSpace(name))
	if c, ok := s.ComposerByID[key]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("composer %q not found", name)
}

// SearchRagasBySwara returns all ragas whose arohana or avarohana contains
// all the given swaras as a subsequence.
func (s *KnowledgeStore) SearchRagasBySwara(swaras []string) []Raga {
	normalized := make([]string, len(swaras))
	for i, sw := range swaras {
		normalized[i] = strings.ToLower(strings.TrimSpace(sw))
	}
	var matches []Raga
	for i := range s.Ragas {
		r := &s.Ragas[i]
		if matchesSwaras(r.Arohana, normalized) || matchesSwaras(r.Avarohana, normalized) {
			matches = append(matches, *r)
		}
	}
	return matches
}

func matchesSwaras(seq, swaras []string) bool {
	if len(swaras) == 0 {
		return true
	}
	if len(seq) < len(swaras) {
		return false
	}
	si := 0
	for _, s := range seq {
		if si < len(swaras) && strings.ToLower(s) == swaras[si] {
			si++
		}
	}
	return si == len(swaras)
}

// SearchRagasByMood returns ragas matching the given rasa and/or time of day.
func (s *KnowledgeStore) SearchRagasByMood(rasa, timeOfDay string) []Raga {
	var matches []Raga
	rasaLower := strings.ToLower(strings.TrimSpace(rasa))
	timeLower := strings.ToLower(strings.TrimSpace(timeOfDay))
	for i := range s.Ragas {
		r := &s.Ragas[i]
		if rasaLower != "" {
			found := false
			for _, rs := range r.Rasa {
				if strings.ToLower(rs) == rasaLower {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if timeLower != "" && strings.ToLower(r.TimeOfDay) != timeLower {
			continue
		}
		matches = append(matches, *r)
	}
	return matches
}

// SearchTalasByBeats returns all talas with the given total beat count.
func (s *KnowledgeStore) SearchTalasByBeats(beats int) []Tala {
	var matches []Tala
	for i := range s.Talas {
		t := &s.Talas[i]
		if t.TotalBeats == beats {
			matches = append(matches, *t)
		}
	}
	return matches
}

// SearchKritis returns all kritis matching the given filter.
func (s *KnowledgeStore) SearchKritis(filter KritiFilter) ([]Kriti, error) {
	if filter.Ragam == "" && filter.Talam == "" && filter.Composer == "" && filter.Language == "" && filter.Tag == "" {
		return nil, fmt.Errorf("at least one filter field must be provided")
	}
	ragamKey := strings.ToLower(filter.Ragam)
	talamKey := strings.ToLower(filter.Talam)
	composerKey := strings.ToLower(filter.Composer)
	languageKey := strings.ToLower(filter.Language)

	ids := s.KritisByRaga[ragamKey]
	if filter.Talam != "" {
		ids = intersect(ids, s.KritisByTala[talamKey])
	}
	if filter.Composer != "" {
		ids = intersect(ids, s.KritisByComposer[composerKey])
	}
	if filter.Language != "" {
		ids = intersect(ids, s.KritisByLanguage[languageKey])
	}

	var results []Kriti
	for _, id := range ids {
		if k, ok := s.KritiByID[strings.ToLower(id)]; ok {
			if filter.Tag != "" {
				found := false
				for _, t := range k.Tags {
					if strings.ToLower(t) == strings.ToLower(filter.Tag) {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}
			results = append(results, *k)
		}
	}
	return results, nil
}

func intersect(a, b []string) []string {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	set := make(map[string]bool)
	for _, v := range a {
		set[v] = true
	}
	var out []string
	for _, v := range b {
		if set[v] {
			out = append(out, v)
		}
	}
	return out
}

// SearchComposersByLanguage returns all composers who composed in the given language.
func (s *KnowledgeStore) SearchComposersByLanguage(language string) []Composer {
	key := strings.ToLower(strings.TrimSpace(language))
	var matches []Composer
	for i := range s.Composers {
		c := &s.Composers[i]
		for _, lang := range c.Language {
			if strings.ToLower(lang) == key {
				matches = append(matches, *c)
				break
			}
		}
	}
	return matches
}
