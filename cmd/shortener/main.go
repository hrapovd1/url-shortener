package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
)

type storage map[string]string

const (
	appPort = `:8080`
	strLen  = 9
)

var (
	urlStor = make(storage, 0)
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *storage) getShort() (string, error) {
	var short string
	for i := 0; i < 3; i++ {
		shortTmp := randSeq(strLen)
		if _, exist := map[string]string(*s)[shortTmp]; exist {
			continue
		}
		short = shortTmp
	}
	if short == "" {
		return short, fmt.Errorf("error get random short")
	}
	return short, nil
}

func (s *storage) getURL(short string) (string, error) {
	url, ok := map[string]string(*s)[short]
	if !ok {
		return url, fmt.Errorf("error get url")
	}
	return url, nil
}

func (s *storage) saveURL(url string, short string) error {
	map[string]string(*s)[short] = url
	return nil
}

func rootPost(rw http.ResponseWriter, req *http.Request) {
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

	short, err := urlStor.getShort()
	if err != nil {
		log.Print(err.Error())
		http.Error(rw, "Convert error", http.StatusInternalServerError)
		return
	}

	if err := urlStor.saveURL(string(body), short); err != nil {
		log.Print(err.Error())
		http.Error(rw, "Convert error", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("content-type", "text/plain")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", short)))
}

func rootGet(rw http.ResponseWriter, req *http.Request) {
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
	url, err := urlStor.getURL(short)
	if err != nil {
		log.Print(err.Error())
		http.Error(rw, "Wrong url path", http.StatusBadRequest)
		return
	}

	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
	rw.Write([]byte(""))
}

func rootHandle(rw http.ResponseWriter, req *http.Request) {
	log.Printf("Got request:\n%+v\n", req)
	switch req.Method {
	case http.MethodGet:
		rootGet(rw, req)
	case http.MethodPost:
		rootPost(rw, req)
	default:
		http.Error(rw, "Method not work", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	logFile, err := os.OpenFile("shortener.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandle)

	if err := http.ListenAndServe(appPort, mux); err != nil {
		log.Fatal(err.Error())
	}
}
