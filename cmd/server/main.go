package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	transportHTTP "github.com/kiraplenkin/go-musthave-devops/internal/transport/http"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var serverPort string

func main() {
	serverCfg := types.ServerConfig{}
	err := env.Parse(&serverCfg)
	if err != nil {
		return
	}

	flag.StringVar(&serverPort, "p", "8080", "port to run server")
	flag.Parse()

	store, err := storage.NewStorage(&serverCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	handler := transportHTTP.NewHandler(*store)
	handler.SetupRouters()

	srv := &http.Server{
		Addr:    serverCfg.ServerAddress + ":" + serverPort,
		Handler: handler.Router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	log.Println("Server Started")

	<-done
	log.Println("Server Stopped")
	ctx, cancel := context.WithCancel(context.Background())
	// TODO try defer
	func() {
		data, err := json.Marshal(&store.Storage)
		if err != nil {
			log.Fatalf("can't marshal json: %+v", err)
		}
		err = storage.SaveToFile(data, serverCfg.FileStoragePath)
		if err != nil {
			log.Fatalf("can't save stats to file: %+v", err)
		}
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}
}
