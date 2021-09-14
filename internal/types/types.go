package types

import (
	"time"
)

type (
	// Config configs of app
	Config struct {
		Endpoint         string
		ServerUpdateTime time.Duration
		RetryCount       int
		RetryWaitTime    time.Duration
		RetryMaxWaitTime time.Duration
	}

	// ServerConfig - config for server app
	ServerConfig struct {
		ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost"`
		FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"test.json"`
	}

	Stats struct {
		TotalAlloc   int
		Sys          int
		Mallocs      int
		Frees        int
		LiveObjects  int
		PauseTotalNs int
		NumGC        int
		NumGoroutine int
	}

	// Storage struct of storage
	Storage map[uint]Stats

	// RequestStats struct to transport by JSON
	RequestStats struct {
		ID           uint `json:"id,omitempty"`
		TotalAlloc   uint `json:"totalAlloc,omitempty"`
		Sys          uint `json:"sys,omitempty"`
		Mallocs      uint `json:"mallocs,omitempty"`
		Frees        uint `json:"frees,omitempty"`
		LiveObjects  uint `json:"liveObjects,omitempty"`
		PauseTotalNs uint `json:"pauseTotalNs,omitempty"`
		NumGC        uint `json:"numGC,omitempty"`
		NumGoroutine uint `json:"numGoroutine,omitempty"`
	}
)

// SenderConfig config for sender service
var SenderConfig = Config{
	Endpoint:         "/update/",
	RetryCount:       5,
	RetryWaitTime:    10,
	RetryMaxWaitTime: 30,
}
