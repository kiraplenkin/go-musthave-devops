package sendingService

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogService_GetStats(t *testing.T) {
	type fields struct {
		client SendingService
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "should return stats",
			fields: fields{
				client: *NewSender(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LogService{
				client: tt.fields.client,
			}
			stats, err := l.GetStats()
			require.NoError(t, err)
			assert.IsType(t, stats, types.Stats{})
			assert.NotEqual(t, stats.StatsType, "")
			assert.NotEqual(t, stats.StatsValue, "")
		})
	}
}
