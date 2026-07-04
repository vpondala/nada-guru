// Package knowledge embeds all knowledge base JSON files into the binary.
package knowledge

import (
	"embed"
)

//go:embed ragas.json
var RagasJSON []byte

//go:embed talas.json
var TalasJSON []byte

//go:embed kritis.json
var KritisJSON []byte

//go:embed composers.json
var ComposersJSON []byte

//go:embed lyrics
var lyricsFS embed.FS

// ReadLyricsFile reads an embedded lyrics file by its filename.
func ReadLyricsFile(name string) ([]byte, error) {
	return lyricsFS.ReadFile("lyrics/" + name)
}

// LyricsFilenames returns all embedded lyrics filenames.
func LyricsFilenames() ([]string, error) {
	entries, err := lyricsFS.ReadDir("lyrics")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
