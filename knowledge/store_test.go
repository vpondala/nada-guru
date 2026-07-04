package knowledge

import (
	"testing"
)

func TestValidate_PassesOnGoodData(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	if err := store.Validate(); err != nil {
		t.Fatalf("Validate() failed: %v", err)
	}
}

func TestMelakarataCount(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	count := 0
	for _, r := range store.Ragas {
		if r.Type == "melakarta" {
			count++
		}
	}
	if count != 72 {
		t.Fatalf("expected 72 Melakarta ragas, got %d", count)
	}
}

func TestJanyaCount(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	count := 0
	for _, r := range store.Ragas {
		if r.Type == "janya" {
			count++
		}
	}
	if count < 20 {
		t.Fatalf("expected at least 20 Janya ragas, got %d", count)
	}
}

func TestAliasLookup_Kalyani(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	r, err := store.LookupRaga("mechakalyani")
	if err != nil {
		t.Fatalf("LookupRaga(\"mechakalyani\") failed: %v", err)
	}
	if r.Name != "Mechakalyani" {
		t.Fatalf("expected Mechakalyani, got %s", r.Name)
	}
	if r.MelakarataNumber != 65 {
		t.Fatalf("expected Melakarta number 65, got %d", r.MelakarataNumber)
	}
}

func TestAliasLookup_Todi(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	r, err := store.LookupRaga("hanumatodi")
	if err != nil {
		t.Fatalf("LookupRaga(\"hanumatodi\") failed: %v", err)
	}
	if r.Name != "Hanumatodi" {
		t.Fatalf("expected Hanumatodi, got %s", r.Name)
	}
	if r.MelakarataNumber != 8 {
		t.Fatalf("expected Melakarta number 8, got %d", r.MelakarataNumber)
	}
}

func TestKritisByRaga_Sri(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	ids := store.KritisByRaga["sri"]
	if len(ids) == 0 {
		t.Fatalf("expected at least 1 kriti for raga 'Sri', got 0")
	}
}

func TestKritisByComposer_Tyagaraja(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	ids := store.KritisByComposer["tyagaraja"]
	if len(ids) < 5 {
		t.Fatalf("expected at least 5 kritis for composer 'tyagaraja', got %d", len(ids))
	}
}

func TestLookupTala(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	tala, err := store.LookupTala("adi")
	if err != nil {
		t.Fatalf("LookupTala(\"adi\") failed: %v", err)
	}
	if tala.Name != "Adi Talam" {
		t.Fatalf("expected Adi Talam, got %s", tala.Name)
	}
}

func TestLookupKriti(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	k, err := store.LookupKriti("endaro_mahanubhavulu")
	if err != nil {
		t.Fatalf("LookupKriti(\"endaro_mahanubhavulu\") failed: %v", err)
	}
	if k.Ragam != "Sri" {
		t.Fatalf("expected ragam Sri, got %s", k.Ragam)
	}
}

func TestLookupComposer(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	c, err := store.LookupComposer("tyagaraja")
	if err != nil {
		t.Fatalf("LookupComposer(\"tyagaraja\") failed: %v", err)
	}
	if c.Name != "Tyagaraja" {
		t.Fatalf("expected Tyagaraja, got %s", c.Name)
	}
}

func TestSearchRagasBySwara(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	results := store.SearchRagasBySwara([]string{"R2", "G3", "M2"})
	found := false
	for _, r := range results {
		if r.Name == "Mechakalyani" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected Mechakalyani in SearchRagasBySwara([\"R2\",\"G3\",\"M2\"]) results")
	}
}

func TestSearchRagasByMood(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	results := store.SearchRagasByMood("bhakti", "morning")
	if len(results) < 3 {
		t.Fatalf("expected at least 3 ragas for mood bhakti/morning, got %d", len(results))
	}
}

func TestSearchTalasByBeats(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	results := store.SearchTalasByBeats(0)
	if len(results) != 38 {
		t.Fatalf("expected 38 talas with 0 beats, got %d", len(results))
	}
}

func TestSearchKritis(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	results, err := store.SearchKritis(KritiFilter{Composer: "tyagaraja"})
	if err != nil {
		t.Fatalf("SearchKritis failed: %v", err)
	}
	if len(results) < 5 {
		t.Fatalf("expected at least 5 kritis for composer tyagaraja, got %d", len(results))
	}
}

func TestSearchKritis_RequiresFilter(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	_, err = store.SearchKritis(KritiFilter{})
	if err == nil {
		t.Fatal("expected error when all filter fields are empty")
	}
}

func TestSearchKritis_ByLanguage(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	results, err := store.SearchKritis(KritiFilter{Language: "Telugu"})
	if err != nil {
		t.Fatalf("SearchKritis by language failed: %v", err)
	}
	if len(results) == 0 {
		t.Fatalf("expected at least 1 Telugu kriti, got 0")
	}
	for _, k := range results {
		if k.Language != "Telugu" {
			t.Fatalf("expected Telugu kriti, got %s", k.Language)
		}
	}
}

func TestSearchKritis_MultiFilter(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	results, err := store.SearchKritis(KritiFilter{Composer: "tyagaraja", Language: "Telugu"})
	if err != nil {
		t.Fatalf("SearchKritis multi-filter failed: %v", err)
	}
	for _, k := range results {
		if k.Composer != "tyagaraja" || k.Language != "Telugu" {
			t.Fatalf("unexpected kriti in results: %+v", k)
		}
	}
}

func TestSearchComposersByLanguage(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	results := store.SearchComposersByLanguage("Klingon")
	if len(results) != 0 {
		t.Fatalf("expected 0 composers for unknown language, got %d", len(results))
	}
}

func TestReadLyricsFile(t *testing.T) {
	data, err := ReadLyricsFile("endaro_mahanubhavulu.json")
	if err != nil {
		t.Fatalf("ReadLyricsFile failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty lyrics data")
	}
}

func TestReadLyricsFile_NotFound(t *testing.T) {
	_, err := ReadLyricsFile("nonexistent.json")
	if err == nil {
		t.Fatal("expected error for nonexistent lyrics file")
	}
}

func TestLyricsFilenames(t *testing.T) {
	names, err := LyricsFilenames()
	if err != nil {
		t.Fatalf("LyricsFilenames failed: %v", err)
	}
	if len(names) != 7 {
		t.Fatalf("expected 7 lyrics files, got %d", len(names))
	}
}

func TestLookupRaga_NotFound(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	_, err = store.LookupRaga("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent raga")
	}
}

func TestLookupTala_NotFound(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	_, err = store.LookupTala("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent tala")
	}
}

func TestLookupKriti_NotFound(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	_, err = store.LookupKriti("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent kriti")
	}
}

func TestLookupComposer_NotFound(t *testing.T) {
	store, err := New()
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	_, err = store.LookupComposer("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent composer")
	}
}
