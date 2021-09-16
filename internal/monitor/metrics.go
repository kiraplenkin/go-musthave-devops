package monitor

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"runtime"
	"time"
)

// Monitor struct of Statistics
type Monitor struct{}

// NewMonitor func to create new Monitoring
func NewMonitor() *Monitor {
	return &Monitor{}
}

// Get func generate types.Stats
func (m *Monitor) Get() (types.RequestStats, error) {
	var s types.RequestStats
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)

	id := uint(time.Now().UnixNano())
	totalAlloc := uint(rtm.TotalAlloc)
	sys := uint(rtm.Sys)
	mallocs := uint(rtm.Mallocs)
	frees := uint(rtm.Frees)
	liveObjects := mallocs - frees
	pauseTotalNs := uint(rtm.PauseTotalNs)
	numGC := uint(rtm.NumGC)
	numGoroutine := uint(runtime.NumGoroutine())

	s.ID = id
	s.TotalAlloc = totalAlloc
	s.Sys = sys
	s.Mallocs = mallocs
	s.Frees = frees
	s.LiveObjects = liveObjects
	s.PauseTotalNs = pauseTotalNs
	s.NumGC = numGC
	s.NumGoroutine = numGoroutine
	return s, nil
}
