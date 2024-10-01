package main

import (
	"log"
	"net/http"
)

func main() {
	proxy := &Proxy{}

	// Start HTTP proxy server
	log.Println("Starting proxy server on :8080")
	log.Fatal(http.ListenAndServe(":8080", proxy))
}
