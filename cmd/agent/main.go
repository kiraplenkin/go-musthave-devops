package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	monitorService "github.com/kiraplenkin/go-musthave-devops/internal/monitor"
	sendingService "github.com/kiraplenkin/go-musthave-devops/internal/sender"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//var (
//	updateFrequency           int
//	serverAddress, serverPort string
//)

func main() {
	agentCfg := types.AgentConfig{}
	err := env.Parse(&agentCfg)
	if err != nil {
		return
	}
	updateFrequency, err := strconv.Atoi(agentCfg.UpdateFrequency)
	if err != nil {
		return
	}
	reportFrequency, err := strconv.Atoi(agentCfg.ReportFrequency)
	if err != nil {
		fmt.Println(err)
		return
	}
	//flag.StringVar(&serverAddress, "s", "", "server address")
	//flag.StringVar(&serverPort, "p", "", "server port")
	//flag.IntVar(&updateFrequency, "f", 0, "update frequency")
	//flag.Parse()
	//if updateFrequency != 0 {
	//	types.SenderConfig.UpdateFrequency = time.Duration(updateFrequency)
	//}
	//if serverAddress != "" {
	//	types.SenderConfig.ServerAddress = serverAddress
	//}
	//if serverPort != "" {
	//	types.SenderConfig.ServerPort = serverPort
	//}

	restyClient := resty.New().
		SetRetryCount(types.SenderConfig.RetryCount).
		SetRetryWaitTime(types.SenderConfig.RetryWaitTime).
		SetRetryMaxWaitTime(types.SenderConfig.RetryMaxWaitTime)

	monitor := monitorService.NewMonitor()
	sender := sendingService.NewSender(restyClient, monitor)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	pollIntervalTicker := time.NewTicker(time.Duration(updateFrequency) * time.Second)
	reportIntervalTicker := time.NewTicker(time.Duration(reportFrequency) * time.Second)

	// update metrics
	go func() {
		for {
			<-pollIntervalTicker.C
			monitor.Update()
		}
	}()

	// report metrics
	go func() {
		for {
			<-reportIntervalTicker.C
			err := sender.Send(agentCfg.ServerAddress)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	<-done
}
