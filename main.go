package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	port = flag.Int("port", 8080, "Port to bind onto")
)

const (
	LocalFSCachePath = "./cache"
)

func main() {
	flag.Parse()

	if err := createCache(); err != nil {
		log.Fatalf("error creating cache dir, err=%v", err)
	}

	// Routes
	http.Handle("/", http.FileServer(http.Dir("./html/")))
	http.HandleFunc("/ping", pingRoute)
	http.Handle("/i/", http.StripPrefix("/i/", http.FileServer(http.Dir(LocalFSCachePath))))
	http.HandleFunc("/upload", uploadRoute)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("error binding to port %d err=%v", *port, err)
	}
}

func createCache() error {
	_, err := os.Stat(LocalFSCachePath)
	if err == nil {
		return nil
	}
	err = os.Mkdir(LocalFSCachePath, 0744)
	if err != nil {
		return err
	}
	return nil
}
