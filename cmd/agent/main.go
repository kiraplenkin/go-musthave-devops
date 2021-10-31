package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/go-resty/resty/v2"
	monitorService "github.com/kiraplenkin/go-musthave-devops/internal/monitor"
	sendingService "github.com/kiraplenkin/go-musthave-devops/internal/sender"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	agentCfg := types.Config{}
	err := env.Parse(&agentCfg)
	if err != nil {
		log.Printf("can't parse env: %+v", err)
		return
	}

	flag.StringVar(&agentCfg.ServerAddress, "a", agentCfg.ServerAddress, "server address")
	flag.StringVar(&agentCfg.UpdateFrequency, "p", agentCfg.UpdateFrequency, "poll interval")
	flag.StringVar(&agentCfg.ReportFrequency, "r", agentCfg.ReportFrequency, "report interval")
	flag.StringVar(&agentCfg.Key, "k", "", "key for hash")
	flag.Parse()

	updateFrequency, err := time.ParseDuration(agentCfg.UpdateFrequency)
	if err != nil {
		log.Printf("can't parse updateFrequency: %+v", err)
		return
	}
	reportFrequency, err := time.ParseDuration(agentCfg.ReportFrequency)
	if err != nil {
		log.Printf("can't parse reportFrequency: %+v", err)
		return
	}

	restyClient := resty.New().
		SetRetryCount(types.SenderConfig.RetryCount).
		SetRetryWaitTime(types.SenderConfig.RetryWaitTime).
		SetRetryMaxWaitTime(types.SenderConfig.RetryMaxWaitTime)

	monitor := monitorService.NewMonitor()
	sender := sendingService.NewSender(restyClient, monitor)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	pollIntervalTicker := time.NewTicker(updateFrequency)
	reportIntervalTicker := time.NewTicker(reportFrequency)

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
			err := sender.Send(agentCfg)
			if err != nil {
				log.Printf("can't send metrics: %+v", err)
				return
			}
		}
	}()

	<-done
}
