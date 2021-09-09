package sender

import (
	"github.com/go-resty/resty/v2"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// SendClientMock mock-struct of SendClient
type SendClientMock struct {
	mock.Mock
}

// Send mock-func for send types.Stats
func (s *SendClientMock) Send(stats types.Stats) error {
	args := s.Called(stats)
	return args.Error(0)
}

// TestSender_Send mock-test for sending types.Stats to server
func TestSender_Send(t *testing.T) {
	mockSender := new(SendClientMock)
	testStats := types.Stats{
		Alloc:        100,
		TotalAlloc:   100,
		Sys:          100,
		Mallocs:      100,
		Frees:        100,
		LiveObjects:  100,
		PauseTotalNs: 100,
		NumGC:        100,
		NumGoroutine: 100,
	}

	mockSender.On("Send", testStats).Return(nil)

	err := mockSender.Send(testStats)
	require.NoError(t, err)
	mockSender.AssertExpectations(t)
}

// TestNewSender test for create SendClient
func TestNewSender(t *testing.T) {
	restyClient := resty.New()
	tests := []struct {
		name string
		want *SendClient
	}{
		{
			name: "Positive test",
			want: NewSender(restyClient),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewSender(restyClient))
		})
	}
}
