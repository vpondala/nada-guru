# Nāda Guru — Implementation Tasks

Derived from [REQUIREMENTS.md](REQUIREMENTS.md) and [DESIGN.md](DESIGN.md).
Tasks are ordered by dependency. Complete each task fully before moving to the next.
Each task maps to one or more requirements via `REQ-XXX` references.

---

## Phase 1 — Project Scaffold

### Task 1.1 — Initialise Go module and directory structure
**Requirements:** REQ-015  
**Files to create:**
- `go.mod` — `module github.com/vpondala/nada-guru`, `go 1.26`
- `go.sum`
- `main.go` — stub with `--mode` flag parsing
- `cmd/cli/runner.go` — empty `Run()` function
- `cmd/server/server.go` — empty `Start()` function
- `agents/`, `tools/`, `knowledge/`, `session/`, `eval/` — empty directories with `.gitkeep`

**Acceptance criteria:**
- [x] `go build ./...` succeeds with no errors
- [x] `go run main.go --mode=cli` prints a startup banner and exits cleanly
- [x] `go run main.go --mode=server` starts an HTTP listener on port 8080 and exits cleanly on SIGINT

---

### Task 1.2 — Add ADK and Gemini Go dependencies
**Requirements:** REQ-011, REQ-013  
**Files to modify:** `go.mod`, `go.sum`

**Commands:**
```bash
go get google.golang.org/adk@latest
go get google.golang.org/genai@latest
```

**Acceptance criteria:**
- [x] `import "google.golang.org/adk/agent/llmagent"` compiles without error
- [x] `import "google.golang.org/genai"` compiles without error
- [x] `go mod tidy` leaves no unused dependencies

---

## Phase 2 — Knowledge Base

### Task 2.1 — Define all domain types
**Requirements:** REQ-012  
**Files to create:** `knowledge/types.go`

Implement all structs exactly as specified in DESIGN.md §4:
- `Raga`
- `Tala`
- `Kriti`
- `Composer`
- `Lyrics`
- `LyricsSection`
- `KritiFilter` (used by search tools)

**Acceptance criteria:**
- [ ] All structs compile with correct JSON tags
- [ ] All pointer fields (`*int`, `*string`) are used where the value may be absent
- [ ] `go vet ./knowledge/...` passes with no warnings

---

### Task 2.2 — Embed knowledge base JSON files
**Requirements:** REQ-012, NFR-002, NFR-005  
**Files to create:** `knowledge/embed.go`

Wire all five `//go:embed` directives:
- `ragas.json`
- `talas.json`
- `kritis.json`
- `composers.json`
- `lyrics/*.json` (directory embed)

**Acceptance criteria:**
- [ ] `go build ./...` embeds all files into the binary
- [ ] `len(ragasJSON) > 0` is true at runtime
- [ ] Binary does not read from the filesystem at runtime for any knowledge base query

---

### Task 2.3 — Implement KnowledgeStore with indexes
**Requirements:** REQ-012  
**Files to create:** `knowledge/store.go`

Implement `KnowledgeStore` as defined in DESIGN.md §4.6:
- `New() (*KnowledgeStore, error)` — unmarshal all embedded JSON, build all indexes
- `Validate() error` — check counts and required fields per REQ-012

Build the following indexes:
- `RagaByID map[string]*Raga`
- `RagaByAlias map[string]*Raga` — normalise aliases to lowercase for case-insensitive lookup
- `TalaByID map[string]*Tala`
- `KritiByID map[string]*Kriti`
- `ComposerByID map[string]*Composer`
- `KritisByRaga map[string][]string`
- `KritisByTala map[string][]string`
- `KritisByComposer map[string][]string`
- `KritisByLanguage map[string][]string`

**Acceptance criteria:**
- [ ] `New()` returns a non-nil store with no error when all JSON files are valid
- [ ] `Validate()` returns an error if `len(Ragas) < 72`
- [ ] `Validate()` returns an error if any Melakarta raga has an empty `Arohana` field
- [ ] `Validate()` returns an error if any Kriti has an empty `Ragam` field
- [ ] `store.RagaByAlias["kalyani"]` returns the Mechakalyani entry
- [ ] `store.RagaByAlias["todi"]` returns the Hanumatodi entry
- [ ] `store.KritisByComposer["tyagaraja"]` returns at least 10 entries

---

### Task 2.4 — Write knowledge base unit tests
**Requirements:** REQ-012  
**Files to create:** `knowledge/store_test.go`

**Test cases:**
- `TestValidate_PassesOnGoodData` — `New()` + `Validate()` succeed on embedded files
- `TestMelakarataCount` — exactly 72 Melakarta ragas present
- `TestJanyaCount` — at least 20 Janya ragas present
- `TestAliasLookup_Kalyani` — alias "Kalyani" resolves to Mechakalyani (#65)
- `TestAliasLookup_Todi` — alias "Todi" resolves to Hanumatodi (#8)
- `TestKritisByRaga_Sri` — at least 1 result for ragam "Sri"
- `TestKritisByComposer_Tyagaraja` — at least 5 results

**Acceptance criteria:**
- [ ] `go test ./knowledge/...` passes with all tests green
- [ ] Test coverage for `store.go` is ≥ 80%

---

## Phase 3 — Session State

### Task 3.1 — Implement session state
**Requirements:** REQ-010  
**Files to create:** `session/state.go`

Implement `State` struct and string constants as defined in DESIGN.md §5.

Provide helpers:
```go
func FromADKSession(s adk.Session) State
func (st State) ToADKSession(s adk.Session)
```

**Acceptance criteria:**
- [ ] `State` struct marshals/unmarshals to JSON with no loss of fields
- [ ] `PreferredScriptTarget` defaults to `"te"` when not set
- [ ] `FromADKSession` returns zero-value `State` (not a panic) when session is empty

---

## Phase 4 — Tool Layer

### Task 4.1 — Implement raga lookup tools
**Requirements:** REQ-001, REQ-002, REQ-009  
**Files to create:** `tools/raga_lookup.go`

Implement:
- `LookupRaga(ctx, name string) (*knowledge.Raga, error)` — exact + alias lookup
- `SearchRagasBySwara(ctx, swaras []string) ([]knowledge.Raga, error)` — subsequence match
- `SearchRagasByMood(ctx, rasa, timeOfDay string) ([]knowledge.Raga, error)` — filter by rasa/time

Each function must:
- Accept a `*knowledge.KnowledgeStore` via dependency injection or closure
- Return a descriptive error (not nil) when no results are found

**Acceptance criteria:**
- [ ] `LookupRaga(ctx, "Kalyani")` returns Mechakalyani without error
- [ ] `LookupRaga(ctx, "kalyani")` (lowercase) returns Mechakalyani without error
- [ ] `LookupRaga(ctx, "NonExistentRaga")` returns a non-nil error
- [ ] `SearchRagasBySwara(ctx, []string{"R2","G3","M2"})` returns at least Mechakalyani (#65)
- [ ] `SearchRagasByMood(ctx, "bhakti", "morning")` returns at least 3 ragas

---

### Task 4.2 — Implement tala lookup tools
**Requirements:** REQ-003  
**Files to create:** `tools/tala_lookup.go`

Implement:
- `LookupTala(ctx, name string) (*knowledge.Tala, error)` — by ID or display name
- `SearchTalasByBeats(ctx, beats int) ([]knowledge.Tala, error)` — by total beat count

**Acceptance criteria:**
- [ ] `LookupTala(ctx, "Adi")` returns the 8-beat Adi Tala entry
- [ ] `LookupTala(ctx, "adi")` (lowercase) returns the same entry
- [ ] `SearchTalasByBeats(ctx, 7)` returns Misra Chapu and Tisra Jati Triputa
- [ ] `LookupTala(ctx, "unknown")` returns a non-nil error

---

### Task 4.3 — Implement kriti search tools
**Requirements:** REQ-004  
**Files to create:** `tools/kriti_search.go`

Implement:
- `SearchKritis(ctx, filter KritiFilter) ([]knowledge.Kriti, error)` — multi-field filter
- `LookupKriti(ctx, id string) (*knowledge.Kriti, error)` — exact ID lookup

**Acceptance criteria:**
- [ ] `SearchKritis(ctx, KritiFilter{Composer:"tyagaraja"})` returns ≥ 5 results
- [ ] `SearchKritis(ctx, KritiFilter{Ragam:"Kalyani", Language:"Telugu"})` returns Telugu Kalyani kritis only
- [ ] `SearchKritis(ctx, KritiFilter{})` returns an error (at least one filter required)
- [ ] `LookupKriti(ctx, "endaro_mahanubhavulu")` returns the correct entry
- [ ] `LookupKriti(ctx, "nonexistent")` returns a non-nil error

---

### Task 4.4 — Implement composer lookup tools
**Requirements:** REQ-005  
**Files to create:** `tools/composer_lookup.go`

Implement:
- `LookupComposer(ctx, name string) (*knowledge.Composer, error)` — by ID or name
- `SearchComposersByLanguage(ctx, language string) ([]knowledge.Composer, error)`

**Acceptance criteria:**
- [ ] `LookupComposer(ctx, "Tyagaraja")` returns the Tyagaraja entry
- [ ] `LookupComposer(ctx, "tyagaraja")` (lowercase) returns the same entry
- [ ] `SearchComposersByLanguage(ctx, "Kannada")` returns Purandaradasa and Kanakadasa
- [ ] `SearchComposersByLanguage(ctx, "Sanskrit")` returns Muthuswami Dikshitar

---

### Task 4.5 — Implement lyrics lookup tool
**Requirements:** REQ-006  
**Files to create:** `tools/lyrics_lookup.go`

Implement:
- `GetLyrics(ctx, kritiID string) (*knowledge.Lyrics, error)`
  1. Check embedded `lyrics/*.json` first (load via `//go:embed`)
  2. If not found, call `ScrapeLyrics` and cache result in session state
- `ScrapeLyrics(ctx, kritiName, composerName string) (*knowledge.Lyrics, error)` — HTTP scrape from known sources

**Acceptance criteria:**
- [ ] `GetLyrics(ctx, "endaro_mahanubhavulu")` returns a non-nil `*Lyrics` from embedded store
- [ ] Returned `Lyrics` has non-empty `Pallavi.Original` and `Pallavi.IAST`
- [ ] `GetLyrics(ctx, "nonexistent_kriti")` calls `ScrapeLyrics` as fallback
- [ ] `GetLyrics(ctx, "nonexistent_kriti")` returns a non-nil error if scraping also fails
- [ ] `ScrapeLyrics` does not panic on HTTP timeout; returns a descriptive error

---

### Task 4.6 — Implement transliteration tool
**Requirements:** REQ-007  
**Files to create:** `tools/transliterate.go`

Implement `TransliterateText(ctx, text, sourceLang, targetLang string) (*TransliterationResult, error)` using Gemini 3.1 Pro.

Prompt template (embed in code):
```
Transliterate the following {{sourceLang}} text into {{targetLang}} script.
Preserve all Carnatic music phonemes accurately, including:
- Anusvara (anunaasika), Visarga, Chandrabindu
- Long vowels (ā, ī, ū)
- Retroflex consonants
For any phoneme with no exact equivalent, use the nearest approximation
and add a note in the Notes array.
Return ONLY valid JSON matching this schema: {original, transliterated, source_lang, target_lang, notes}
Input text:
{{text}}
```

**Acceptance criteria:**
- [ ] `TransliterateText(ctx, "వాతాపి గణపతిం", "te", "te")` returns the same text unchanged with a note
- [ ] `TransliterateText(ctx, "वातापि गणपतिं", "sa", "te")` returns a Telugu-script transliteration
- [ ] `TransliterateText(ctx, "...", "ta", "te")` includes footnotes for Tamil-specific phonemes (ழ, ற)
- [ ] `TransliterateText(ctx, "text", "en", "te")` returns an error — unsupported source language
- [ ] Returns error on Gemini API failure after 3 retries with exponential backoff

---

### Task 4.7 — Implement YouTube fetch and audio tools
**Requirements:** REQ-008  
**Files to create:** `tools/youtube_fetch.go`, `tools/gemini_audio.go`

Implement:
- `FetchYouTubeMetadata(ctx, url string) (*YouTubeMetadata, error)` — calls YouTube Data API v3
- `ExtractAudio(ctx, url string, maxSeconds int) ([]byte, error)` — shells out to `yt-dlp`
- `AnalyseAudioWithGemini(ctx, audioWAV []byte, hints *YouTubeMetadata) (*AudioAnalysisResult, error)` — Gemini 3.1 Pro multimodal

Gemini audio analysis prompt:
```
You are an expert in Carnatic classical music.
Listen to this audio clip and identify:
1. Ragam (raga name) — be specific, e.g. "Mechakalyani" not just "Kalyani"
2. Talam (rhythmic cycle) — e.g. "Adi", "Misra Chapu"
3. Composition name if recognisable
4. Performing artist if identifiable
Video title hint: {{title}}
Channel hint: {{channel}}
Return ONLY valid JSON: {ragam, talam, kriti, artist, confidence, candidates}
```

**Acceptance criteria:**
- [ ] `FetchYouTubeMetadata(ctx, "https://youtube.com/watch?v=dQw4w9WgXcQ")` returns non-nil metadata (title must be non-empty)
- [ ] `FetchYouTubeMetadata(ctx, "not-a-url")` returns a descriptive error
- [ ] `ExtractAudio` returns `[]byte` of length > 0 for a valid public YouTube URL
- [ ] `ExtractAudio` returns an error for a private or invalid URL
- [ ] `AnalyseAudioWithGemini` returns `Confidence` of `"high"`, `"medium"`, or `"low"` — never empty
- [ ] All three functions are skipped in unit tests when `CI=true`

---

### Task 4.8 — Write tool unit tests
**Requirements:** REQ-001 through REQ-009  
**Files to create:** `tools/*_test.go` for each tool file

**Acceptance criteria:**
- [ ] `go test ./tools/...` passes for all non-network tests
- [ ] YouTube and Gemini audio tests are gated with `t.Skip` when `os.Getenv("CI") != ""`
- [ ] All lookup tools are tested with both a known-good and a not-found case

---

## Phase 5 — Agents

### Task 5.1 — Implement Raga Agent
**Requirements:** REQ-001, REQ-002  
**Files to create:** `agents/raga_agent.go`

```go
func NewRagaAgent(store *knowledge.KnowledgeStore) (*llmagent.Agent, error)
```

- Model: `gemini-3.5-flash`
- Tools: `LookupRaga`, `SearchRagasBySwara`, `SearchRagasByMood`, `google_search`
- System instruction: Explain ragas in structured format — arohana, avarohana, vadi, samvadi, rasa, time of day, Melakarta number, janya ragas, famous compositions. For comparisons, render a side-by-side table.

**Acceptance criteria:**
- [ ] Agent compiles and registers all tools without panic
- [ ] Manual prompt "Tell me about Kalyani Ragam" returns arohana `S R2 G3 M2 P D2 N3 Ṡ`
- [ ] Manual prompt "Compare Bhairavi and Kharaharapriya" returns a response mentioning both ragas and their swara differences
- [ ] Manual prompt "Morning ragas with devotional mood" returns at least 2 raga names

---

### Task 5.2 — Implement Tala Agent
**Requirements:** REQ-003  
**Files to create:** `agents/tala_agent.go`

```go
func NewTalaAgent(store *knowledge.KnowledgeStore) (*llmagent.Agent, error)
```

- Model: `gemini-3.5-flash`
- Tools: `LookupTala`, `SearchTalasByBeats`, `google_search`
- System instruction: Explain talas with anga breakdown, total beats, clap pattern, and examples. Use structured layout.

**Acceptance criteria:**
- [ ] Agent compiles and registers all tools without panic
- [ ] Manual prompt "How does Adi Talam work?" returns the structure `Laghu(4) + Drutam + Drutam` and beat count 8
- [ ] Manual prompt "Which tala has 7 beats?" returns Misra Chapu

---

### Task 5.3 — Implement Kriti Agent
**Requirements:** REQ-004  
**Files to create:** `agents/kriti_agent.go`

```go
func NewKritiAgent(store *knowledge.KnowledgeStore) (*llmagent.Agent, error)
```

- Model: `gemini-3.5-flash`
- Tools: `SearchKritis`, `LookupKriti`, `google_search`
- System instruction: Return kritis in tabular format (Name | Ragam | Talam | Composer | Language). Offer to show lyrics or transliterate.

**Acceptance criteria:**
- [ ] Agent compiles and registers all tools without panic
- [ ] Manual prompt "Which Tyagaraja kritis are in Bhairavi?" returns ≥ 1 result
- [ ] Manual prompt "Tell me about Endaro Mahanubhavulu" returns ragam Sri, talam Adi, and composer Tyagaraja
- [ ] Response offers follow-up options: lyrics or transliteration

---

### Task 5.4 — Implement Composer Agent
**Requirements:** REQ-005  
**Files to create:** `agents/composer_agent.go`

```go
func NewComposerAgent(store *knowledge.KnowledgeStore) (*llmagent.Agent, error)
```

- Model: `gemini-3.5-flash`
- Tools: `LookupComposer`, `SearchComposersByLanguage`, `google_search`
- System instruction: Return full composer profile — name, dates, region, languages, deity, notable works, famous kritis. For "Who are the Trinity?" return all three as a group.

**Acceptance criteria:**
- [ ] Agent compiles and registers all tools without panic
- [ ] Manual prompt "Tell me about Muthuswami Dikshitar" returns birth year 1775, language Sanskrit, and at least 2 famous kritis
- [ ] Manual prompt "Who are the Carnatic Trinity?" returns Tyagaraja, Dikshitar, and Syama Sastri
- [ ] Manual prompt "Which composers wrote in Kannada?" returns Purandaradasa and Kanakadasa

---

### Task 5.5 — Implement Lyrics Agent
**Requirements:** REQ-006  
**Files to create:** `agents/lyrics_agent.go`

```go
func NewLyricsAgent(store *knowledge.KnowledgeStore) (*llmagent.Agent, error)
```

- Model: `gemini-3.5-flash`
- Tools: `GetLyrics`, `ScrapeLyrics`, `google_search`
- System instruction: Display lyrics section by section (Pallavi → Anupallavi → Charanams). Always show original script alongside IAST. Offer word-by-word meaning on request. Offer transliteration after displaying lyrics.

**Acceptance criteria:**
- [ ] Agent compiles and registers all tools without panic
- [ ] Manual prompt "Show me the lyrics of Endaro Mahanubhavulu" returns Pallavi in Telugu script and IAST
- [ ] Response includes section labels: Pallavi, Anupallavi, Charanam 1
- [ ] Response offers follow-up: "Would you like word-by-word meaning or Telugu transliteration?"

---

### Task 5.6 — Implement Transliteration Agent
**Requirements:** REQ-007  
**Files to create:** `agents/transliteration_agent.go`

```go
func NewTransliterationAgent(store *knowledge.KnowledgeStore) (*llmagent.Agent, error)
```

- Model: `gemini-3.5-flash` (falls back to `gemini-3.1-pro` for Sanskrit input)
- Tools: `TransliterateText`, `GetLyrics` (as AgentTool wrapping LyricsAgent)
- System instruction: For each lyrics section, produce a two-column table: Original | Telugu Transliteration. Preserve section structure. Add phoneme approximation footnotes for Tamil source.

**Acceptance criteria:**
- [ ] Agent compiles and registers all tools without panic
- [ ] Manual prompt "Transliterate Vatapi Ganapatim into Telugu" returns Telugu script output in a table format
- [ ] Output preserves Pallavi / Anupallavi / Charanam structure
- [ ] Manual prompt requesting a non-Telugu target returns: "Only Telugu transliteration is currently supported"

---

### Task 5.7 — Implement YouTube Analyser Agent
**Requirements:** REQ-008  
**Files to create:** `agents/youtube_analyser_agent.go`

```go
func NewYouTubeAnalyserAgent(store *knowledge.KnowledgeStore) (*llmagent.Agent, error)
```

- Model: `gemini-3.1-pro`
- Tools: `FetchYouTubeMetadata`, `ExtractAudio`, `AnalyseAudioWithGemini`, `LookupRaga`, `SearchKritis`
- System instruction: When a YouTube URL is detected, run the full analysis pipeline automatically. Return: identified ragam (with full profile), identified talam, composition name (if known), artist (if known), and a list of other compositions in the same raga.

**Acceptance criteria:**
- [ ] Agent compiles and registers all tools without panic
- [ ] Agent detects a YouTube URL in the input without requiring an explicit instruction
- [ ] When audio analysis confidence is "low", response presents all candidate ragas and asks user to confirm
- [ ] When URL is invalid, response states the error and offers to accept a free-text description
- [ ] Test is skipped in CI (`t.Skip` when `os.Getenv("CI") != ""`)

---

### Task 5.8 — Implement Root Orchestrator Agent
**Requirements:** REQ-010, REQ-011  
**Files to create:** `agents/orchestrator.go`

```go
func New(store *knowledge.KnowledgeStore) (*llmagent.Agent, error)
```

- Model: `gemini-3.1-pro`
- SubAgents: all 7 agents from Tasks 5.1–5.7
- Tools: `google_search` (ADK built-in)
- Session state: read/write `session.State` on every turn

System instruction must include:
- Intent classification rules (one primary sub-agent per query)
- Multi-domain chaining order (Kriti → Lyrics → Transliteration)
- Session reference resolution ("this raga", "their kritis", etc.)
- Error fallback: if sub-agent fails, return partial response with explanation

**Acceptance criteria:**
- [ ] Agent compiles with all 7 sub-agents registered, no panic
- [ ] Single-domain query routes to exactly one sub-agent (check agent_used in response)
- [ ] Multi-domain query "Show me lyrics for a Tyagaraja kriti in Kalyani and transliterate into Telugu" invokes Kriti → Lyrics → Transliteration in sequence
- [ ] Follow-up "give me more songs in this raga" correctly resolves to the last discussed raga from session state
- [ ] Full response is returned in ≤ 15 seconds for a text query on standard network (NFR-001)

---

## Phase 6 — Interfaces

### Task 6.1 — Implement CLI runner
**Requirements:** REQ-015  
**Files to create:** `cmd/cli/runner.go`

Implement an interactive read-eval-print loop:
- Print banner on startup
- Read user input from stdin line by line
- Send each line to the root orchestrator agent
- Print response with a `Guru >` prefix
- Exit cleanly on `"quit"` or `"exit"` input or SIGINT

**Acceptance criteria:**
- [ ] `go run main.go --mode=cli` starts and prints the banner
- [ ] Typing a question returns a response from the agent
- [ ] Typing `quit` exits with status code 0
- [ ] SIGINT (Ctrl-C) exits with status code 0 and prints a farewell message
- [ ] Empty input (pressing Enter) is ignored gracefully

---

### Task 6.2 — Implement HTTP server
**Requirements:** REQ-015  
**Files to create:** `cmd/server/server.go`

Implement two endpoints as per DESIGN.md §8:

**`POST /chat`**
- Parse `{"message": "...", "session_id": "..."}` from request body
- Return `400` if message is empty
- Invoke root orchestrator agent with session ID
- Return `{"response": "...", "session_id": "...", "agent_used": "...", "latency_ms": N}`

**`GET /health`**
- Return `{"status": "ok", "knowledge_base": {...}, "version": "0.1.0"}`

**Acceptance criteria:**
- [ ] `POST /chat` with a valid message returns 200 and a non-empty `response` field
- [ ] `POST /chat` with empty `message` returns 400
- [ ] `POST /chat` with the same `session_id` in two sequential requests maintains session context
- [ ] `GET /health` returns 200 with correct knowledge base counts
- [ ] Server logs one structured JSON line per request (NFR-006)

---

### Task 6.3 — Implement structured logging
**Requirements:** NFR-006  
**Files to modify:** `cmd/server/server.go`, `agents/orchestrator.go`

Use the standard `log/slog` package. Each agent invocation must emit:
```json
{
  "ts":           "2026-07-04T10:23:01Z",
  "agent":        "raga_agent",
  "query":        "Tell me about Kalyani Ragam",
  "tools_called": ["lookup_raga"],
  "latency_ms":   843,
  "session_id":   "usr-abc123",
  "error":        null
}
```

**Acceptance criteria:**
- [ ] `LOG_FORMAT=json go run main.go --mode=server` emits newline-delimited JSON logs
- [ ] `LOG_FORMAT=text go run main.go --mode=server` emits human-readable logs
- [ ] `latency_ms` field is always present
- [ ] `error` is `null` on success and a string on failure

---

## Phase 7 — Evaluation

### Task 7.1 — Create evaluation test cases
**Requirements:** REQ-016  
**Files to create:** `eval/test_cases.json`

Write at least 20 test cases covering:
- Raga lookup by name and alias (5 cases)
- Tala lookup and beat count query (3 cases)
- Kriti search by composer, raga, and language (4 cases)
- Composer info including Trinity query (3 cases)
- Lyrics retrieval (2 cases)
- Transliteration (2 cases)
- YouTube analysis (1 case, `skip_in_ci: true`)

Each case follows the schema defined in DESIGN.md §14.

**Acceptance criteria:**
- [ ] `eval/test_cases.json` is valid JSON
- [ ] At least 20 entries present
- [ ] All YouTube-related cases have `"skip_in_ci": true`
- [ ] All non-YouTube cases have `"skip_in_ci": false`

---

### Task 7.2 — Implement eval test runner
**Requirements:** REQ-016  
**Files to create:** `eval/eval_test.go`

Implement `TestEvalSuite` that:
1. Loads `eval/test_cases.json`
2. For each test case, skips if `skip_in_ci: true` and `CI` env var is set
3. Sends `input` to the root orchestrator agent
4. Asserts that all strings in `expect_contains` appear in the response (case-insensitive)
5. Prints a summary: total / passed / failed / pass-rate

**Acceptance criteria:**
- [ ] `go test ./eval/...` runs without panic
- [ ] Non-CI test cases all pass (≥ 90% pass rate)
- [ ] `CI=true go test ./eval/...` skips YouTube cases cleanly
- [ ] Summary line is printed at the end: `Pass rate: N/M (X%)`

---

## Phase 8 — Polish and Submission Prep

### Task 8.1 — Write README.md
**Requirements:** REQ-015  
**Files to create:** `README.md`

Cover:
- What Nāda Guru is (2-3 sentences)
- Prerequisites (`go 1.26+`, `yt-dlp`, environment variables)
- Setup: `go mod download`, set `GEMINI_API_KEY`, set `YOUTUBE_API_KEY`
- Running CLI: `go run main.go --mode=cli`
- Running server: `go run main.go --mode=server`
- Running tests: `go test ./...`
- Running eval: `go test ./eval/...`
- Architecture overview (link to DESIGN.md)

**Acceptance criteria:**
- [ ] A new developer can clone, set env vars, and run `go run main.go --mode=cli` following only the README
- [ ] All commands in README are copy-paste executable

---

### Task 8.2 — Kaggle notebook writeup
**Requirements:** Kaggle Capstone Submission  
**Files to create:** `KAGGLE_WRITEUP.md`

Structure:
1. Introduction and personal motivation
2. What makes Nāda Guru special (multimodal YouTube analysis, multi-agent depth, Go on ADK 2.0)
3. Architecture diagram (ASCII from DESIGN.md §2)
4. Sub-agent and tool overview table
5. Knowledge base stats (103 ragas, 38 talas, 72 kritis, 12 composers)
6. 5 annotated demo interactions (screenshots or code output)
7. Eval results (pass rate from Task 7.2)
8. Learnings and future work (voice input, quiz mode, more transliteration targets)
9. Conclusion

**Acceptance criteria:**
- [ ] Writeup is between 800–2000 words
- [ ] Includes at least one architecture diagram
- [ ] Includes at least 5 real agent interaction examples
- [ ] Includes eval pass-rate result
- [ ] Links to GitHub repo with all source code

---

### Task 8.3 — Record demo video
**Requirements:** Kaggle Capstone Submission  
**Target duration:** 3–5 minutes

Script (from plan):
- `[0:00–0:30]` Hook and problem statement
- `[0:30–1:00]` Architecture overview (show diagram)
- `[1:00–2:30]` Live demo: YouTube URL → raga identified; Kalyani query; Tyagaraja kritis in Adi Talam; lyrics + Telugu transliteration
- `[2:30–3:30]` Code walkthrough: orchestrator, one sub-agent, one tool
- `[3:30–4:00]` Eval results (pass rate)
- `[4:00–4:30]` Closing: rationale, future roadmap

**Acceptance criteria:**
- [ ] Video is between 3:00 and 5:30 minutes
- [ ] Video is uploaded to YouTube or another publicly accessible URL
- [ ] Audio is clear throughout
- [ ] All 4 demo scenarios are shown live (not slides)

---

## Task Completion Checklist

| Task | Phase | Status |
|------|-------|--------|
| 1.1 — Go module scaffold | Project | ☐ |
| 1.2 — ADK + Gemini dependencies | Project | ☐ |
| 2.1 — Domain types | Knowledge Base | ☐ |
| 2.2 — Embed JSON files | Knowledge Base | ☐ |
| 2.3 — KnowledgeStore + indexes | Knowledge Base | ☐ |
| 2.4 — Knowledge base tests | Knowledge Base | ☐ |
| 3.1 — Session state | Session | ☐ |
| 4.1 — Raga lookup tools | Tools | ☐ |
| 4.2 — Tala lookup tools | Tools | ☐ |
| 4.3 — Kriti search tools | Tools | ☐ |
| 4.4 — Composer lookup tools | Tools | ☐ |
| 4.5 — Lyrics lookup tool | Tools | ☐ |
| 4.6 — Transliteration tool | Tools | ☐ |
| 4.7 — YouTube tools | Tools | ☐ |
| 4.8 — Tool unit tests | Tools | ☐ |
| 5.1 — Raga Agent | Agents | ☐ |
| 5.2 — Tala Agent | Agents | ☐ |
| 5.3 — Kriti Agent | Agents | ☐ |
| 5.4 — Composer Agent | Agents | ☐ |
| 5.5 — Lyrics Agent | Agents | ☐ |
| 5.6 — Transliteration Agent | Agents | ☐ |
| 5.7 — YouTube Analyser Agent | Agents | ☐ |
| 5.8 — Root Orchestrator Agent | Agents | ☐ |
| 6.1 — CLI runner | Interfaces | ☐ |
| 6.2 — HTTP server | Interfaces | ☐ |
| 6.3 — Structured logging | Interfaces | ☐ |
| 7.1 — Eval test cases | Evaluation | ☐ |
| 7.2 — Eval test runner | Evaluation | ☐ |
| 8.1 — README.md | Polish | ☐ |
| 8.2 — Kaggle writeup | Polish | ☐ |
| 8.3 — Demo video | Polish | ☐ |