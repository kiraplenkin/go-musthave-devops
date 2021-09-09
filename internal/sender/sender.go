package sender

import (
	"bytes"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var ErrCantSaveData = errors.New("sent data not saved")

// SendClient struct of client
type SendClient struct {
	Client *resty.Client
}

// NewSender func to create new client for send types.Stats
func NewSender(resty *resty.Client) *SendClient {
	resty.
		SetRetryCount(types.SenderConfig.RetryCount).
		SetRetryWaitTime(types.SenderConfig.RetryWaitTime).
		SetRetryMaxWaitTime(types.SenderConfig.RetryMaxWaitTime)
	return &SendClient{
		Client: resty,
	}
}

// Send func send data with sender.SendClient
func (s *SendClient) Send(stats types.Stats) error {
	data := url.Values{}
	id := time.Now().UnixNano()
	data.Set("id", strconv.Itoa(int(id)))
	data.Set("Alloc", strconv.Itoa(stats.Alloc))
	data.Set("TotalAlloc", strconv.Itoa(stats.TotalAlloc))
	data.Set("Sys", strconv.Itoa(stats.Sys))
	data.Set("Mallocs", strconv.Itoa(stats.Mallocs))
	data.Set("Frees", strconv.Itoa(stats.Frees))
	data.Set("LiveObjects", strconv.Itoa(stats.LiveObjects))
	data.Set("PauseTotalNs", strconv.Itoa(stats.PauseTotalNs))
	data.Set("NumGC", strconv.Itoa(stats.NumGC))
	data.Set("NumGoroutine", strconv.Itoa(stats.NumGoroutine))
	post, err := s.Client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(bytes.NewBufferString(data.Encode())).
		Post(types.SenderConfig.Endpoint)
	if err != nil {
		return nil
	}
	if post.StatusCode() != http.StatusCreated {
		return ErrCantSaveData
	}
	return nil
}
