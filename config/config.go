package config

type PodcastinationConfig struct {
	PostgresDatasource string `json:"postgres_datasource"`
	PullDirectory      string `json:"pull_directory"`
	PodcastDirectory   string `json:"podcast_directory"`
}
