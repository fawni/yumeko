package common

type Video struct {
	URL      string `json:"url,omitempty"`
	Uploader string `json:"uploaderName,omitempty"`
	Title    string `json:"title,omitempty"`
	Duration int    `json:"duration,omitempty"`
	Uploaded int64  `json:"uploaded,omitempty"`
	Error    string `json:"error,omitempty"`

	Details Stream
	Channel Channel
}

type Stream struct {
	Views       int64  `json:"views,omitempty"`
	Description string `json:"description,omitempty"`
	Likes       int64  `json:"likes,omitempty"`
	Dislikes    int64  `json:"dislikes,omitempty"`
	ChannelURL  string `json:"uploaderUrl,omitempty"`
}

type Channel struct {
	Subscibers float64 `json:"subscriberCount,omitempty"`
}

type Search struct {
	Items []Video `json:"items,omitempty"`
}
