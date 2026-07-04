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

const transliterationModelName = "gemini-2.5-flash"

const transliterationInstruction = `You are Nāda Guru's Transliteration specialist.
For each lyrics section, produce a two-column table: Original | Telugu Transliteration.
Preserve section structure.
Add phoneme approximation footnotes for Tamil source.
Only Telugu transliteration is currently supported.`

// NewTransliterationAgent creates the Transliteration specialist agent.
func NewTransliterationAgent(store *knowledge.KnowledgeStore) (agent.Agent, error) {
	tools.Init(store)

	ctx := context.Background()
	model, err := gemini.NewModel(ctx, transliterationModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	searchAgent, err := NewSearchAgent("transliteration_search_agent")
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "transliteration_agent",
		Model:       model,
		Description: "Transliterates Carnatic lyrics from Sanskrit, Tamil, or Kannada into Telugu script. Presents output as a side-by-side table.",
		Instruction: transliterationInstruction,
		SubAgents: []agent.Agent{
			searchAgent,
		},
		Tools: []tool.Tool{
			tools.TransliterateTextTool,
		},
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				slog.Info("agent invocation", "agent", "transliteration_agent")
				return nil, nil
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}
