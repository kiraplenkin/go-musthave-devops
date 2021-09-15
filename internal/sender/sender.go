package sender

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"net/http"
)

// SendClient struct of client
type SendClient struct {
}

// NewSender func to create new client for send types.Stats
func NewSender() *SendClient {
	return &SendClient{}
}

// Send func send data with sender.SendClient
func (s *SendClient) Send(client resty.Client, stats types.RequestStats, serverAddress, serverPort string) error {
	r, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	post, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bytes.NewBufferString(string(r))).
		Post(serverAddress + ":" + serverPort + types.SenderConfig.Endpoint)
	if err != nil {
		return nil
	}
	if post.StatusCode() != http.StatusCreated {
		return types.ErrCantSaveData
	}
	return nil
}
