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

const youtubeModelName = "gemini-3.1-pro"

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

	fetchTool, err := functiontool.New(functiontool.Config{
		Name:        "fetch_youtube_metadata",
		Description: "Fetch metadata for a YouTube video",
	}, func(ctx agent.ToolContext, input struct{ URL string `json:"url"` }) (struct {
		VideoID     string `json:"video_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		ChannelName string `json:"channel_name"`
		DurationSec int    `json:"duration_seconds"`
		Error       string `json:"error,omitempty"`
	}, error) {
		m, err := tools.FetchYouTubeMetadata(ctx, input.URL)
		if err != nil {
			return struct {
				VideoID     string `json:"video_id"`
				Title       string `json:"title"`
				Description string `json:"description"`
				ChannelName string `json:"channel_name"`
				DurationSec int    `json:"duration_seconds"`
				Error       string `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		return struct {
			VideoID     string `json:"video_id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			ChannelName string `json:"channel_name"`
			DurationSec int    `json:"duration_seconds"`
			Error       string `json:"error,omitempty"`
		}{
			VideoID:     m.VideoID,
			Title:       m.Title,
			Description: m.Description,
			ChannelName: m.ChannelName,
			DurationSec: m.DurationSec,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	extractTool, err := functiontool.New(functiontool.Config{
		Name:        "extract_audio",
		Description: "Extract audio from a YouTube URL",
	}, func(ctx agent.ToolContext, input struct {
		URL       string `json:"url"`
		MaxSeconds int   `json:"max_seconds"`
	}) (struct {
		Size  int    `json:"size"`
		Error string `json:"error,omitempty"`
	}, error) {
		_, err := tools.ExtractAudio(ctx, input.URL, input.MaxSeconds)
		if err != nil {
			return struct {
				Size  int    `json:"size"`
				Error string `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		return struct {
			Size  int    `json:"size"`
			Error string `json:"error,omitempty"`
		}{}, nil
	})
	if err != nil {
		return nil, err
	}

	analyseTool, err := functiontool.New(functiontool.Config{
		Name:        "analyse_audio_with_gemini",
		Description: "Analyse audio with Gemini multimodal",
	}, func(ctx agent.ToolContext, input struct {
		AudioWAV []byte `json:"audio_wav"`
		Title    string `json:"title,omitempty"`
		Channel  string `json:"channel,omitempty"`
	}) (struct {
		Ragam      string   `json:"ragam"`
		Talam      string   `json:"talam"`
		Kriti      string   `json:"kriti"`
		Artist     string   `json:"artist"`
		Confidence string   `json:"confidence"`
		Candidates []string `json:"candidates"`
		Error      string   `json:"error,omitempty"`
	}, error) {
	var hints *tools.YouTubeMetadata
	if input.Title != "" || input.Channel != "" {
		hints = &tools.YouTubeMetadata{
			Title:       input.Title,
			ChannelName: input.Channel,
		}
	}
		result, err := tools.AnalyseAudioWithGemini(ctx, input.AudioWAV, hints)
		if err != nil {
			return struct {
				Ragam      string   `json:"ragam"`
				Talam      string   `json:"talam"`
				Kriti      string   `json:"kriti"`
				Artist     string   `json:"artist"`
				Confidence string   `json:"confidence"`
				Candidates []string `json:"candidates"`
				Error      string   `json:"error,omitempty"`
			}{Error: err.Error()}, nil
		}
		return struct {
			Ragam      string   `json:"ragam"`
			Talam      string   `json:"talam"`
			Kriti      string   `json:"kriti"`
			Artist     string   `json:"artist"`
			Confidence string   `json:"confidence"`
			Candidates []string `json:"candidates"`
			Error      string   `json:"error,omitempty"`
		}{
			Ragam:      result.Ragam,
			Talam:      result.Talam,
			Kriti:      result.Kriti,
			Artist:     result.Artist,
			Confidence: result.Confidence,
			Candidates: result.Candidates,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	lookupRagaTool, err := functiontool.New(functiontool.Config{
		Name:        "lookup_raga",
		Description: "Lookup a raga by name",
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

	a, err := llmagent.New(llmagent.Config{
		Name:        "youtube_analyser_agent",
		Model:       model,
		Description: "Accepts a YouTube URL, extracts audio, and uses Gemini multimodal to identify the Ragam, Talam, composition, and artist.",
		Instruction: youtubeInstruction,
		Tools: []tool.Tool{
			fetchTool,
			extractTool,
			analyseTool,
			lookupRagaTool,
			searchKritisTool,
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
