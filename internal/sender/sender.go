package sender

import (
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

// Send ...
func (s *SendClient) Send(serverAddress, serverPort string) error {
	//s.monitor.Mu.Lock()
	//defer s.monitor.Mu.Unlock()
	for metric, stat := range s.monitor.MonitorStorage {
		r := stat.Type + "/" + metric + "/" + fmt.Sprintf("%f", stat.Value)
		post, err := s.resty.R().
			SetHeader("Content-Type", "text/plain").
			Post(serverAddress + ":" + serverPort + types.SenderConfig.Endpoint + r)
		if err != nil {
			return nil
		}
		if post.StatusCode() != http.StatusCreated {
			return types.ErrCantSaveData
		}
	}
	return nil
}
