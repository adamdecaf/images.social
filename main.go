package main

import (
	"flag"
	"fmt"
	"github.com/adamdecaf/images.social/cache"
	"log"
	"net/http"
)

var (
	port = flag.Int("port", 8080, "Port to bind onto")
)

func main() {
	flag.Parse()

	if err := cache.Init(); err != nil {
		log.Fatalf("error creating cache, err=%v", err)
	}

	http.Handle("/", http.FileServer(http.Dir("./html/")))
	http.HandleFunc("/ping", pingRoute)
	http.Handle("/i/", http.StripPrefix("/i/", http.FileServer(http.Dir(cache.Dir()))))
	http.HandleFunc("/upload", uploadRoute)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatalf("error binding to port %d err=%v", *port, err)
	}
}
