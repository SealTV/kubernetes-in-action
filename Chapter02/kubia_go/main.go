package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved request from: %s", r.RemoteAddr)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World! %s\n", time.Now())
	fmt.Fprintf(w, "You've hit %s\n", r.Host)
}

func main() {
	http.HandleFunc("/", greet)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
