package sender

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
)

// SendingService - struct of client
type SendingService struct {
	Client *resty.Client
}

// NewSender - func to create new client for send types.Stats
func NewSender() *SendingService {
	restyClient := resty.New()
	restyClient.
		SetRetryCount(types.SenderCfg.RetryCount).
		SetRetryWaitTime(types.SenderCfg.RetryWaitTime).
		SetRetryMaxWaitTime(types.SenderCfg.RetryMaxWaitTime)
	return &SendingService{
		Client: restyClient,
	}
}

// Logger - interface for get and send types.Stats
type Logger interface {
	GetStats() (types.Stats, error)
	SendStats() error
}

// LogService - struct of LogService
type LogService struct {
	client SendingService
}

// NewLogger - func to create new LogService
func NewLogger(sender SendingService) *LogService {
	return &LogService{
		client: sender,
	}
}

// GetStats - func generate types.Stats
func (l *LogService) GetStats() (types.Stats, error) {
	var newStats types.Stats
	newStats.StatsType = "Counter"
	newStats.StatsValue = strconv.Itoa(runtime.NumGoroutine())
	return newStats, nil
}

// SendStats - func run GetStats and send data with SendingService
func (l *LogService) SendStats() error {
	newStats, err := l.GetStats()
	if err != nil {
		return errors.New("can't get stats")
	}

	requestData := types.RequestStats{
		ID:    strconv.Itoa(rand.Intn(1000)),
		Type:  newStats.StatsType,
		Value: newStats.StatsValue,
	}
	r, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	post, err := l.client.Client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(bytes.NewBufferString(string(r))).
		Post(types.SenderCfg.Endpoint)
	if err != nil {
		return err
	}
	if post.StatusCode() != http.StatusCreated {
		return errors.New("sent data not saved")
	}
	return nil
}
