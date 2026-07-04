// Package tools provides lookup and search tools for the Nāda Guru knowledge base.
package tools

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// YouTubeMetadata holds video information from the YouTube Data API v3.
type YouTubeMetadata struct {
	VideoID     string `json:"video_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ChannelName string `json:"channel_name"`
	DurationSec int    `json:"duration_seconds"`
}

// FetchYouTubeMetadata retrieves video metadata from the YouTube Data API v3.
//
// ADK tool name: "fetch_youtube_metadata"
// Parameters:    url (string) — full YouTube URL
// Env required:  YOUTUBE_API_KEY
func FetchYouTubeMetadata(ctx context.Context, rawURL string) (*YouTubeMetadata, error) {
	if strings.TrimSpace(rawURL) == "" {
		return nil, fmt.Errorf("url must not be empty")
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid url: %w", err)
	}
	if !strings.Contains(u.Host, "youtube.com") && !strings.Contains(u.Host, "youtu.be") {
		return nil, fmt.Errorf("not a YouTube URL: %q", rawURL)
	}
	return nil, fmt.Errorf("YouTube Data API integration not yet implemented")
}

// ExtractAudio downloads and extracts the first 90 seconds of audio from a
// YouTube URL using yt-dlp, returning raw audio bytes in WAV format.
//
// ADK tool name: "extract_audio"  (internal)
// Parameters:    url (string), maxSeconds (int)
func ExtractAudio(ctx context.Context, rawURL string, maxSeconds int) ([]byte, error) {
	if strings.TrimSpace(rawURL) == "" {
		return nil, fmt.Errorf("url must not be empty")
	}
	if maxSeconds <= 0 {
		return nil, fmt.Errorf("maxSeconds must be positive")
	}
	return nil, fmt.Errorf("audio extraction requires yt-dlp integration")
}
