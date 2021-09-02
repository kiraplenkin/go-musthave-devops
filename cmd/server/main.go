package main

import (
	"fmt"
	"github.com/kiraplenkin/go-musthave-devops/internal/stats"
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	transportHTTP "github.com/kiraplenkin/go-musthave-devops/internal/transport/http"
	"log"
	"net/http"
)

// App - the struct of app
type App struct{}

// Run - handles the startup application
func (a *App) Run() error {
	fmt.Println("Setting Up App")

	store := storage.New()

	statsService := stats.NewService(store)

	handler := transportHTTP.NewHandler(statsService)
	handler.SetupRouters()

	log.Fatal(http.ListenAndServe(":8080", handler.Router))

	return nil
}

func main() {
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting the App")
		fmt.Println(err)
	}
}
