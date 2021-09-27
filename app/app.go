package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/life-unlimited/podcastination-server/config"
	"github.com/life-unlimited/podcastination-server/stores"
	"github.com/life-unlimited/podcastination-server/tasks"
	"github.com/life-unlimited/podcastination-server/web_server"
	"github.com/pkg/errors"
	"log"
	"time"
)

type App struct {
	config    config.PodcastinationConfig
	db        *sql.DB
	scheduler *tasks.Scheduler
	webServer *web_server.WebServer
	Stores    stores.Stores
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
	db, err := connectDB(a.config.PostgresDatasource, defaultMaxDBConnections)
	if err != nil {
		panic(fmt.Errorf("could not open db connection: %v", err))
	}
	err = testDBConnection(db)
	if err != nil {
		return errors.Wrap(err, "test db connection")
	}
	// Perform database migrations if needed.
	err = performDBMigrations(db)
	if err != nil {
		return errors.Wrap(err, "perform database migrations")
	}
	a.db = db
	// Setup stores.Stores.
	a.Stores = stores.Stores{
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
	// Refresh all podcast.xml files.

	// Let's go.
	a.scheduler.ScheduleJob(&tasks.ImportJob{
		StaticContentURL: a.config.StaticContentURL,
		PullDir:          a.config.PullDir,
		PodcastDir:       a.config.PodcastDir,
		ImportInterval:   time.Duration(a.config.ImportInterval) * time.Minute,
		Store: tasks.ImportJobStores{
			Podcasts: a.Stores.Podcasts,
			Owners:   a.Stores.Owners,
			Seasons:  a.Stores.Seasons,
			Episodes: a.Stores.Episodes,
		},
	}, true)
	// Start web web_server.
	a.webServer = web_server.NewServer(web_server.Config{
		StaticDir: a.config.PodcastDir,
		Addr:      a.config.ServerAddr,
	}, &a.Stores)
	err = a.webServer.Start()
	if err != nil {
		log.Fatalf("could not start web server: %v", err)
	}
	return nil
}

// Shutdown shuts down the app.
func (a *App) Shutdown() error {
	a.scheduler.Stop()
	if err := a.webServer.Stop(); err != nil {
		return fmt.Errorf("stop web server: %v", err)
	}
	if a.db != nil {
		return closeDB(a.db)
	}
	return nil
}

// closeDB closes the database connection.
func closeDB(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return fmt.Errorf("could not close db connection: %v", err)
	}
	return nil
}
