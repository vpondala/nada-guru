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
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

const lyricsModelName = "gemini-3.5-flash"

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

	getLyricsTool, err := functiontool.New(functiontool.Config{
		Name:        "get_lyrics",
		Description: "Get lyrics for a kriti by ID",
	}, func(ctx agent.ToolContext, input struct{ KritiID string `json:"kriti_id"` }) (struct {
		KritiID string `json:"kriti_id"`
		Pallavi struct {
			Original string `json:"original"`
			IAST     string `json:"iast"`
		} `json:"pallavi"`
		Anupallavi struct {
			Original string `json:"original"`
			IAST     string `json:"iast"`
		} `json:"anupallavi"`
		Charanams []struct {
			Original string `json:"original"`
			IAST     string `json:"iast"`
		} `json:"charanams"`
		Error string `json:"error,omitempty"`
	}, error) {
		lyrics, err := tools.GetLyrics(ctx, input.KritiID)
		if err != nil {
			return struct {
				KritiID string `json:"kriti_id"`
				Pallavi struct {
					Original string `json:"original"`
					IAST     string `json:"iast"`
				} `json:"pallavi"`
				Anupallavi struct {
					Original string `json:"original"`
					IAST     string `json:"iast"`
				} `json:"anupallavi"`
				Charanams []struct {
					Original string `json:"original"`
					IAST     string `json:"iast"`
				} `json:"charanams"`
				Error string `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		charanams := make([]struct {
			Original string `json:"original"`
			IAST     string `json:"iast"`
		}, len(lyrics.Charanams))
		for i, c := range lyrics.Charanams {
			charanams[i] = struct {
				Original string `json:"original"`
				IAST     string `json:"iast"`
			}{Original: c.Original, IAST: c.IAST}
		}
		return struct {
			KritiID string `json:"kriti_id"`
			Pallavi struct {
				Original string `json:"original"`
				IAST     string `json:"iast"`
			} `json:"pallavi"`
			Anupallavi struct {
				Original string `json:"original"`
				IAST     string `json:"iast"`
			} `json:"anupallavi"`
			Charanams []struct {
				Original string `json:"original"`
				IAST     string `json:"iast"`
			} `json:"charanams"`
			Error string `json:"error,omitempty"`
		}{
			KritiID: lyrics.KritiID,
			Pallavi: struct {
				Original string `json:"original"`
				IAST     string `json:"iast"`
			}{Original: lyrics.Pallavi.Original, IAST: lyrics.Pallavi.IAST},
			Anupallavi: struct {
				Original string `json:"original"`
				IAST     string `json:"iast"`
			}{Original: lyrics.Anupallavi.Original, IAST: lyrics.Anupallavi.IAST},
			Charanams: charanams,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	scrapeLyricsTool, err := functiontool.New(functiontool.Config{
		Name:        "scrape_lyrics",
		Description: "Scrape lyrics from external sources (internal)",
	}, func(ctx agent.ToolContext, input struct {
		KritiName    string `json:"kriti_name"`
		ComposerName string `json:"composer_name"`
	}) (struct {
		Error string `json:"error,omitempty"`
	}, error) {
		_, err := tools.ScrapeLyrics(ctx, input.KritiName, input.ComposerName)
		return struct {
			Error string `json:"error,omitempty"`
		}{Error: err.Error()}, nil
	})
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "lyrics_agent",
		Model:       model,
		Description: "Retrieves full lyrics (Pallavi, Anupallavi, Charanams) for a Carnatic composition in original script with IAST and meaning.",
		Instruction: lyricsInstruction,
		Tools: []tool.Tool{
			getLyricsTool,
			scrapeLyricsTool,
			geminitool.GoogleSearch{},
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
