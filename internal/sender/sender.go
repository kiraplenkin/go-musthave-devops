package sender

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kiraplenkin/go-musthave-devops/internal/monitor"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"net/http"
)

// SendClient struct of client
type SendClient struct {
	resty   *resty.Client
	monitor *monitor.Monitor
}

// NewSender func to create new client for send types.Stats
func NewSender(resty *resty.Client, monitor *monitor.Monitor) *SendClient {
	return &SendClient{resty: resty, monitor: monitor}
}

// SendURL ...
func (s *SendClient) SendURL(serverAddress, serverPort string) error {
	//s.monitor.Mu.Lock()
	//defer s.monitor.Mu.Unlock()
	for metric, stat := range s.monitor.MonitorStorage {
		r := stat.Type + "/" + metric + "/" + fmt.Sprintf("%f", stat.Value)
		post, err := s.resty.R().
			SetHeader("Content-Type", "text/plain").
			Post(serverAddress + types.SenderConfig.Endpoint + r)
		if err != nil {
			return nil
		}
		if post.StatusCode() != http.StatusCreated {
			return types.ErrCantSaveData
		}
	}
	return nil
}

// Send ...
func (s *SendClient) Send(agentConfig types.Config) error {
	//s.monitor.Mu.Lock()
	//defer s.monitor.Mu.Unlock()
	for id, stat := range s.monitor.MonitorStorage {
		if stat.Type != "gauge" && stat.Type != "counter" {
			return types.ErrUnknownStat
		}

		requestStat := types.Metrics{}
		requestStat.ID = id
		requestStat.MType = stat.Type
		if stat.Type == "gauge" {
			requestStat.Value = &stat.Value
			if agentConfig.Key != "" {
				h := hmac.New(sha256.New, []byte(agentConfig.Key))
				h.Write([]byte(fmt.Sprintf("%s:gauge:%f", id, stat.Value)))
				requestStat.Hash = string(h.Sum(nil))
			}
		} else {
			value := int64(stat.Value)
			requestStat.Delta = &value
			if agentConfig.Key != "" {
				h := hmac.New(sha256.New, []byte(agentConfig.Key))
				h.Write([]byte(fmt.Sprintf("%s:counter:%d", id, value)))
				requestStat.Hash = string(h.Sum(nil))
			}
		}
		rawRequest, err := json.Marshal(requestStat)
		if err != nil {
			return err
		}

		post, err := s.resty.R().
			SetHeader("Content-Type", "application/json").
			SetBody(bytes.NewBufferString(string(rawRequest))).
			Post("http://" + agentConfig.ServerAddress + types.SenderConfig.Endpoint)
		if err != nil {
			return nil
		}
		if post.StatusCode() != http.StatusOK {
			return types.ErrCantSaveData
		}
	}
	return nil
}
