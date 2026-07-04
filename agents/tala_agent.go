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

const talaModelName = "gemini-2.5-flash"

const talaInstruction = `You are Nāda Guru's Tala specialist.
Explain talas with anga breakdown, total beats, clap pattern, and examples.
Use structured layout.`

// NewTalaAgent creates the Tala specialist agent.
func NewTalaAgent(store *knowledge.KnowledgeStore) (agent.Agent, error) {
	tools.Init(store)

	ctx := context.Background()
	model, err := gemini.NewModel(ctx, talaModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	lookupTalaTool, err := functiontool.New(functiontool.Config{
		Name:        "lookup_tala",
		Description: "Lookup a tala by name",
	}, func(ctx agent.ToolContext, input struct{ Name string `json:"name"` }) (struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Error    string `json:"error,omitempty"`
	}, error) {
		t, err := tools.LookupTala(ctx, input.Name)
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
		}{ID: t.ID, Name: t.Name}, nil
	})
	if err != nil {
		return nil, err
	}

	searchBeatsTool, err := functiontool.New(functiontool.Config{
		Name:        "search_talas_by_beats",
		Description: "Search talas by beat count",
	}, func(ctx agent.ToolContext, input struct{ Beats int `json:"beats"` }) (struct {
		Results []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"results"`
		Error string `json:"error,omitempty"`
	}, error) {
		results, err := tools.SearchTalasByBeats(ctx, input.Beats)
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
		for i, t := range results {
			out[i] = struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{ID: t.ID, Name: t.Name}
		}
		return struct {
			Results []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"results"`
			Error string `json:"error,omitempty"`
		}{Results: out}, nil
	})
	searchAgent, err := NewSearchAgent("tala_search_agent")
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "tala_agent",
		Model:       model,
		Description: "Answers questions about Carnatic talas — structure, angas, beat counts, clap patterns, and example compositions.",
		Instruction: talaInstruction,
		SubAgents: []agent.Agent{
			searchAgent,
		},
		Tools: []tool.Tool{
			lookupTalaTool,
			searchBeatsTool,
		},
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				slog.Info("agent invocation", "agent", "tala_agent")
				return nil, nil
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}
