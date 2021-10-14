package monitor

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestNewMonitor test for creating Monitor
func TestNewMonitor(t *testing.T) {
	tests := []struct {
		name string
		want *Monitor
	}{
		{
			name: "Creating new monitor storage",
			want: NewMonitor(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewMonitor())
		})
	}
}

// TestUpdateMonitor test for generating types.Stats
func TestUpdateMonitor(t *testing.T) {
	type want struct {
		statName string
		statType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Test Alloc stat",
			want: want{
				statName: "Alloc",
				statType: "gauge",
			},
		},
		{
			name: "Test BuckHashSys stat",
			want: want{
				statName: "BuckHashSys",
				statType: "gauge",
			},
		},
		{
			name: "Test Frees stat",
			want: want{
				statName: "Frees",
				statType: "gauge",
			},
		},
		{
			name: "Test GCCPUFraction stat",
			want: want{
				statName: "GCCPUFraction",
				statType: "gauge",
			},
		},
		{
			name: "Test GCSys stat",
			want: want{
				statName: "GCSys",
				statType: "gauge",
			},
		},
		{
			name: "Test HeapAlloc stat",
			want: want{
				statName: "HeapAlloc",
				statType: "gauge",
			},
		},
		{
			name: "Test HeapIdle stat",
			want: want{
				statName: "HeapIdle",
				statType: "gauge",
			},
		},
		{
			name: "Test HeapInuse stat",
			want: want{
				statName: "HeapInuse",
				statType: "gauge",
			},
		},
		{
			name: "Test HeapObjects stat",
			want: want{
				statName: "HeapObjects",
				statType: "gauge",
			},
		},
		{
			name: "Test TotalAlloc stat",
			want: want{
				statName: "TotalAlloc",
				statType: "gauge",
			},
		},
		{
			name: "Test HeapReleased stat",
			want: want{
				statName: "HeapReleased",
				statType: "gauge",
			},
		},
		{
			name: "Test HeapSys stat",
			want: want{
				statName: "HeapSys",
				statType: "gauge",
			},
		},
		{
			name: "Test LastGC stat",
			want: want{
				statName: "LastGC",
				statType: "gauge",
			},
		},
		{
			name: "Test Lookups stat",
			want: want{
				statName: "Lookups",
				statType: "gauge",
			},
		},
		{
			name: "Test HeapSys stat",
			want: want{
				statName: "HeapSys",
				statType: "gauge",
			},
		},
		{
			name: "Test MCacheInuse stat",
			want: want{
				statName: "MCacheInuse",
				statType: "gauge",
			},
		},
		{
			name: "Test MCacheSys stat",
			want: want{
				statName: "MCacheSys",
				statType: "gauge",
			},
		},
		{
			name: "Test MSpanInuse stat",
			want: want{
				statName: "MSpanInuse",
				statType: "gauge",
			},
		},
		{
			name: "Test MSpanSys stat",
			want: want{
				statName: "MSpanSys",
				statType: "gauge",
			},
		},
		{
			name: "Test Mallocs stat",
			want: want{
				statName: "Mallocs",
				statType: "gauge",
			},
		},
		{
			name: "Test NextGC stat",
			want: want{
				statName: "NextGC",
				statType: "gauge",
			},
		},
		{
			name: "Test NumForcedGC stat",
			want: want{
				statName: "NumForcedGC",
				statType: "gauge",
			},
		},
		{
			name: "Test NumGC stat",
			want: want{
				statName: "NumGC",
				statType: "gauge",
			},
		},
		{
			name: "Test OtherSys stat",
			want: want{
				statName: "OtherSys",
				statType: "gauge",
			},
		},
		{
			name: "Test PauseTotalNs stat",
			want: want{
				statName: "PauseTotalNs",
				statType: "gauge",
			},
		},
		{
			name: "Test StackInuse stat",
			want: want{
				statName: "StackInuse",
				statType: "gauge",
			},
		},
		{
			name: "Test StackSys stat",
			want: want{
				statName: "StackSys",
				statType: "gauge",
			},
		},
		{
			name: "Test Sys stat",
			want: want{
				statName: "Sys",
				statType: "gauge",
			},
		},
		{
			name: "Test RandomValue stat",
			want: want{
				statName: "RandomValue",
				statType: "gauge",
			},
		},
		{
			name: "Test PollCount stat",
			want: want{
				statName: "PollCount",
				statType: "counter",
			},
		},
	}
	monitor := NewMonitor()
	monitor.Update()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stat, ok := monitor.MonitorStorage[tt.want.statName]
			require.Equal(t, true, ok)
			assert.IsType(t, types.Stats{}, stat)
			require.NotNil(t, stat.Value)
			require.Equal(t, tt.want.statType, stat.Type)
		})
	}
}

// TestCheckCounter test increase counter
func TestCheckCounter(t *testing.T) {
	tests := []struct{
		name string
		want float64
 	}{
		{
			name: "First value",
			want: 1.0,
		},
		{
			name: "Update counter ones",
			want: 2.0,
		},
		{
			name: "Update counter twice",
			want: 3.0,
		},
	}
	monitor := NewMonitor()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor.Update()
			assert.Equal(t, tt.want, monitor.MonitorStorage["PollCount"].Value)
		})
	}
}