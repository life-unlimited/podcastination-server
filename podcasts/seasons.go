package podcasts

type Season struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	Description   string `json:"description"`
	ImageLocation string `json:"image_location"`
	PodcastId     int    `json:"podcast_id"`
	Num           int    `json:"num"`
	Key           string `json:"key"`
}
