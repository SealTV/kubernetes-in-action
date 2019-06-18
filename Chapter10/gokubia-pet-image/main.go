package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const dataFile = "/var/data/gokubia.txt"

func greet(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler request")

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Some error on read request body: %v", err)
			return
		}
		defer r.Body.Close()

		file, err := os.OpenFile(dataFile, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Some error on open file: %v", err)
			return
		}
		defer file.Close()

		if _, err := fmt.Fprintln(file, string(data)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Some error on write data in file: %v", err)
			return
		}

		log.Println("new data has been received and sotred.")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Data stored on pod %s\n", hostname)
		return
	}

	if _, err := os.Stat(dataFile); err != nil {
		if !os.IsNotExist(err) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Some error on check file exist: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "No data posted yet\n")
		return
	}

	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Some error on read file: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "You've hit %s\n", hostname)
	fmt.Fprintf(w, "Data stored on this pod: %s \n", string(data))
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

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT)
	signal.Notify(stop, syscall.SIGTERM)

	<-stop
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println("Stop server")
}
