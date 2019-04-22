package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved request from: %s", r.Host)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World! %s\n", time.Now())
	fmt.Fprintf(w, "You've hit %s\n", r.RemoteAddr)
}

func main() {
	http.HandleFunc("/", greet)

	go func() {
		log.Println("Start listen")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	cancel := make(chan os.Signal)
	signal.Notify(cancel, syscall.SIGTERM)
	signal.Notify(cancel, syscall.SIGINT)

	sig := <-cancel
	log.Printf("Stop service by signal: %v", sig)
}
