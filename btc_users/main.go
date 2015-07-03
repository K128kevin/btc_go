package main

import (
	"log"
	"net/http"
)

func main() {
	// globalSessions = NewManager("memory", "gosessionid", 3600)
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}