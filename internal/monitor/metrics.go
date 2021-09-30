package monitor

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"math/rand"
	"runtime"
)

// Monitor struct of Statistics
type Monitor struct{}

// NewMonitor func to create new Monitoring
func NewMonitor() *Monitor {
	return &Monitor{}
}

// Get func generate types.Stats
func (m *Monitor) Get() ([]types.Metric, error) {
	var s []types.Metric
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	s = append(s, types.Metric{
		ID:    "Alloc",
		Type:  "gauge",
		Value: float64(rtm.Alloc),
	})
	s = append(s, types.Metric{
		ID:    "BuckHashSys",
		Type:  "gauge",
		Value: float64(rtm.BuckHashSys),
	})
	s = append(s, types.Metric{
		ID:    "Frees",
		Type:  "gauge",
		Value: float64(rtm.Frees),
	})
	s = append(s, types.Metric{
		ID:    "GCCPUFraction",
		Type:  "gauge",
		Value: rtm.GCCPUFraction,
	})
	s = append(s, types.Metric{
		ID:    "GCSys",
		Type:  "gauge",
		Value: float64(rtm.GCSys),
	})
	s = append(s, types.Metric{
		ID:    "HeapAlloc",
		Type:  "gauge",
		Value: float64(rtm.HeapAlloc),
	})
	s = append(s, types.Metric{
		ID:    "HeapIdle",
		Type:  "gauge",
		Value: float64(rtm.HeapIdle),
	})
	s = append(s, types.Metric{
		ID:    "HeapInuse",
		Type:  "gauge",
		Value: float64(rtm.HeapInuse),
	})
	s = append(s, types.Metric{
		ID:    "HeapObjects",
		Type:  "gauge",
		Value: float64(rtm.HeapObjects),
	})
	s = append(s, types.Metric{
		ID:    "HeapReleased",
		Type:  "gauge",
		Value: float64(rtm.HeapReleased),
	})
	s = append(s, types.Metric{
		ID:    "HeapSys",
		Type:  "gauge",
		Value: float64(rtm.HeapSys),
	})
	s = append(s, types.Metric{
		ID:    "LastGC",
		Type:  "gauge",
		Value: float64(rtm.LastGC),
	})
	s = append(s, types.Metric{
		ID:    "Lookups",
		Type:  "gauge",
		Value: float64(rtm.Lookups),
	})
	s = append(s, types.Metric{
		ID:    "HeapSys",
		Type:  "gauge",
		Value: float64(rtm.HeapSys),
	})
	s = append(s, types.Metric{
		ID:    "MCacheInuse",
		Type:  "gauge",
		Value: float64(rtm.MCacheInuse),
	})
	s = append(s, types.Metric{
		ID:    "MCacheSys",
		Type:  "gauge",
		Value: float64(rtm.MCacheSys),
	})
	s = append(s, types.Metric{
		ID:    "MSpanInuse",
		Type:  "counter",
		Value: float64(rtm.MSpanInuse),
	})
	s = append(s, types.Metric{
		ID:    "MSpanSys",
		Type:  "gauge",
		Value: float64(rtm.MSpanSys),
	})
	s = append(s, types.Metric{
		ID:    "Mallocs",
		Type:  "gauge",
		Value: float64(rtm.Mallocs),
	})
	s = append(s, types.Metric{
		ID:    "NextGC",
		Type:  "gauge",
		Value: float64(rtm.NextGC),
	})
	s = append(s, types.Metric{
		ID:    "NumForcedGC",
		Type:  "gauge",
		Value: float64(rtm.NumForcedGC),
	})
	s = append(s, types.Metric{
		ID:    "NumGCv",
		Type:  "gauge",
		Value: float64(rtm.NumGC),
	})
	s = append(s, types.Metric{
		ID:    "OtherSys",
		Type:  "gauge",
		Value: float64(rtm.OtherSys),
	})
	s = append(s, types.Metric{
		ID:    "PauseTotalNs",
		Type:  "gauge",
		Value: float64(rtm.PauseTotalNs),
	})
	s = append(s, types.Metric{
		ID:    "StackInuse",
		Type:  "gauge",
		Value: float64(rtm.StackInuse),
	})
	s = append(s, types.Metric{
		ID:    "StackSys",
		Type:  "gauge",
		Value: float64(rtm.StackSys),
	})
	s = append(s, types.Metric{
		ID:    "Sys",
		Type:  "gauge",
		Value: float64(rtm.Sys),
	})
	s = append(s, types.Metric{
		ID:    "PollCount",
		Type:  "counter",
		Value: 1,
	})
	s = append(s, types.Metric{
		ID:    "RandomValue",
		Type:  "gauge",
		Value: float64(rand.Intn(10000)),
	})
	return s, nil
}
