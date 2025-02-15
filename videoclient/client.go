package videoclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	apiKey string
}

type PageToken string

const FirstToken PageToken = ""

func New(apiKey string) Client {
	return Client{apiKey: apiKey}
}

func (c Client) GetResponse(nextPageToken PageToken) (Response, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=UU3tNpTOHsTnkmbwztCs30sA&maxResults=10&key=%s&maxResults=50", c.apiKey)
	if nextPageToken != FirstToken {
		url = fmt.Sprintf("%s&pageToken=%s", url, nextPageToken)
	}
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

func (c Client) GetVideoLength(videoId string) (string, error) {
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

	rawLength := x.Items[0].ContentDetails.Duration

	length := formatLength(rawLength)

	return length, nil
}

func formatLength(raw string) string {
	//PT#H#M#S
	//h:m:s
	s := strings.TrimPrefix(raw, "PT")
	h, s, found := strings.Cut(s, "H")
	if found {
		h = h + ":"
	} else {
		s = h
		h = ""
	}
	m, s, found := strings.Cut(s, "M")
	if !found {
		m = "0"
		s = m
	}
	sec, _, _ := strings.Cut(s, "S")
	return fmt.Sprintf("%s%s:%s", h, m, sec)

}
