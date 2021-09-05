package types

import (
	"time"
)

type (
	// Config - configs of app
	Config struct {
		Endpoint         string
		ServerUpdateTime time.Duration
		RetryCount       int
		RetryWaitTime    time.Duration
		RetryMaxWaitTime time.Duration
	}

	// Stats - struct of stats
	Stats struct {
		StatsType  string
		StatsValue string
	}
)

// SenderConfig - config for sender service
var SenderConfig = Config{
	Endpoint:         "http://localhost:8080/api/stat/",
	ServerUpdateTime: 20,
	RetryCount:       5,
	RetryWaitTime:    10,
	RetryMaxWaitTime: 30,
}
