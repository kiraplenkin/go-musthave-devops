package monitor

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"runtime"
	"strconv"
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
		Value: strconv.FormatUint(rtm.Alloc, 10),
	})
	//s = append(s, types.Metric{
	//	ID:    "Sys",
	//	Type:  "gauge",
	//	Value: strconv.FormatUint(rtm.Sys, 10),
	//})
	//s = append(s, types.Metric{
	//	ID:    "GCSys",
	//	Type:  "gauge",
	//	Value: strconv.FormatUint(rtm.GCSys, 10),
	//})
	//s = append(s, types.Metric{
	//	ID:    "OtherSys",
	//	Type:  "gauge",
	//	Value: strconv.FormatUint(rtm.OtherSys, 10),
	//})
	//s = append(s, types.Metric{
	//	ID:    "Mallocs",
	//	Type:  "gauge",
	//	Value: strconv.FormatUint(rtm.Mallocs, 10),
	//})
	//s = append(s, types.Metric{
	//	ID:    "Frees",
	//	Type:  "counter",
	//	Value: strconv.FormatUint(rtm.Frees, 10),
	//})
	//s = append(s, types.Metric{
	//	ID:    "HeapObjects",
	//	Type:  "counter",
	//	Value: strconv.FormatUint(rtm.HeapObjects, 10),
	//})
	//s = append(s, types.Metric{
	//	ID:    "LiveObjects",
	//	Type:  "counter",
	//	Value: strconv.FormatUint(rtm.Mallocs-rtm.Frees, 10),
	//})
	//s = append(s, types.Metric{
	//	ID:    "LiveObjects",
	//	Type:  "counter",
	//	Value: strconv.FormatUint(rtm.PauseTotalNs, 10),
	//})
	return s, nil
}
