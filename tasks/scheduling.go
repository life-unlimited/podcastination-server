package tasks

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Scheduler struct {
	config SchedulingConfig
	db     *sql.DB
	stop   chan struct{}
	jobs   []*job
}

type job struct {
	job  *SchedulingJob
	stop chan struct{}
}

type SchedulingConfig struct {
	PullDir        string
	PodcastDir     string
	ImportInterval time.Duration
}

type SchedulingJob interface {
	name() string
	run() error
	importInterval() time.Duration
}

func NewScheduler(config SchedulingConfig, db *sql.DB) *Scheduler {
	return &Scheduler{
		config: config,
		db:     db,
		stop:   make(chan struct{}),
	}
}

// ScheduleJob schedules and runs a SchedulingJob.
func (s *Scheduler) ScheduleJob(j SchedulingJob, initialRun bool) {
	newJob := &job{
		job:  &j,
		stop: make(chan struct{}),
	}
	s.jobs = append(s.jobs, newJob)
	// Run.
	go func(myJob *job) {
		if initialRun {
			if err := (*myJob.job).run(); err != nil {
				log.Printf("%s initial run failed: %v", jobLogPrefix(myJob), err)
			}
		}
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
	}(newJob)
}

func (s *Scheduler) Stop() {
	s.stop <- struct{}{}
	// TODO
}

func jobLogPrefix(job *job) string {
	name := (*(*job).job).name()
	return fmt.Sprintf("[JOB] %s:", name)
}
