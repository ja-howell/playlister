package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	PlaylistID         string   `json:"playlist_id,omitempty"`
	APIKey             string   `json:"api_key`
	IgnoredPlaylists   []string `json:"ignored_playlists,omitempty"`
	LastCollectionDate string   `json:"last_collection_date,omitempty"`
}

func newConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read file: %w", err)
	}

	var config Config

	if err := json.Unmarshal(contents, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}
	return config, nil
}
