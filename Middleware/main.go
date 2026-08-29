package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const backendURL = "http://api:8080"

// const usersURL = "http://users:8080"
// const authURL = "http://auth:8080"
// const bookingsURL = "http://bookings:8080"
// const officesURL = "http://offices:8080"

func routeTarget(path string) string {
	switch {
	case strings.HasPrefix(path, "/users"):
		return backendURL
		//return usersURL
	case path == "/login", path == "/register", path == "/update", strings.HasPrefix(path, "/delete"):
		return backendURL
		//return authURL
	case strings.HasPrefix(path, "/offices"):
		return backendURL
		//return officesURL
	case strings.HasPrefix(path, "/bookings"):
		return backendURL
		//return bookingsURL
	default:
		return backendURL
		//return backendURL
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	target := routeTarget(r.URL.Path)
	forwardURL := target + r.URL.Path
	if r.URL.RawQuery != "" {
		forwardURL += "?" + r.URL.RawQuery
	}

	forwardReq, err := http.NewRequest(r.Method, forwardURL, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to build upstream request: %v", err), http.StatusBadGateway)
		return
	}
	forwardReq.Header = make(http.Header)
	for key, values := range r.Header {
		for _, value := range values {
			forwardReq.Header.Add(key, value)
		}
	}
	forwardReq.Host = "api"

	resp, err := http.DefaultClient.Do(forwardReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to reach backend service: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("failed to write response body: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirectHandler)

	log.Println("Middleware running on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}
