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
	"google.golang.org/adk/tool/functiontool"
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

	lookupComposerTool, err := functiontool.New(functiontool.Config{
		Name:        "lookup_composer",
		Description: "Lookup a composer by name",
	}, func(ctx agent.ToolContext, input struct{ Name string `json:"name"` }) (struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Error    string `json:"error,omitempty"`
	}, error) {
		c, err := tools.LookupComposer(ctx, input.Name)
		if err != nil {
			return struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Error    string `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		return struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Error    string `json:"error,omitempty"`
		}{ID: c.ID, Name: c.Name}, nil
	})
	if err != nil {
		return nil, err
	}

	searchLangTool, err := functiontool.New(functiontool.Config{
		Name:        "search_composers_by_language",
		Description: "Search composers by language",
	}, func(ctx agent.ToolContext, input struct{ Language string `json:"language"` }) (struct {
		Results []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"results"`
		Error string `json:"error,omitempty"`
	}, error) {
		results, err := tools.SearchComposersByLanguage(ctx, input.Language)
		if err != nil {
			return struct {
				Results []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"results"`
				Error string `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		out := make([]struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}, len(results))
		for i, c := range results {
			out[i] = struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{ID: c.ID, Name: c.Name}
		}
		return struct {
			Results []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"results"`
			Error string `json:"error,omitempty"`
		}{Results: out}, nil
	})
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
			lookupComposerTool,
			searchLangTool,
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
