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
	interval() time.Duration
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
		alive := true
		for alive {
			select {
			case <-newJob.stop:
				alive = false
				break
			case <-time.After(j.interval()):
				if err := j.run(); err != nil {
					log.Fatalf("%s could not run job: %v", jobLogPrefix(myJob), err)
				}
				break
			}
		}
	}(newJob)
	log.Printf("scheduled job %s (every %v, initial run: %v)", j.name(), j.interval(), initialRun)
}

// Stop stops all registered jobs.
func (s *Scheduler) Stop() {
	totalJobCount := len(s.jobs)
	for _, j := range s.jobs {
		log.Printf("scheduler stopping %d/%d jobs...", len(s.jobs), totalJobCount)
		j.stop <- struct{}{}
	}
}

func jobLogPrefix(job *job) string {
	name := (*(*job).job).name()
	return fmt.Sprintf("[JOB] %s:", name)
}
