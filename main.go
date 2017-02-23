package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	port = flag.Int("port", 8080, "Port to bind onto")
)

const (
	LocalFSCachePath = "./cache"
)

func main() {
	flag.Parse()

	// Routes
	http.Handle("/", http.FileServer(http.Dir("./html/")))
	http.HandleFunc("/ping", pingRoute)
	http.Handle("/i/", http.StripPrefix("/i/", http.FileServer(http.Dir(LocalFSCachePath))))

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatalf("error binding to port %d err=%v", port, err)
	}
}
