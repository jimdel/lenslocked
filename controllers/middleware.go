package controllers

import (
	"fmt"
	"net/http"
	"time"
)

func Performance(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip := r.RemoteAddr
		h(w, r)
		fmt.Printf("Request from %v took: %v\n", ip, time.Since(start))
	}
}
