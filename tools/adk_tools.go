// Package tools provides lookup and search tools for the Nāda Guru knowledge base.
package tools

import (
	"fmt"
	"sync"

	"github.com/vpondala/nada-guru/knowledge"
	"github.com/vpondala/nada-guru/session"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

var (
	LookupRagaTool             tool.Tool
	SearchRagasBySwaraTool     tool.Tool
	SearchRagasByMoodTool      tool.Tool
	LookupTalaTool             tool.Tool
	SearchTalasByBeatsTool     tool.Tool
	SearchKritisTool           tool.Tool
	LookupKritiTool            tool.Tool
	LookupComposerTool         tool.Tool
	SearchComposersByLangTool  tool.Tool
	GetLyricsTool              tool.Tool
	ScrapeLyricsTool           tool.Tool
	TransliterateTextTool      tool.Tool
	FetchYouTubeMetadataTool   tool.Tool
	ExtractAudioTool           tool.Tool
	AnalyseAudioWithGeminiTool tool.Tool
)

var adkToolsInitOnce sync.Once

// InitADKTools registers all function tools.
func InitADKTools() {
	adkToolsInitOnce.Do(func() {
		var err error

		// 1. lookup_raga
		LookupRagaTool, err = functiontool.New(functiontool.Config{
			Name:        "lookup_raga",
			Description: "Lookup a raga by name or alias",
		}, func(ctx agent.ToolContext, input struct{ Name string `json:"name"` }) (struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Error string `json:"error,omitempty"`
		}, error) {
			r, err := LookupRaga(ctx, input.Name)
			if err != nil {
				return struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Error string `json:"error,omitempty"`
				}{Error: err.Error()}, nil
			}
			_ = ctx.State().Set(session.KeyLastRaga, r.ID)
			return struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Error string `json:"error,omitempty"`
			}{ID: r.ID, Name: r.Name}, nil
		})
		if err != nil {
			panic(fmt.Errorf("failed to create lookup_raga tool: %w", err))
		}

		// 2. search_ragas_by_swara
		SearchRagasBySwaraTool, err = functiontool.New(functiontool.Config{
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
			results, err := SearchRagasBySwara(ctx, input.Swaras)
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
			// If there's a unique match, set it in the session
			if len(results) == 1 {
				_ = ctx.State().Set(session.KeyLastRaga, results[0].ID)
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
			panic(fmt.Errorf("failed to create search_ragas_by_swara tool: %w", err))
		}

		// 3. search_ragas_by_mood
		SearchRagasByMoodTool, err = functiontool.New(functiontool.Config{
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
			results, err := SearchRagasByMood(ctx, input.Rasa, input.TimeOfDay)
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
			if len(results) == 1 {
				_ = ctx.State().Set(session.KeyLastRaga, results[0].ID)
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
			panic(fmt.Errorf("failed to create search_ragas_by_mood tool: %w", err))
		}

		// 4. lookup_tala
		LookupTalaTool, err = functiontool.New(functiontool.Config{
			Name:        "lookup_tala",
			Description: "Lookup a tala by name",
		}, func(ctx agent.ToolContext, input struct{ Name string `json:"name"` }) (struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Error string `json:"error,omitempty"`
		}, error) {
			t, err := LookupTala(ctx, input.Name)
			if err != nil {
				return struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Error string `json:"error,omitempty"`
				}{Error: err.Error()}, nil
			}
			_ = ctx.State().Set(session.KeyLastTala, t.ID)
			return struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Error string `json:"error,omitempty"`
			}{ID: t.ID, Name: t.Name}, nil
		})
		if err != nil {
			panic(fmt.Errorf("failed to create lookup_tala tool: %w", err))
		}

		// 5. search_talas_by_beats
		SearchTalasByBeatsTool, err = functiontool.New(functiontool.Config{
			Name:        "search_talas_by_beats",
			Description: "Search talas by beat count",
		}, func(ctx agent.ToolContext, input struct{ Beats int `json:"beats"` }) (struct {
			Results []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"results"`
			Error string `json:"error,omitempty"`
		}, error) {
			results, err := SearchTalasByBeats(ctx, input.Beats)
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
			if len(results) == 1 {
				_ = ctx.State().Set(session.KeyLastTala, results[0].ID)
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
			panic(fmt.Errorf("failed to create search_talas_by_beats tool: %w", err))
		}

		// 6. search_kritis
		SearchKritisTool, err = functiontool.New(functiontool.Config{
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
			results, err := SearchKritis(ctx, knowledge.KritiFilter{
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
			if len(results) == 1 {
				_ = ctx.State().Set(session.KeyLastKriti, results[0].ID)
				_ = ctx.State().Set(session.KeyLastComposer, results[0].Composer)
				if r, err := store.LookupRaga(results[0].Ragam); err == nil {
					_ = ctx.State().Set(session.KeyLastRaga, r.ID)
				}
				if t, err := store.LookupTala(results[0].Talam); err == nil {
					_ = ctx.State().Set(session.KeyLastTala, t.ID)
				}
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
			panic(fmt.Errorf("failed to create search_kritis tool: %w", err))
		}

		// 7. lookup_kriti
		LookupKritiTool, err = functiontool.New(functiontool.Config{
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
			k, err := LookupKriti(ctx, input.ID)
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
			_ = ctx.State().Set(session.KeyLastKriti, k.ID)
			_ = ctx.State().Set(session.KeyLastComposer, k.Composer)
			if r, err := store.LookupRaga(k.Ragam); err == nil {
				_ = ctx.State().Set(session.KeyLastRaga, r.ID)
			}
			if t, err := store.LookupTala(k.Talam); err == nil {
				_ = ctx.State().Set(session.KeyLastTala, t.ID)
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
		if err != nil {
			panic(fmt.Errorf("failed to create lookup_kriti tool: %w", err))
		}

		// 8. lookup_composer
		LookupComposerTool, err = functiontool.New(functiontool.Config{
			Name:        "lookup_composer",
			Description: "Lookup a composer by name",
		}, func(ctx agent.ToolContext, input struct{ Name string `json:"name"` }) (struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Error string `json:"error,omitempty"`
		}, error) {
			c, err := LookupComposer(ctx, input.Name)
			if err != nil {
				return struct {
					ID    string `json:"id"`
					Name  string `json:"name"`
					Error string `json:"error,omitempty"`
				}{Error: err.Error()}, nil
			}
			_ = ctx.State().Set(session.KeyLastComposer, c.ID)
			return struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Error string `json:"error,omitempty"`
			}{ID: c.ID, Name: c.Name}, nil
		})
		if err != nil {
			panic(fmt.Errorf("failed to create lookup_composer tool: %w", err))
		}

		// 9. search_composers_by_language
		SearchComposersByLangTool, err = functiontool.New(functiontool.Config{
			Name:        "search_composers_by_language",
			Description: "Search composers by language",
		}, func(ctx agent.ToolContext, input struct{ Language string `json:"language"` }) (struct {
			Results []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"results"`
			Error string `json:"error,omitempty"`
		}, error) {
			results, err := SearchComposersByLanguage(ctx, input.Language)
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
			if len(results) == 1 {
				_ = ctx.State().Set(session.KeyLastComposer, results[0].ID)
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
			panic(fmt.Errorf("failed to create search_composers_by_language tool: %w", err))
		}

		// 10. get_lyrics
		GetLyricsTool, err = functiontool.New(functiontool.Config{
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
			lyrics, err := GetLyrics(ctx, input.KritiID)
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
			_ = ctx.State().Set(session.KeyLastLyrics, lyrics.KritiID)
			_ = ctx.State().Set(session.KeyLastKriti, lyrics.KritiID)
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
			panic(fmt.Errorf("failed to create get_lyrics tool: %w", err))
		}

		// 11. scrape_lyrics
		ScrapeLyricsTool, err = functiontool.New(functiontool.Config{
			Name:        "scrape_lyrics",
			Description: "Scrape lyrics from external sources (internal)",
		}, func(ctx agent.ToolContext, input struct {
			KritiName    string `json:"kriti_name"`
			ComposerName string `json:"composer_name"`
		}) (struct {
			Error string `json:"error,omitempty"`
		}, error) {
			_, err := ScrapeLyrics(ctx, input.KritiName, input.ComposerName)
			if err != nil {
				return struct {
					Error string `json:"error,omitempty"`
				}{Error: err.Error()}, nil
			}
			return struct {
				Error string `json:"error,omitempty"`
			}{}, nil
		})
		if err != nil {
			panic(fmt.Errorf("failed to create scrape_lyrics tool: %w", err))
		}

		// 12. transliterate_text
		TransliterateTextTool, err = functiontool.New(functiontool.Config{
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
			result, err := TransliterateText(ctx, input.Text, input.SourceLang, input.TargetLang)
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
			_ = ctx.State().Set(session.KeyScriptTarget, input.TargetLang)
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
			panic(fmt.Errorf("failed to create transliterate_text tool: %w", err))
		}

		// 13. fetch_youtube_metadata
		FetchYouTubeMetadataTool, err = functiontool.New(functiontool.Config{
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
			m, err := FetchYouTubeMetadata(ctx, input.URL)
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
			panic(fmt.Errorf("failed to create fetch_youtube_metadata tool: %w", err))
		}

		// 14. extract_audio
		ExtractAudioTool, err = functiontool.New(functiontool.Config{
			Name:        "extract_audio",
			Description: "Extract audio from a YouTube URL",
		}, func(ctx agent.ToolContext, input struct {
			URL        string `json:"url"`
			MaxSeconds int    `json:"max_seconds"`
		}) (struct {
			Size  int    `json:"size"`
			Error string `json:"error,omitempty"`
		}, error) {
			_, err := ExtractAudio(ctx, input.URL, input.MaxSeconds)
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
			panic(fmt.Errorf("failed to create extract_audio tool: %w", err))
		}

		// 15. analyse_audio_with_gemini
		AnalyseAudioWithGeminiTool, err = functiontool.New(functiontool.Config{
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
			var hints *YouTubeMetadata
			if input.Title != "" || input.Channel != "" {
				hints = &YouTubeMetadata{
					Title:       input.Title,
					ChannelName: input.Channel,
				}
			}
			result, err := AnalyseAudioWithGemini(ctx, input.AudioWAV, hints)
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
			panic(fmt.Errorf("failed to create analyse_audio_with_gemini tool: %w", err))
		}
	})
}
