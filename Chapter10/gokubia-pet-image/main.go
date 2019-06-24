package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	dataFilePath = flag.String("df", "/var/data/gokubia.txt", "data file path")
	domain       = flag.String("domain", "gokubia.default.svc.cluster.local", "data file path")

	dataFile = "/var/data/gokubia.txt"
)

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
			fmt.Fprintf(w, "Some error on read request body: %v\n", err)
			return
		}
		defer r.Body.Close()

		file, err := os.OpenFile(dataFile, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Some error on open file: %v\n", err)
			return
		}
		defer file.Close()

		if _, err := fmt.Fprintln(file, string(data)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Some error on write data in file: %v\n", err)
			return
		}

		log.Println("new data has been received and sotred.")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Data stored on pod %s\n", hostname)
		return
	}

	w.WriteHeader(http.StatusOK)
	if r.URL.Path == "/data" {
		data, err := ioutil.ReadFile(dataFile)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Some error on read file: %v\n", err)
			return
		}

		fmt.Fprintf(w, "You've hit %s\n", hostname)
		fmt.Fprintf(w, "Data stored on this pod: %s \n", string(data))
		return
	}
	fmt.Fprintf(w, "You've hit %s\n", hostname)
	fmt.Fprintf(w, "Data stored on this cluster:\n")

	cname, addresses, err := net.LookupSRV("http", "tcp", *domain)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Some error on lookup DNS SRV: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Cname value: %s\n", cname)
	if len(addresses) == 0 {
		fmt.Fprintf(w, "No peers descovered\n")
		return
	}

	for _, item := range addresses {
		host := strings.Trim(item.Target, ".")
		fmt.Fprintf(w, "target: %s port: %v priority: %s weight: %v\n", item.Target, item.Port, item.Priority, item.Weight)
		fmt.Fprintf(w, "REQUEST TO: http://%s:%d/data\n%+v\n", host, item.Port, *item)
		resp, err := http.Get(fmt.Sprintf("http://%s:%d/data", host, item.Port))
		if err != nil {
			fmt.Fprintf(w, "Error on request data: %v\n", err)
			return
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintf(w, "Error on read responce data: %v\n", err)
			return
		}
		resp.Body.Close()
		fmt.Fprintf(w, "- %s: %s\n", item.Target, string(data))
	}
}

func main() {
	flag.Parse()
	dataFile = *dataFilePath

	if _, err := os.Stat(dataFile); err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		if _, err := os.Create(dataFile); err != nil {
			panic(err)
		}
	}

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
