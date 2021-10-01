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
		ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"127.0.0.1"`
		FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"test.json"`
	}

	// Stats struct to save gauge metric
	Stats struct {
		Type  string
		Value float64
	}

	Metric struct {
		ID    string  `json:"id"`
		Type  string  `json:"type"`
		Value float64 `json:"value"`
	}

	// Storage struct of storage
	Storage map[string]Stats
)

var (
	// SenderConfig config for sender service
	SenderConfig = Config{
		Endpoint:         "/update/",
		UpdateFrequency:  2,
		ReportFrequency:  10,
		ServerAddress:    "http://127.0.0.1",
		ServerPort:       "8080",
		RetryCount:       10,
		RetryWaitTime:    5,
		RetryMaxWaitTime: 30,
	}

	ErrCantGetStats = errors.New("can't get stats by ID")
	ErrCantSaveData = errors.New("sent data not saved")

	Metrics = []string{
		"Alloc",
		"BuckHashSys",
		"Frees",
		"GCCPUFraction",
		"GCSys",
		"HeapAlloc",
		"HeapIdle",
		"HeapInuse",
		"HeapObjects",
		"HeapReleased",
		"HeapSys",
		"LastGC",
		"Lookups",
		"MCacheInuse",
		"MCacheSys",
		"MSpanInuse",
		"MSpanSys",
		"Mallocs",
		"NextGC",
		"NumForcedGC",
		"NumGC",
		"OtherSys",
		"PauseTotalNs",
		"StackInuse",
		"StackSys",
		"Sys",
		"PollCount",
		"RandomValue",
	}
)
