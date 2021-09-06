package main

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	transportHTTP "github.com/kiraplenkin/go-musthave-devops/internal/transport/http"
	"github.com/kiraplenkin/go-musthave-devops/internal/types"
	"log"
	"net/http"
)

// App - the struct of app
type App struct{}

// Run - function that startup application
func (a *App) Run() error {
	fmt.Println("Setting Up App")

	store := storage.NewStorage()
	handler := transportHTTP.NewHandler(*store)
	handler.SetupRouters()

	serverCfg := types.ServerConfig{}
	err := env.Parse(&serverCfg)
	if err != nil {
		return err
	}

	log.Fatal(http.ListenAndServe(serverCfg.ServerAddress+serverCfg.ServerPort, handler.Router))
	return nil
}

func main() {
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting the App")
		fmt.Println(err)
	}
}
