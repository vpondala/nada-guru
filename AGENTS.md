# AGENTS.md — Nāda Guru Project Constitution

## Project Overview
Nāda Guru is a multi-agent Carnatic music learning AI system.
It is built on Google ADK 2.0 (Go), Gemini 3.1 Pro / 3.5 Flash, and an embedded knowledge base.
Users interact via CLI or HTTP REST to ask questions about Ragas, Talams, Composers,
Kritis, Lyrics, and get Telugu transliterations.

## Specification Documents (Source of Truth)
ALWAYS read these before generating any code:
- REQUIREMENTS.md — 16 functional requirements with GIVEN/WHEN/THEN acceptance criteria
- DESIGN.md — Go structs, agent configs, tool signatures, API contracts, data-flow sequences
- TASKS.md — 31 ordered implementation tasks across 8 phases

DO NOT deviate from these specs without explicit user approval.

## Tech Stack (Non-Negotiable)
- Language: Go 1.26+ (NOT Python, NOT TypeScript)
- Agent framework: google.golang.org/adk v2 (ADK 2.0 Go)
- LLM: Gemini 3.1 Pro / 3.5 Flash
- Knowledge base: //go:embed (JSON files compiled into binary — NO external database)
- Audio extraction: yt-dlp subprocess (shell out from Go)
- YouTube API: YouTube Data API v3 (REST, NOT YouTube client library)
- Logging: log/slog (structured JSON or text, NOT logrus, NOT zap)
- Testing: standard go test (NOT testify unless already imported)

## Module Name
github.com/vpondala/nada-guru

## Code Conventions
- All Go types are defined in knowledge/types.go — do not duplicate types elsewhere
- All //go:embed directives live in knowledge/embed.go only
- Agent constructors are named New<AgentName>Agent(store *knowledge.KnowledgeStore)
- Tool functions accept ctx context.Context as first argument
- Return (result, error) — never panic on user input
- All errors are descriptive strings, never "error occurred"
- Session state keys are constants defined in session/state.go (KeyLastRaga, etc.)
- Log one structured slog line per agent invocation

## Environment Variables
Required: GEMINI_API_KEY, YOUTUBE_API_KEY
Optional: PORT (default 8080), LOG_FORMAT (default "text"), CI (skip YouTube tests)

## Knowledge Base Facts
- ragas.json: 103 entries (72 Melakarta + 31 Janya) — DO NOT modify the schema
- talas.json: 38 entries — DO NOT modify the schema
- kritis.json: 72 entries — DO NOT modify the schema
- composers.json: 12 entries — DO NOT modify the schema
- lyrics/*.json: individual kriti lyrics files — schema defined in DESIGN.md §4.5

## What NOT to Do
- Do NOT use Python in any file
- Do NOT use an external database (Postgres, Redis, ChromaDB, etc.)
- Do NOT use a Go ORM
- Do NOT add HTTP middleware frameworks (Gin, Echo, Fiber) — use net/http only
- Do NOT import testify unless already in go.mod
- Do NOT generate code for tasks not yet listed as in-scope in TASKS.md
