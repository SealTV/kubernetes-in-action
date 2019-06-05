package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var counter int

func greet(w http.ResponseWriter, r *http.Request) {
	counter++

	if counter > 5 {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	fmt.Fprintf(w, "Hello World! %s\n", time.Now())
	hostname, _ := os.Hostname()
	fmt.Fprintf(w, "Value %s\n", hostname)
}

func main() {
	http.HandleFunc("/", greet)
	http.ListenAndServe(":8080", nil)
}
