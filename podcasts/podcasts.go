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
	Episodes      []Episode
}

type Owner struct {
	copyright string
	creator   string
	email     string
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
