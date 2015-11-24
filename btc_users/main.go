package main

import (
	"log"
	"net/http"
)

func main() {
	go StartSessionUpdates()
	go ThrottleLoginAttempts()
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}