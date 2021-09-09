package monitor

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"runtime"
)

// Monitor struct of Statistics
type Monitor struct{}

// NewMonitor func to create new Monitoring
func NewMonitor() *Monitor {
	return &Monitor{}
}

// Get func generate types.Stats
func (m *Monitor) Get() (types.Stats, error) {
	//var newStats types.Stats
	//newStats.StatsType = "Counter"
	//newStats.StatsValue = strconv.Itoa(runtime.NumGoroutine())
	//return newStats, nil
	var s types.Stats
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	// Number of goroutines
	s.NumGoroutine = runtime.NumGoroutine()

	// Misc mem stats
	s.Alloc = int(rtm.Alloc)
	s.TotalAlloc = int(rtm.TotalAlloc)
	s.Sys = int(rtm.Sys)
	s.Mallocs = int(rtm.Mallocs)
	s.Frees = int(rtm.Frees)
	s.LiveObjects = int(s.Mallocs - s.Frees)

	// GC stats
	s.PauseTotalNs = int(rtm.PauseTotalNs)
	s.NumGC = int(rtm.NumGC)
	return s, nil
}
