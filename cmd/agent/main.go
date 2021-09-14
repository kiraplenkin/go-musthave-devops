package main

import (
	"flag"
	"fmt"
	"github.com/go-resty/resty/v2"
	monitorService "github.com/kiraplenkin/go-musthave-devops/internal/monitor"
	sendingService "github.com/kiraplenkin/go-musthave-devops/internal/sender"
	"time"
)

var (
	RestyClient = resty.New()
)

func main() {
	serverAddress := flag.String("s", "localhost", "server address")
	serverPort := flag.String("p", "8080", "server port")
	updateFrequency := flag.Int("f", 5, "update frequency")
	flag.Parse()

	sender := sendingService.NewSender(RestyClient)
	monitor := monitorService.NewMonitor()

	for {
		ticker := time.NewTicker(time.Duration(*updateFrequency) * time.Second)
		<-ticker.C
		stats, err := monitor.Get()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = sender.Send(stats, *serverAddress, *serverPort)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
