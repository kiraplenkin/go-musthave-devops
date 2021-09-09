package monitor

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewMonitor test for creating Monitor
func TestNewMonitor(t *testing.T) {
	tests := []struct {
		name string
		want *Monitor
	}{
		{
			name: "Positive test",
			want: NewMonitor(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewMonitor())
		})
	}
}

// TestMonitor_Get test for generate types.Stats
func TestMonitor_Get(t *testing.T) {
	type want struct {
		stats types.Stats
		error error
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Positive test",
			want: want{
				stats: types.Stats{},
				error: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monitor := NewMonitor()
			stats, err := monitor.Get()
			assert.Equal(t, tt.want.error, err)
			assert.IsType(t, tt.want.stats, stats)
		})
	}
}
