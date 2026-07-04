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

const kritiModelName = "gemini-2.5-flash"

const kritiInstruction = `You are Nāda Guru's Kriti specialist.
Return kritis in tabular format (Name | Ragam | Talam | Composer | Language).
Offer to show lyrics or transliterate.`

// NewKritiAgent creates the Kriti specialist agent.
func NewKritiAgent(store *knowledge.KnowledgeStore) (agent.Agent, error) {
	tools.Init(store)

	ctx := context.Background()
	model, err := gemini.NewModel(ctx, kritiModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	searchAgent, err := NewSearchAgent("kriti_search_agent")
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "kriti_agent",
		Model:       model,
		Description: "Finds Carnatic compositions by raga, tala, composer, or language. Returns metadata and offers to fetch lyrics or transliterate.",
		Instruction: kritiInstruction,
		SubAgents: []agent.Agent{
			searchAgent,
		},
		Tools: []tool.Tool{
			tools.SearchKritisTool,
			tools.LookupKritiTool,
		},
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				slog.Info("agent invocation", "agent", "kriti_agent")
				return nil, nil
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}
