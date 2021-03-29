package app

import (
	"database/sql"
	"fmt"
	"life-unlimited/podcastination/config"
)

type App struct {
	config config.PodcastinationConfig
	db     *sql.DB
}

// NewApp creates a new App.
func NewApp(config config.PodcastinationConfig) *App {
	return &App{
		config: config,
	}
}

// Boot boots the App.
func (a *App) Boot() error {
	// Connect to database.
	db, err := connectToDB(a.config.PostgresDatasource)
	if err != nil {
		panic(fmt.Errorf("could not open db connection: %v", err))
	}
	a.db = db
	defer closeDB(a.db)
	// TODO
	return nil
}

func connectToDB(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}

func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		panic(fmt.Errorf("could not close db connection: %v", err))
	}
}
