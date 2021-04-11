package tasks

import (
	"database/sql"
	"log"
	"time"
)

type Scheduler struct {
	config SchedulingConfig
	db     *sql.DB
	stop   chan bool
}

type SchedulingConfig struct {
	PullDir        string
	PodcastDir     string
	ImportInterval time.Duration
}

type SchedulingTask interface {
	run()
	importInterval() time.Duration
}

func NewScheduler(config SchedulingConfig, db *sql.DB) *Scheduler {
	return &Scheduler{
		config: config,
		db:     db,
		stop:   make(chan bool),
	}
}

func (s *Scheduler) Run() {
	// Run initial imports.
	log.Println("Performing initial import...")
	if err := s.runImports(); err != nil {
		log.Fatalf("could not run initial import tasks: %v", err)
	}
	log.Println("Done.")
	// TODO: Setup scheduling and interrupt handling.
	//alive := true
	//for alive {
	//	select {
	//	case <-s.stop:
	//		alive = false
	//		break
	//	case <-time.After(s.config.ImportInterval):
	//		if err := s.runImports(); err != nil {
	//			log.Fatalf("could not run import tasks: %v", err)
	//		}
	//		break
	//	}
	//}
}

func (s *Scheduler) Stop() {
	s.stop <- true
}
