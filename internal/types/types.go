package types

import (
	"errors"
	"time"
)

type (
	// AgentConfig configs of app
	AgentConfig struct {
		Endpoint         string
		RetryCount       int
		RetryWaitTime    time.Duration
		RetryMaxWaitTime time.Duration
	}

	// Config for server and agent apps
	Config struct {
		ServerAddress   string `env:"ADDRESS" envDefault:"localhost:8080"`
		FileStoragePath string `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
		StoreInterval   string `env:"STORE_INTERVAL" envDefault:"5m"`
		Restore         bool   `env:"RESTORE" envDefault:"true"`
		UpdateFrequency string `env:"POLL_INTERVAL" envDefault:"2s"`
		ReportFrequency string `env:"REPORT_INTERVAL" envDefault:"10s"`
		Key             string `env:"KEY"`
		Database        string `env:"DATABASE_DSN" envDefault:"localhost:5432"`
	}

	// Stats ...
	Stats struct {
		Type  string
		Value float64
	}

	// Metrics ...
	Metrics struct {
		ID    string   `json:"id"`              // Имя метрики
		MType string   `json:"type"`            // Параметр принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
		Hash  string   `json:"hash,omitempty"`  // Значение hash-функции
	}

	// Storage struct of storage
	Storage struct {
		GaugeStorage   map[string]Stats
		CounterStorage map[string]int64
	}
)

var (
	// SenderConfig config for sender service
	SenderConfig = AgentConfig{
		Endpoint:         "/update/",
		RetryCount:       10,
		RetryWaitTime:    5,
		RetryMaxWaitTime: 30,
	}

	ErrCantGetStats = errors.New("can't get stats by ID")
	ErrCantSaveData = errors.New("sent data not saved")
	ErrUnknownStat  = errors.New("unknown stat")
)
