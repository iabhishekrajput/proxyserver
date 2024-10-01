package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

// Handle HTTPS (CONNECT method) requests
func (p *Proxy) handleHTTPS(w http.ResponseWriter, r *http.Request) {
	// Establish a TCP connection to the target
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second) // Added timeout for connection
	if err != nil {
		http.Error(w, "Unable to connect to destination", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)

	// Hijack the client connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, "Hijacking failed", http.StatusServiceUnavailable)
		return
	}

	// Relay data between client and destination with pooled connections
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

// Handle HTTP requests (non-CONNECT)
func (p *Proxy) handleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Host == "localhost:8080" || r.Host == "127.0.0.1:8080" {
		http.Error(w, "Proxying to itself is not allowed", http.StatusBadRequest)
		return
	}

	// Modify request for proxying
	r.RequestURI = ""
	r.URL.Scheme = "http"
	r.URL.Host = r.Host

	// Set up a custom transport with connection pooling
	transport := &http.Transport{
		MaxIdleConns:        100,              // Maximum idle connections
		MaxIdleConnsPerHost: 10,               // Max idle connections per host
		IdleConnTimeout:     90 * time.Second, // Timeout for idle connections
		DisableCompression:  true,             // Disable compression
	}

	// Set up a reverse proxy with custom transport
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = r.URL.Scheme
			req.URL.Host = r.URL.Host
			req.URL.Path = r.URL.Path
		},
		Transport: transport,
	}

	// Log the request and proxy it
	log.Printf("Proxying request to: %s", r.URL.String())
	proxy.ServeHTTP(w, r)
}
