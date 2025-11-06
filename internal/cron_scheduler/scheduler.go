package cronscheduler

import (
	"context"
	"log"
	"log/slog"

	"github.com/durid-ah/nmap-api/internal/config"
	"github.com/durid-ah/nmap-api/internal/nmap_scanner"
	"github.com/go-co-op/gocron/v2"
)

type contextKey string

const (
	ownerContextKey contextKey = "owner"
)

type BackgroundScheduler struct {
	scheduler gocron.Scheduler
	job       gocron.Job
	cancel    context.CancelFunc
}

func NewBackgroundScheduler(config *config.Config) *BackgroundScheduler {
	ctx, cancel := context.WithCancel(
		context.WithValue(context.Background(), ownerContextKey, "scheduler"))

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("error creating scheduler", "error", err)
		log.Fatal(err)
	}

	// add a job to the scheduler
	job, err := scheduler.NewJob(
		gocron.CronJob(config.NmapCronTab, false),
		gocron.NewTask(
			nmapscanner.CreateScannerTask(config),
		),
		gocron.WithContext(ctx),
	)

	if err != nil {
		slog.Error("error creating job", "error", err)
		log.Fatal(err)
	}

	slog.Info("job created", "job_id", job.ID())

	return &BackgroundScheduler{scheduler: scheduler, job: job, cancel: cancel}
}

func (s *BackgroundScheduler) Start() {
	s.scheduler.Start()
	s.job.RunNow()
}

func (s *BackgroundScheduler) Shutdown() error {
	s.cancel()
	return s.scheduler.Shutdown()
}
