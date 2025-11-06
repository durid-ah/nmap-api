package db

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Host struct {
	Hostname string `gorm:"primaryKey;not null"`
	IP       string `gorm:"unique;not null"`
}

type Storage struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewStorage(logger *slog.Logger) (*Storage, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Host{})

	return &Storage{db: db, logger: logger}, nil
}

func (s *Storage) CreateHost(ctx context.Context, host *Host) error {

	return gorm.G[Host](s.db).Create(ctx, host)
}

func (s *Storage) UpdateHost(ctx context.Context, host *Host) error {
	affectedRows, err := gorm.G[Host](s.db).
		Where("hostname = ?", host.Hostname).
		Updates(ctx, *host)

	if err != nil {
		s.logger.Error("failed to update host", "error", err, "host", host)
		return err
	}
	if affectedRows == 0 {
		s.logger.Error("no host found to update", "host", host)
		return fmt.Errorf("no host found to update: %w", err)
	}
	return nil
}

func (s *Storage) GetHostIPMap(ctx context.Context) (map[string]string, error) {
	hosts, err := gorm.G[Host](s.db).Find(ctx)
	if err != nil {
		s.logger.Error("failed to get all hosts", "error", err)
		return nil, err
	}
	hostIPMap := make(map[string]string, len(hosts))
	for _, host := range hosts {
		hostIPMap[host.Hostname] = host.IP
	}
	return hostIPMap, nil
}
