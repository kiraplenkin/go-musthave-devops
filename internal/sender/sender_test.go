package sender

import (
	"github.com/go-resty/resty/v2"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var restyClient = resty.New()

// SendClientMock mock-struct of SendClient
type SendClientMock struct {
	mock.Mock
}

// Send mock-func for send types.Stats
func (s *SendClientMock) Send(stats types.RequestStats) error {
	args := s.Called(stats)
	return args.Error(0)
}

// TestSender_Send mock-test for sending types.Stats to server
func TestSender_Send(t *testing.T) {
	mockSender := new(SendClientMock)
	testStats := types.RequestStats{
		ID:           1,
		TotalAlloc:   101,
		Sys:          102,
		Mallocs:      103,
		Frees:        104,
		LiveObjects:  105,
		NumGoroutine: 108,
	}

	mockSender.On("Send", testStats).Return(nil)

	err := mockSender.Send(testStats)
	require.NoError(t, err)
	mockSender.AssertExpectations(t)
}

// TestNewSender test for create SendClient
func TestNewSender(t *testing.T) {
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
