package tools

import (
	"context"
	"os"
	"testing"
)

func TestFetchYouTubeMetadata_InvalidURL(t *testing.T) {
	_, err := FetchYouTubeMetadata(context.Background(), "not-a-url")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestExtractAudio_InvalidURL(t *testing.T) {
	_, err := ExtractAudio(context.Background(), "not-a-url", 90)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestExtractAudio_InvalidDuration(t *testing.T) {
	_, err := ExtractAudio(context.Background(), "https://youtube.com/watch?v=123", 0)
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestAnalyseAudioWithGemini_EmptyAudio(t *testing.T) {
	_, err := AnalyseAudioWithGemini(context.Background(), nil, nil)
	if err == nil {
		t.Fatal("expected error for empty audio")
	}
}

func TestFetchYouTubeMetadata_NotYouTube(t *testing.T) {
	_, err := FetchYouTubeMetadata(context.Background(), "https://example.com")
	if err == nil {
		t.Fatal("expected error for non-YouTube URL")
	}
}

func TestYouTubeIntegration(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping network-dependent test in CI")
	}
	t.Skip("YouTube API and yt-dlp integration not yet implemented")
}
