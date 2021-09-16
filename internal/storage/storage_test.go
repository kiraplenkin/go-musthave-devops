package storage

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	store = &Store{
		Storage: types.Storage{
			1: types.Stats{
				TotalAlloc:   100,
				Sys:          100,
				Mallocs:      100,
				Frees:        100,
				LiveObjects:  100,
				PauseTotalNs: 100,
				NumGC:        100,
				NumGoroutine: 100,
			},
		},
	}
)

// TestStore_GetStatsByID test for getting types.Stats by ID from types.Storage
func TestStore_GetStatsByID(t *testing.T) {

	tests := []struct {
		name string
		args uint
		want *types.Stats
		err  error
	}{
		{
			name: "Positive test",
			args: 1,
			want: &types.Stats{
				TotalAlloc:   100,
				Sys:          100,
				Mallocs:      100,
				Frees:        100,
				LiveObjects:  100,
				PauseTotalNs: 100,
				NumGC:        100,
				NumGoroutine: 100,
			},
			err: nil,
		},
		{
			name: "Negative test",
			args: 2,
			want: nil,
			err:  types.ErrCantGetStats,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := store.GetStatsByID(tt.args)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, stats)
		})
	}
}

// TestStore_SaveStats test for saving types.Stats to Store
func TestStore_SaveStats(t *testing.T) {
	type args struct {
		ID    uint
		stats types.Stats
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Positive test",
			args: args{
				ID: 2,
				stats: types.Stats{
					TotalAlloc:   100,
					Sys:          100,
					Mallocs:      100,
					Frees:        100,
					LiveObjects:  100,
					PauseTotalNs: 100,
					NumGC:        100,
					NumGoroutine: 100,
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.SaveStats(tt.args.ID, tt.args.stats)
			assert.Equal(t, tt.err, err)
		})
	}
}

// TestStore_GetAllStats test for getting all types.Stats from types.Storage
func TestStore_GetAllStats(t *testing.T) {
	tests := []struct {
		name string
		want *types.Storage
		err  error
	}{
		{
			name: "Positive test",
			want: &store.Storage,
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testedStore, err := store.GetAllStats()
			assert.Equal(t, tt.want, testedStore)
			assert.Equal(t, tt.err, err)
			//err = os.Remove("test_file")
			//if err != nil {
			//	require.NoError(t, err)
			//}
		})
	}
}
