package types

import (
	"time"
)

type (
	// AgentConfig - config for agent app
	AgentConfig struct {
		Endpoint         string
		ServerUpdateTime time.Duration
		RetryCount       int
		RetryWaitTime    time.Duration
		RetryMaxWaitTime time.Duration
	}

	// ServerConfig - config for server app
	ServerConfig struct {
		ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost"`
		ServerPort      string `env:"SERVER_PORT" envDefault:":8080"`
		FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/"`
	}

	// Stats - struct of stats
	Stats struct {
		StatsType  string `json:"stats_type"`
		StatsValue string `json:"stats_value"`
	}

	Storage map[uint]Stats

	// RequestStats - struct of stats that transport on json
	RequestStats struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Value string `json:"value"`
	}
)

// SenderCfg - config for sender service
var SenderCfg = AgentConfig{
	Endpoint:         "http://localhost:8080/api/stat/",
	ServerUpdateTime: 20,
	RetryCount:       5,
	RetryWaitTime:    10,
	RetryMaxWaitTime: 30,
}
