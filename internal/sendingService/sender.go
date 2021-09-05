package sendingService

import (
	"bytes"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"math/rand"
	"net/http"
	"net/url"
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
		SetRetryCount(types.SenderConfig.RetryCount).
		SetRetryWaitTime(types.SenderConfig.RetryWaitTime).
		SetRetryMaxWaitTime(types.SenderConfig.RetryMaxWaitTime)
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
	data := url.Values{}
	data.Set("id", strconv.Itoa(rand.Intn(1000)))
	data.Add("type", newStats.StatsType)
	data.Add("value", newStats.StatsValue)

	post, err := l.client.Client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(bytes.NewBufferString(data.Encode())).
		Post(types.SenderConfig.Endpoint)
	if err != nil {
		return err
	}
	if post.StatusCode() != http.StatusCreated {
		return errors.New("sent data not saved")
	}
	return nil
}
