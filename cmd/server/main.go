package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/kiraplenkin/go-musthave-devops/internal/compress"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	transportHTTP "github.com/kiraplenkin/go-musthave-devops/internal/transport/http"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	serverCfg := types.Config{}

	err := env.Parse(&serverCfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	flag.StringVar(&serverCfg.ServerAddress, "a", serverCfg.ServerAddress, "server address")
	flag.BoolVar(&serverCfg.Restore, "r", serverCfg.Restore, "restore storage")
	flag.StringVar(&serverCfg.StoreInterval, "i", serverCfg.StoreInterval, "store interval")
	flag.StringVar(&serverCfg.FileStoragePath, "f", serverCfg.FileStoragePath, "file storage")
	flag.StringVar(&serverCfg.Key, "k", "", "key for hash")
	flag.StringVar(&serverCfg.Database, "d", serverCfg.Database, "database connection string")
	flag.Parse()

	storeInterval, err := time.ParseDuration(serverCfg.StoreInterval)
	if err != nil {
		fmt.Println(err)
		return
	}

	storeIntervalTicker := time.NewTicker(storeInterval)

	store, err := storage.NewStorage(&serverCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	handler := transportHTTP.NewHandler(store, serverCfg)
	handler.SetupRouters()

	srv := &http.Server{
		Addr:    serverCfg.ServerAddress,
		Handler: compress.GzipHandle(handler.Router),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// http server
	go func() {
		log.Println("Server started")
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = store.WriteToFile()
	if err != nil {
		return
	}
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}
}
