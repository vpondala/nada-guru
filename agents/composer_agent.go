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

const composerModelName = "gemini-2.5-flash"

const composerInstruction = `You are Nāda Guru's Composer specialist.
Provide full composer profile — name, dates, region, languages, deity, notable works, and famous kritis.
For "Who are the Trinity?" return all three as a group.`

// NewComposerAgent creates the Composer specialist agent.
func NewComposerAgent(store *knowledge.KnowledgeStore) (agent.Agent, error) {
	tools.Init(store)

	ctx := context.Background()
	model, err := gemini.NewModel(ctx, composerModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	searchAgent, err := NewSearchAgent("composer_search_agent")
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "composer_agent",
		Model:       model,
		Description: "Provides biographical and compositional details about Carnatic music composers.",
		Instruction: composerInstruction,
		SubAgents: []agent.Agent{
			searchAgent,
		},
		Tools: []tool.Tool{
			tools.LookupComposerTool,
			tools.SearchComposersByLangTool,
		},
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				slog.Info("agent invocation", "agent", "composer_agent")
				return nil, nil
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}
