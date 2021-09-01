package main

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetStats(t *testing.T) {
	stats := getStats()
	t.Run("should return Stats", func(t *testing.T) {
		assert.IsType(t, stats, storage.Stats{})
	})
	t.Run("stats must contain StatsType", func(t *testing.T) {
		assert.NotEqual(t, stats.StatsType, "")
	})
	t.Run("stats must contain StatsValue", func(t *testing.T) {
		assert.NotEqual(t, stats.StatsValue, "")
	})
}

type SendServiceMock struct {
	mock.Mock
}

func (s *SendServiceMock) SendStats(stats storage.Stats) error {
	args := s.Called(stats)

	return args.Error(0)
}

func TestSaveStats(t *testing.T) {
	logService := new(SendServiceMock)

	stats := storage.Stats{
		StatsType:  "Test",
		StatsValue: "Test",
	}

	logService.On("SendStats", stats).Return(nil)

	myService := MyService{logService}
	err := myService.SaveStats(stats)
	if err != nil {
		// TODO return error
		return 
	}

	logService.AssertExpectations(t)
}