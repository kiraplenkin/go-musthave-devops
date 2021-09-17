package storage

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
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
				NumGoroutine: 100,
			},
		},
	}
	cfg            = types.ServerConfig{ServerAddress: "localhost", FileStoragePath: "test_file"}
	testStorage, _ = NewStorage(&cfg)
)

// TestStore_GetStatsByID test for getting types.Stats by ID from types.Storage
func TestStore_GetStatsByID(t *testing.T) {

	tests := []struct {
		name string
		args int
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
		ID    int
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
		})
	}
}

// TestNewStorage test for creating new storage
func TestNewStorage(t *testing.T) {
	tests := []struct {
		name string
		want *Store
	}{
		{
			name: "Positive test",
			want: testStorage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newStorage, err := NewStorage(&cfg)
			require.NoError(t, err)
			assert.Equal(t, tt.want, newStorage)
		})
	}
}

func TestSaveToFile(t *testing.T) {
	type args struct {
		data     []byte
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Positive test",
			args: args{
				data: []byte(`{"id":1,"value":2"}`),
				fileName: "test_file.json",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveToFile(tt.args.data, tt.args.fileName)
			assert.Equal(t, tt.wantErr, err)
			_, err = os.Stat(tt.args.fileName)
			require.NoError(t, err)
			defer os.Remove(tt.args.fileName)
		})
	}
}