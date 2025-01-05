package videoclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ja-howell/playlister/models"
)

type Client struct {
	apiKey string
}

func New(apiKey string) Client {
	return Client{apiKey: apiKey}
}

//TODO GetVideosSince(date)
//Get all the videos since the last collection date

func (c Client) getResponse() (Response, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=UU3tNpTOHsTnkmbwztCs30sA&maxResults=10&key=%s&maxResults=50", c.apiKey)
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

func (c Client) getVideo(raw RawVideo) models.Video {
	video := convertRawtoVideo(raw)
	videoLength, err := c.getVideoLength(raw.Snippet.ResourceId.VideoId)
	if err != nil {
		log.Printf("Failed to create video length: %v", err)
	}
	video.VideoLength = videoLength
	return video
}

func convertRawtoVideo(raw RawVideo) models.Video {
	return models.Video{
		Name:        raw.Snippet.Title,
		Url:         "https://www.youtube.com/watch?v=" + raw.Snippet.ResourceId.VideoId,
		Thumbnail:   raw.Snippet.Thumbnails["standard"].Url,
		PublishedAt: raw.Snippet.PublishedAt,
	}
}

func (c Client) getVideoLength(videoId string) (string, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=contentDetails&id=%s&key=%s", videoId, c.apiKey)
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

	x := struct {
		Items []struct {
			ContentDetails struct {
				Duration string `json:"duration,omitempty"`
			} `json:"contentDetails,omitempty"`
		} `json:"items,omitempty"`
	}{}

	json.Unmarshal(body, &x)

	return x.Items[0].ContentDetails.Duration, nil

}
