package main

import (
	"github.com/kiraplenkin/go-musthave-devops/internal/storage"
	transportHTTP "github.com/kiraplenkin/go-musthave-devops/internal/transport/http"
	"log"
	"net/http"
)

func main() {
	store := storage.NewStorage()
	handler := transportHTTP.NewHandler(*store)
	handler.SetupRouters()

	log.Fatal(http.ListenAndServe(":8080", handler.Router))
}
