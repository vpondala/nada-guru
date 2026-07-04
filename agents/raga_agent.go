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

const ragaModelName = "gemini-2.5-flash"

const ragaInstruction = `You are Nāda Guru's Raga specialist.
Explain ragas in structured format: arohana, avarohana, vadi, samvadi, rasa, time of day, Melakarta number, janya ragas, and famous compositions.
For comparisons, render a side-by-side table.
Always preserve Carnatic music phonemes accurately.`

// NewRagaAgent creates the Raga specialist agent.
func NewRagaAgent(store *knowledge.KnowledgeStore) (agent.Agent, error) {
	tools.Init(store)

	ctx := context.Background()
	model, err := gemini.NewModel(ctx, ragaModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	searchAgent, err := NewSearchAgent("raga_search_agent")
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "raga_agent",
		Model:       model,
		Description: "Answers questions about Carnatic ragas — arohana, avarohana, vadi, samvadi, rasa, time of day, Melakarta classification, and related compositions.",
		Instruction: ragaInstruction,
		SubAgents: []agent.Agent{
			searchAgent,
		},
		Tools: []tool.Tool{
			tools.LookupRagaTool,
			tools.SearchRagasBySwaraTool,
			tools.SearchRagasByMoodTool,
		},
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				slog.Info("agent invocation", "agent", "raga_agent")
				return nil, nil
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}
