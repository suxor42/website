package main

import (
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"
	"encoding/json"
	"bytes"
	"log"
)

func main() {
	port := os.Getenv("PORT")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := &http.Server{
		Addr:           ":"+port,
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		http.Handle("/", http.FileServer(http.Dir("./wwwroot")))
		http.HandleFunc("/blub", blub)
		s.ListenAndServe()
	}()

	<-signals
	fmt.Println("Stopping server")

	go func() {
		time.Sleep(11 * time.Second)
		log.Println("Server didn't stop within 10 seconds. Force stop server")
		cancel()
	}()

	err := s.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}

}


func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "Hello World!")
	fmt.Fprintln(writer, request.URL.Path)
}

func blub(writer http.ResponseWriter, request *http.Request) {
	headers := request.Header
	switch request.Method {
	case http.MethodGet:
		responseJson, err := json.Marshal(headers)
		if err != nil {
			returnStatus(writer, http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(writer, bytes.NewBuffer(responseJson))
	default:
		returnStatus(writer, http.StatusMethodNotAllowed)
	}

}

func returnStatus(w http.ResponseWriter, statuscode int) {
	w.WriteHeader(statuscode)
	fmt.Fprintln(w, http.StatusText(statuscode))
}
