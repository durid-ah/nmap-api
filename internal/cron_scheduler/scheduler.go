package cronscheduler

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/durid-ah/nmap-api/internal/config"
	"github.com/durid-ah/nmap-api/internal/nmap_scanner"
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
			func(_ctx context.Context) {
				scannerCtx, scannerCtxCancel := context.WithTimeout(
					context.WithValue(_ctx, ownerContextKey, "scanner"), 5*time.Minute)
				defer scannerCtxCancel()
				slog.Info("running scanner", "context", scannerCtx)
				scanner, err := nmapscanner.NewNmapScanner(scannerCtx, config)
				if err != nil {
					slog.Error("unable to create nmap scanner", "error", err)
					return
				}
				err = scanner.Run(scannerCtx)
				if err != nil {
					slog.Error("unable to run nmap scan", "error", err)
				}
			},
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
	// s.job.RunNow()
}

func (s *BackgroundScheduler) Shutdown() error {
	s.cancel()
	return s.scheduler.Shutdown()
}
