// Package agents implements specialist agents for the Nāda Guru system.
package agents

import (
	"context"
	"log/slog"
	"os"

	"github.com/vpondala/nada-guru/knowledge"
	"github.com/vpondala/nada-guru/tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

const lyricsModelName = "gemini-2.5-flash"

const lyricsInstruction = `You are Nāda Guru's Lyrics specialist.
Display lyrics section by section (Pallavi → Anupallavi → Charanams).
Always show original script alongside IAST.
Offer word-by-word meaning on request.
Offer transliteration after displaying lyrics.`

// NewLyricsAgent creates the Lyrics specialist agent.
func NewLyricsAgent(store *knowledge.KnowledgeStore) (agent.Agent, error) {
	tools.Init(store)

	ctx := context.Background()
	model, err := gemini.NewModel(ctx, lyricsModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	searchAgent, err := NewSearchAgent("lyrics_search_agent")
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "lyrics_agent",
		Model:       model,
		Description: "Retrieves full lyrics (Pallavi, Anupallavi, Charanams) for a Carnatic composition in original script with IAST and meaning.",
		Instruction: lyricsInstruction,
		SubAgents: []agent.Agent{
			searchAgent,
		},
		Tools: []tool.Tool{
			tools.GetLyricsTool,
			tools.ScrapeLyricsTool,
		},
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				slog.Info("agent invocation", "agent", "lyrics_agent")
				return nil, nil
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}
