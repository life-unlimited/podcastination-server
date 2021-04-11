package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"life-unlimited/podcastination/config"
	"life-unlimited/podcastination/tasks"
	"time"
)

type App struct {
	config    config.PodcastinationConfig
	db        *sql.DB
	scheduler *tasks.Scheduler
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
	// Create scheduler.
	a.scheduler = tasks.NewScheduler(tasks.SchedulingConfig{
		PullDir:        a.config.PullDir,
		PodcastDir:     a.config.PodcastDir,
		ImportInterval: time.Duration(a.config.ImportInterval) * time.Minute,
	}, a.db)
	// Perform integrity check.
	// TODO: Perform integrity check.
	// Let's go.
	a.scheduler.Run()
	return nil
}

// connectToDB connects to the database with postgres driver.
func connectToDB(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}

// closeDB closes the database connection.
func closeDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		panic(fmt.Errorf("could not close db connection: %v", err))
	}
}
