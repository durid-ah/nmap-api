package main

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/durid-ah/nmap-api/internal/config"
	"github.com/durid-ah/nmap-api/internal/cron_scheduler"
)

func main() {
	cfg := config.NewConfig()

	// TODO: Child logger
	opts := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	// TODO: pass in individual logger?
	handler := slog.NewJSONHandler(os.Stdout, &opts)
	slog.SetDefault(slog.New(handler))

	// scheduler := initJob()
	scheduler := cronscheduler.NewBackgroundScheduler(cfg)
	scheduler.Start()

	defer func() {
		slog.Info("shutting down scheduler")
		err := scheduler.Shutdown()
		if err != nil {
			slog.Error("error shutting down scheduler", "error", err)
			log.Fatal(err)
		}
	}()

	time.Sleep(time.Minute * 2)

}
