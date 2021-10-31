package storage

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	cfg = types.Config{
		Restore:         false,
		FileStoragePath: "test.json",
	}
)

// TestNewStorage test creating Store
func TestNewStorage(t *testing.T) {
	testedStorage, err := NewStorage(&cfg)
	defer os.Remove(cfg.FileStoragePath)
	require.NoError(t, err)
	assert.IsType(t, &Store{}, testedStorage)
}

// TestUpdateGaugeStats test creating and updating gouge types.Stats
func TestGaugeStats(t *testing.T) {
	testedStorage, err := NewStorage(&cfg)
	defer os.Remove(cfg.FileStoragePath)
	require.NoError(t, err)

	type args struct {
		ID    string
		stats types.Stats
	}

	tests := []struct {
		name string
		args args
		want *types.Stats
		err  error
	}{
		{
			name: "create gouge stats",
			args: args{
				ID: "testGauge",
				stats: types.Stats{
					Type: "gauge",
					Value: 1.0,
				},
			},
			want: &types.Stats{
				Type: "gauge",
				Value: 1.0,
			},
			err: nil,
		},
		{
			name: "update gouge stats",
			args: args{
				ID: "testGauge",
				stats: types.Stats{
					Type: "gauge",
					Value: 2.0,
				},
			},
			want: &types.Stats{
				Type: "gauge",
				Value: 2.0,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testedStorage.UpdateGaugeStats(tt.args.ID, tt.args.stats)
			require.NoError(t, err)
			stats, err := testedStorage.GetGaugeStatsByID(tt.args.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.want, stats)
		})
	}
}

// TestCounterStats test creating and updating counter types.Stats
func TestCounterStats(t *testing.T) {
	testedStorage, err := NewStorage(&cfg)
	defer os.Remove(cfg.FileStoragePath)
	require.NoError(t, err)

	type args struct {
		ID    string
		stats types.Stats
	}

	tests := []struct {
		name string
		args args
		want int64
		err  error
	}{
		{
			name: "create counter stats",
			args: args{
				ID: "testCounter",
				stats: types.Stats{
					Type: "gauge",
					Value: 1.0,
				},
			},
			want: 1,
			err: nil,
		},
		{
			name: "update gouge stats",
			args: args{
				ID: "testCounter",
				stats: types.Stats{
					Type: "gauge",
					Value: 2.0,
				},
			},
			want: 3,
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testedStorage.UpdateCounterStats(tt.args.ID, tt.args.stats)
			require.NoError(t, err)
			value, err := testedStorage.GetCounterStatsByID(tt.args.ID)
			require.NoError(t, err)
			assert.Equal(t, tt.want, value)
		})
	}
}
