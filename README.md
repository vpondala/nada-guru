# Nāda Guru

Nāda Guru is a multi-agent Carnatic music learning AI system built on Google ADK 2.0 (Go) and Gemini. Users can ask questions about Ragas, Talams, Kritis, Composers, Lyrics, and get Telugu transliterations. It also supports YouTube audio analysis to identify ragas and compositions from performances.

## Prerequisites

- Go 1.26+
- yt-dlp (for YouTube audio extraction)
- `GEMINI_API_KEY` — Google AI Studio API key
- `YOUTUBE_API_KEY` — YouTube Data API v3 key (for YouTube analysis)

## Setup

```bash
go mod download
export GEMINI_API_KEY="your-key"
export YOUTUBE_API_KEY="your-key"
```

## Running

```bash
# Interactive CLI
go run main.go --mode=cli

# HTTP REST server
go run main.go --mode=server
```

The server listens on port 8080 by default (override with `PORT` env var).

## Endpoints

- `POST /chat` — Send a message to the root orchestrator agent
- `GET /health` — Health check with knowledge base counts

## Running Tests

```bash
go test ./...
```

## Running Evaluation

```bash
go test ./eval/...
```

Set `CI=true` to skip YouTube-dependent eval cases.

## Architecture

See [DESIGN.md](DESIGN.md) for the full technical design, agent definitions, tool signatures, and data-flow sequences.
