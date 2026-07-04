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

const youtubeModelName = "gemini-2.5-pro"

const youtubeInstruction = `You are Nāda Guru's YouTube Analyser.
When a YouTube URL is detected, run the full analysis pipeline automatically.
Return: identified ragam (with full profile), identified talam, composition name (if known), artist (if known), and a list of other compositions in the same raga.`

// NewYouTubeAnalyserAgent creates the YouTube Analyser agent.
func NewYouTubeAnalyserAgent(store *knowledge.KnowledgeStore) (agent.Agent, error) {
	tools.Init(store)

	ctx := context.Background()
	model, err := gemini.NewModel(ctx, youtubeModelName, &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "youtube_analyser_agent",
		Model:       model,
		Description: "Accepts a YouTube URL, extracts audio, and uses Gemini multimodal to identify the Ragam, Talam, composition, and artist.",
		Instruction: youtubeInstruction,
		Tools: []tool.Tool{
			tools.FetchYouTubeMetadataTool,
			tools.ExtractAudioTool,
			tools.AnalyseAudioWithGeminiTool,
			tools.LookupRagaTool,
			tools.SearchKritisTool,
		},
		BeforeAgentCallbacks: []agent.BeforeAgentCallback{
			func(ctx agent.CallbackContext) (*genai.Content, error) {
				slog.Info("agent invocation", "agent", "youtube_analyser_agent")
				return nil, nil
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return a, nil
}
