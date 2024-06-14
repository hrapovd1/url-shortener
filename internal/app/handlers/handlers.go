package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/hrapovd1/url-shortener/internal/app/storage"
)

func rootPost(s storage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			log.Printf("got url: %s\n", req.URL.Path)
			http.Error(rw, "Wrong url path", http.StatusBadRequest)
			return
		}

		if !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
			log.Printf("got header: %s\n", req.Header.Get("Content-Type"))
			http.Error(rw, "Wrong content type.", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Printf("Got error, when read body: %s\n", err.Error())
		}

		short, err := s.GetShort()
		if err != nil {
			log.Print(err.Error())
			http.Error(rw, "Convert error", http.StatusInternalServerError)
			return
		}

		if err := s.SaveURL(string(body), short); err != nil {
			log.Print(err.Error())
			http.Error(rw, "Convert error", http.StatusInternalServerError)
			return
		}

		rw.Header().Set("content-type", "text/plain")
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", short)))
	}
}

func rootGet(s storage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			log.Printf("got url: %s\n", req.URL.Path)
			http.Error(rw, "Wrong url path", http.StatusBadRequest)
			return
		}

		/*
			if !strings.Contains(req.Header.Get("Content-Type"), "text/plain") {
				log.Printf("got header: %s\n", req.Header.Get("Content-Type"))
				http.Error(rw, "Wrong content type.", http.StatusBadRequest)
				return
			}
		*/

		short := req.URL.Path[1:]
		url, err := s.GetURL(short)
		if err != nil {
			log.Print(err.Error())
			http.Error(rw, "Wrong url path", http.StatusBadRequest)
			return
		}

		rw.Header().Set("Location", url)
		rw.WriteHeader(http.StatusTemporaryRedirect)
		rw.Write([]byte(""))
	}
}

func RootHandler(s storage.Storage) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		log.Printf("Got request:\n%+v\n", req)
		switch req.Method {
		case http.MethodGet:
			rootGet(s)(rw, req)
		case http.MethodPost:
			rootPost(s)(rw, req)
		default:
			http.Error(rw, "Method not work", http.StatusMethodNotAllowed)
			return
		}
	}
}
