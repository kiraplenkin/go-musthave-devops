package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	transportHTTP "github.com/kiraplenkin/go-musthave-devops/internal/transport/http"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

//var serverPort string

func main() {
	serverCfg := types.ServerConfig{}
	err := env.Parse(&serverCfg)
	if err != nil {
		return
	}
	storeInterval, err := strconv.Atoi(serverCfg.StoreInterval)
	if err != nil {
		return
	}

	storeIntervalTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)

	//flag.StringVar(&serverPort, "p", "8080", "port to run server")
	//flag.Parse()

	store, err := storage.NewStorage(&serverCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	handler := transportHTTP.NewHandler(*store)
	handler.SetupRouters()

	srv := &http.Server{
		Addr:    serverCfg.ServerAddress,
		Handler: handler.Router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// http server
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	// save to file
	go func() {
		for {
			<-storeIntervalTicker.C
			err := store.WriteToFile()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = store.WriteToFile()
	if err != nil {
		return
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}
}
