package types

import (
	"errors"
	"time"
)

type (
	// Config configs of app
	Config struct {
		Endpoint         string
		RetryCount       int
		RetryWaitTime    time.Duration
		RetryMaxWaitTime time.Duration
	}

	// ServerConfig config for server app
	ServerConfig struct {
		ServerAddress   string `env:"ADDRESS" envDefault:"localhost:8080"`
		FileStoragePath string `env:"STORE_FILE" envDefault:"test.json"`
		StoreInterval string `env:"STORE_INTERVAL" envDefault:"300"`
		Restore bool `env:"RESTORE" envDefault:"true"`
	}

	// AgentConfig ...
	AgentConfig struct {
		ServerAddress   string `env:"ADDRESS" envDefault:"localhost:8080"`
		UpdateFrequency string `env:"POLL_INTERVAL" envDefault:"2"`
		ReportFrequency string `env:"REPORT_INTERVAL" envDefault:"10"`
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
	}

	// Storage struct of storage
	Storage struct {
		GaugeStorage   map[string]Stats
		CounterStorage map[string]int64
	}
)

var (
	// SenderConfig config for sender service
	SenderConfig = Config{
		Endpoint:         "/update/",
		RetryCount:       10,
		RetryWaitTime:    5,
		RetryMaxWaitTime: 30,
	}

	ErrCantGetStats = errors.New("can't get stats by ID")
	ErrCantSaveData = errors.New("sent data not saved")
	ErrUnknownStat  = errors.New("unknown stat")
)
