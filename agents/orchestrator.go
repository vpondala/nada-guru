// Package agents implements specialist agents for the Nāda Guru system.
package agents

import (
	"context"
	"log/slog"
	"os"

	"github.com/vpondala/nada-guru/knowledge"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

const orchestratorModelName = "gemini-3.1-pro-preview"

const orchestratorInstruction = `You are Nāda Guru, an expert AI guide for Carnatic classical music.
You help learners understand Ragas, Talams, Kritis, Composers, Lyrics, and Transliterations.
Route each user query to the most appropriate specialist agent.
Use session state to resolve pronouns and references from prior turns.
If a query spans multiple domains (e.g. lyrics + transliteration), invoke specialist agents sequentially and combine their responses.
Always respond in the same language the user used.`

// New creates the fully wired root orchestrator agent with all sub-agents registered.
func New(store *knowledge.KnowledgeStore) (agent.Agent, error) {
	ctx := context.Background()
	model, err := gemini.NewModel(ctx, orchestratorModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	ragaAgent, err := NewRagaAgent(store)
	if err != nil {
		return nil, err
	}
	talaAgent, err := NewTalaAgent(store)
	if err != nil {
		return nil, err
	}
	kritiAgent, err := NewKritiAgent(store)
	if err != nil {
		return nil, err
	}
	composerAgent, err := NewComposerAgent(store)
	if err != nil {
		return nil, err
	}
	lyricsAgent, err := NewLyricsAgent(store)
	if err != nil {
		return nil, err
	}
	transliterationAgent, err := NewTransliterationAgent(store)
	if err != nil {
		return nil, err
	}
	youtubeAgent, err := NewYouTubeAnalyserAgent(store)
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "nada_guru",
		Model:       model,
		Description: "Root orchestrator for Carnatic music learning queries. Routes to specialist agents for raga, tala, kriti, composer, lyrics, transliteration, and YouTube analysis.",
		Instruction: orchestratorInstruction,
		SubAgents: []agent.Agent{
			ragaAgent,
			talaAgent,
			kritiAgent,
			composerAgent,
			lyricsAgent,
			transliterationAgent,
			youtubeAgent,
		},
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				slog.Info("agent invocation", "agent", "nada_guru")
				return nil, nil
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}
