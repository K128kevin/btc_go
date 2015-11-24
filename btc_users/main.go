package main

import (
	"log"
	"net/http"
)

func main() {
	go startUpdates()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}