package models

type Video struct {
	Name        string `json:"name,omitempty"`
	Url         string `json:"url,omitempty"`
	Thumbnail   string `json:"thumbnail,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
	VideoLength string `json:"video_length,omitempty"`
	Playlist    string `json:"playlist,omitempty"`
}
