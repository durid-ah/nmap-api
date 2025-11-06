package nmapscanner

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ullaakut/nmap/v3"
	"github.com/durid-ah/nmap-api/internal/config"
)

type contextKey string

const (
	ownerContextKey contextKey = "owner"
)

type NmapScanner struct {
	scanner *nmap.Scanner
}

func NewNmapScanner(ctx context.Context, config *config.Config) (*NmapScanner, error) {
	scanner, err := nmap.NewScanner(ctx,
		nmap.WithTargets(config.NmapTarget),
		nmap.WithPingScan(),
	)
	if err != nil {
		slog.Error("unable to create nmap scanner", "error", err)
		return nil, err
	}

	return &NmapScanner{scanner: scanner}, nil
}

func (s *NmapScanner) Run(ctx context.Context) error {
	result, warnings, err := s.scanner.Run()
	if len(*warnings) > 0 {
		slog.Warn("run finished with warnings") // Warnings are non-critical errors from nmap.
		for _, warning := range *warnings {
			slog.Warn("nmap warning", "warning", warning)
		}
	}

	if err != nil {
		slog.Error("unable to run nmap scan", "error", err)
		return fmt.Errorf("unable to run nmap scan: %v", err)
	}

	// Use the results to print an example output
	for _, host := range result.Hosts {

		if len(host.Hostnames) == 0 || len(host.Addresses) == 0 {
			continue
		}

		slog.Info("Hostname", "hostname", host.Hostnames[0].Name, "ip", host.Addresses[0].Addr)
	}

	slog.Info("Nmap done", "hosts_up", len(result.Hosts), "elapsed", result.Stats.Finished.Elapsed)
	return nil
}

func CreateScannerTask(config *config.Config) func(ctx context.Context) {
	return func(_ctx context.Context) {
		scannerCtx, scannerCtxCancel := context.WithTimeout(
			context.WithValue(_ctx, ownerContextKey, "scanner"), 5*time.Minute)
		defer scannerCtxCancel()
		slog.Info("running scanner", "owner", scannerCtx.Value(ownerContextKey))
		scanner, err := NewNmapScanner(scannerCtx, config)
		if err != nil {
			slog.Error("unable to create nmap scanner", "owner", scannerCtx.Value(ownerContextKey), "error", err)
			return
		}
		err = scanner.Run(scannerCtx)
		if err != nil {
			slog.Error("unable to run nmap scan", "owner", scannerCtx.Value(ownerContextKey), "error", err)
		}
	}
}
