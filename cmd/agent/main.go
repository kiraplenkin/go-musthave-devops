package main

import (
	"bytes"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	"math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

const endpoint = "http://localhost:8080/api/stat/"
const updateTime = 2

var statsTypes = [2]string{"Counter", "Gauge"}

// SendingService ...
type SendingService interface {
	SendStats(storage.Stats) error
}

// SendService ...
type SendService struct{}

// MyService ...
type MyService struct {
	sendingService SendingService
}

// SendStats ...
func (s SendService) SendStats(stats storage.Stats) error {
	client := resty.New()

	data := url.Values{}
	data.Set("id", strconv.Itoa(rand.Intn(1000)))
	data.Add("type", stats.StatsType)
	data.Add("value", stats.StatsValue)

	_, err := client.
		SetRetryCount(5).
		SetRetryWaitTime(5*time.Second).
		SetRetryMaxWaitTime(20*time.Second).
		R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(bytes.NewBufferString(data.Encode())).
		Post(endpoint)
	if err != nil {
		return err
	}
	return nil
}

// SaveStats ...
func (a MyService) SaveStats(stats storage.Stats) error {
	err := a.sendingService.SendStats(stats)
	if err != nil {
		return err
	}
	return nil
}

// getStats ...
func getStats() storage.Stats {
	var stats storage.Stats
	n := rand.Int() % len(statsTypes)
	stats.StatsType = statsTypes[n]
	stats.StatsValue = strconv.Itoa(rand.Intn(100000))
	return stats
}

func main() {
	logService := SendService{}
	myService := MyService{logService}
	for {
		ticker := time.NewTicker(updateTime * time.Second)
		<-ticker.C
		err := myService.SaveStats(getStats())
		if err != nil {
			return
		}
	}
}
