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

	// Stats struct to save stats
	Stats struct {
		TotalAlloc   int
		Sys          int
		Mallocs      int
		Frees        int
		LiveObjects  int
		NumGoroutine int
	}

	// Storage struct of storage
	Storage map[int]Stats

	// RequestStats struct to transport by JSON
	RequestStats struct {
		ID           int `json:"id"`
		TotalAlloc   int `json:"totalAlloc"`
		Sys          int `json:"sys"`
		Mallocs      int `json:"mallocs"`
		Frees        int `json:"frees"`
		LiveObjects  int `json:"liveObjects"`
		NumGoroutine int `json:"numGoroutine"`
	}
)

var (
	// SenderConfig config for sender service
	SenderConfig = Config{
		Endpoint:         "/update/",
		UpdateFrequency:  5,
		ServerAddress:    "http://localhost",
		ServerPort:       "8080",
		RetryCount:       5,
		RetryWaitTime:    10,
		RetryMaxWaitTime: 30,
	}

	ErrCantGetStats = errors.New("can't get stats by ID")
	ErrCantSaveData = errors.New("sent data not saved")
)
