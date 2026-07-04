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

	searchKritisTool, err := functiontool.New(functiontool.Config{
		Name:        "search_kritis",
		Description: "Search kritis by raga, tala, composer, or language",
	}, func(ctx agent.ToolContext, input struct {
		Ragam    string `json:"ragam,omitempty"`
		Talam    string `json:"talam,omitempty"`
		Composer string `json:"composer,omitempty"`
		Language string `json:"language,omitempty"`
		Tag      string `json:"tag,omitempty"`
	}) (struct {
		Results []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Ragam    string `json:"ragam"`
			Talam    string `json:"talam"`
			Composer string `json:"composer"`
			Language string `json:"language"`
		} `json:"results"`
		Error string `json:"error,omitempty"`
	}, error) {
		results, err := tools.SearchKritis(ctx, knowledge.KritiFilter{
			Ragam:    input.Ragam,
			Talam:    input.Talam,
			Composer: input.Composer,
			Language: input.Language,
			Tag:      input.Tag,
		})
		if err != nil {
			return struct {
				Results []struct {
					ID       string `json:"id"`
					Name     string `json:"name"`
					Ragam    string `json:"ragam"`
					Talam    string `json:"talam"`
					Composer string `json:"composer"`
					Language string `json:"language"`
				} `json:"results"`
				Error string `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		out := make([]struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Ragam    string `json:"ragam"`
			Talam    string `json:"talam"`
			Composer string `json:"composer"`
			Language string `json:"language"`
		}, len(results))
		for i, k := range results {
			out[i] = struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Ragam    string `json:"ragam"`
				Talam    string `json:"talam"`
				Composer string `json:"composer"`
				Language string `json:"language"`
			}{ID: k.ID, Name: k.Name, Ragam: k.Ragam, Talam: k.Talam, Composer: k.Composer, Language: k.Language}
		}
		return struct {
			Results []struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Ragam    string `json:"ragam"`
				Talam    string `json:"talam"`
				Composer string `json:"composer"`
				Language string `json:"language"`
			} `json:"results"`
			Error string `json:"error,omitempty"`
		}{Results: out}, nil
	})
	if err != nil {
		return nil, err
	}

	lookupKritiTool, err := functiontool.New(functiontool.Config{
		Name:        "lookup_kriti",
		Description: "Lookup a kriti by ID",
	}, func(ctx agent.ToolContext, input struct{ ID string `json:"id"` }) (struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Ragam    string `json:"ragam"`
		Talam    string `json:"talam"`
		Composer string `json:"composer"`
		Error    string `json:"error,omitempty"`
	}, error) {
		k, err := tools.LookupKriti(ctx, input.ID)
		if err != nil {
			return struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Ragam    string `json:"ragam"`
				Talam    string `json:"talam"`
				Composer string `json:"composer"`
				Error    string `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		return struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Ragam    string `json:"ragam"`
			Talam    string `json:"talam"`
			Composer string `json:"composer"`
			Error    string `json:"error,omitempty"`
		}{ID: k.ID, Name: k.Name, Ragam: k.Ragam, Talam: k.Talam, Composer: k.Composer}, nil
	})
	searchAgent, err := NewSearchAgent()
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
			searchKritisTool,
			lookupKritiTool,
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
