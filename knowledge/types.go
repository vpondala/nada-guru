// Package knowledge holds all domain types for the Nāda Guru knowledge base.
package knowledge

// Raga represents either a Melakarta (parent) or Janya (derived) raga.
type Raga struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	MelakarataNumber int     `json:"melakarta_number"`
	ParentMelakarta int      `json:"parent_melakarta"`
	Chakra          string   `json:"chakra"`
	Madhyama        string   `json:"madhyama"`
	Rishabha        string   `json:"rishabha"`
	Gandhara        string   `json:"gandhara"`
	Dhaivata        string   `json:"dhaivata"`
	Nishada         string   `json:"nishada"`
	Arohana         []string `json:"arohana"`
	Avarohana       []string `json:"avarohana"`
	Aliases         []string `json:"aliases"`
	Rasa            []string `json:"rasa"`
	TimeOfDay       string   `json:"time_of_day"`
	Description     string   `json:"description"`
	JanyaRagas      []string `json:"janya_ragas"`
}

// Tala represents a rhythmic cycle used in Carnatic music.
type Tala struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Family              string   `json:"family"`
	Jati                string   `json:"jati"`
	Structure           string   `json:"structure"`
	Angas               []string `json:"angas"`
	TotalBeats          int      `json:"total_beats"`
	ClapPattern         string   `json:"clap_pattern"`
	Description         string   `json:"description"`
	CommonCompositions  []string `json:"common_compositions"`
}

// Kriti represents a single Carnatic composition.
type Kriti struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Composer    string   `json:"composer"`
	Ragam       string   `json:"ragam"`
	Talam       string   `json:"talam"`
	Language    string   `json:"language"`
	Script      string   `json:"script"`
	LyricsFile  string   `json:"lyrics_file"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

// Composer represents a Carnatic music composer.
type Composer struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	FullName           string   `json:"full_name"`
	Born               *int     `json:"born"`
	Died               *int     `json:"died"`
	Era                string   `json:"era"`
	Language           []string `json:"language"`
	Region             string   `json:"region"`
	Deity              string   `json:"deity"`
	NotableWorks       []string `json:"notable_works"`
	TotalCompositions  string   `json:"total_compositions"`
	Description        string   `json:"description"`
	FamousKritis       []string `json:"famous_kritis"`
}

// Lyrics holds the full structured text of a Carnatic composition.
type Lyrics struct {
	KritiID     string         `json:"kriti_id"`
	Ragam       string         `json:"ragam"`
	Talam       string         `json:"talam"`
	Composer    string         `json:"composer"`
	Language    string         `json:"language"`
	Script      string         `json:"script"`
	Pallavi     LyricsSection  `json:"pallavi"`
	Anupallavi  LyricsSection  `json:"anupallavi"`
	Charanams   []LyricsSection `json:"charanams"`
}

// LyricsSection is a single structural section of a kriti.
type LyricsSection struct {
	Number   int    `json:"number,omitempty"`
	Original string `json:"original"`
	IAST     string `json:"iast"`
	TeluguTr string `json:"transliteration_te,omitempty"`
	Meaning  string `json:"meaning"`
}

// KritiFilter defines optional search filters for kritis.
type KritiFilter struct {
	Ragam    string `json:"ragam,omitempty"`
	Talam    string `json:"talam,omitempty"`
	Composer string `json:"composer,omitempty"`
	Language string `json:"language,omitempty"`
	Tag      string `json:"tag,omitempty"`
}
