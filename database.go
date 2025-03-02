package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ja-howell/playlister/models"
)

type Database map[string][]models.Video

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
