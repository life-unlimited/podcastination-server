package podcasts

type Podcast struct {
	Title         string
	Subtitle      string
	Language      Language
	Owner         Owner
	Description   string
	Keywords      []string
	Link          string
	ImageLocation string
	Type          Type
	Seasons       []Season
}

type Owner struct {
	Name      string
	Email     string
	Copyright string
}

type Type string

const (
	TypeSermon  Type = "sermon"
	TypeService Type = "service"
	TypeEvent   Type = "event"
)

type Language string

const (
	LangDE Language = "de-de"
	LangEN Language = "en-us"
)
