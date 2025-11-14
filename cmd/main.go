package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/durid-ah/nmap-api/config"
	"github.com/durid-ah/nmap-api/cron_scheduler"
	"github.com/durid-ah/nmap-api/db"
	"github.com/durid-ah/nmap-api/nmap_scanner"
)

func main() {
	cfg := config.NewConfig()

	_, err := db.NewStorage(slog.Default())
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Child logger
	opts := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stdout, &opts)
	slog.SetDefault(slog.New(handler))

	// run the scanner once at startup to populate the db
	slog.Info("running intial scan to populate the db...")
	scanTask := nmapscanner.CreateScannerTask(cfg)
	scanTask(context.Background())
	slog.Info("intial scan completed")

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
