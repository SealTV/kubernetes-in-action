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

func greet(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler request")
	fmt.Fprintf(w, "Hello World! %v\n", time.Now())
	fmt.Fprintf(w, "Host: %v\n", r.Host)
	fmt.Fprintf(w, "Method: %v\n", r.Method)
	fmt.Fprintf(w, "RequestURI: %v\n", r.RequestURI)
	fmt.Fprintf(w, "RemoteAddr: %v\n", r.RemoteAddr)
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
	}
	fmt.Fprintf(w, "OS Hostname: %v\n", hostname)
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
