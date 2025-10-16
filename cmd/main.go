package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "log/slog"
    "time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/go-co-op/gocron/v2"
    "github.com/durid-ah/nmap-api/config"
)


func initJob() gocron.Scheduler {
	// create a scheduler
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		slog.Error("error creating scheduler", "error", err)
		log.Fatal(err)
	}
	
	// add a job to the scheduler
	j, err := scheduler.NewJob(
		gocron.CronJob("* */1 * * * *", true),
		gocron.NewTask(
			func(_ctx context.Context) {
				scanner, err := initScanner(_ctx)
				if err != nil {
					slog.Error("error initializing scanner", "error", err)
					return
				}
				runScanner(scanner)
			},
		),
		gocron.WithContext(ctx),
	)

	if err != nil {
		slog.Error("error creating job", "error", err)
		log.Fatal(err)
	}

	// each job has a unique id
	slog.Info("job created", "job_id", j.ID())

	// start the scheduler
	scheduler.Start()

	return scheduler
}

func initScanner(ctx context.Context) (*nmap.Scanner, error) {
	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets("192.168.1.*"),
		nmap.WithPingScan(),
	)

	if err != nil {
		slog.Error("unable to create nmap scanner", "error", err)
		return nil, err
	}

	return scanner, err
}

func runScanner(scanner *nmap.Scanner) {
	result, warnings, err := scanner.Run()
	if len(*warnings) > 0 {
		log.Printf("run finished with warnings: %s\n", *warnings) // Warnings are non-critical errors from nmap.
	}
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	// Use the results to print an example output
	for _, host := range result.Hosts {

		if len(host.Hostnames) == 0 || len(host.Addresses) == 0 {
			continue
		}

		slog.Info("Hostname", "hostname", host.Hostnames[0], "ip", host.Addresses[0])
	}

	slog.Info("Nmap done", "hosts_up", len(result.Hosts), "elapsed", result.Stats.Finished.Elapsed)
}

func main() {
	cfg := config.NewConfig()
	fmt.Printf("config: %+v", cfg)
	// TODO: Config
	// TODO: Child logger
    opts := slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	// TODO: pass in individual logger?
	handler := slog.NewJSONHandler(os.Stdout, &opts)
	slog.SetDefault(slog.New(handler))


	scheduler := initJob()
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
