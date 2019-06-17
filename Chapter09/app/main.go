package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var counter = 0

func greet(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler request")

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if counter >= 5 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Some internal server error has occurred! This is pod: %s\n", hostname)
		return
	}
	counter++

	fmt.Fprintf(w, "Hello World! %v\n", time.Now())
	fmt.Fprintf(w, "This is v4 running in pod %v\n", hostname)
	w.WriteHeader(http.StatusOK)
}

func main() {
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/", greet)

	server := &http.Server{Addr: ":8080", Handler: serverMux}

	go func() {
		log.Println("Start server")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGTERM)

	<-stop
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("Stop server")
}
