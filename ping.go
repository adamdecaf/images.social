package main

import (
	"fmt"
	"net/http"
)

func pingRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}
