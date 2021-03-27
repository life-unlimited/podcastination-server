package app

import "life-unlimited/podcastination/config"

type App struct {
	config config.PodcastinationConfig
}

// NewApp creates a new App.
func NewApp(config config.PodcastinationConfig) *App {
	return &App{
		config: config,
	}
}

// Boot boots the App.
func (a *App) Boot() error {
	// TODO
	return nil
}
