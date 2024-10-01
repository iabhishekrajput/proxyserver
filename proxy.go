package main

import (
	"net/http"
)

// Proxy struct
type Proxy struct{}

// ServeHTTP is the entry point to handle HTTP and HTTPS requests
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleHTTPS(w, r)
	} else {
		p.handleHTTP(w, r)
	}
}
