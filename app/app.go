package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"life-unlimited/podcastination/config"
	"life-unlimited/podcastination/stores"
	"life-unlimited/podcastination/tasks"
	"log"
	"time"
)

type App struct {
	config    config.PodcastinationConfig
	db        *sql.DB
	scheduler *tasks.Scheduler
	Stores    Stores
}

type Stores struct {
	Podcasts stores.PodcastStore
	Owners   stores.OwnerStore
	Seasons  stores.SeasonStore
	Episodes stores.EpisodeStore
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
	// Setup stores.
	a.Stores = Stores{
		Podcasts: stores.PodcastStore{DB: a.db},
		Owners:   stores.OwnerStore{DB: a.db},
		Seasons:  stores.SeasonStore{DB: a.db},
		Episodes: stores.EpisodeStore{DB: a.db},
	}
	// Check database connection.
	_, err = a.Stores.Podcasts.All()
	if err != nil {
		return fmt.Errorf("could not connect to db: %v", err)
	}
	log.Println("connection to database established.")
	// Create scheduler.
	a.scheduler = tasks.NewScheduler(tasks.SchedulingConfig{
		PullDir:        a.config.PullDir,
		PodcastDir:     a.config.PodcastDir,
		ImportInterval: time.Duration(a.config.ImportInterval) * time.Minute,
	}, a.db)
	// Perform integrity check.
	// TODO: Perform integrity check.
	// Let's go.
	a.scheduler.ScheduleJob(&tasks.ImportJob{
		PullDir:    a.config.PullDir,
		PodcastDir: a.config.PodcastDir,
		Store: tasks.ImportJobStores{
			Podcasts: a.Stores.Podcasts,
			Owners:   a.Stores.Owners,
			Seasons:  a.Stores.Seasons,
			Episodes: a.Stores.Episodes,
		},
	}, true)
	return nil
}

func (a *App) Shutdown() error {
	if a.db != nil {
		return closeDB(a.db)
	}
	return nil
}

// connectToDB connects to the database with postgres driver.
func connectToDB(dataSourceName string) (*sql.DB, error) {
	return sql.Open("postgres", dataSourceName)
}

// closeDB closes the database connection.
func closeDB(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("could not close db connection: %v", err)
	}
	return nil
}
