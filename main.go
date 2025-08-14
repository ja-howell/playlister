package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ja-howell/playlister/models"
	"github.com/ja-howell/playlister/videoclient"
)

type Client interface {
	GetResponse(nextPageToken videoclient.PageToken) (videoclient.Response, error)
	GetVideoLength(videoId string) (string, error)
}

// Helper function to get environment variable with default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Get configuration from environment variables with defaults
	dataDir := getEnv("DATA_DIR", "/data")
	dbFilename := getEnv("DATABASE_FILENAME", "db.json")
	configFilename := getEnv("CONFIG_FILENAME", "config.json")
	apiKeyFilename := getEnv("API_KEY_FILENAME", "API_KEY")

	databasePath := fmt.Sprintf("%s/%s", dataDir, dbFilename)
	configPath := fmt.Sprintf("%s/%s", dataDir, configFilename)
	apiKeyFilepath := apiKeyFilename // Keep as is for backward compatibility

	err := run(databasePath, configPath, apiKeyFilepath)
	if err != nil {
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}

func run(databasePath, configPath, apiKeyFilepath string) error {
	config, err := newConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	apiKey := config.APIKey
	if apiKey == "" {
		apiKey, err = getAPIKey(apiKeyFilepath)
		if err != nil {
			return fmt.Errorf("failed to get API Key: %w", err)
		}
	}

	client := videoclient.New(apiKey)

	database, err := newDatabase(databasePath)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	videos, err := GetVideosSince(client, config.LastCollectionDate)
	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}
	log.Printf("Fetched %d videos", len(videos))

	for _, video := range videos {
		database[video.Playlist] = append(database[video.Playlist], video)
		log.Printf("Added '%s' to playlist: %s", video.Name, video.Playlist)
	}

	err = writeToFile(databasePath, database)
	if err != nil {
		return fmt.Errorf("failed to write database: %w", err)
	}

	config.LastCollectionDate = time.Now().Format("2006-01-02T15:04:05Z")
	err = writeToFile(configPath, config)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	log.Print("Finished updating")

	return nil
}

func getAPIKey(apiKeyFilepath string) (string, error) {
	f, err := os.Open(apiKeyFilepath)
	if err != nil {
		return "", fmt.Errorf("failed to open API Key file: %w", err)
	}

	defer f.Close()

	apiKey, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read API Key: %w", err)
	}

	return string(apiKey), nil
}

// Get all the videos since the last collection date
func GetVideosSince(c Client, lastCollectionDate string) ([]models.Video, error) {
	log.Printf("Fetching videos since %s", lastCollectionDate)
	videos := []models.Video{}

	next := videoclient.FirstToken

	//process the videos
	done := false
	for !done {
		response, err := c.GetResponse(next)
		if err != nil {
			return []models.Video{}, fmt.Errorf("failed to retrieve videos: %w", err)
		}
		for _, item := range response.Items {
			snippet := item.Snippet
			if snippet.PublishedAt < lastCollectionDate {
				done = true
				break
			}
			newVideo := getVideo(c, snippet)
			// filters out videos shorter than 10 minutes
			if len(newVideo.VideoLength) > 4 {
				videos = append(videos, newVideo)
			}
		}
		next = videoclient.PageToken(response.NextPageToken)
	}

	return videos, nil
}

func getVideo(c Client, snippet videoclient.Snippet) models.Video {
	video := convertSnippettoVideo(snippet)
	videoLength, err := c.GetVideoLength(snippet.ResourceId.VideoId)
	if err != nil {
		log.Printf("Failed to create video length: %v", err)
	}
	video.VideoLength = videoLength

	video.Playlist, video.Name = parsePlaylistFromName(video.Name)
	return video
}

func convertSnippettoVideo(snippet videoclient.Snippet) models.Video {
	return models.Video{
		Name:        snippet.Title,
		Url:         "https://www.youtube.com/watch?v=" + snippet.ResourceId.VideoId,
		Thumbnail:   snippet.Thumbnails["standard"].Url,
		PublishedAt: snippet.PublishedAt,
	}
}

func parsePlaylistFromName(rawName string) (playlist, name string) {
	lastParen := strings.LastIndex(rawName, "(")
	if lastParen == -1 {
		return "[Missing]", rawName
	}

	name = strings.TrimSpace(rawName[:lastParen])
	playlist = strings.TrimSpace(rawName[lastParen+1:])
	playlist = playlist[:len(playlist)-1]
	return playlist, name
}

func writeToFile(path string, obj any) error {
	b, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	err = os.WriteFile(path, b, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}
