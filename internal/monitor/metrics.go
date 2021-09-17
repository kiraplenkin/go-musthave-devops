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

	s.ID = int(time.Now().UnixNano())
	s.TotalAlloc = int(rtm.TotalAlloc)
	s.Sys = int(rtm.Sys)
	s.Mallocs = int(rtm.Mallocs)
	s.Frees = int(rtm.Frees)
	s.LiveObjects = int(rtm.Mallocs - rtm.Frees)
	s.NumGoroutine = runtime.NumGoroutine()
	return s, nil
}
