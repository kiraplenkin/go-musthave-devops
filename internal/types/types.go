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

	// Stats struct of one stat
	//Stats struct {
	//	StatsType  string
	//	StatsValue string
	//}

	Stats struct {
		Alloc        int
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
)

// SenderConfig config for sender service
var SenderConfig = Config{
	Endpoint:         "http://localhost:8080/api/stat/",
	ServerUpdateTime: 5,
	RetryCount:       5,
	RetryWaitTime:    10,
	RetryMaxWaitTime: 30,
}
