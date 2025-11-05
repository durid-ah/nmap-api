package db

import (
	"log/slog"

	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type Host struct {
	Hostname  string `gorm:"primaryKey;not null"`
	IP        string `gorm:"unique;not null"`
}

type Storage struct {
	db *gorm.DB
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