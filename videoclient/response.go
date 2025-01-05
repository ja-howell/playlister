package videoclient

type Response struct {
	NextPageToken string     `json:"nextPageToken,omitempty"`
	Items         []RawVideo `json:"items,omitempty"`
}

type RawVideo struct {
	Snippet Snippet `json:"snippet,omitempty"`
}

type Snippet struct {
	PublishedAt string               `json:"publishedAt,omitempty"`
	Title       string               `json:"title,omitempty"`
	Thumbnails  map[string]Thumbnail `json:"thumbnails,omitempty"`
	ResourceId  ResourceId           `json:"resourceId,omitempty"`
}

type Thumbnail struct {
	Url    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

type ResourceId struct {
	VideoId string `json:"videoId,omitempty"`
}
