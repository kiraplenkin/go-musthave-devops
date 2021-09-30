package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	monitorService "github.com/kiraplenkin/go-musthave-devops/internal/monitor"
	sendingService "github.com/kiraplenkin/go-musthave-devops/internal/sender"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"time"
)

var (
	updateFrequency           int
	serverAddress, serverPort string
)

func main() {
	flag.StringVar(&serverAddress, "s", "", "server address")
	flag.StringVar(&serverPort, "p", "", "server port")
	flag.IntVar(&updateFrequency, "f", 0, "update frequency")
	flag.Parse()
	if updateFrequency != 0 {
		types.SenderConfig.UpdateFrequency = time.Duration(updateFrequency)
	}
	if serverAddress != "" {
		types.SenderConfig.ServerAddress = serverAddress
	}
	if serverPort != "" {
		types.SenderConfig.ServerPort = serverPort
	}

	restyClient := resty.New().
		SetRetryCount(types.SenderConfig.RetryCount).
		SetRetryWaitTime(types.SenderConfig.RetryWaitTime).
		SetRetryMaxWaitTime(types.SenderConfig.RetryMaxWaitTime)

	sender := sendingService.NewSender(restyClient)
	monitor := monitorService.NewMonitor()

	for {
		ticker := time.NewTicker(types.SenderConfig.UpdateFrequency * time.Second)
		<-ticker.C
		stats, err := monitor.Get()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, stat := range stats {
			fmt.Println(stat)
			err = sender.Send(stat, types.SenderConfig.ServerAddress, types.SenderConfig.ServerPort)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	}
}
