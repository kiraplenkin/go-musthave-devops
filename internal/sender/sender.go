package sender

import (
	"bytes"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/kiraplenkin/go-musthave-devops/internal/transformation"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"net/http"
)

// SendClient struct of client
type SendClient struct {
	resty *resty.Client
}

// NewSender func to create new client for send types.Stats
func NewSender(resty *resty.Client) *SendClient {
	return &SendClient{resty: resty}
}

// Send func send data with sender.SendClient
func (s *SendClient) Send(stats types.RequestStats, serverAddress, serverPort string) error {
	rawRequest, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	compressRequest, err := transformation.Compress(rawRequest)
	if err != nil {
		return err
	}

	encodedRequest, err := transformation.EncodeDecode(compressRequest, "encode")
	if err != nil {
		return err
	}

	post, err := s.resty.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(bytes.NewBufferString(string(encodedRequest))).
		Post(serverAddress + ":" + serverPort + types.SenderConfig.Endpoint)
	if err != nil {
		return nil
	}
	if post.StatusCode() != http.StatusCreated {
		return types.ErrCantSaveData
	}
	return nil
}
