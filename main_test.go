package main

import (
	"testing"

	"github.com/ja-howell/playlister/models"
	"github.com/ja-howell/playlister/videoclient"
	"github.com/stretchr/testify/assert"
)

func TestGetVideo(t *testing.T) {
	snippet := videoclient.Snippet{
		PublishedAt: "2009-10-21T18:00:00Z",
		Title:       "Fake Title",
		Thumbnails: map[string]videoclient.Thumbnail{
			"standard": videoclient.Thumbnail{
				Url:    "www.testing.com",
				Width:  100,
				Height: 100,
			},
		},
		ResourceId: videoclient.ResourceId{
			VideoId: "fakeid",
		},
	}

	expected := models.Video{
		Name:        "Fake Title",
		Url:         "https://www.youtube.com/watch?v=fakeid",
		Thumbnail:   "www.testing.com",
		PublishedAt: "2009-10-21T18:00:00Z",
		VideoLength: "1:05:00",
	}

	assert := assert.New(t)
	got := getVideo(testClient{}, snippet)
	assert.Equal(expected, got)
}

type testClient struct{}

func (t testClient) GetResponse(nextPageToken videoclient.PageToken) (videoclient.Response, error) {
	return videoclient.Response{}, nil
}

func (t testClient) GetVideoLength(videoId string) (string, error) {
	return "1:05:00", nil
}
