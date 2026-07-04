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

	transliterateTool, err := functiontool.New(functiontool.Config{
		Name:        "transliterate_text",
		Description: "Transliterate text from source script to Telugu",
	}, func(ctx agent.ToolContext, input struct {
		Text       string `json:"text"`
		SourceLang string `json:"source_lang"`
		TargetLang string `json:"target_lang"`
	}) (struct {
		Original       string   `json:"original"`
		Transliterated string   `json:"transliterated"`
		SourceLang     string   `json:"source_lang"`
		TargetLang     string   `json:"target_lang"`
		Notes          []string `json:"notes"`
		Error          string   `json:"error,omitempty"`
	}, error) {
		result, err := tools.TransliterateText(ctx, input.Text, input.SourceLang, input.TargetLang)
		if err != nil {
			return struct {
				Original       string   `json:"original"`
				Transliterated string   `json:"transliterated"`
				SourceLang     string   `json:"source_lang"`
				TargetLang     string   `json:"target_lang"`
				Notes          []string `json:"notes"`
				Error          string   `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		return struct {
			Original       string   `json:"original"`
			Transliterated string   `json:"transliterated"`
			SourceLang     string   `json:"source_lang"`
			TargetLang     string   `json:"target_lang"`
			Notes          []string `json:"notes"`
			Error          string   `json:"error,omitempty"`
		}{
			Original:       result.Original,
			Transliterated: result.Transliterated,
			SourceLang:     result.SourceLang,
			TargetLang:     result.TargetLang,
			Notes:          result.Notes,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "transliteration_agent",
		Model:       model,
		Description: "Transliterates Carnatic lyrics from Sanskrit, Tamil, or Kannada into Telugu script. Presents output as a side-by-side table.",
		Instruction: transliterationInstruction,
		Tools: []tool.Tool{
			transliterateTool,
			geminitool.GoogleSearch{},
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
