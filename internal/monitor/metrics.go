package monitor

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"math/rand"
	"runtime"
	"sync"
)

// Monitor struct of Statistics
type Monitor struct {
	Mu sync.Mutex
	MonitorStorage map[string]types.Stats
}

// NewMonitor func to create new Monitoring
func NewMonitor() *Monitor {
	m := make(map[string]types.Stats)
	m["PollCount"] = types.Stats{
		Type:  "counter",
		Value: 0.0,
	}
	return &Monitor{MonitorStorage: m}
}

// Update ...
func (m *Monitor) Update() {
	//m.Mu.Lock()
	//defer m.Mu.Unlock()

	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	m.MonitorStorage["Alloc"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.Alloc),
	}
	m.MonitorStorage["BuckHashSys"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.BuckHashSys),
	}
	m.MonitorStorage["Frees"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.Frees),
	}
	m.MonitorStorage["GCCPUFraction"] = types.Stats{
		Type:  "gauge",
		Value: rtm.GCCPUFraction,
	}
	m.MonitorStorage["GCSys"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.GCSys),
	}
	m.MonitorStorage["HeapAlloc"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.HeapAlloc),
	}
	m.MonitorStorage["HeapIdle"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.HeapIdle),
	}
	m.MonitorStorage["HeapInuse"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.HeapInuse),
	}
	m.MonitorStorage["HeapObjects"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.HeapObjects),
	}
	m.MonitorStorage["HeapReleased"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.HeapReleased),
	}
	m.MonitorStorage["HeapSys"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.HeapSys),
	}
	m.MonitorStorage["LastGC"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.LastGC),
	}
	m.MonitorStorage["Lookups"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.Lookups),
	}
	m.MonitorStorage["HeapSys"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.HeapSys),
	}
	m.MonitorStorage["MCacheInuse"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.MCacheInuse),
	}
	m.MonitorStorage["MCacheSys"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.MCacheSys),
	}
	m.MonitorStorage["MSpanInuse"] = types.Stats{
		Type:  "counter",
		Value: float64(rtm.MSpanInuse),
	}
	m.MonitorStorage["MSpanSys"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.MSpanSys),
	}
	m.MonitorStorage["Mallocs"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.Mallocs),
	}
	m.MonitorStorage["NextGC"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.NextGC),
	}
	m.MonitorStorage["NumForcedGC"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.NumForcedGC),
	}
	m.MonitorStorage["NumGCv"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.NumGC),
	}
	m.MonitorStorage["OtherSys"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.OtherSys),
	}
	m.MonitorStorage["PauseTotalNs"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.PauseTotalNs),
	}
	m.MonitorStorage["StackInuse"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.StackInuse),
	}
	m.MonitorStorage["StackSys"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.StackSys),
	}
	m.MonitorStorage["Sys"] = types.Stats{
		Type:  "gauge",
		Value: float64(rtm.Sys),
	}
	m.MonitorStorage["RandomValue"] = types.Stats{
		Type:  "gauge",
		Value: float64(rand.Intn(10000)),
	}
	if count, ok := m.MonitorStorage["PollCount"]; ok {
		count.Value++
		m.MonitorStorage["PollCount"] = count
	}
}
