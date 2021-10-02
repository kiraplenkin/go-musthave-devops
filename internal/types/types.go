package types

import (
	"errors"
	"time"
)

type (
	// Config configs of app
	Config struct {
		Endpoint         string
		UpdateFrequency  time.Duration
		ReportFrequency  time.Duration
		ServerAddress    string
		ServerPort       string
		RetryCount       int
		RetryWaitTime    time.Duration
		RetryMaxWaitTime time.Duration
	}

	// ServerConfig config for server app
	ServerConfig struct {
		ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost"`
		FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"test.json"`
	}

	// Stats struct to save gauge metric
	Stats struct {
		Type  string
		Value float64
	}

	Metrics struct {
		ID    string   `json:"id"`              // Имя метрики
		MType string   `json:"type"`            // Параметр принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
	}

	// Storage struct of storage
	Storage struct {
		GougeStorage   map[string]Stats
		CounterStorage map[string]int64
	}
)

var (
	// SenderConfig config for sender service
	SenderConfig = Config{
		Endpoint:         "/update/",
		UpdateFrequency:  2,
		ReportFrequency:  10,
		ServerAddress:    "http://localhost",
		ServerPort:       "8080",
		RetryCount:       10,
		RetryWaitTime:    5,
		RetryMaxWaitTime: 30,
	}

	ErrCantGetStats = errors.New("can't get stats by ID")
	ErrCantSaveData = errors.New("sent data not saved")
	ErrUnknownStat  = errors.New("unknown stat")
)
