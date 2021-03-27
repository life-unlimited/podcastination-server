package main

import (
	"life-unlimited/podcastination/app"
	"life-unlimited/podcastination/config"
)

func main() {
	// TODO: Read the config from file.
	podcastination := app.NewApp(config.PodcastinationConfig{
		PullDirectory:    "",
		PodcastDirectory: "",
	})
	if err := podcastination.Boot(); err != nil {
		panic(err)
	}
}
