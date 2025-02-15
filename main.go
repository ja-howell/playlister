package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ja-howell/playlister/models"
	"github.com/ja-howell/playlister/videoclient"
)

const apiKeyFilepath = "API_KEY"

func main() {
	err := run()
	if err != nil {
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}

func run() error {
	apiKey, err := getAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API Key: %w", err)
	}
	client := videoclient.New(apiKey)

	config, err := newConfig("config.json")
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// database, err := newDatabase("../db.json")
	// if err != nil {
	// 	return fmt.Errorf("failed to create database: %w", err)
	// }
	videos, err := GetVideosSince(client, config.LastCollectionDate)
	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}

	for _, video := range videos {
		fmt.Printf("Length: %v     Name: %v\n", video.VideoLength, video.Name)
	}
	fmt.Println(len(videos))

	// _ = database
	return nil
}

func getAPIKey() (string, error) {
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

// TODO GetVideosSince(date)
// Get all the videos since the last collection date
func GetVideosSince(c videoclient.Client, lastCollectionDate string) ([]models.Video, error) {
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
			videos = append(videos, getVideo(c, snippet))
		}
		next = videoclient.PageToken(response.NextPageToken)
	}

	return videos, nil
}

func getVideo(c videoclient.Client, snippet videoclient.Snippet) models.Video {
	video := convertSnippettoVideo(snippet)
	videoLength, err := c.GetVideoLength(snippet.ResourceId.VideoId)
	if err != nil {
		log.Printf("Failed to create video length: %v", err)
	}
	video.VideoLength = videoLength
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
