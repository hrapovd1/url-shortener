package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hrapovd1/url-shortener/internal/app/handlers"
	"github.com/hrapovd1/url-shortener/internal/app/storage"
)

const appPort = `:8080`

func main() {
	logFile, err := os.OpenFile("shortener.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	urlStorage := storage.NewMemStorage()
	rootHandler := handlers.RootHandler(urlStorage)

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	if err := http.ListenAndServe(appPort, mux); err != nil {
		log.Fatal(err.Error())
	}
}
