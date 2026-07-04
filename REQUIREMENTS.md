# Nāda Guru — Requirements Specification

## Overview

Nāda Guru is a multi-agent AI system for learning Carnatic classical music. It is built using Google ADK 2.0 (Go), Gemini 3.1 Pro / 3.5 Flash, and a curated embedded knowledge base. Users can ask questions in natural language or paste YouTube links to receive contextual knowledge about Ragas, Talams, Composers, Artists, Kritis/Keertanas, Lyrics, and Transliterations.

## Scope

This specification covers the Nāda Guru backend agent system and its REST/CLI interface. Out of scope: a graphical front-end UI, audio playback, user authentication, and persistent user databases (beyond in-session memory).

---

## Requirements

### REQ-001 — Raga Information Lookup

**User Story**
As a Carnatic music learner, I want to ask the agent about a specific Raga by name, so that I can learn its arohana, avarohana, vadi, samvadi, rasa, time of performance, and related compositions.

**Acceptance Criteria**

* GIVEN a user asks "Tell me about Kalyani Ragam", WHEN the Raga Agent processes the query, THEN the response must include the raga's arohana, avarohana, vadi, samvadi, associated rasa(s), recommended time of performance, Melakarta number and chakra (if applicable), and a list of at least 3 famous compositions in that raga.
* GIVEN a user provides a partial or alternate name (e.g. "Yaman" for Mechakalyani, "Todi" for Hanumatodi), WHEN the agent resolves the query, THEN it must correctly identify the raga using its `aliases` field and return complete information.
* GIVEN a user asks about a Janya (derived) raga, WHEN the agent responds, THEN it must also identify the parent Melakarta and describe how the Janya differs from it (e.g. missing swaras, vakra patterns).
* GIVEN the raga is not found in the embedded knowledge base, WHEN the agent cannot match the name, THEN it must use the Google Search tool to retrieve information and clearly indicate the source is external.

---

### REQ-002 — Raga Comparison

**User Story**
As a Carnatic music learner, I want to compare two ragas side by side, so that I can understand the similarities and differences in their swara structures and emotional character.

**Acceptance Criteria**

* GIVEN a user asks "What is the difference between Bhairavi and Kharaharapriya?", WHEN the Raga Agent processes the query, THEN the response must display arohana, avarohana, vadi, samvadi, and rasa for both ragas in a structured side-by-side format.
* GIVEN two ragas share common swaras, WHEN the comparison is rendered, THEN the response must highlight the shared and differing swaras explicitly.
* GIVEN two ragas belong to different Melakartas, WHEN the comparison is rendered, THEN the response must note both Melakarta parent numbers.

---

### REQ-003 — Tala Information Lookup

**User Story**
As a Carnatic music learner, I want to ask the agent about a specific Talam, so that I can understand its structure, anga breakdown, beat count, and clap pattern.

**Acceptance Criteria**

* GIVEN a user asks "How does Adi Talam work?", WHEN the Tala Agent processes the query, THEN the response must include the tala's family (Suladi Sapta or Chapu), jati, anga breakdown (Laghu/Drutam/Anudrutam), total beat count, and clap/wave pattern.
* GIVEN a user asks about a tala by its beat count (e.g. "Which tala has 7 beats?"), WHEN the agent processes the query, THEN it must return all matching talas (e.g. Misra Chapu, Tisra Jati Triputa) with their structures.
* GIVEN a user asks "Give me examples of kritis in Misra Chapu", WHEN the Tala Agent responds, THEN it must return at least 2 example compositions from the kritis knowledge base with that tala.

---

### REQ-004 — Kriti / Keertana Search

**User Story**
As a Carnatic music learner, I want to search for Carnatic compositions by raga, tala, composer, or language, so that I can discover songs to listen to and study.

**Acceptance Criteria**

* GIVEN a user asks "Which kritis did Tyagaraja compose in Bhairavi?", WHEN the Kriti Agent processes the query, THEN the response must return all matching kritis from the knowledge base filtered by both composer and raga, displaying the kriti name, talam, and language.
* GIVEN a user asks "Show me Telugu kritis in Adi talam", WHEN the Kriti Agent processes the query, THEN the response must return a filtered list of kritis with language = "Telugu" and talam = "Adi".
* GIVEN a user asks about a kriti by name (e.g. "Tell me about Endaro Mahanubhavulu"), WHEN the Kriti Agent processes the query, THEN the response must return the full metadata record: ragam, talam, composer, language, and a brief description.
* GIVEN a kriti is found in the knowledge base, WHEN the Kriti Agent responds, THEN it must offer to show the full lyrics and/or transliteration as follow-up actions.
* GIVEN no kritis match the search criteria in the local knowledge base, WHEN the Kriti Agent cannot find results, THEN it must fall back to Google Search and indicate the source.

---

### REQ-005 — Composer Information

**User Story**
As a Carnatic music learner, I want to learn about a Carnatic composer, so that I can understand their era, language, deity of devotion, compositional style, and landmark works.

**Acceptance Criteria**

* GIVEN a user asks "Tell me about Muthuswami Dikshitar", WHEN the Composer Agent processes the query, THEN the response must include the composer's full name, dates, era, language(s) of composition, region, primary deity, notable work collections, estimated total compositions, a descriptive summary, and a list of famous kritis.
* GIVEN a user asks "Who are the Carnatic Trinity?", WHEN the Composer Agent processes the query, THEN it must return details for Tyagaraja, Muthuswami Dikshitar, and Syama Sastri, along with a brief explanation of what makes them the Trinity.
* GIVEN a user asks "Which composers wrote in Kannada?", WHEN the Composer Agent processes the query, THEN it must return all composers from the knowledge base whose `language` field includes "Kannada".
* GIVEN a composer is not found in the knowledge base, WHEN the agent cannot find a match, THEN it must use Google Search and indicate the source.

---

### REQ-006 — Lyrics Retrieval

**User Story**
As a Carnatic music learner, I want to retrieve the full lyrics of a Carnatic composition, so that I can read, study, and practise singing the text.

**Acceptance Criteria**

* GIVEN a user asks "Show me the lyrics of Endaro Mahanubhavulu", WHEN the Lyrics Agent processes the query, THEN the response must present the lyrics structured into Pallavi, Anupallavi, and all Charanam sections in the original script.
* GIVEN a lyrics file exists for the requested kriti in the embedded `lyrics/` knowledge base, WHEN the Lyrics Agent retrieves it, THEN it must serve it from the embedded store without making any external network calls.
* GIVEN a lyrics file does not exist in the embedded knowledge base, WHEN the Lyrics Agent cannot find it locally, THEN it must attempt to retrieve lyrics from known Carnatic lyrics sources (karnatik.com, shivkumar.org) via web scraping and cache the result in session memory for the duration of the session.
* GIVEN lyrics have been retrieved (from any source), WHEN displaying them, THEN the agent must also show the IAST romanisation alongside the original script to aid pronunciation.
* GIVEN the user asks for word meanings, WHEN the Lyrics Agent receives this follow-up, THEN it must invoke Gemini to provide a word-by-word or phrase-by-phrase meaning in English for the displayed lyrics section.
* GIVEN lyrics cannot be found from any source, WHEN the Lyrics Agent exhausts all retrieval paths, THEN it must inform the user clearly and suggest alternative compositions with available lyrics.

---

### REQ-007 — Telugu Transliteration

**User Story**
As a Carnatic music learner who reads Telugu script, I want to transliterate lyrics from Sanskrit, Tamil, or Kannada into Telugu script, so that I can read and sing compositions regardless of their source language.

**Acceptance Criteria**

* GIVEN a user asks "Transliterate this kriti into Telugu" for a Sanskrit composition, WHEN the Transliteration Agent processes the request, THEN it must identify the source language from the kriti's metadata and invoke the transliteration pipeline accordingly.
* GIVEN the source language of the kriti is already Telugu, WHEN the Transliteration Agent processes the request, THEN it must return the existing Telugu text with a note confirming no transliteration is required.
* GIVEN a Sanskrit kriti is to be transliterated, WHEN the Transliteration Agent processes it, THEN it must produce Telugu script output that preserves Carnatic phonemes: anusvara (ం), visarga (ః), chandrabindu, and long vowels (ā, ī, ū) accurately.
* GIVEN a Tamil kriti is to be transliterated into Telugu, WHEN the Transliteration Agent processes it, THEN it must handle Tamil-specific phonemes (ழ, ற, ண etc.) by mapping them to the nearest Telugu equivalents and adding a footnote about approximations made.
* GIVEN transliteration is complete, WHEN the Transliteration Agent returns the output, THEN it must present the result as a side-by-side table: Original Script | Telugu Transliteration, with the section structure (Pallavi / Anupallavi / Charanam) preserved.
* GIVEN a user requests transliteration into a language other than Telugu, WHEN the Transliteration Agent receives the request, THEN it must respond that only Telugu transliteration is currently supported but indicate that additional languages (Kannada, Tamil, Devanagari, IAST Roman) are on the roadmap.

---

### REQ-008 — YouTube Audio Analysis

**User Story**
As a Carnatic music learner, I want to paste a YouTube link into the agent, so that it can identify the Raga, Talam, composition name, and artist from the audio.

**Acceptance Criteria**

* GIVEN a user pastes a YouTube URL, WHEN the YouTube Analyser Agent detects the URL in the input, THEN it must automatically invoke the YouTube analysis pipeline without requiring the user to specify they want analysis.
* GIVEN a YouTube URL is provided, WHEN the YouTube Analyser Agent processes it, THEN it must first fetch the video's title, description, and channel name from the YouTube Data API v3 to use as contextual hints.
* GIVEN the YouTube metadata has been fetched, WHEN the agent proceeds to audio analysis, THEN it must extract an audio stream (MP3/WAV format) using yt-dlp, limited to the first 90 seconds of the recording.
* GIVEN audio has been extracted, WHEN the YouTube Analyser Agent invokes Gemini 3.1 Pro multimodal analysis, THEN the prompt must ask Gemini to identify: (a) the Ragam, (b) the Talam, (c) the composition name if recognisable, and (d) the performing artist if identifiable.
* GIVEN Gemini returns an identified Raga, WHEN the YouTube Analyser Agent has a raga name, THEN it must look up the raga in the knowledge base and return the full raga profile along with other compositions in the same raga.
* GIVEN Gemini cannot confidently identify the raga, WHEN Gemini returns low-confidence or multiple candidates, THEN the agent must present all candidate ragas with their likelihood and ask the user for confirmation.
* GIVEN the YouTube URL is invalid, private, age-restricted, or unavailable, WHEN the audio extraction fails, THEN the agent must report the specific error reason and offer to accept free-text description of the music instead.

---

### REQ-009 — Free-Text Song Description Analysis

**User Story**
As a Carnatic music learner, I want to describe a song in free text (e.g. hum the swara pattern or describe what I remember), so that the agent can suggest possible ragas or compositions that match my description.

**Acceptance Criteria**

* GIVEN a user describes a melody or swara sequence in free text (e.g. "It goes Sa Ri Ga Ma Pa, very serene, evening raga"), WHEN the Root Orchestrator Agent processes the input, THEN it must route the request to the Raga Agent for swara-pattern matching.
* GIVEN a swara sequence is described, WHEN the Raga Agent searches the knowledge base, THEN it must return ragas whose arohana or avarohana contains the described swara sequence as a subsequence.
* GIVEN a user describes the mood or occasion (e.g. "a morning raga with a devotional feel"), WHEN the Raga Agent searches the knowledge base, THEN it must filter by the `rasa` and `time_of_day` fields and return matching ragas ranked by relevance.

---

### REQ-010 — Conversational Session Memory

**User Story**
As a Carnatic music learner, I want the agent to remember the context of our ongoing conversation, so that I can ask follow-up questions without repeating context (e.g. "give me more songs in this raga" after discussing Kalyani).

**Acceptance Criteria**

* GIVEN a user has just asked about a specific raga, WHEN the user sends a follow-up such as "give me more compositions in this raga", THEN the Root Orchestrator Agent must resolve "this raga" from the session state without asking the user to re-specify.
* GIVEN a user has just discussed a composer, WHEN the user asks "what are their most famous kritis?", THEN the agent must resolve "their" to the composer discussed in the immediately preceding exchange.
* GIVEN a conversation session begins, WHEN the agent initialises, THEN it must create an ADK session state object that persists: last discussed raga, last discussed tala, last discussed composer, last retrieved lyrics (kriti ID), and user's preferred transliteration language (default: Telugu).
* GIVEN a session ends or a new session is started, WHEN the agent initialises a new session, THEN all prior session state must be cleared and not leaked between sessions.

---

### REQ-011 — Multi-Agent Orchestration

**User Story**
As a developer, I want the system to route user queries to the correct specialised sub-agent automatically, so that each agent can focus on its domain and the system handles complex multi-domain queries gracefully.

**Acceptance Criteria**

* GIVEN any user input, WHEN the Root Orchestrator Agent processes it, THEN it must classify the intent and route to exactly one primary sub-agent: Raga Agent, Tala Agent, Kriti Agent, Composer Agent, YouTube Analyser Agent, Lyrics Agent, or Transliteration Agent.
* GIVEN a query spans multiple domains (e.g. "Show me lyrics for a Tyagaraja kriti in Kalyani and transliterate into Telugu"), WHEN the Root Orchestrator Agent processes it, THEN it must invoke sub-agents in the correct sequential order (Kriti Agent → Lyrics Agent → Transliteration Agent) and aggregate the results into a single coherent response.
* GIVEN a sub-agent fails or returns an error, WHEN the orchestrator detects the failure, THEN it must log the error, attempt a fallback strategy (e.g. Google Search), and return a partial response with a clear explanation rather than a blank or error response.
* GIVEN the system is operational, WHEN any query is received, THEN the full response must be returned within 15 seconds under normal network conditions.

---

### REQ-012 — Knowledge Base Integrity

**User Story**
As a developer, I want the embedded knowledge base to be validated at startup, so that the system does not silently serve corrupt or incomplete data.

**Acceptance Criteria**

* GIVEN the application starts, WHEN the `knowledge` package initialises, THEN it must unmarshal all four embedded JSON files (ragas.json, talas.json, kritis.json, composers.json) and verify that: (a) ragas.json contains exactly 72 Melakarta entries and at least 20 Janya entries, (b) talas.json contains at least 8 entries, (c) kritis.json contains at least 50 entries, (d) composers.json contains at least 10 entries.
* GIVEN a Melakarta raga entry is loaded, WHEN validated, THEN it must have non-empty `arohana`, `avarohana`, `melakarta_number` (1–72), and `chakra` fields.
* GIVEN a kriti entry is loaded, WHEN validated, THEN it must have non-empty `id`, `name`, `ragam`, `talam`, `composer`, and `language` fields.
* GIVEN validation fails for any entry, WHEN the application starts, THEN it must log a structured error identifying the malformed entry by ID and continue with valid entries only (fail-open, not fail-closed).

---

### REQ-013 — Gemini API Integration

**User Story**
As a developer, I want the system to interact with Gemini 3.1 Pro for complex reasoning and Gemini 3.5 Flash for fast agentic domain Q&A, so that response quality and latency are optimised per use case.

**Acceptance Criteria**

* GIVEN the Root Orchestrator Agent and YouTube Analyser Agent are invoked, WHEN they call Gemini, THEN they must use the `gemini-3.1-pro` model.
* GIVEN the Raga, Tala, Kriti, Composer, Lyrics, and Transliteration sub-agents are invoked, WHEN they call Gemini, THEN they must use the `gemini-3.5-flash` model.
* GIVEN any Gemini API call is made, WHEN the call fails due to rate limiting or network error, THEN the agent must retry up to 3 times with exponential backoff before returning an error.
* GIVEN a Gemini API key is not set in the environment, WHEN the application starts, THEN it must exit with a descriptive error message indicating which environment variable is missing (`GEMINI_API_KEY`).

---

### REQ-014 — Google Search Grounding

**User Story**
As a user, I want the agent to use live web search when its knowledge base does not contain the information I need, so that I always receive a useful answer rather than "I don't know."

**Acceptance Criteria**

* GIVEN a query cannot be answered from the embedded knowledge base, WHEN the relevant sub-agent determines the knowledge base is insufficient, THEN it must invoke the ADK built-in Google Search tool to retrieve external information.
* GIVEN Google Search results are used, WHEN the agent composes its response, THEN it must clearly indicate the information came from web search (e.g. "Based on web search: ...") and include the source URL(s).
* GIVEN Google Search returns no relevant results, WHEN the agent cannot find information from any source, THEN it must respond with a polite, specific message explaining what was searched and suggesting alternative approaches.

---

### REQ-015 — CLI and HTTP Server Interface

**User Story**
As a developer and demo presenter, I want to interact with Nāda Guru via both a CLI and an HTTP REST endpoint, so that I can demonstrate it in a terminal and integrate it into a Kaggle notebook.

**Acceptance Criteria**

* GIVEN the application is launched with `go run main.go --mode=cli`, WHEN the CLI is active, THEN it must present an interactive prompt that accepts free-text input, sends it to the Root Orchestrator Agent, and prints the formatted response.
* GIVEN the application is launched with `go run main.go --mode=server`, WHEN the HTTP server is active, THEN it must listen on `PORT` (default 8080) and expose a `POST /chat` endpoint accepting `{"message": "...", "session_id": "..."}` and returning `{"response": "...", "session_id": "..."}`.
* GIVEN a `session_id` is provided in the HTTP request, WHEN the server processes the request, THEN it must load and update the corresponding ADK session state to maintain conversational continuity across HTTP calls.
* GIVEN the HTTP server is running, WHEN a `GET /health` request is received, THEN it must return `{"status": "ok", "knowledge_base": {"ragas": N, "kritis": N, "talas": N, "composers": N}}` confirming the knowledge base is loaded.

---

### REQ-016 — Agent Evaluation

**User Story**
As a developer, I want a test suite of golden Q&A pairs that can be run against the agent system, so that I can validate response quality before submission.

**Acceptance Criteria**

* GIVEN the evaluation suite is run with `go test ./eval/...`, WHEN the test runner executes, THEN it must load `eval/test_cases.json` and run each test case against the live agent system.
* GIVEN a test case specifies an expected raga name in the answer, WHEN the agent response is evaluated, THEN the test must pass if the raga name appears in the response (case-insensitive substring match).
* GIVEN a test case specifies an expected composer name, WHEN the agent response is evaluated, THEN the test must pass if the composer name or any of their known aliases appear in the response.
* GIVEN the full test suite is executed, WHEN all cases have been evaluated, THEN the runner must print a summary report: total cases, passed, failed, and pass-rate percentage.
* GIVEN the test suite includes YouTube analysis test cases, WHEN those cases run, THEN they must be marked `skip_in_ci: true` and skipped when the `CI` environment variable is set, to avoid external network dependency in automated runs.

---

## Non-Functional Requirements

### NFR-001 — Response Latency
The system must return a complete response to any text-based query within **15 seconds** under normal network conditions (excluding YouTube audio download time).

### NFR-002 — Knowledge Base Size
The embedded knowledge base must remain under **10 MB** total to stay within Kaggle notebook memory limits and compile-time `//go:embed` constraints.

### NFR-003 — Language Support
The system must correctly render and process Unicode text in Telugu (te), Sanskrit/Devanagari (sa), Tamil (ta), and Kannada (kn) scripts without garbling characters.

### NFR-004 — Extensibility
Adding a new transliteration target language must require changes only to the Transliteration Agent's instruction string — no new agents, tools, or knowledge base files.

### NFR-005 — Stateless Knowledge Base
The knowledge base must be fully embedded at compile time via `//go:embed` and require no external database, file system access at runtime, or network calls to serve knowledge base queries.

### NFR-006 — Observability
Every agent invocation must emit a structured log line containing: timestamp, agent name, user query (truncated to 100 chars), response latency in milliseconds, and tool calls made.