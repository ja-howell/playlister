package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Video struct {
	Name        string `json:"name,omitempty"`
	Url         string `json:"url,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	VideoLength string `json:"video_length,omitempty"`
}

type Database map[string][]Video

func newDatabase(path string) (Database, error) {
	f, err := os.Open(path)
	if err != nil {
		return Database{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return Database{}, fmt.Errorf("failed to read file: %w", err)
	}

	var database Database

	if err := json.Unmarshal(contents, &database); err != nil {
		return Database{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}
	return database, nil
}
