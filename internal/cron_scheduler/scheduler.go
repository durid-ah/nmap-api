package cronscheduler

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/durid-ah/nmap-api/internal/config"
	"github.com/go-co-op/gocron/v2"
)

type contextKey string

const (
	ownerContextKey contextKey = "owner"
)

type BackgroundScheduler struct {
	scheduler gocron.Scheduler
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
	j, err := scheduler.NewJob(
		gocron.CronJob(config.NmapCronTab, false),
		gocron.NewTask(
			func(_ctx context.Context) {
				scannerCtx, scannerCtxCancel := context.WithTimeout(_ctx, 5*time.Minute)
				defer scannerCtxCancel()
				slog.Info("running scanner", "context", scannerCtx)
				slog.Info("running job")
			},
		),
		gocron.WithContext(ctx),
	)

	if err != nil {
		slog.Error("error creating job", "error", err)
		log.Fatal(err)
	}

	slog.Info("job created", "job_id", j.ID())

	return &BackgroundScheduler{scheduler: scheduler, cancel: cancel}
}

func (s *BackgroundScheduler) Start() {
	s.scheduler.Start()
}

func (s *BackgroundScheduler) Shutdown() error {
	s.cancel()
	return s.scheduler.Shutdown()
}
