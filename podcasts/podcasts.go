package podcasts

type Podcast struct {
	Id            int
	Title         string
	Subtitle      string
	Language      Language
	OwnerId       int
	Description   string
	Keywords      []string
	Link          string
	FeedLink      string
	ImageLocation string
	PodcastType   PodcastType
	Key           string
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
