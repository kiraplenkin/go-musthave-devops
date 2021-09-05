package main

import (
	"fmt"
	sendingService "github.com/kiraplenkin/go-musthave-devops/internal/sender"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"time"
)

func main() {
	sender := sendingService.NewSender()
	logger := sendingService.NewLogger(*sender)

	for {
		ticker := time.NewTicker(types.SenderConfig.ServerUpdateTime)
		<-ticker.C
		err := logger.SendStats()
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
	}
}
