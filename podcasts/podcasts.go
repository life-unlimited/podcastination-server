package podcasts

type Podcast struct {
	Id            int         `json:"id"`
	Title         string      `json:"title"`
	Subtitle      string      `json:"subtitle"`
	Language      Language    `json:"language"`
	OwnerId       int         `json:"owner_id"`
	Description   string      `json:"description"`
	Keywords      []string    `json:"keywords"`
	Link          string      `json:"link"`
	FeedLink      string      `json:"feed_link"`
	ImageLocation string      `json:"image_location"`
	PodcastType   PodcastType `json:"podcast_type"`
	Key           string      `json:"key"`
}

type PodcastType string

const (
	TypeSermon  PodcastType = "sermon"
	TypeService PodcastType = "service"
	TypeEvent   PodcastType = "event"
)

type Language string

const (
	LangDE Language = "de-de"
	LangEN Language = "en-us"
)
