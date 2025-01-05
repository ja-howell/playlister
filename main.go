package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	resp, err := getResponse(apiKey)
	if err != nil {
		return fmt.Errorf("failed to get a response: %w", err)
	}

	config, err := newConfig("config.json")
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// database, err := newDatabase("../db.json")
	// if err != nil {
	// 	return fmt.Errorf("failed to create database: %w", err)
	// }

	fmt.Println(resp)

	_ = resp
	_ = config
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

func getResponse(apiKey string) (Response, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=UU3tNpTOHsTnkmbwztCs30sA&maxResults=10&key=%s&maxResults=50", apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return Response{}, fmt.Errorf("failed to fetch endpoint: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("status not ok: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, fmt.Errorf("failed to read body: %w", err)
	}
	response := Response{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}
	return response, nil
}

func getVideo(raw RawVideo, apiKey string) Video {
	video := convertRawtoVideo(raw)
	videoLength, err := getVideoLength(raw.Snippet.ResourceId.VideoId, apiKey)
	if err != nil {
		log.Printf("Failed to create video length: %v", err)
	}
	video.VideoLength = videoLength
	return video
}

func convertRawtoVideo(raw RawVideo) Video {
	return Video{
		Name:        raw.Snippet.Title,
		Url:         "https://www.youtube.com/watch?v=" + raw.Snippet.ResourceId.VideoId,
		Thumbnail:   raw.Snippet.Thumbnails["standard"].Url,
		PublishedAt: raw.Snippet.PublishedAt,
	}
}

func getVideoLength(videoId string, apiKey string) (string, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=contentDetails&id=%s&key=%s", videoId, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to collect video length: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status not ok: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	// TODO: Unmarshal into an anonymous struct
	// x := struct {
	// 	foo string `json:"Foo"`
	// }{}

	_ = body
	// json.Unmarshal(body, &x)

	return "", nil

}
