package main

import (
	"fmt"
	"github.com/kiraplenkin/go-musthave-devops/internal/stats"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	transportHTTP "github.com/kiraplenkin/go-musthave-devops/internal/transport/http"
	"net/http"
)

// App - the struct which contains ...
type App struct{}

// Run - handles the startup application
func (a *App) Run() error {
	fmt.Println("Setting Up App")

	store := storage.New()

	statsService := stats.NewService(store)

	handler := transportHTTP.NewHandler(statsService)
	handler.SetupRouters()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		fmt.Println("Failed to set up server")
		return err
	}

	return nil
}

func main() {
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting the App")
		fmt.Println(err)
	}
}
