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

	lookupRagaTool, err := functiontool.New(functiontool.Config{
		Name:        "lookup_raga",
		Description: "Lookup a raga by name or alias",
	}, func(ctx agent.ToolContext, input struct{ Name string `json:"name"` }) (struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Error string `json:"error,omitempty"`
	}, error) {
		r, err := tools.LookupRaga(ctx, input.Name)
		if err != nil {
			return struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Error string `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		return struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Error string `json:"error,omitempty"`
		}{ID: r.ID, Name: r.Name}, nil
	})
	if err != nil {
		return nil, err
	}

	searchSwarasTool, err := functiontool.New(functiontool.Config{
		Name:        "search_ragas_by_swara",
		Description: "Search ragas by swara pattern",
	}, func(ctx agent.ToolContext, input struct {
		Swaras []string `json:"swaras"`
	}) (struct {
		Results []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"results"`
		Error string `json:"error,omitempty"`
	}, error) {
		results, err := tools.SearchRagasBySwara(ctx, input.Swaras)
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
		for i, r := range results {
			out[i] = struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{ID: r.ID, Name: r.Name}
		}
		return struct {
			Results []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"results"`
			Error string `json:"error,omitempty"`
		}{Results: out}, nil
	})
	if err != nil {
		return nil, err
	}

	searchMoodTool, err := functiontool.New(functiontool.Config{
		Name:        "search_ragas_by_mood",
		Description: "Search ragas by mood and time of day",
	}, func(ctx agent.ToolContext, input struct {
		Rasa      string `json:"rasa"`
		TimeOfDay string `json:"time_of_day"`
	}) (struct {
		Results []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"results"`
		Error string `json:"error,omitempty"`
	}, error) {
		results, err := tools.SearchRagasByMood(ctx, input.Rasa, input.TimeOfDay)
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
		for i, r := range results {
			out[i] = struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}{ID: r.ID, Name: r.Name}
		}
		return struct {
			Results []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"results"`
			Error string `json:"error,omitempty"`
		}{Results: out}, nil
	})
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
			lookupRagaTool,
			searchSwarasTool,
			searchMoodTool,
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
