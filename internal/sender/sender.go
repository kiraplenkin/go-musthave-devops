package sender

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"net/http"
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
func (s *SendClient) Send(stats types.RequestStats, serverAddress, serverPort string) error {
	r, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	post, err := s.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes.NewBufferString(string(r))).
		Post(serverAddress + ":" + serverPort + types.SenderConfig.Endpoint)
	if err != nil {
		return nil
	}
	if post.StatusCode() != http.StatusCreated {
		return ErrCantSaveData
	}
	return nil
}
