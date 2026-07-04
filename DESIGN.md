# Nāda Guru — Technical Design

## 1. Overview

This document describes the technical design of the Nāda Guru multi-agent Carnatic music learning system. It is derived from [REQUIREMENTS.md](REQUIREMENTS.md) and serves as the authoritative reference for implementation. All Go types, agent configurations, tool signatures, API contracts, and data-flow sequences are defined here.

**Tech Stack:** Google ADK 2.0 (Go) · Gemini 3.1 Pro / 3.5 Flash · Google AI Studio · Antigravity · YouTube Data API v3 · yt-dlp

---

## 2. System Architecture

```
┌──────────────────────────────────────────────────────────────────────────┐
│                          NĀDA GURU SYSTEM                                │
│                                                                          │
│   ┌──────────┐    ┌──────────┐                                           │
│   │  CLI     │    │  HTTP    │  ← User Entry Points                      │
│   │ (stdin)  │    │ :8080    │                                           │
│   └────┬─────┘    └────┬─────┘                                           │
│        └───────┬────────┘                                                │
│                ▼                                                          │
│   ┌────────────────────────────┐                                         │
│   │   ROOT ORCHESTRATOR AGENT  │  gemini-3.1-pro                        │
│   │   (agents/orchestrator.go) │  Session State Manager                  │
│   └──┬───┬───┬───┬───┬───┬────┘                                         │
│      │   │   │   │   │   │                                               │
│   ┌──▼─┐ │ ┌─▼──┐│ ┌─▼──┐│ ┌──▼─┐ ┌──▼──┐ ┌──▼──┐ ┌──▼──────────┐    │
│   │Raga│ │ │Tala││ │Krti││ │Comp│ │Lyrc │ │Trns │ │YouTube      │    │
│   │Agt │ │ │Agt ││ │Agt ││ │Agt │ │Agt  │ │Agt  │ │Analyser Agt │    │
│   └──┬─┘ │ └─┬──┘│ └─┬──┘│ └──┬─┘ └──┬──┘ └──┬──┘ └──┬──────────┘    │
│      │   │   │   │   │   │    │      │       │       │                  │
│   ┌──▼───▼───▼───▼───▼───▼────▼──────▼───────▼───────▼──────────────┐  │
│   │                        TOOL LAYER                                 │  │
│   │  raga_lookup · tala_lookup · kriti_search · composer_lookup      │  │
│   │  lyrics_lookup · transliterate · youtube_fetch · gemini_audio    │  │
│   │  google_search (ADK built-in)                                    │  │
│   └──────────────────────────────┬───────────────────────────────────┘  │
│                                  │                                        │
│   ┌──────────────────────────────▼───────────────────────────────────┐  │
│   │                     KNOWLEDGE LAYER                               │  │
│   │           (//go:embed — compiled into binary)                     │  │
│   │   ragas.json · talas.json · kritis.json · composers.json         │  │
│   │   lyrics/*.json                                                   │  │
│   └──────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
          │                              │
          ▼                              ▼
   Gemini 3.1 Pro / 3.5 Flash API      YouTube Data API v3
   (Google AI Studio key)        yt-dlp (audio extraction)
```

---

## 3. Package Structure

```
nada-guru/
├── go.mod                          # module github.com/vpondala/nada-guru
├── go.sum
├── main.go                         # Entry: parses --mode flag, starts CLI or server
├── cmd/
│   ├── cli/runner.go               # Interactive CLI loop
│   └── server/server.go            # HTTP server, /chat and /health handlers
├── agents/
│   ├── orchestrator.go             # Root agent — llmagent.New() with all sub-agents
│   ├── raga_agent.go
│   ├── tala_agent.go
│   ├── kriti_agent.go
│   ├── composer_agent.go
│   ├── lyrics_agent.go
│   ├── transliteration_agent.go
│   └── youtube_analyser_agent.go
├── tools/
│   ├── raga_lookup.go
│   ├── tala_lookup.go
│   ├── kriti_search.go
│   ├── composer_lookup.go
│   ├── lyrics_lookup.go
│   ├── transliterate.go
│   ├── youtube_fetch.go
│   └── gemini_audio.go
├── knowledge/
│   ├── embed.go                    # //go:embed directives
│   ├── store.go                    # KnowledgeStore — load, validate, index
│   ├── types.go                    # All shared Go struct definitions
│   ├── ragas.json
│   ├── talas.json
│   ├── kritis.json
│   ├── composers.json
│   └── lyrics/
│       └── *.json
├── session/
│   └── state.go                    # SessionState struct + ADK session helpers
├── eval/
│   ├── eval_test.go                # go test runner
│   └── test_cases.json
└── README.md
```

---

## 4. Knowledge Base Data Models

All types are defined in `knowledge/types.go`.

### 4.1 Raga

```go
// Raga represents either a Melakarta (parent) or Janya (derived) raga.
type Raga struct {
    ID              string   `json:"id"`                // snake_case unique key, e.g. "mechakalyani"
    Name            string   `json:"name"`              // Display name, e.g. "Mechakalyani"
    Type            string   `json:"type"`              // "melakarta" | "janya"
    MelakarataNumber int     `json:"melakarta_number"`  // 1–72 (Melakarta only; 0 for Janya)
    ParentMelakarta int      `json:"parent_melakarta"`  // Melakarta number (Janya only; 0 for Melakarta)
    Chakra          string   `json:"chakra"`            // e.g. "Rudra" (Melakarta only)
    Madhyama        string   `json:"madhyama"`          // "M1" | "M2"
    Rishabha        string   `json:"rishabha"`          // "R1" | "R2" | "R3"
    Gandhara        string   `json:"gandhara"`          // "G1" | "G2" | "G3"
    Dhaivata        string   `json:"dhaivata"`          // "D1" | "D2" | "D3"
    Nishada         string   `json:"nishada"`           // "N1" | "N2" | "N3"
    Arohana         []string `json:"arohana"`           // ascending swara sequence
    Avarohana       []string `json:"avarohana"`         // descending swara sequence
    Aliases         []string `json:"aliases"`           // alternate names
    Rasa            []string `json:"rasa"`              // emotional qualities
    TimeOfDay       string   `json:"time_of_day"`       // "morning" | "evening" | "any" etc.
    Description     string   `json:"description"`       // prose summary
    JanyaRagas      []string `json:"janya_ragas"`       // names of derived ragas (Melakarta only)
}
```

### 4.2 Tala

```go
// Tala represents a rhythmic cycle used in Carnatic music.
type Tala struct {
    ID                  string   `json:"id"`                   // e.g. "adi"
    Name                string   `json:"name"`                 // e.g. "Adi Talam"
    Family              string   `json:"family"`               // "suladi_sapta" | "chapu"
    Jati                string   `json:"jati"`                 // "chatusra" | "tisra" | "misra" | "khanda"
    Structure           string   `json:"structure"`            // human-readable, e.g. "Laghu(4) + Drutam + Drutam"
    Angas               []string `json:"angas"`                // e.g. ["L4", "D", "D"]
    TotalBeats          int      `json:"total_beats"`          // e.g. 8
    ClapPattern         string   `json:"clap_pattern"`         // e.g. "1 2 3 4 | wave | wave"
    Description         string   `json:"description"`
    CommonCompositions  []string `json:"common_compositions"`  // kriti IDs
}
```

### 4.3 Kriti

```go
// Kriti represents a single Carnatic composition.
type Kriti struct {
    ID          string   `json:"id"`           // snake_case unique key, e.g. "endaro_mahanubhavulu"
    Name        string   `json:"name"`         // Display name
    Composer    string   `json:"composer"`     // composer ID (FK to composers.json)
    Ragam       string   `json:"ragam"`        // raga name (display name, not ID)
    Talam       string   `json:"talam"`        // tala name (display name, not ID)
    Language    string   `json:"language"`     // "Telugu" | "Sanskrit" | "Tamil" | "Kannada"
    Script      string   `json:"script"`       // BCP-47 script tag: "te" | "sa" | "ta" | "kn"
    LyricsFile  string   `json:"lyrics_file"`  // relative path within knowledge/, empty if unavailable
    Description string   `json:"description"`
    Tags        []string `json:"tags"`         // e.g. ["pancharatna", "popular", "bhakti"]
}
```

### 4.4 Composer

```go
// Composer represents a Carnatic music composer.
type Composer struct {
    ID                 string   `json:"id"`                  // snake_case, e.g. "tyagaraja"
    Name               string   `json:"name"`                // common name
    FullName           string   `json:"full_name"`
    Born               *int     `json:"born"`                // year, nil if unknown
    Died               *int     `json:"died"`                // year, nil if unknown
    Era                string   `json:"era"`                 // e.g. "18th–19th century"
    Language           []string `json:"language"`            // languages of composition
    Region             string   `json:"region"`
    Deity              string   `json:"deity"`
    NotableWorks       []string `json:"notable_works"`
    TotalCompositions  string   `json:"total_compositions"`  // prose estimate
    Description        string   `json:"description"`
    FamousKritis       []string `json:"famous_kritis"`       // kriti IDs
}
```

### 4.5 Lyrics

```go
// Lyrics holds the full structured text of a Carnatic composition.
type Lyrics struct {
    KritiID  string         `json:"kriti_id"`
    Ragam    string         `json:"ragam"`
    Talam    string         `json:"talam"`
    Composer string         `json:"composer"`
    Language string         `json:"language"`
    Script   string         `json:"script"`
    Pallavi  LyricsSection  `json:"pallavi"`
    Anupallavi LyricsSection `json:"anupallavi"`
    Charanams []LyricsSection `json:"charanams"`
}

// LyricsSection is a single structural section of a kriti.
type LyricsSection struct {
    Number   int    `json:"number,omitempty"`   // Charanam index (1-based); 0 for Pallavi/Anupallavi
    Original string `json:"original"`            // text in original script (Unicode)
    IAST     string `json:"iast"`                // romanised transliteration (IAST standard)
    TeluguTr string `json:"transliteration_te,omitempty"` // Telugu script transliteration (if pre-computed)
    Meaning  string `json:"meaning"`             // English meaning
}
```

### 4.6 KnowledgeStore

```go
// KnowledgeStore is the in-memory index built at startup from embedded JSON.
// Defined in knowledge/store.go.
type KnowledgeStore struct {
    Ragas     []Raga
    Talas     []Tala
    Kritis    []Kriti
    Composers []Composer

    // Indexes for O(1) lookup
    RagaByID      map[string]*Raga
    RagaByAlias   map[string]*Raga     // normalised alias → *Raga
    TalaByID      map[string]*Tala
    KritiByID     map[string]*Kriti
    ComposerByID  map[string]*Composer

    // Inverted indexes for search
    KritisByRaga     map[string][]string  // ragam name → []kriti IDs
    KritisByTala     map[string][]string  // tala name  → []kriti IDs
    KritisByComposer map[string][]string  // composer ID → []kriti IDs
    KritisByLanguage map[string][]string  // language → []kriti IDs
}

// New loads and validates all embedded JSON, builds indexes, and returns a Store.
// Returns an error if any required file is missing or unparseable.
func New() (*KnowledgeStore, error)

// Validate checks counts and required fields per REQ-012.
func (s *KnowledgeStore) Validate() error
```

---

## 5. Session State

Defined in `session/state.go`. Stored in ADK session state and serialised as JSON.

```go
// State holds all per-session conversational context (REQ-010).
type State struct {
    SessionID            string `json:"session_id"`
    LastRagaID           string `json:"last_raga_id"`            // last discussed raga
    LastTalaID           string `json:"last_tala_id"`
    LastComposerID       string `json:"last_composer_id"`
    LastKritiID          string `json:"last_kriti_id"`
    LastLyricsKritiID    string `json:"last_lyrics_kriti_id"`    // last lyrics retrieved
    PreferredScriptTarget string `json:"preferred_script_target"` // default "te" (Telugu)
    TurnCount            int    `json:"turn_count"`
}

// Key constants for ADK session state map
const (
    KeyLastRaga     = "nada.last_raga_id"
    KeyLastTala     = "nada.last_tala_id"
    KeyLastComposer = "nada.last_composer_id"
    KeyLastKriti    = "nada.last_kriti_id"
    KeyLastLyrics   = "nada.last_lyrics_kriti_id"
    KeyScriptTarget = "nada.preferred_script_target"
)
```

---

## 6. Agent Definitions

All agents are defined in `agents/` using `google.golang.org/adk/agent/llmagent`.

### 6.1 Root Orchestrator Agent

```go
// agents/orchestrator.go

// New returns the fully wired root orchestrator agent with all sub-agents
// registered as sub-agents and tools.
func New(store *knowledge.KnowledgeStore) (*llmagent.Agent, error)

// Model:       gemini-3.1-pro
// SubAgents:   RagaAgent, TalaAgent, KritiAgent, ComposerAgent,
//              LyricsAgent, TransliterationAgent, YouTubeAnalyserAgent
// Tools:       google_search (ADK built-in)
// Instruction: See prompt constant in orchestrator.go
```

**System Instruction (excerpt):**
```
You are Nāda Guru, an expert AI guide for Carnatic classical music.
You help learners understand Ragas, Talams, Kritis, Composers, Lyrics, and Transliterations.
Route each user query to the most appropriate specialist agent.
Use session state to resolve pronouns and references from prior turns.
If a query spans multiple domains (e.g. lyrics + transliteration), 
invoke specialist agents sequentially and combine their responses.
Always respond in the same language the user used.
```

### 6.2 Raga Agent

```go
// agents/raga_agent.go
// Model:  gemini-3.5-flash
// Tools:  LookupRaga, SearchRagasBySwara, SearchRagasByMood, google_search
// Description: "Answers questions about Carnatic ragas — arohana, avarohana,
//               vadi, samvadi, rasa, time of day, Melakarta classification,
//               and related compositions."
```

### 6.3 Tala Agent

```go
// agents/tala_agent.go
// Model:  gemini-3.5-flash
// Tools:  LookupTala, SearchTalasByBeats, google_search
// Description: "Answers questions about Carnatic talas — structure, angas,
//               beat counts, clap patterns, and example compositions."
```

### 6.4 Kriti Agent

```go
// agents/kriti_agent.go
// Model:  gemini-3.5-flash
// Tools:  SearchKritis, LookupKriti, google_search
// Description: "Finds Carnatic compositions by raga, tala, composer, or language.
//               Returns metadata and offers to fetch lyrics or transliterate."
```

### 6.5 Composer Agent

```go
// agents/composer_agent.go
// Model:  gemini-3.5-flash
// Tools:  LookupComposer, SearchComposersByLanguage, google_search
// Description: "Provides biographical and compositional details about
//               Carnatic music composers."
```

### 6.6 Lyrics Agent

```go
// agents/lyrics_agent.go
// Model:  gemini-3.5-flash
// Tools:  GetLyrics, ScrapeLyrics, google_search
// Description: "Retrieves full lyrics (Pallavi, Anupallavi, Charanams) for a
//               Carnatic composition in original script with IAST and meaning."
```

### 6.7 Transliteration Agent

```go
// agents/transliteration_agent.go
// Model:  gemini-3.5-flash  (gemini-3.1-pro for complex Sanskrit→Telugu)
// Tools:  TransliterateText, GetLyrics (as AgentTool calling LyricsAgent)
// Description: "Transliterates Carnatic lyrics from Sanskrit, Tamil, or Kannada
//               into Telugu script. Presents output as a side-by-side table."
```

### 6.8 YouTube Analyser Agent

```go
// agents/youtube_analyser_agent.go
// Model:  gemini-3.1-pro  (multimodal audio required)
// Tools:  FetchYouTubeMetadata, ExtractAudio, AnalyseAudioWithGemini,
//         LookupRaga, SearchKritis
// Description: "Accepts a YouTube URL, extracts audio, and uses Gemini
//               multimodal to identify the Ragam, Talam, composition, and artist."
```

---

## 7. Tool Function Signatures

All tools are plain Go functions registered with ADK via its tool-registration API. Defined in `tools/`.

### 7.1 Raga Tools (`tools/raga_lookup.go`)

```go
// LookupRaga returns a Raga by exact ID or by alias (case-insensitive).
// Returns an error if not found, triggering Google Search fallback in the agent.
//
// ADK tool name: "lookup_raga"
// Parameters:    name (string) — raga name or alias
func LookupRaga(ctx context.Context, name string) (*knowledge.Raga, error)

// SearchRagasBySwara returns all ragas whose arohana or avarohana contains
// all the given swaras as a subsequence.
//
// ADK tool name: "search_ragas_by_swara"
// Parameters:    swaras ([]string) — e.g. ["S","R2","G3","P"]
func SearchRagasBySwara(ctx context.Context, swaras []string) ([]knowledge.Raga, error)

// SearchRagasByMood returns ragas matching the given rasa and/or time of day.
//
// ADK tool name: "search_ragas_by_mood"
// Parameters:    rasa (string), timeOfDay (string) — either may be empty
func SearchRagasByMood(ctx context.Context, rasa, timeOfDay string) ([]knowledge.Raga, error)
```

### 7.2 Tala Tools (`tools/tala_lookup.go`)

```go
// LookupTala returns a Tala by ID or name (case-insensitive).
//
// ADK tool name: "lookup_tala"
// Parameters:    name (string)
func LookupTala(ctx context.Context, name string) (*knowledge.Tala, error)

// SearchTalasByBeats returns all talas with the given total beat count.
//
// ADK tool name: "search_talas_by_beats"
// Parameters:    beats (int)
func SearchTalasByBeats(ctx context.Context, beats int) ([]knowledge.Tala, error)
```

### 7.3 Kriti Tools (`tools/kriti_search.go`)

```go
// KritiFilter defines optional search filters; all fields are optional.
type KritiFilter struct {
    Ragam    string `json:"ragam,omitempty"`
    Talam    string `json:"talam,omitempty"`
    Composer string `json:"composer,omitempty"`  // composer ID
    Language string `json:"language,omitempty"`
    Tag      string `json:"tag,omitempty"`
}

// SearchKritis returns all kritis matching the given filter.
// At least one filter field must be non-empty.
//
// ADK tool name: "search_kritis"
// Parameters:    filter (KritiFilter)
func SearchKritis(ctx context.Context, filter KritiFilter) ([]knowledge.Kriti, error)

// LookupKriti returns a single Kriti by exact ID (case-insensitive).
//
// ADK tool name: "lookup_kriti"
// Parameters:    id (string)
func LookupKriti(ctx context.Context, id string) (*knowledge.Kriti, error)
```

### 7.4 Composer Tools (`tools/composer_lookup.go`)

```go
// LookupComposer returns a Composer by ID or name (case-insensitive).
//
// ADK tool name: "lookup_composer"
// Parameters:    name (string)
func LookupComposer(ctx context.Context, name string) (*knowledge.Composer, error)

// SearchComposersByLanguage returns all composers who composed in the given language.
//
// ADK tool name: "search_composers_by_language"
// Parameters:    language (string) — e.g. "Telugu", "Sanskrit"
func SearchComposersByLanguage(ctx context.Context, language string) ([]knowledge.Composer, error)
```

### 7.5 Lyrics Tools (`tools/lyrics_lookup.go`)

```go
// GetLyrics retrieves lyrics for a kriti, first from the embedded store,
// then from external sources if not found locally.
//
// ADK tool name: "get_lyrics"
// Parameters:    kritiID (string)
func GetLyrics(ctx context.Context, kritiID string) (*knowledge.Lyrics, error)

// ScrapeLyrics attempts to retrieve lyrics from karnatik.com or shivkumar.org.
// Results are cached in session state for the duration of the session.
// Only called by GetLyrics when embedded lyrics are unavailable.
//
// ADK tool name: "scrape_lyrics"   (internal; not directly exposed to agents)
// Parameters:    kritiName (string), composerName (string)
func ScrapeLyrics(ctx context.Context, kritiName, composerName string) (*knowledge.Lyrics, error)
```

### 7.6 Transliteration Tools (`tools/transliterate.go`)

```go
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
type TransliterationResult struct {
    Original       string `json:"original"`         // source text unchanged
    Transliterated string `json:"transliterated"`   // target script output
    SourceLang     string `json:"source_lang"`
    TargetLang     string `json:"target_lang"`
    Notes          []string `json:"notes"`          // approximation footnotes
}

func TransliterateText(ctx context.Context, text, sourceLang, targetLang string) (*TransliterationResult, error)
```

### 7.7 YouTube Tools (`tools/youtube_fetch.go`, `tools/gemini_audio.go`)

```go
// YouTubeMetadata holds video information from the YouTube Data API v3.
type YouTubeMetadata struct {
    VideoID     string `json:"video_id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    ChannelName string `json:"channel_name"`
    DurationSec int    `json:"duration_seconds"`
}

// FetchYouTubeMetadata retrieves video metadata from the YouTube Data API v3.
//
// ADK tool name: "fetch_youtube_metadata"
// Parameters:    url (string) — full YouTube URL
// Env required:  YOUTUBE_API_KEY
func FetchYouTubeMetadata(ctx context.Context, url string) (*YouTubeMetadata, error)

// ExtractAudio downloads and extracts the first 90 seconds of audio from a
// YouTube URL using yt-dlp, returning raw audio bytes in WAV format.
//
// ADK tool name: "extract_audio"  (internal)
// Parameters:    url (string), maxSeconds (int)
func ExtractAudio(ctx context.Context, url string, maxSeconds int) ([]byte, error)

// AudioAnalysisResult holds Gemini's multimodal analysis of an audio clip.
type AudioAnalysisResult struct {
    Ragam      string   `json:"ragam"`        // identified raga name
    Talam      string   `json:"talam"`        // identified tala name
    Kriti      string   `json:"kriti"`        // composition name (may be empty)
    Artist     string   `json:"artist"`       // performing artist (may be empty)
    Confidence string   `json:"confidence"`   // "high" | "medium" | "low"
    Candidates []string `json:"candidates"`   // alternate ragas if confidence < high
}

// AnalyseAudioWithGemini sends audio bytes to Gemini 3.1 Pro multimodal and
// returns a structured analysis of the Carnatic music content.
//
// ADK tool name: "analyse_audio_with_gemini"
// Parameters:    audioWAV ([]byte), hints (*YouTubeMetadata) — optional context
func AnalyseAudioWithGemini(ctx context.Context, audioWAV []byte, hints *YouTubeMetadata) (*AudioAnalysisResult, error)
```

---

## 8. HTTP API Contract

Defined in `cmd/server/server.go`.

### POST /chat

**Request:**
```json
{
  "message":    "Tell me about Kalyani Ragam",
  "session_id": "usr-abc123"
}
```

**Response `200 OK`:**
```json
{
  "response":   "Kalyani (Mechakalyani, #65) is a Melakarta raga in the Rudra chakra...",
  "session_id": "usr-abc123",
  "agent_used": "raga_agent",
  "latency_ms": 1240
}
```

**Response `400 Bad Request`** (empty message):
```json
{
  "error": "message must not be empty"
}
```

**Response `500 Internal Server Error`:**
```json
{
  "error":      "agent invocation failed",
  "detail":     "gemini api rate limit exceeded",
  "session_id": "usr-abc123"
}
```

---

### GET /health

**Response `200 OK`:**
```json
{
  "status": "ok",
  "knowledge_base": {
    "ragas":     103,
    "talas":      38,
    "kritis":     72,
    "composers":  12
  },
  "version": "0.1.0"
}
```

---

## 9. CLI Interface

Defined in `cmd/cli/runner.go`.

```
$ go run main.go --mode=cli

🎵 Nāda Guru — Carnatic Music Learning Agent
Type your question, or 'quit' to exit.
────────────────────────────────────────────────
You > Tell me about Kalyani Ragam
Guru > Kalyani (Mechakalyani #65) belongs to the Rudra chakra...
       Arohana:   S R2 G3 M2 P D2 N3 Ṡ
       Avarohana: Ṡ N3 D2 P M2 G3 R2 S
       ...
You > Show me some kritis in this raga
Guru > [resolves "this raga" → Kalyani from session state]
       Here are famous compositions in Kalyani...
You > quit
Bye! Subham astu 🙏
```

**Flags:**
```
--mode     string  "cli" | "server" (default: "cli")
--port     int     HTTP server port (default: 8080, server mode only)
--log      string  "json" | "text" (default: "text")
```

---

## 10. Data Flow Sequences

### 10.1 Text Query — Single Domain

```
User ──► POST /chat {"message": "What is Adi Talam?"}
              │
              ▼
         Orchestrator Agent
              │  classifies intent → tala query
              ▼
         Tala Agent
              │  calls LookupTala(ctx, "Adi")
              ▼
         tools.LookupTala
              │  searches KnowledgeStore.TalaByID["adi"]
              │  returns *Tala{...}
              ▼
         Tala Agent formats response
              ▼
         Orchestrator returns response
              │  updates session.LastTalaID = "adi"
              ▼
         POST /chat response {"response": "..."}
```

### 10.2 Text Query — Multi-Domain (Lyrics + Transliteration)

```
User ──► "Show me lyrics of Vatapi Ganapatim in Telugu"
              │
              ▼
         Orchestrator Agent
              │  classifies: multi-domain (lyrics + transliterate)
              │
              ├─ Step 1: Kriti Agent → LookupKriti("vatapi_ganapatim")
              │          returns Kriti{Language:"Sanskrit", Script:"sa", ...}
              │          updates session.LastKritiID
              │
              ├─ Step 2: Lyrics Agent → GetLyrics("vatapi_ganapatim")
              │          found in embedded store → returns *Lyrics
              │          updates session.LastLyricsKritiID
              │
              └─ Step 3: Transliteration Agent
                         → TransliterateText(pallavi.Original, "sa", "te")
                         → TransliterateText(anupallavi.Original, "sa", "te")
                         → TransliterateText(charanam[0].Original, "sa", "te")
                         assembles side-by-side table
                         │
                         ▼
         Orchestrator aggregates and returns final response
```

### 10.3 YouTube URL Analysis

```
User ──► "What raga is this? https://youtube.com/watch?v=XYZ"
              │
              ▼
         Orchestrator Agent
              │  detects YouTube URL → routes to YouTube Analyser Agent
              │
              ├─ Step 1: FetchYouTubeMetadata(url)
              │          → YouTube Data API v3 → *YouTubeMetadata
              │
              ├─ Step 2: ExtractAudio(url, maxSeconds=90)
              │          → yt-dlp subprocess → []byte (WAV)
              │
              ├─ Step 3: AnalyseAudioWithGemini(audioWAV, metadata)
              │          → Gemini 3.1 Pro (multimodal audio)
              │          → *AudioAnalysisResult{Ragam:"Kalyani", Confidence:"high"}
              │
              ├─ Step 4: LookupRaga("Kalyani")
              │          → full Raga profile from KnowledgeStore
              │
              └─ Step 5: SearchKritis(KritiFilter{Ragam:"Kalyani"})
                         → list of compositions in Kalyani
                         │
                         ▼
         Orchestrator returns: raga profile + identified composition + related kritis
```

---

## 11. Error Handling Strategy

| Scenario | Handling |
|---|---|
| Raga / tala / kriti not in knowledge base | Sub-agent invokes `google_search`; response notes source |
| Gemini API rate limit / 429 | Exponential backoff: 1s → 2s → 4s, max 3 retries (REQ-013) |
| YouTube URL invalid / private | Return descriptive error; offer free-text description fallback |
| Audio extraction (yt-dlp) failure | Return error with reason; skip audio analysis |
| Lyrics scraping fails | Return error; suggest alternative kritis with embedded lyrics |
| Knowledge base validation failure | Log structured error per entry; continue with valid entries (fail-open) |
| Sub-agent timeout (>15s) | Orchestrator returns partial response with timeout notice |

---

## 12. Environment Variables

| Variable | Required | Description |
|---|---|---|
| `GEMINI_API_KEY` | Yes | Google AI Studio API key for all Gemini calls |
| `YOUTUBE_API_KEY` | Yes (for YT analysis) | YouTube Data API v3 key |
| `PORT` | No (default: 8080) | HTTP server port |
| `LOG_FORMAT` | No (default: text) | `"json"` for structured logging |
| `CI` | No | If set, skips YouTube analysis eval test cases |

---

## 13. Observability

Every agent invocation emits a structured log line (REQ-016 NFR-006):

```json
{
  "ts":          "2026-07-04T10:23:01Z",
  "agent":       "raga_agent",
  "query":       "Tell me about Kalyani Ragam",
  "tools_called": ["lookup_raga"],
  "latency_ms":  843,
  "session_id":  "usr-abc123",
  "error":       null
}
```

---

## 14. Eval Test Case Schema

Defined in `eval/test_cases.json`:

```json
[
  {
    "id":           "tc-001",
    "description":  "Raga lookup by display name",
    "input":        "Tell me about Kalyani Ragam",
    "expect_contains": ["Mechakalyani", "M2", "R2", "G3", "Rudra"],
    "agent_expected": "raga_agent",
    "skip_in_ci":   false
  },
  {
    "id":           "tc-010",
    "description":  "YouTube raga identification",
    "input":        "https://www.youtube.com/watch?v=EXAMPLE",
    "expect_contains": ["ragam", "talam"],
    "agent_expected": "youtube_analyser_agent",
    "skip_in_ci":   true
  }
]
```

**Test case fields:**

| Field | Type | Description |
|---|---|---|
| `id` | string | Unique test identifier |
| `description` | string | Human-readable intent |
| `input` | string | User message or YouTube URL |
| `expect_contains` | []string | All strings must appear in response (case-insensitive) |
| `agent_expected` | string | Expected agent name in response metadata |
| `skip_in_ci` | bool | If true, skip when `CI` env var is set |