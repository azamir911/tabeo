package main

import (
	"log"
	"net/http"
)

func main() {
	// Define routes and start the server
	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
