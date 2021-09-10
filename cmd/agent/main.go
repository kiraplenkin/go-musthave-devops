package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	monitorService "github.com/kiraplenkin/go-musthave-devops/internal/monitor"
	sendingService "github.com/kiraplenkin/go-musthave-devops/internal/sender"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"time"
)

var RestyClient = resty.New()

func main() {
	sender := sendingService.NewSender(RestyClient)
	monitor := monitorService.NewMonitor()

	for {
		ticker := time.NewTicker(types.SenderConfig.ServerUpdateTime * time.Second)
		<-ticker.C
		stats, err := monitor.Get()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = sender.Send(stats)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
